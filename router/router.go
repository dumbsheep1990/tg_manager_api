package router

import (
	"github.com/gin-gonic/gin"
	"tg_manager_api/middleware"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	router := gin.Default()

	// 使用中间件
	router.Use(middleware.Cors()) // 跨域中间件
	
	// API v1 版本路由组
	apiV1 := router.Group("/api/v1")
	
	// 初始化Telegram管理相关路由
	InitTgManagerRouter(apiV1)
	
	// 初始化任务和工作节点相关路由
	InitTaskWorkerRouter(apiV1)

	return router
}
