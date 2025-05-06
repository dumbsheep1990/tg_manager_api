package core

import (
	"fmt"
	"net/http"
	"time"
	"tg_manager_api/global"
	"tg_manager_api/initialize"
	"tg_manager_api/router"

	"go.uber.org/zap"
)

// RunServer 启动HTTP服务器
func RunServer() {
	Router := router.InitRouter()
	
	address := fmt.Sprintf(":%d", global.Config.System.Port)
	server := &http.Server{
		Addr:           address,
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	// 初始化消息消费者
	go initialize.InitResultConsumer()
	
	// 启动服务器
	global.Logger.Info("服务器启动成功", zap.String("address", address))
	if err := server.ListenAndServe(); err != nil {
		global.Logger.Error("服务器启动失败", zap.Error(err))
	}
}
