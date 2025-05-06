package initialize

import (
	"tg_manager_api/global"
	"tg_manager_api/model"
)

// initMigrate handles database migrations for all models
func initMigrate() {
	err := global.DB.AutoMigrate(
		// System models
		&model.Account{},
		&model.AccountGroup{},
	)
	
	if err != nil {
		global.DB = nil
		panic("migrate table failed")
	}
}
