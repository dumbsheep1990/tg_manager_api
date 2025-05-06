package message_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tg_manager_api/services/message/service"
)

// 模拟消息模板数据库接口
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(message interface{}) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) Update(message interface{}) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByID(id uint) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockMessageRepository) FindAll(query interface{}) ([]interface{}, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]interface{}), int64(args.Int(1)), args.Error(2)
}

// 模拟消息发送接口
type MockMessageSender struct {
	mock.Mock
}

func (m *MockMessageSender) SendMessage(message interface{}, recipients []string) (interface{}, error) {
	args := m.Called(message, recipients)
	return args.Get(0), args.Error(1)
}

// 测试创建消息模板
func TestCreateMessageTemplate(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	messageService := service.NewMessageService()
	
	// 模拟模板数据
	template := map[string]interface{}{
		"template_name": "测试模板",
		"message_type": "text",
		"content": "这是一条测试消息",
		"has_media": false,
	}
	
	// 设置模拟行为
	mockRepo.On("Create", mock.Anything).Return(nil)
	
	// 测试创建模板
	result, err := messageService.CreateMessageTemplate(template)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

// 测试获取消息模板列表
func TestGetMessageTemplates(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	messageService := service.NewMessageService()
	
	// 模拟查询参数
	query := map[string]interface{}{
		"page": 1,
		"page_size": 10,
		"template_name": "测试",
		"message_type": "text",
	}
	
	// 模拟返回数据
	mockTemplates := []interface{}{
		map[string]interface{}{
			"id": uint(1),
			"template_name": "测试模板1",
			"message_type": "text",
			"content": "这是测试模板1",
			"has_media": false,
			"created_at": time.Now(),
		},
		map[string]interface{}{
			"id": uint(2),
			"template_name": "测试模板2",
			"message_type": "text",
			"content": "这是测试模板2",
			"has_media": false,
			"created_at": time.Now(),
		},
	}
	
	// 设置模拟行为
	mockRepo.On("FindAll", mock.Anything).Return(mockTemplates, 2, nil)
	
	// 测试获取模板列表
	templates, total, err := messageService.GetMessageTemplates(query)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(templates))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// 测试更新消息模板
func TestUpdateMessageTemplate(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	messageService := service.NewMessageService()
	
	// 模拟模板数据
	template := map[string]interface{}{
		"id": uint(1),
		"template_name": "更新的模板",
		"message_type": "text",
		"content": "这是更新后的测试消息",
		"has_media": false,
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(template, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	
	// 测试更新模板
	result, err := messageService.UpdateMessageTemplate(uint(1), template)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

// 测试删除消息模板
func TestDeleteMessageTemplate(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	messageService := service.NewMessageService()
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(map[string]interface{}{"id": uint(1)}, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	
	// 测试删除模板
	err := messageService.DeleteMessageTemplate(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// 测试发送消息
func TestSendMessage(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockSender := new(MockMessageSender)
	messageService := service.NewMessageService()
	
	// 模拟消息数据
	messageData := map[string]interface{}{
		"template_id": uint(1),
		"recipients": []string{"user1", "user2", "user3"},
		"scheduled_time": time.Now().Add(time.Hour).Format(time.RFC3339),
	}
	
	// 模拟模板数据
	templateData := map[string]interface{}{
		"id": uint(1),
		"template_name": "测试模板",
		"message_type": "text",
		"content": "这是测试消息内容",
		"has_media": false,
	}
	
	// 模拟发送结果
	sendResult := map[string]interface{}{
		"total": 3,
		"success": 3,
		"failed": 0,
		"message_id": "msg123456",
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(templateData, nil)
	mockSender.On("SendMessage", mock.Anything, []string{"user1", "user2", "user3"}).Return(sendResult, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)
	
	// 测试发送消息
	result, err := messageService.SendMessage(messageData)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 3, result["total"])
	assert.Equal(t, 3, result["success"])
	assert.Equal(t, 0, result["failed"])
	mockRepo.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

// 测试获取消息发送历史
func TestGetMessageHistory(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	messageService := service.NewMessageService()
	
	// 模拟查询参数
	query := map[string]interface{}{
		"page": 1,
		"page_size": 10,
		"message_type": "text",
		"status": "success",
	}
	
	// 模拟返回数据
	mockHistory := []interface{}{
		map[string]interface{}{
			"id": uint(1),
			"template_id": uint(1),
			"template_name": "测试模板1",
			"message_type": "text",
			"total_count": 10,
			"success_count": 10,
			"fail_count": 0,
			"status": "success",
			"sent_at": time.Now(),
		},
		map[string]interface{}{
			"id": uint(2),
			"template_id": uint(1),
			"template_name": "测试模板1",
			"message_type": "text",
			"total_count": 5,
			"success_count": 5,
			"fail_count": 0,
			"status": "success",
			"sent_at": time.Now(),
		},
	}
	
	// 设置模拟行为
	mockRepo.On("FindAll", mock.Anything).Return(mockHistory, 2, nil)
	
	// 测试获取消息历史
	history, total, err := messageService.GetMessageHistory(query)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(history))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// 测试获取消息状态
func TestGetMessageStatus(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	messageService := service.NewMessageService()
	
	// 模拟消息数据
	mockMessage := map[string]interface{}{
		"id": uint(1),
		"status": "success",
		"total_count": 10,
		"success_count": 9,
		"fail_count": 1,
		"progress": 100,
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockMessage, nil)
	
	// 测试获取消息状态
	status, err := messageService.GetMessageStatus(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "success", status["status"])
	assert.Equal(t, 10, status["total_count"])
	assert.Equal(t, 9, status["success_count"])
	assert.Equal(t, 1, status["fail_count"])
	mockRepo.AssertExpectations(t)
}
