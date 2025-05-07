package account

import (
	"tg_manager_api/services/account/service"
)

// ServiceFactory 账号服务工厂
type ServiceFactory struct{}

// NewServiceFactory 创建服务工厂实例
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{}
}

// AccountService 获取账号服务实例
func (f *ServiceFactory) AccountService() service.AccountService {
	return service.NewAccountService()
}
