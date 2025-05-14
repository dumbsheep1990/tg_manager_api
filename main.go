package main

import (
	"os"
	"os/signal"
	"syscall"
	"tg_manager_api/core"
	"tg_manager_api/global"
	"tg_manager_api/initialize"
	"tg_manager_api/services/core"
	"tg_manager_api/services/task/scheduler"

	"go.uber.org/zap"
)

// 全局任务调度器
var taskScheduler *scheduler.TaskScheduler

func main() {
	// 初始化配置、数据库、缓存等组件
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitRabbitMQ()
	initialize.InitEtcd()
	initialize.InitNacos()
	
	// 初始化服务模块
	core.InitializeServices()
	
	// 初始化任务调度器
	taskScheduler = initialize.InitTaskScheduler()
	if taskScheduler == nil {
		global.Logger.Warn("任务调度器初始化失败，将无法处理分布式任务")
	} else {
		global.Logger.Info("任务调度器启动成功")
	}
	
	// 优雅关闭服务
	go gracefulShutdown()
	
	// 启动服务器
	global.Logger.Info("服务启动完成")
	core.RunServer()
}

// gracefulShutdown 优雅关闭服务
func gracefulShutdown() {
	// 监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	
	<-c
	
	// 收到退出信号后执行清理操作
	global.Logger.Info("收到退出信号，正在关闭服务...")
	
	// 停止任务调度器
	if taskScheduler != nil {
		taskScheduler.Stop()
		global.Logger.Info("任务调度器已停止")
	}
	
	// 从nacos注销服务
	initialize.DeregisterService()
	
	// 关闭etcd连接
	if global.EtcdClient != nil {
		global.EtcdClient.Close()
	}
	
	// 关闭RabbitMQ连接
	if global.RabbitMQChannel != nil {
		global.RabbitMQChannel.Close()
	}
	if global.RabbitMQConn != nil {
		global.RabbitMQConn.Close()
	}
	
	// 关闭Redis连接
	if global.Redis != nil {
		global.Redis.Close()
	}
	
	// 关闭数据库连接
	if global.DB != nil {
		db, _ := global.DB.DB()
		db.Close()
	}
	
	global.Logger.Info("服务已安全关闭")
	os.Exit(0)
}
