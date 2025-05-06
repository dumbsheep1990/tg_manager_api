package initialize

import (
	"os"
	"tg_manager_api/global"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB() {
	m := global.Config.MySQL
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
	
	var logMode logger.LogLevel
	if m.LogMode {
		logMode = logger.Info
	} else {
		logMode = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	
	if err != nil {
		global.Logger.Error("MySQL连接失败", zap.Error(err))
		os.Exit(1)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	global.DB = db

	// Auto migrate the database
	initMigrate()
}
