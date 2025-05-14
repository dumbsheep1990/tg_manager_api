package initialize

import (
	"fmt"
	
	"tg_manager_api/global"
	"tg_manager_api/services/rabbitmq"
	"tg_manager_api/services/task/scheduler"
	taskService "tg_manager_api/services/task/service"
	workerService "tg_manager_api/services/worker/service"
)

// InitTaskScheduler 初始化任务调度器
func InitTaskScheduler() *scheduler.TaskScheduler {
	global.LOG.Info("正在初始化任务调度器...")
	
	// 创建RabbitMQ连接
	rabbitMQService, err := rabbitmq.NewRabbitMQService(global.Config.RabbitMQ.URL)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}
	
	// 获取服务实例
	taskSvc := taskService.NewTaskService()
	workerSvc := workerService.NewWorkerService()
	
	// 创建任务调度器
	taskScheduler := scheduler.NewTaskScheduler(taskSvc, workerSvc, rabbitMQService)
	
	// 启动调度器
	err = taskScheduler.Start()
	if err != nil {
		global.LOG.Error(fmt.Sprintf("任务调度器启动失败: %v", err))
		return nil
	}
	
	global.LOG.Info("任务调度器初始化完成")
	return taskScheduler
}
