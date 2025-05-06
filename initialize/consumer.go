package initialize

import (
	"context"
	"encoding/json"
	"tg_manager_api/global"
	"tg_manager_api/model"
	
	"go.uber.org/zap"
)

// InitResultConsumer 初始化消息消费者
func InitResultConsumer() {
	config := global.Config.RabbitMQ
	
	// 创建新的channel来消费消息
	ch, err := global.RabbitMQConn.Channel()
	if err != nil {
		global.Logger.Fatal("消费者创建channel失败", zap.Error(err))
		return
	}
	defer ch.Close()
	
	// 声明要消费的队列
	q, err := ch.QueueDeclare(
		config.Queue.TelegramResults, // 队列名称
		true,                         // 持久化
		false,                        // 不使用时删除
		false,                        // 非排他
		false,                        // 非阻塞
		nil,                          // 附加参数
	)
	if err != nil {
		global.Logger.Fatal("消费者声明队列失败", zap.Error(err))
		return
	}
	
	// 消费消息
	msgs, err := ch.Consume(
		q.Name, // 队列
		"",     // 消费者标识
		true,   // 自动确认
		false,  // 非排他
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 参数
	)
	if err != nil {
		global.Logger.Fatal("注册消费者失败", zap.Error(err))
		return
	}
	
	global.Logger.Info("消息消费者启动成功")
	
	// 处理收到的消息
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			handleResultMessage(d.Body)
		}
	}()
	
	<-forever
}

// handleResultMessage 处理来自Python工作者的结果消息
func handleResultMessage(body []byte) {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	if err != nil {
		global.Logger.Error("解析结果消息失败", zap.Error(err))
		return
	}
	
	// 检查消息类型并处理
	messageType, ok := result["type"].(string)
	if !ok {
		global.Logger.Error("消息类型缺失或无效")
		return
	}
	
	ctx := context.Background()
	global.Logger.Info("收到消息", zap.String("类型", messageType))
	
	switch messageType {
	case "tdata.import.result":
		handleTdataImportResult(ctx, result)
	case "telegram.action.result":
		handleTelegramActionResult(ctx, result)
	default:
		global.Logger.Warn("未知消息类型", zap.String("类型", messageType))
	}
}

// handleTdataImportResult 处理tdata导入操作的结果
func handleTdataImportResult(ctx context.Context, result map[string]interface{}) {
	taskID, ok := result["task_id"].(string)
	if !ok {
		global.Logger.Error("tdata导入结果缺少task_id")
		return
	}
	
	// 更新数据库中的账号状态
	success, ok := result["success"].(bool)
	if !ok {
		global.Logger.Error("tdata导入结果缺少success标记")
		return
	}
	
	var account model.Account
	if err := global.DB.Where("task_id = ?", taskID).First(&account).Error; err != nil {
		global.Logger.Error("无法查找到对应任务的账号", 
			zap.String("task_id", taskID), 
			zap.Error(err))
		return
	}
	
	if success {
		// 更新成功导入的账号数据
		accountInfo, ok := result["account_info"].(map[string]interface{})
		if !ok {
			global.Logger.Error("成功的tdata导入结果缺少account_info")
			return
		}
		
		phone, _ := accountInfo["phone"].(string)
		username, _ := accountInfo["username"].(string)
		
		account.Status = "ACTIVE"
		account.Phone = phone
		account.Username = username
		
		if err := global.DB.Save(&account).Error; err != nil {
			global.Logger.Error("无法在成功导入后更新账号", zap.Error(err))
		}
		
		global.Logger.Info("账号导入成功", zap.String("phone", phone))
	} else {
		// 更新失败状态
		errorMsg, _ := result["error"].(string)
		
		account.Status = "IMPORT_FAILED"
		account.ErrorMessage = errorMsg
		
		if err := global.DB.Save(&account).Error; err != nil {
			global.Logger.Error("无法在失败导入后更新账号", zap.Error(err))
		}
		
		global.Logger.Warn("账号导入失败", zap.String("error", errorMsg))
	}
}

// handleTelegramActionResult 处理Telegram操作的结果
func handleTelegramActionResult(ctx context.Context, result map[string]interface{}) {
	taskID, ok := result["task_id"].(string)
	if !ok {
		global.Logger.Error("Telegram操作结果缺少task_id")
		return
	}
	
	// 这里将是根据操作结果更新任务记录的逻辑
	// 这将取决于特定操作（发送消息、加入群组等）
	global.Logger.Info("收到任务结果", zap.String("task_id", taskID))
}
