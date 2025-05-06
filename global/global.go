package global

import (
	"tg_manager_api/config"
	
	"github.com/go-redis/redis/v8"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/streadway/amqp"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config         config.Configuration
	DB             *gorm.DB
	Redis          *redis.Client
	RabbitMQConn   *amqp.Connection
	RabbitMQChannel *amqp.Channel
	Logger         *zap.Logger       // zap日志
	EtcdClient     *clientv3.Client   // etcd客户端
	NacosClient    naming_client.INamingClient // nacos客户端
)
