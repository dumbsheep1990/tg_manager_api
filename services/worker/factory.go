package worker

import (
	"sync"
	
	"github.com/gin-gonic/gin"
	
	"tg_manager_api/services/worker/service"
)

var (
	workerServiceInstance service.WorkerServiceI
	once                  sync.Once
)

// GetWorkerService 返回工作节点服务的单例实例
func GetWorkerService() service.WorkerServiceI {
	once.Do(func() {
		workerServiceInstance = service.NewWorkerService()
	})
	return workerServiceInstance
}

// InjectWorkerService 将工作节点服务注入到gin上下文中
func InjectWorkerService(c *gin.Context) {
	c.Set("workerService", GetWorkerService())
	c.Next()
}

// GetWorkerServiceFromContext 从gin上下文中检索工作节点服务
func GetWorkerServiceFromContext(c *gin.Context) service.WorkerServiceI {
	return c.MustGet("workerService").(service.WorkerServiceI)
}
