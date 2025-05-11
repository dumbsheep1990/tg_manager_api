package task

import (
	"sync"
	
	"github.com/gin-gonic/gin"
	
	"tg_manager_api/services/task/service"
)

var (
	taskServiceInstance service.TaskServiceI
	once                sync.Once
)

// GetTaskService 返回任务服务的单例实例
func GetTaskService() service.TaskServiceI {
	once.Do(func() {
		taskServiceInstance = service.NewTaskService()
	})
	return taskServiceInstance
}

// InjectTaskService 将任务服务注入到gin上下文中
func InjectTaskService(c *gin.Context) {
	c.Set("taskService", GetTaskService())
	c.Next()
}

// GetTaskServiceFromContext 从gin上下文中检索任务服务
func GetTaskServiceFromContext(c *gin.Context) service.TaskServiceI {
	return c.MustGet("taskService").(service.TaskServiceI)
}
