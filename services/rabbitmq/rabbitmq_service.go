package rabbitmq

import (
	"encoding/json"
	"fmt"
	
	"github.com/streadway/amqp"
	
	"tg_manager_api/global"
)

// RabbitMQService RabbitMQ服务接口
type RabbitMQService interface {
	// 发布消息到队列
	PublishMessage(exchange, routingKey string, body interface{}) error
	
	// 创建消费者
	CreateConsumer(exchange, queueName, bindingKey string, handler func([]byte) error) error
	
	// 关闭连接
	Close() error
}

// rabbitMQService RabbitMQ服务实现
type rabbitMQService struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQService 创建RabbitMQ服务实例
func NewRabbitMQService(url string) (RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	
	return &rabbitMQService{
		connection: conn,
		channel:    ch,
	}, nil
}

// PublishMessage 发布消息到队列
func (s *rabbitMQService) PublishMessage(exchange, routingKey string, body interface{}) error {
	// 确保交换机存在
	err := s.channel.ExchangeDeclare(
		exchange, // 交换机名称
		"topic",  // 交换机类型
		true,     // 持久化
		false,    // 自动删除
		false,    // 内部使用
		false,    // 非阻塞
		nil,      // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}
	
	// 将消息体编码为JSON
	var messageBody []byte
	var marshalErr error
	
	switch v := body.(type) {
	case []byte:
		messageBody = v
	case string:
		messageBody = []byte(v)
	default:
		messageBody, marshalErr = json.Marshal(body)
		if marshalErr != nil {
			return global.ErrorJsonMarshalFailed
		}
	}
	
	// 发布消息
	err = s.channel.Publish(
		exchange,   // 交换机
		routingKey, // 路由键
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         messageBody,
			DeliveryMode: amqp.Persistent, // 持久化消息
		},
	)
	if err != nil {
		return global.ErrorQueuePublishFailed
	}
	
	return nil
}

// CreateConsumer 创建消费者
func (s *rabbitMQService) CreateConsumer(exchange, queueName, bindingKey string, handler func([]byte) error) error {
	// 声明交换机
	err := s.channel.ExchangeDeclare(
		exchange, // 交换机名称
		"topic",  // 交换机类型
		true,     // 持久化
		false,    // 自动删除
		false,    // 内部使用
		false,    // 非阻塞
		nil,      // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}
	
	// 声明队列
	q, err := s.channel.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 非阻塞
		nil,       // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	
	// 绑定队列到交换机
	err = s.channel.QueueBind(
		q.Name,     // 队列名称
		bindingKey, // 绑定键
		exchange,   // 交换机
		false,      // 非阻塞
		nil,        // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}
	
	// 开始消费消息
	msgs, err := s.channel.Consume(
		q.Name, // 队列名称
		"",     // 消费者标签
		false,  // 自动应答
		false,  // 独占
		false,  // 不接收同一个连接的投递
		false,  // 非阻塞
		nil,    // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}
	
	// 异步处理消息
	go func() {
		for msg := range msgs {
			err := handler(msg.Body)
			if err != nil {
				// 处理失败，将消息重新放回队列
				msg.Nack(false, true)
			} else {
				// 处理成功，确认消息
				msg.Ack(false)
			}
		}
	}()
	
	return nil
}

// Close 关闭连接
func (s *rabbitMQService) Close() error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.connection != nil {
		return s.connection.Close()
	}
	return nil
}
