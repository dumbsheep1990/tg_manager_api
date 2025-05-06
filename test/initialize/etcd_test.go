package initialize_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tg_manager_api/initialize"
)

// 模拟 etcd 客户端
type MockEtcdClient struct {
	mock.Mock
}

func (m *MockEtcdClient) Put(key, value string, opts ...interface{}) error {
	args := m.Called(key, value, opts)
	return args.Error(0)
}

func (m *MockEtcdClient) Get(key string, opts ...interface{}) (string, error) {
	args := m.Called(key, opts)
	return args.String(0), args.Error(1)
}

func (m *MockEtcdClient) Delete(key string, opts ...interface{}) error {
	args := m.Called(key, opts)
	return args.Error(0)
}

func TestSaveTaskStatus(t *testing.T) {
	mockClient := new(MockEtcdClient)
	
	// 设置模拟行为
	mockClient.On("Put", "tasks/123", "running", mock.Anything).Return(nil)
	
	// 测试保存任务状态
	err := initialize.SaveTaskStatus("123", "running")
	
	// 验证结果
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestGetTaskStatus(t *testing.T) {
	mockClient := new(MockEtcdClient)
	
	// 设置模拟行为
	mockClient.On("Get", "tasks/123", mock.Anything).Return("running", nil)
	
	// 测试获取任务状态
	status, err := initialize.GetTaskStatus("123")
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "running", status)
	mockClient.AssertExpectations(t)
}

func TestAcquireLock(t *testing.T) {
	mockClient := new(MockEtcdClient)
	
	// 设置模拟行为
	mockClient.On("Put", "locks/resource1", mock.Anything, mock.Anything).Return(nil)
	
	// 测试获取锁
	lock, err := initialize.AcquireLock("resource1", 30*time.Second)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, lock)
	mockClient.AssertExpectations(t)
}

func TestReleaseLock(t *testing.T) {
	mockClient := new(MockEtcdClient)
	
	// 设置模拟行为
	mockClient.On("Delete", "locks/resource1", mock.Anything).Return(nil)
	
	// 测试释放锁
	err := initialize.ReleaseLock("resource1")
	
	// 验证结果
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
