package service

import (
	"context"
	"time"
	
	"tg_manager_api/global"
	"tg_manager_api/model"
)

// WorkerServiceI worker服务接口
type WorkerServiceI interface {
	// 注册工作节点
	RegisterWorker(ctx context.Context, hostname, ip string, maxTasks int, tags string) (string, error)
	
	// 更新工作节点心跳
	UpdateHeartbeat(ctx context.Context, workerID string) error
	
	// 获取可用的工作节点
	GetAvailableWorker(ctx context.Context, tags string) (*model.Worker, error)
	
	// 分配任务到工作节点
	AssignTaskToWorker(ctx context.Context, taskID, workerID string) error
	
	// 获取工作节点状态
	GetWorkerStatus(ctx context.Context, workerID string) (*model.Worker, error)
	
	// 获取所有工作节点
	GetAllWorkers(ctx context.Context, page, pageSize int) ([]*model.Worker, int64, error)
	
	// 获取工作节点执行的任务列表
	GetWorkerTasks(ctx context.Context, workerID string, page, pageSize int) ([]*model.Task, int64, error)
}

// NewWorkerService 创建worker服务实例
func NewWorkerService() WorkerServiceI {
	return &workerService{}
}

// workerService worker服务实现
type workerService struct{}

// RegisterWorker 注册工作节点
func (s *workerService) RegisterWorker(ctx context.Context, hostname, ip string, maxTasks int, tags string) (string, error) {
	// 检查是否已存在相同hostname和IP的节点
	var existingWorker model.Worker
	result := global.DB.Where("hostname = ? AND ip = ?", hostname, ip).First(&existingWorker)
	
	// 如果找到现有节点且状态为offline，则重用该节点
	if result.Error == nil {
		existingWorker.Status = "online"
		existingWorker.LastHeartbeat = time.Now()
		existingWorker.MaxTasks = maxTasks
		existingWorker.CurrentTasks = 0
		existingWorker.Tags = tags
		
		if err := global.DB.Save(&existingWorker).Error; err != nil {
			return "", err
		}
		
		return existingWorker.WorkerID, nil
	}
	
	// 创建新工作节点
	workerID := generateWorkerID()
	worker := model.Worker{
		WorkerID:      workerID,
		Hostname:      hostname,
		IP:            ip,
		Status:        "online",
		LastHeartbeat: time.Now(),
		MaxTasks:      maxTasks,
		CurrentTasks:  0,
		Tags:          tags,
	}
	
	if err := global.DB.Create(&worker).Error; err != nil {
		return "", err
	}
	
	return workerID, nil
}

// UpdateHeartbeat 更新工作节点心跳
func (s *workerService) UpdateHeartbeat(ctx context.Context, workerID string) error {
	// 更新工作节点心跳时间
	result := global.DB.Model(&model.Worker{}).
		Where("worker_id = ?", workerID).
		Updates(map[string]interface{}{
			"last_heartbeat": time.Now(),
			"status":         "online",
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return global.ErrorWorkerNotFound
	}
	
	return nil
}

// GetAvailableWorker 获取可用的工作节点
func (s *workerService) GetAvailableWorker(ctx context.Context, tags string) (*model.Worker, error) {
	var worker model.Worker
	
	// 查询条件：状态为online且当前任务数小于最大任务数
	query := global.DB.Where("status = ? AND current_tasks < max_tasks", "online")
	
	// 如果指定了标签，则按标签筛选
	if tags != "" {
		query = query.Where("FIND_IN_SET(?, tags) > 0", tags)
	}
	
	// 按当前任务负载排序，选择负载最小的节点
	result := query.Order("current_tasks ASC").First(&worker)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &worker, nil
}

// AssignTaskToWorker 分配任务到工作节点
func (s *workerService) AssignTaskToWorker(ctx context.Context, taskID, workerID string) error {
	// 创建任务分配记录
	assignment := model.TaskAssignment{
		TaskID:     taskID,
		WorkerID:   workerID,
		Status:     "pending",
		AssignedAt: time.Now(),
		Priority:   0,
	}
	
	if err := global.DB.Create(&assignment).Error; err != nil {
		return err
	}
	
	// 更新工作节点的当前任务数
	if err := global.DB.Model(&model.Worker{}).
		Where("worker_id = ?", workerID).
		Update("current_tasks", global.DB.Raw("current_tasks + 1")).
		Error; err != nil {
		return err
	}
	
	return nil
}

// GetWorkerStatus 获取工作节点状态
func (s *workerService) GetWorkerStatus(ctx context.Context, workerID string) (*model.Worker, error) {
	var worker model.Worker
	
	if err := global.DB.Where("worker_id = ?", workerID).First(&worker).Error; err != nil {
		return nil, err
	}
	
	return &worker, nil
}

// GetAllWorkers 获取所有工作节点
func (s *workerService) GetAllWorkers(ctx context.Context, page, pageSize int) ([]*model.Worker, int64, error) {
	var workers []*model.Worker
	var total int64
	
	// 获取总记录数
	if err := global.DB.Model(&model.Worker{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := global.DB.Offset(offset).Limit(pageSize).Find(&workers).Error; err != nil {
		return nil, 0, err
	}
	
	return workers, total, nil
}

// GetWorkerTasks 获取工作节点执行的任务列表
func (s *workerService) GetWorkerTasks(ctx context.Context, workerID string, page, pageSize int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64
	
	// 获取总记录数
	if err := global.DB.Model(&model.Task{}).
		Joins("JOIN task_assignments ON task_assignments.task_id = tasks.task_id").
		Where("task_assignments.worker_id = ?", workerID).
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := global.DB.
		Joins("JOIN task_assignments ON task_assignments.task_id = tasks.task_id").
		Where("task_assignments.worker_id = ?", workerID).
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).
		Error; err != nil {
		return nil, 0, err
	}
	
	return tasks, total, nil
}

// 生成工作节点ID
func generateWorkerID() string {
	return "wrk_" + time.Now().Format("20060102150405") + randomString(8)
}

// 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // 确保随机性
	}
	return string(result)
}
