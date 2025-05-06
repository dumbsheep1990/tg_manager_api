package initialize

import (
	"fmt"
	"os"
	"tg_manager_api/global"
	
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig() {
	var config string
	// 默认使用config.toml作为配置文件
	config = "config.toml"
	
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("toml")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("读取配置文件失败: %s \n", err)
		os.Exit(1)
	}
	
	// 解析配置到结构体
	if err := v.Unmarshal(&global.Config); err != nil {
		fmt.Printf("解析配置文件失败: %s \n", err)
		os.Exit(1)
	}
	
	// 初始化日志
	InitLogger()
}
