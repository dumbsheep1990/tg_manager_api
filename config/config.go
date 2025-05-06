package config

// Configuration 系统全局配置结构
type Configuration struct {
	System   System   `mapstructure:"system" json:"system" toml:"system"`
	MySQL    MySQL    `mapstructure:"mysql" json:"mysql" toml:"mysql"`
	Redis    Redis    `mapstructure:"redis" json:"redis" toml:"redis"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq" json:"rabbitmq" toml:"rabbitmq"`
	Etcd     Etcd     `mapstructure:"etcd" json:"etcd" toml:"etcd"`
	Nacos    Nacos    `mapstructure:"nacos" json:"nacos" toml:"nacos"`
	Zap      Zap      `mapstructure:"zap" json:"zap" toml:"zap"`
}

// System 系统基础配置
type System struct {
	Env     string `mapstructure:"env" json:"env" toml:"env"`         // 环境模式: dev, test, prod
	Port    int    `mapstructure:"port" json:"port" toml:"port"`       // 服务端口
	DbType  string `mapstructure:"db-type" json:"dbType" toml:"db-type"` // 数据库类型
	OssType string `mapstructure:"oss-type" json:"ossType" toml:"oss-type"` // 对象存储类型
}

// MySQL 数据库配置
type MySQL struct {
	Path         string `mapstructure:"path" json:"path" toml:"path"`                   // 服务器地址
	Port         string `mapstructure:"port" json:"port" toml:"port"`                   // 端口
	Config       string `mapstructure:"config" json:"config" toml:"config"`             // 高级配置
	Dbname       string `mapstructure:"db-name" json:"dbname" toml:"db-name"`           // 数据库名
	Username     string `mapstructure:"username" json:"username" toml:"username"`       // 用户名
	Password     string `mapstructure:"password" json:"password" toml:"password"`       // 密码
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" toml:"max-idle-conns"` // 最大空闲连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" toml:"max-open-conns"` // 最大连接数
	LogMode      bool   `mapstructure:"log-mode" json:"logMode" toml:"log-mode"`         // 是否开启Gorm全局日志
}

// Redis 缓存配置
type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" toml:"addr"`           // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" toml:"password"` // 密码
	DB       int    `mapstructure:"db" json:"db" toml:"db"`                 // 数据库
}

// RabbitMQ 消息队列配置
type RabbitMQ struct {
	URL      string `mapstructure:"url" json:"url" toml:"url"` // RabbitMQ连接字符串
	Exchange struct {
		Tasks   string `mapstructure:"tasks" json:"tasks" toml:"tasks"`     // 任务交换机
		Results string `mapstructure:"results" json:"results" toml:"results"` // 结果交换机
	} `mapstructure:"exchange" json:"exchange" toml:"exchange"`
	Queue struct {
		TdataImport     string `mapstructure:"tdata-import" json:"tdataImport" toml:"tdata-import"`         // tdata导入队列
		TelegramAction  string `mapstructure:"telegram-action" json:"telegramAction" toml:"telegram-action"`   // Telegram操作队列
		TelegramResults string `mapstructure:"telegram-results" json:"telegramResults" toml:"telegram-results"` // 结果队列
	} `mapstructure:"queue" json:"queue" toml:"queue"`
}

// Zap 日志配置
type Zap struct {
	Level         string `mapstructure:"level" json:"level" toml:"level"`                   // 日志级别: debug, info, warn, error, dpanic, panic, fatal
	Format        string `mapstructure:"format" json:"format" toml:"format"`                 // 日志输出格式: console, json
	Prefix        string `mapstructure:"prefix" json:"prefix" toml:"prefix"`                 // 日志前缀
	Director      string `mapstructure:"director" json:"director" toml:"director"`           // 日志文件夹
	LinkName      string `mapstructure:"link-name" json:"linkName" toml:"link-name"`         // 软链接名称
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" toml:"show-line"`         // 显示行号
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" toml:"encode-level"` // 编码级别
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" toml:"stacktrace-key"` // 栈名称
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" toml:"log-in-console"` // 输出控制台
}

// Etcd 任务状态存储服务配置
type Etcd struct {
	Endpoints      []string `mapstructure:"endpoints" json:"endpoints" toml:"endpoints"`             // etcd节点地址
	DialTimeout    int      `mapstructure:"dial-timeout" json:"dialTimeout" toml:"dial-timeout"`       // 连接超时时间(秒)
	Username       string   `mapstructure:"username" json:"username" toml:"username"`                 // 用户名
	Password       string   `mapstructure:"password" json:"password" toml:"password"`                 // 密码
	TaskPrefix     string   `mapstructure:"task-prefix" json:"taskPrefix" toml:"task-prefix"`         // 任务状态key前缀
	LockPrefix     string   `mapstructure:"lock-prefix" json:"lockPrefix" toml:"lock-prefix"`         // 分布式锁前缀
	TaskTTL        int      `mapstructure:"task-ttl" json:"taskTTL" toml:"task-ttl"`                 // 任务状态过期时间(秒)
}

// Nacos 服务注册发现配置
type Nacos struct {
	IPAddr        string `mapstructure:"ip-addr" json:"ipAddr" toml:"ip-addr"`                   // Nacos服务器IP
	Port          uint64 `mapstructure:"port" json:"port" toml:"port"`                           // Nacos服务器端口
	NamespaceID   string `mapstructure:"namespace-id" json:"namespaceId" toml:"namespace-id"`       // 命名空间ID
	LogDir        string `mapstructure:"log-dir" json:"logDir" toml:"log-dir"`                     // 日志目录
	CacheDir      string `mapstructure:"cache-dir" json:"cacheDir" toml:"cache-dir"`               // 缓存目录
	ServiceName   string `mapstructure:"service-name" json:"serviceName" toml:"service-name"`       // 服务名称
	GroupName     string `mapstructure:"group-name" json:"groupName" toml:"group-name"`             // 服务分组
	Weight        float64 `mapstructure:"weight" json:"weight" toml:"weight"`                       // 服务权重
	ClusterName   string `mapstructure:"cluster-name" json:"clusterName" toml:"cluster-name"`       // 所属集群
	Username      string `mapstructure:"username" json:"username" toml:"username"`                 // 用户名
	Password      string `mapstructure:"password" json:"password" toml:"password"`                 // 密码
}
