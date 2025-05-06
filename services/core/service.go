package core

import (
	"tg_manager_api/global"
	accountService "tg_manager_api/services/account/service"
	dashboardService "tg_manager_api/services/dashboard/service"
	messageService "tg_manager_api/services/message/service"
	taskService "tg_manager_api/services/task/service"
	tdataService "tg_manager_api/services/tdata/service"
)

// ServiceGroup 服务集合
type ServiceGroup struct {
	AccountService   accountService.AccountService
	TaskService      taskService.TaskService
	TdataService     tdataService.TdataService
	MessageService   messageService.MessageService
	DashboardService dashboardService.DashboardService
}

// ServiceGroupApp 全局服务组实例
var ServiceGroupApp = new(ServiceGroup)

// InitializeServices 初始化所有服务
func InitializeServices() {
	global.Logger.Info("初始化服务组件...")
	
	// 初始化各服务模块
	ServiceGroupApp.AccountService = accountService.NewAccountService()
	ServiceGroupApp.TaskService = taskService.NewTaskService()
	ServiceGroupApp.TdataService = tdataService.NewTdataService()
	ServiceGroupApp.MessageService = messageService.NewMessageService()
	ServiceGroupApp.DashboardService = dashboardService.NewDashboardService()
	
	global.Logger.Info("服务组件初始化完成")
}
