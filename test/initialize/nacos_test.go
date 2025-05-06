package initialize_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tg_manager_api/initialize"
)

// 模拟 Nacos 客户端
type MockNacosClient struct {
	mock.Mock
}

func (m *MockNacosClient) RegisterInstance(param interface{}) (bool, error) {
	args := m.Called(param)
	return args.Bool(0), args.Error(1)
}

func (m *MockNacosClient) DeregisterInstance(param interface{}) (bool, error) {
	args := m.Called(param)
	return args.Bool(0), args.Error(1)
}

func (m *MockNacosClient) SelectInstances(param interface{}) ([]interface{}, error) {
	args := m.Called(param)
	return args.Get(0).([]interface{}), args.Error(1)
}

func TestRegisterService(t *testing.T) {
	mockClient := new(MockNacosClient)
	
	// 设置模拟行为
	mockClient.On("RegisterInstance", mock.Anything).Return(true, nil)
	
	// 测试服务注册
	success, err := initialize.RegisterService("account-service", "127.0.0.1", 8080)
	
	// 验证结果
	assert.NoError(t, err)
	assert.True(t, success)
	mockClient.AssertExpectations(t)
}

func TestDeregisterService(t *testing.T) {
	mockClient := new(MockNacosClient)
	
	// 设置模拟行为
	mockClient.On("DeregisterInstance", mock.Anything).Return(true, nil)
	
	// 测试服务注销
	success, err := initialize.DeregisterService("account-service", "127.0.0.1", 8080)
	
	// 验证结果
	assert.NoError(t, err)
	assert.True(t, success)
	mockClient.AssertExpectations(t)
}

func TestGetServiceInstances(t *testing.T) {
	mockClient := new(MockNacosClient)
	
	// 模拟服务实例
	mockInstances := []interface{}{
		map[string]interface{}{
			"ip": "127.0.0.1",
			"port": 8080,
			"serviceName": "account-service",
			"healthy": true,
		},
		map[string]interface{}{
			"ip": "127.0.0.1",
			"port": 8081,
			"serviceName": "account-service",
			"healthy": true,
		},
	}
	
	// 设置模拟行为
	mockClient.On("SelectInstances", mock.Anything).Return(mockInstances, nil)
	
	// 测试获取服务实例
	instances, err := initialize.GetServiceInstances("account-service")
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(instances))
	mockClient.AssertExpectations(t)
}
