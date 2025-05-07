package router

import (
	"github.com/gin-gonic/gin"
	
	v1 "tg_manager_api/api/v1"
	"tg_manager_api/api/v1/account_group"
	"tg_manager_api/api/v1/tdata_account"
)

// InitTgManagerRouter 初始化Telegram管理相关路由
func InitTgManagerRouter(Router *gin.RouterGroup) {
	// 账号分组管理路由
	accountGroupRouter := Router.Group("account-groups")
	{
		accountGroupRouter.POST("", account_group.CreateAccountGroup)        // 创建账号分组
		accountGroupRouter.GET("", account_group.ListAccountGroups)          // 获取账号分组列表
		accountGroupRouter.GET("/:id", account_group.GetAccountGroup)        // 获取账号分组详情
		accountGroupRouter.PUT("/:id", account_group.UpdateAccountGroup)     // 更新账号分组
		accountGroupRouter.DELETE("/:id", account_group.DeleteAccountGroup)  // 删除账号分组
	}

	// tdata账号管理路由
	tdataAccountRouter := Router.Group("tdata-accounts")
	{
		tdataAccountRouter.POST("/import", tdata_account.ImportTdataAccount)     // 导入tdata账号
		tdataAccountRouter.GET("", tdata_account.ListTdataAccounts)              // 获取账号列表
		tdataAccountRouter.GET("/:id", tdata_account.GetTdataAccount)            // 获取账号详情
		tdataAccountRouter.PUT("/:id", tdata_account.UpdateTdataAccount)         // 更新账号信息
		tdataAccountRouter.DELETE("/:id", tdata_account.DeleteTdataAccount)      // 删除账号
	}
}
