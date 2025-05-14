package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	
	"tg_manager_api/global"
	"tg_manager_api/model"
	"tg_manager_api/services/rabbitmq"
	"tg_manager_api/services/task/service"
	workerSvc "tg_manager_api/services/worker/service"
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	taskService   service.TaskServiceI
	workerService workerSvc.WorkerServiceI
	rabbitMQ      rabbitmq.RabbitMQService
	running       bool
	mutex         sync.Mutex
	stopChan      chan struct{}
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler(taskService service.TaskServiceI, workerService workerSvc.WorkerServiceI, rabbitMQ rabbitmq.RabbitMQService) *TaskScheduler {
	return &TaskScheduler{
		taskService:   taskService,
		workerService: workerService,
		rabbitMQ:      rabbitMQ,
		running:       false,
		stopChan:      make(chan struct{}),
	}
}

// Start 启动任务调度器
func (s *TaskScheduler) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if s.running {
		return fmt.Errorf("task scheduler is already running")
	}
	
	// 启动任务调度循环
	go s.scheduleLoop()
	
	// 启动任务结果处理器
	err := s.startResultProcessor()
	if err != nil {
		return fmt.Errorf("failed to start result processor: %w", err)
	}
	
	s.running = true
	global.LOG.Info("Task scheduler started")
	return nil
}

// Stop 停止任务调度器
func (s *TaskScheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if !s.running {
		return
	}
	
	close(s.stopChan)
	s.running = false
	global.LOG.Info("Task scheduler stopped")
}

// 任务调度循环
func (s *TaskScheduler) scheduleLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.schedulePendingTasks()
		case <-s.stopChan:
			return
		}
	}
}

// 调度待处理任务
func (s *TaskScheduler) schedulePendingTasks() {
	// 获取所有待处理的任务
	var pendingTasks []model.Task
	if err := global.DB.Where("status = ?", "pending").
		Order("priority DESC, created_at ASC").
		Find(&pendingTasks).Error; err != nil {
		global.LOG.Error(fmt.Sprintf("Failed to fetch pending tasks: %v", err))
		return
	}
	
	if len(pendingTasks) == 0 {
		return
	}
	
	// 获取可用的工作节点
	availableWorkers, err := s.workerService.GetAvailableWorkers(context.Background())
	if err != nil || len(availableWorkers) == 0 {
		global.LOG.Warn("No available workers for pending tasks")
		return
	}
	
	// 轮询分配任务
	workerIndex := 0
	for _, task := range pendingTasks {
		// 检查是否有可用工作节点
		if len(availableWorkers) == 0 {
			break
		}
		
		// 选择工作节点
		worker := availableWorkers[workerIndex]
		workerIndex = (workerIndex + 1) % len(availableWorkers)
		
		// 分配任务
		if err := s.assignTaskToWorker(context.Background(), &task, worker.WorkerID); err != nil {
			global.LOG.Error(fmt.Sprintf("Failed to assign task %s to worker %s: %v", 
				task.TaskID, worker.WorkerID, err))
			continue
		}
		
		// 更新工作节点当前任务数
		worker.CurrentTasks++
		
		// 如果工作节点任务数达到最大值，从可用列表移除
		if worker.CurrentTasks >= worker.MaxTasks {
			availableWorkers = append(availableWorkers[:workerIndex], availableWorkers[workerIndex+1:]...)
			if workerIndex > 0 {
				workerIndex--
			}
		}
	}
}

// 分配任务给工作节点
func (s *TaskScheduler) assignTaskToWorker(ctx context.Context, task *model.Task, workerID string) error {
	// 创建任务分配记录
	assignment := model.TaskAssignment{
		TaskID:     task.TaskID,
		WorkerID:   workerID,
		Status:     "assigned",
		AssignedAt: time.Now(),
		Priority:   task.Priority,
	}
	
	if err := global.DB.Create(&assignment).Error; err != nil {
		return fmt.Errorf("failed to create task assignment: %w", err)
	}
	
	// 更新任务状态
	if err := global.DB.Model(task).Updates(map[string]interface{}{
		"status": "assigned",
	}).Error; err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	
	// 更新工作节点当前任务数
	if err := global.DB.Model(&model.Worker{}).
		Where("worker_id = ?", workerID).
		Update("current_tasks", global.DB.Raw("current_tasks + 1")).
		Error; err != nil {
		return fmt.Errorf("failed to update worker task count: %w", err)
	}
	
	// 发送任务到队列
	taskMessage := map[string]interface{}{
		"task_id":    task.TaskID,
		"task_type":  task.TaskType,
		"account_id": task.AccountID,
		"params":     task.Params,
		"worker_id":  workerID,
		"timeout":    task.TimeoutSec,
		"created_at": task.CreatedAt,
	}
	
	// 序列化任务消息
	taskData, err := json.Marshal(taskMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal task message: %w", err)
	}
	
	// 发送到RabbitMQ
	if err := s.rabbitMQ.PublishTask(task.TaskType, taskData); err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}
	
	global.LOG.Info(fmt.Sprintf("Task %s assigned to worker %s", task.TaskID, workerID))
	return nil
}

// 启动任务结果处理器
func (s *TaskScheduler) startResultProcessor() error {
	// 创建任务结果消费者
	handler := func(data []byte) error {
		return s.processTaskResult(data)
	}
	
	if err := s.rabbitMQ.CreateTaskResultConsumer(handler); err != nil {
		return fmt.Errorf("failed to create task result consumer: %w", err)
	}
	
	global.LOG.Info("Task result processor started")
	return nil
}

// 处理任务结果
func (s *TaskScheduler) processTaskResult(data []byte) error {
	// 解析任务结果
	var result struct {
		TaskID      string                 `json:"task_id"`
		AccountID   uint                   `json:"account_id"`
		WorkerID    string                 `json:"worker_id"`
		Status      string                 `json:"status"`
		Result      map[string]interface{} `json:"result"`
		Error       string                 `json:"error"`
		CompletedAt time.Time              `json:"completed_at"`
	}
	
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("failed to unmarshal task result: %w", err)
	}
	
	// 更新任务状态
	status := service.TaskStatusCompleted
	if result.Status == "failed" {
		status = service.TaskStatusFailed
	}
	
	// 序列化结果数据
	resultData, _ := json.Marshal(result.Result)
	resultStr := string(resultData)
	
	// 更新任务状态
	ctx := context.Background()
	if err := s.taskService.UpdateTaskStatus(ctx, result.TaskID, status, resultStr, result.Error); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	
	// 更新任务分配记录
	if err := global.DB.Model(&model.TaskAssignment{}).
		Where("task_id = ? AND worker_id = ?", result.TaskID, result.WorkerID).
		Updates(map[string]interface{}{
			"status":       result.Status,
			"completed_at": result.CompletedAt,
		}).Error; err != nil {
		return fmt.Errorf("failed to update task assignment: %w", err)
	}
	
	// 更新工作节点当前任务数
	if err := global.DB.Model(&model.Worker{}).
		Where("worker_id = ?", result.WorkerID).
		Update("current_tasks", global.DB.Raw("current_tasks - 1")).
		Error; err != nil {
		return fmt.Errorf("failed to update worker task count: %w", err)
	}
	
	global.LOG.Info(fmt.Sprintf("Task %s processing completed with status: %s", result.TaskID, result.Status))
	return nil
}
