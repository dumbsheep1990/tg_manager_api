package account_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tg_manager_api/services/account/service"
)

// 模拟账号数据库接口
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(account interface{}) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) Update(account interface{}) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAccountRepository) FindByID(id uint) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockAccountRepository) FindAll(query interface{}) ([]interface{}, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]interface{}), int64(args.Int(1)), args.Error(2)
}

// 测试创建账号
func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	accountService := service.NewAccountService()
	
	// 模拟账号数据
	account := map[string]interface{}{
		"account_name": "test_account",
		"phone": "13800138000",
		"username": "test_user",
		"group_id": uint(1),
	}
	
	// 设置模拟行为
	mockRepo.On("Create", mock.Anything).Return(nil)
	
	// 测试创建账号
	result, err := accountService.CreateAccount(account)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

// 测试获取账号列表
func TestGetAccounts(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	accountService := service.NewAccountService()
	
	// 模拟查询参数
	query := map[string]interface{}{
		"page": 1,
		"page_size": 10,
		"account_name": "test",
		"group_id": uint(1),
	}
	
	// 模拟返回数据
	mockAccounts := []interface{}{
		map[string]interface{}{
			"id": uint(1),
			"account_name": "test_account1",
			"phone": "13800138001",
			"username": "test_user1",
			"group_id": uint(1),
		},
		map[string]interface{}{
			"id": uint(2),
			"account_name": "test_account2",
			"phone": "13800138002",
			"username": "test_user2",
			"group_id": uint(1),
		},
	}
	
	// 设置模拟行为
	mockRepo.On("FindAll", mock.Anything).Return(mockAccounts, 2, nil)
	
	// 测试获取账号列表
	accounts, total, err := accountService.GetAccounts(query)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(accounts))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// 测试更新账号
func TestUpdateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	accountService := service.NewAccountService()
	
	// 模拟账号数据
	account := map[string]interface{}{
		"id": uint(1),
		"account_name": "updated_account",
		"phone": "13800138000",
		"username": "test_user",
		"group_id": uint(2),
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(account, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	
	// 测试更新账号
	result, err := accountService.UpdateAccount(uint(1), account)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

// 测试删除账号
func TestDeleteAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	accountService := service.NewAccountService()
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(map[string]interface{}{"id": uint(1)}, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	
	// 测试删除账号
	err := accountService.DeleteAccount(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// 测试获取账号详情
func TestGetAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	accountService := service.NewAccountService()
	
	// 模拟账号数据
	mockAccount := map[string]interface{}{
		"id": uint(1),
		"account_name": "test_account",
		"phone": "13800138000",
		"username": "test_user",
		"group_id": uint(1),
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockAccount, nil)
	
	// 测试获取账号详情
	account, err := accountService.GetAccount(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, mockAccount, account)
	mockRepo.AssertExpectations(t)
}
