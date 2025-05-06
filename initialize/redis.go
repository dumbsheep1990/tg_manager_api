package initialize

import (
	"context"
	"os"
	"tg_manager_api/global"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// InitRedis 初始化Redis连接
func InitRedis() {
	redisConfig := global.Config.Redis
	
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Error("Redis连接失败", zap.Error(err))
		os.Exit(1)
	}
	
	global.Logger.Info("Redis连接成功")
	global.Redis = client
}
