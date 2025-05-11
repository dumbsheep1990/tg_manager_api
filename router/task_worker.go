package router

import (
	"github.com/gin-gonic/gin"
	
	"tg_manager_api/api/v1/task"
	"tg_manager_api/api/v1/worker"
	taskService "tg_manager_api/services/task"
	workerService "tg_manager_api/services/worker"
)

// InitTaskWorkerRouter 初始化任务和工作节点相关路由
func InitTaskWorkerRouter(Router *gin.RouterGroup) {
	// 注册服务中间件
	Router.Use(taskService.InjectTaskService)
	Router.Use(workerService.InjectWorkerService)
	
	// 实例化控制器
	taskController := task.TaskController{}
	workerController := worker.WorkerController{}
	
	// 任务管理路由
	taskRouter := Router.Group("tasks")
	{
		taskRouter.POST("", taskController.CreateTask)                         // 创建任务
		taskRouter.GET("", taskController.GetTaskList)                         // 获取任务列表
		taskRouter.GET("/:id", taskController.GetTaskDetail)                   // 获取任务详情
		taskRouter.POST("/:id/cancel", taskController.CancelTask)              // 取消任务
	}
	
	// 账号关联任务路由
	Router.GET("/accounts/:account_id/tasks", taskController.GetTasksByAccount) // 获取账号关联的任务
	
	// 工作节点管理路由
	workerRouter := Router.Group("workers")
	{
		workerRouter.POST("/register", workerController.RegisterWorker)        // 注册工作节点
		workerRouter.POST("/heartbeat", workerController.Heartbeat)            // 工作节点心跳
		workerRouter.GET("", workerController.GetWorkerList)                   // 获取工作节点列表
		workerRouter.GET("/:id", workerController.GetWorkerDetail)             // 获取工作节点详情
		workerRouter.GET("/:id/tasks", workerController.GetWorkerTasks)        // 获取工作节点任务列表
	}
}
