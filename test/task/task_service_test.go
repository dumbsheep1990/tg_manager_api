package task_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tg_manager_api/services/task/service"
)

// 模拟任务数据库接口
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(task interface{}) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(task interface{}) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskRepository) FindByID(id uint) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockTaskRepository) FindAll(query interface{}) ([]interface{}, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]interface{}), int64(args.Int(1)), args.Error(2)
}

// 模拟etcd任务状态管理接口
type MockTaskStatusManager struct {
	mock.Mock
}

func (m *MockTaskStatusManager) SaveTaskStatus(taskID, status string) error {
	args := m.Called(taskID, status)
	return args.Error(0)
}

func (m *MockTaskStatusManager) GetTaskStatus(taskID string) (string, error) {
	args := m.Called(taskID)
	return args.String(0), args.Error(1)
}

// 测试创建任务
func TestCreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockStatusManager := new(MockTaskStatusManager)
	taskService := service.NewTaskService()
	
	// 模拟任务数据
	task := map[string]interface{}{
		"task_name": "测试任务",
		"task_type": "message",
		"account_group_id": uint(1),
		"message_template_id": uint(1),
		"target_count": 100,
	}
	
	// 设置模拟行为
	mockRepo.On("Create", mock.Anything).Return(nil)
	mockStatusManager.On("SaveTaskStatus", mock.Anything, "waiting").Return(nil)
	
	// 测试创建任务
	result, err := taskService.CreateTask(task)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
	mockStatusManager.AssertExpectations(t)
}

// 测试获取任务列表
func TestGetTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := service.NewTaskService()
	
	// 模拟查询参数
	query := map[string]interface{}{
		"page": 1,
		"page_size": 10,
		"task_name": "测试",
		"task_type": "message",
		"status": "running",
	}
	
	// 模拟返回数据
	mockTasks := []interface{}{
		map[string]interface{}{
			"id": uint(1),
			"task_name": "测试任务1",
			"task_type": "message",
			"status": "running",
			"progress": 45,
			"created_at": time.Now(),
		},
		map[string]interface{}{
			"id": uint(2),
			"task_name": "测试任务2",
			"task_type": "message",
			"status": "running",
			"progress": 60,
			"created_at": time.Now(),
		},
	}
	
	// 设置模拟行为
	mockRepo.On("FindAll", mock.Anything).Return(mockTasks, 2, nil)
	
	// 测试获取任务列表
	tasks, total, err := taskService.GetTasks(query)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

// 测试开始任务
func TestStartTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockStatusManager := new(MockTaskStatusManager)
	taskService := service.NewTaskService()
	
	// 模拟任务数据
	mockTask := map[string]interface{}{
		"id": uint(1),
		"task_name": "测试任务",
		"task_type": "message",
		"status": "waiting",
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockTask, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	mockStatusManager.On("SaveTaskStatus", "1", "running").Return(nil)
	
	// 测试开始任务
	err := taskService.StartTask(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockStatusManager.AssertExpectations(t)
}

// 测试暂停任务
func TestPauseTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockStatusManager := new(MockTaskStatusManager)
	taskService := service.NewTaskService()
	
	// 模拟任务数据
	mockTask := map[string]interface{}{
		"id": uint(1),
		"task_name": "测试任务",
		"task_type": "message",
		"status": "running",
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockTask, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	mockStatusManager.On("SaveTaskStatus", "1", "paused").Return(nil)
	
	// 测试暂停任务
	err := taskService.PauseTask(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockStatusManager.AssertExpectations(t)
}

// 测试停止任务
func TestStopTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockStatusManager := new(MockTaskStatusManager)
	taskService := service.NewTaskService()
	
	// 模拟任务数据
	mockTask := map[string]interface{}{
		"id": uint(1),
		"task_name": "测试任务",
		"task_type": "message",
		"status": "running",
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockTask, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	mockStatusManager.On("SaveTaskStatus", "1", "stopped").Return(nil)
	
	// 测试停止任务
	err := taskService.StopTask(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockStatusManager.AssertExpectations(t)
}

// 测试更新任务进度
func TestUpdateTaskProgress(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := service.NewTaskService()
	
	// 模拟任务数据
	mockTask := map[string]interface{}{
		"id": uint(1),
		"task_name": "测试任务",
		"task_type": "message",
		"status": "running",
		"progress": 45,
	}
	
	// 设置模拟行为
	mockRepo.On("FindByID", uint(1)).Return(mockTask, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)
	
	// 测试更新任务进度
	err := taskService.UpdateTaskProgress(uint(1), 75)
	
	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// 测试获取任务状态
func TestGetTaskStatus(t *testing.T) {
	mockStatusManager := new(MockTaskStatusManager)
	taskService := service.NewTaskService()
	
	// 设置模拟行为
	mockStatusManager.On("GetTaskStatus", "1").Return("running", nil)
	
	// 测试获取任务状态
	status, err := taskService.GetTaskStatus(uint(1))
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "running", status)
	mockStatusManager.AssertExpectations(t)
}
