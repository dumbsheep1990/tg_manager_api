package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
	
	"tg_manager_api/global"
	"tg_manager_api/model"
	"tg_manager_api/services/rabbitmq"
	"tg_manager_api/services/worker/service"
	"tg_manager_api/utils"
)

// TaskServiceI 任务服务接口
type TaskServiceI interface {
	// 创建任务
	CreateTask(ctx context.Context, taskType string, accountID uint, params map[string]interface{}) (string, error)
	
	// 获取任务列表
	GetTasks(ctx context.Context, page, pageSize int) ([]*Task, int64, error)
	
	// 获取任务详情
	GetTask(ctx context.Context, taskID string) (*Task, error)
	
	// 更新任务状态
	UpdateTaskStatus(ctx context.Context, taskID string, status TaskStatus, result string, errorMsg string) error
	
	// 取消任务
	CancelTask(ctx context.Context, taskID string) error
	
	// 获取任务日志
	GetTaskLogs(ctx context.Context, taskID string, page, pageSize int) ([]string, int64, error)
}

// taskService 任务服务实现
type taskService struct{}

// taskServiceImpl 任务服务实现
type taskServiceImpl struct {
	workerService service.WorkerServiceI
}

// CreateTask 创建任务
func (s *taskServiceImpl) CreateTask(ctx context.Context, taskType string, accountID uint, params map[string]interface{}) (*model.Task, error) {
	// 检查账号是否存在
	var account model.Account
	if err := global.DB.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, global.ErrorAccountNotFound
		}
		return nil, err
	}
	
	// 生成任务ID
	taskID := generateTaskID()
	
	// 创建任务记录
	task := &model.Task{
		TaskID:    taskID,
		TaskType:  taskType,
		AccountID: accountID,
		Params:    params,
		Status:    "pending",
		Priority:  0, // 默认优先级
		TimeoutSec: 300, // 默认5分钟超时
	}
	
	// 保存到数据库
	if err := global.DB.Create(task).Error; err != nil {
		return nil, err
	}
	
	// 尝试分配任务给可用的工作节点
	go s.AssignTask(context.Background(), task)
	
	return task, nil
}

// GetTasks 获取任务列表
func (s *taskServiceImpl) GetTasks(ctx context.Context, page, pageSize int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64
	
	// 获取总记录数
	if err := global.DB.Model(&model.Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := global.DB.Preload("Account").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	
	return tasks, total, nil
}

// GetTask 获取任务详情
func (s *taskServiceImpl) GetTask(ctx context.Context, taskID string) (*model.Task, error) {
	var task model.Task
	
	if err := global.DB.Preload("Account").
		Preload("TaskRecords").
		Where("task_id = ?", taskID).
		First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, global.ErrorTaskNotFound
		}
		return nil, err
	}
	
	return &task, nil
}

// UpdateTaskStatus 更新任务状态
func (s *taskServiceImpl) UpdateTaskStatus(ctx context.Context, taskID, status string, result map[string]interface{}, errorMsg string) error {
	// 查找任务
	var task model.Task
	if err := global.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return global.ErrorTaskNotFound
		}
		return err
	}
	
	// 更新状态
	updateFields := map[string]interface{}{
		"status": status,
	}
	
	// 根据状态设置其他字段
	switch status {
	case "processing":
		now := time.Now()
		updateFields["started_at"] = now
	case "completed", "failed":
		now := time.Now()
		updateFields["completed_at"] = now
		if errorMsg != "" {
			updateFields["error_message"] = errorMsg
		}
	}
	
	// 更新数据库
	if err := global.DB.Model(&task).Updates(updateFields).Error; err != nil {
		return err
	}
	
	return nil
}

// CancelTask 取消任务
func (s *taskServiceImpl) CancelTask(ctx context.Context, taskID string) error {
	// 查找任务
	var task model.Task
	if err := global.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return global.ErrorTaskNotFound
		}
		return err
	}
	
	// 只有待处理或处理中的任务可以取消
	if task.Status != "pending" && task.Status != "processing" {
		return global.ErrorInvalidTaskStatus
	}
	
	// 更新任务状态为取消
	if err := global.DB.Model(&task).Updates(map[string]interface{}{
		"status":        "canceled",
		"completed_at":  time.Now(),
		"error_message": "Task canceled by user",
	}).Error; err != nil {
		return err
	}
	
	// 如果已分配给工作节点，还需要通知工作节点取消任务
	// 这里简化处理，仅更新数据库状态
	// 实际实现中应该通过消息队列通知工作节点
	
	return nil
}

// ProcessTaskResult 处理任务结果
func (s *taskServiceImpl) ProcessTaskResult(ctx context.Context, taskID string, result map[string]interface{}, success bool, errorMsg string) error {
	// 查找任务
	var task model.Task
	if err := global.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return global.ErrorTaskNotFound
		}
		return err
	}
	
	// 查找任务分配记录，获取worker信息
	var assignment model.TaskAssignment
	if err := global.DB.Where("task_id = ?", taskID).First(&assignment).Error; err != nil {
		// 如果没有分配记录，仅更新任务状态
		status := "completed"
		if !success {
			status = "failed"
		}
		
		return s.UpdateTaskStatus(ctx, taskID, status, result, errorMsg)
	}
	
	// 创建任务记录
	now := time.Now()
	taskRecord := model.TaskRecord{
		TaskID:       taskID,
		WorkerID:     assignment.WorkerID,
		Status:       "completed",
		Result:       result,
		ErrorMessage: errorMsg,
		StartedAt:    assignment.AssignedAt,
		CompletedAt:  &now,
	}
	
	// 计算执行时间（毫秒）
	taskRecord.ExecutionTime = int(now.Sub(assignment.AssignedAt).Milliseconds())
	
	// 保存任务记录
	if err := global.DB.Create(&taskRecord).Error; err != nil {
		return err
	}
	
	// 更新任务状态
	status := "completed"
	if !success {
		status = "failed"
	}
	
	// 更新任务分配状态
	if err := global.DB.Model(&assignment).Updates(map[string]interface{}{
		"status":       status,
		"completed_at": now,
	}).Error; err != nil {
		return err
	}
	
	// 更新工作节点任务计数
	if err := global.DB.Model(&model.Worker{}).
		Where("worker_id = ?", assignment.WorkerID).
		Update("current_tasks", gorm.Expr("current_tasks - 1")).
		Error; err != nil {
		return err
	}
	
	// 更新任务状态
	return s.UpdateTaskStatus(ctx, taskID, status, result, errorMsg)
}

// AssignTask 分配任务
func (s *taskServiceImpl) AssignTask(ctx context.Context, task *model.Task) error {
	// 获取可用的工作节点
	worker, err := s.workerService.GetAvailableWorker(ctx, "")
	if err != nil {
		// 如果没有可用的工作节点，将任务状态设置为pending，等待后续分配
		return global.ErrorNoAvailableWorker
	}
	
	// 分配任务给工作节点
	if err := s.workerService.AssignTaskToWorker(ctx, task.TaskID, worker.WorkerID); err != nil {
		return err
	}
	
	// 更新任务状态为assigned
	if err := s.UpdateTaskStatus(ctx, task.TaskID, "assigned", nil, ""); err != nil {
		return err
	}
	
	// 发送任务到消息队列
	if err := s.sendTaskToQueue(task, worker.WorkerID); err != nil {
		// 如果发送失败，更新任务状态
		s.UpdateTaskStatus(ctx, task.TaskID, "failed", nil, fmt.Sprintf("Failed to send task to queue: %v", err))
		return err
	}
	
	return nil
}

// 发送任务到消息队列
func (s *taskServiceImpl) sendTaskToQueue(task *model.Task, workerID string) error {
	// 创建RabbitMQ连接
	rabbitmqService, err := rabbitmq.NewRabbitMQService(global.Config.RabbitMQ.URL)
	if err != nil {
		return err
	}
	defer rabbitmqService.Close()
	
	// 准备任务消息
	taskMessage := map[string]interface{}{
		"task_id":    task.TaskID,
		"task_type":  task.TaskType,
		"account_id": task.AccountID,
		"params":     task.Params,
		"worker_id":  workerID,
		"timeout":    task.TimeoutSec,
		"created_at": task.CreatedAt,
	}
	
	// 发送任务到队列
	return rabbitmqService.PublishMessage(
		global.Config.RabbitMQ.Exchange.Tasks,         // 交换机
		fmt.Sprintf("task.%s", task.TaskType),  // 路由键
		taskMessage,                           // 消息体
	)
}

// 生成任务ID
func generateTaskID() string {
	return fmt.Sprintf("task_%s", uuid.New().String())
}
