package initialize

import (
	"os"
	"tg_manager_api/global"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// InitRabbitMQ 初始化RabbitMQ连接
func InitRabbitMQ() {
	config := global.Config.RabbitMQ
	
	// 建立连接
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		global.Logger.Error("RabbitMQ连接失败", zap.Error(err))
		os.Exit(1)
	}
	
	// 创建channel
	ch, err := conn.Channel()
	if err != nil {
		global.Logger.Error("RabbitMQ创建channel失败", zap.Error(err))
		os.Exit(1)
	}
	
	// 声明交换机
	err = ch.ExchangeDeclare(
		config.Exchange.Tasks,  // name
		"direct",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		global.Logger.Error("声明任务交换机失败", zap.Error(err))
		os.Exit(1)
	}
	
	err = ch.ExchangeDeclare(
		config.Exchange.Results, // name
		"direct",                // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		global.Logger.Error("声明结果交换机失败", zap.Error(err))
		os.Exit(1)
	}
	
	// 声明队列
	_, err = ch.QueueDeclare(
		config.Queue.TdataImport,     // name
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		global.Logger.Error("声明tdata导入队列失败", zap.Error(err))
		os.Exit(1)
	}
	
	_, err = ch.QueueDeclare(
		config.Queue.TelegramAction,  // name
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		global.Logger.Error("声明Telegram操作队列失败", zap.Error(err))
		os.Exit(1)
	}
	
	_, err = ch.QueueDeclare(
		config.Queue.TelegramResults, // name
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		global.Logger.Error("声明Telegram结果队列失败", zap.Error(err))
		os.Exit(1)
	}
	
	// 绑定队列到交换机
	err = ch.QueueBind(
		config.Queue.TdataImport,    // queue name
		"tdata.import",              // routing key
		config.Exchange.Tasks,       // exchange
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("绑定tdata导入队列失败", zap.Error(err))
		os.Exit(1)
	}
	
	err = ch.QueueBind(
		config.Queue.TelegramAction, // queue name
		"telegram.action.#",         // routing key
		config.Exchange.Tasks,       // exchange
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("绑定Telegram操作队列失败", zap.Error(err))
		os.Exit(1)
	}
	
	err = ch.QueueBind(
		config.Queue.TelegramResults, // queue name
		"telegram.results",           // routing key
		config.Exchange.Results,      // exchange
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("绑定Telegram结果队列失败", zap.Error(err))
		os.Exit(1)
	}
	
	global.Logger.Info("RabbitMQ初始化成功")
	
	global.RabbitMQConn = conn
	global.RabbitMQChannel = ch
}
