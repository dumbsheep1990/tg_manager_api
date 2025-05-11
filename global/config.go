package global

import (
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

// Config 全局配置
var (
	DB     *gorm.DB
	Config SystemConfig
	RabbitMQChannel *amqp.Channel
)

// SystemConfig 系统配置
type SystemConfig struct {
	System   SystemSettings   `mapstructure:"system" json:"system"`
	Database DatabaseSettings `mapstructure:"database" json:"database"`
	RabbitMQ RabbitMQSettings `mapstructure:"rabbitmq" json:"rabbitmq"`
	Redis    RedisSettings    `mapstructure:"redis" json:"redis"`
	Etcd     EtcdSettings     `mapstructure:"etcd" json:"etcd"`
	Nacos    NacosSettings    `mapstructure:"nacos" json:"nacos"`
	Log      LogSettings      `mapstructure:"log" json:"log"`
}

// SystemSettings 系统设置
type SystemSettings struct {
	Env      string `mapstructure:"env" json:"env"`           // 环境
	Port     int    `mapstructure:"port" json:"port"`         // 端口
	DbType   string `mapstructure:"db-type" json:"dbType"`    // 数据库类型
	UseRedis bool   `mapstructure:"use-redis" json:"useRedis"` // 使用Redis
}

// DatabaseSettings 数据库设置
type DatabaseSettings struct {
	Host         string `mapstructure:"host" json:"host"`
	Port         string `mapstructure:"port" json:"port"`
	Config       string `mapstructure:"config" json:"config"`
	Dbname       string `mapstructure:"db-name" json:"dbname"`
	Username     string `mapstructure:"username" json:"username"`
	Password     string `mapstructure:"password" json:"password"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns"`
	LogMode      bool   `mapstructure:"log-mode" json:"logMode"`
	LogZap       bool   `mapstructure:"log-zap" json:"logZap"`
}

// RabbitMQSettings RabbitMQ设置
type RabbitMQSettings struct {
	URL      string        `mapstructure:"url" json:"url"`
	Exchange ExchangeSettings `mapstructure:"exchange" json:"exchange"`
	Queue    QueueSettings    `mapstructure:"queue" json:"queue"`
}

// ExchangeSettings 交换机设置
type ExchangeSettings struct {
	Tasks      string `mapstructure:"tasks" json:"tasks"`           // 任务交换机
	TaskResults string `mapstructure:"task-results" json:"taskResults"` // 任务结果交换机
	System     string `mapstructure:"system" json:"system"`         // 系统交换机
}

// QueueSettings 队列设置
type QueueSettings struct {
	Tasks      string `mapstructure:"tasks" json:"tasks"`           // 任务队列
	TaskResults string `mapstructure:"task-results" json:"taskResults"` // 任务结果队列
	System     string `mapstructure:"system" json:"system"`         // 系统队列
}

// RedisSettings Redis设置
type RedisSettings struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

// EtcdSettings Etcd设置
type EtcdSettings struct {
	Endpoints []string `mapstructure:"endpoints" json:"endpoints"`
	Username  string   `mapstructure:"username" json:"username"`
	Password  string   `mapstructure:"password" json:"password"`
}

// NacosSettings Nacos设置
type NacosSettings struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}

// LogSettings 日志设置
type LogSettings struct {
	Level      string `mapstructure:"level" json:"level"`
	Path       string `mapstructure:"path" json:"path"`
	MaxSize    int    `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int    `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int    `mapstructure:"max-age" json:"maxAge"`
	Compress   bool   `mapstructure:"compress" json:"compress"`
}
