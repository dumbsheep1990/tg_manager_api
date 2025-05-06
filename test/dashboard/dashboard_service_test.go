package dashboard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tg_manager_api/services/dashboard/service"
)

// 模拟仪表盘数据访问层
type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetSystemOverview() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockDashboardRepository) GetAccountStats() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockDashboardRepository) GetTaskStats() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockDashboardRepository) GetMessageStats() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockDashboardRepository) GetSystemStatus() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockDashboardRepository) GetRecentActivities(query interface{}) ([]interface{}, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]interface{}), int64(args.Int(1)), args.Error(2)
}

func (m *MockDashboardRepository) GetTaskTrend(params interface{}) ([]interface{}, error) {
	args := m.Called(params)
	return args.Get(0).([]interface{}), args.Error(1)
}

// 测试获取系统概览
func TestGetDashboardOverview(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	dashboardService := service.NewDashboardService()
	
	// 模拟返回数据
	mockOverview := map[string]interface{}{
		"accountCount": 150,
		"activeAccounts": 120,
		"taskCount": 48,
		"runningTasks": 5,
		"messageCount": 1240,
		"successMessages": 1180,
		"systemStatus": "normal",
		"workerCount": 3,
	}
	
	// 设置模拟行为
	mockRepo.On("GetSystemOverview").Return(mockOverview, nil)
	
	// 测试获取系统概览
	overview, err := dashboardService.GetDashboardOverview()
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 150, overview["accountCount"])
	assert.Equal(t, 120, overview["activeAccounts"])
	assert.Equal(t, 48, overview["taskCount"])
	assert.Equal(t, 5, overview["runningTasks"])
	assert.Equal(t, "normal", overview["systemStatus"])
	mockRepo.AssertExpectations(t)
}

// 测试获取账号统计
func TestGetAccountStats(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	dashboardService := service.NewDashboardService()
	
	// 模拟返回数据
	mockStats := map[string]interface{}{
		"totalAccounts": 150,
		"onlineAccounts": 80,
		"offlineAccounts": 40,
		"disabledAccounts": 30,
		"groupDistribution": []map[string]interface{}{
			{"name": "分组1", "value": 50},
			{"name": "分组2", "value": 60},
			{"name": "分组3", "value": 40},
		},
		"accountGrowth": []map[string]interface{}{
			{"date": "2025-04-01", "count": 100},
			{"date": "2025-04-15", "count": 120},
			{"date": "2025-05-01", "count": 150},
		},
	}
	
	// 设置模拟行为
	mockRepo.On("GetAccountStats").Return(mockStats, nil)
	
	// 测试获取账号统计
	stats, err := dashboardService.GetAccountStats()
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 150, stats["totalAccounts"])
	assert.Equal(t, 80, stats["onlineAccounts"])
	assert.Equal(t, 3, len(stats["groupDistribution"].([]map[string]interface{})))
	assert.Equal(t, 3, len(stats["accountGrowth"].([]map[string]interface{})))
	mockRepo.AssertExpectations(t)
}

// 测试获取任务统计
func TestGetTaskStats(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	dashboardService := service.NewDashboardService()
	
	// 模拟返回数据
	mockStats := map[string]interface{}{
		"totalTasks": 48,
		"runningTasks": 5,
		"completedTasks": 35,
		"failedTasks": 8,
		"taskTypeDistribution": []map[string]interface{}{
			{"type": "message", "count": 30},
			{"type": "addFriend", "count": 10},
			{"type": "joinGroup", "count": 8},
		},
		"taskStatusDistribution": []map[string]interface{}{
			{"status": "running", "count": 5},
			{"status": "completed", "count": 35},
			{"status": "failed", "count": 8},
		},
		"successRate": 87.5,
	}
	
	// 设置模拟行为
	mockRepo.On("GetTaskStats").Return(mockStats, nil)
	
	// 测试获取任务统计
	stats, err := dashboardService.GetTaskStats()
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 48, stats["totalTasks"])
	assert.Equal(t, 5, stats["runningTasks"])
	assert.Equal(t, 3, len(stats["taskTypeDistribution"].([]map[string]interface{})))
	assert.Equal(t, 87.5, stats["successRate"])
	mockRepo.AssertExpectations(t)
}

// 测试获取消息统计
func TestGetMessageStats(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	dashboardService := service.NewDashboardService()
	
	// 模拟返回数据
	mockStats := map[string]interface{}{
		"totalMessages": 1240,
		"successMessages": 1180,
		"failedMessages": 60,
		"messageTypeDistribution": []map[string]interface{}{
			{"type": "text", "count": 800},
			{"type": "image", "count": 300},
			{"type": "mixed", "count": 140},
		},
		"messageTrend": []map[string]interface{}{
			{"date": "2025-04-01", "count": 400},
			{"date": "2025-04-15", "count": 800},
			{"date": "2025-05-01", "count": 1240},
		},
		"successRate": 95.2,
	}
	
	// 设置模拟行为
	mockRepo.On("GetMessageStats").Return(mockStats, nil)
	
	// 测试获取消息统计
	stats, err := dashboardService.GetMessageStats()
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 1240, stats["totalMessages"])
	assert.Equal(t, 1180, stats["successMessages"])
	assert.Equal(t, 3, len(stats["messageTypeDistribution"].([]map[string]interface{})))
	assert.Equal(t, 95.2, stats["successRate"])
	mockRepo.AssertExpectations(t)
}

// 测试获取系统状态
func TestGetSystemStatus(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	dashboardService := service.NewDashboardService()
	
	// 模拟返回数据
	mockStatus := map[string]interface{}{
		"status": "normal",
		"cpuUsage": 35.2,
		"memoryUsage": 42.8,
		"diskUsage": 58.6,
		"workerNodes": []map[string]interface{}{
			{"id": "worker1", "status": "active", "tasks": 2, "lastHeartbeat": "2025-05-06T12:55:00Z"},
			{"id": "worker2", "status": "active", "tasks": 3, "lastHeartbeat": "2025-05-06T12:54:30Z"},
			{"id": "worker3", "status": "active", "tasks": 0, "lastHeartbeat": "2025-05-06T12:55:15Z"},
		},
		"uptimeHours": 168,
	}
	
	// 设置模拟行为
	mockRepo.On("GetSystemStatus").Return(mockStatus, nil)
	
	// 测试获取系统状态
	status, err := dashboardService.GetSystemStatus()
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "normal", status["status"])
	assert.Equal(t, 35.2, status["cpuUsage"])
	assert.Equal(t, 3, len(status["workerNodes"].([]map[string]interface{})))
	mockRepo.AssertExpectations(t)
}

// 测试获取最近活动
func TestGetRecentActivities(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	dashboardService := service.NewDashboardService()
	
	// 模拟查询参数
	query := map[string]interface{}{
		"page": 1,
		"pageSize": 10,
	}
	
	// 模拟返回数据
	mockActivities := []interface{}{
		map[string]interface{}{
			"id": uint(1),
			"time": "2025-05-06T12:50:00Z",
			"type": "任务",
			"content": "创建了任务'测试任务1'",
			"operator": "admin",
			"status": "成功",
		},
		map[string]interface{}{
			"id": uint(2),
			"time": "2025-05-06T12:45:00Z",
			"type": "账号",
			"content": "导入了15个TData账号",
			"operator": "admin",
			"status": "成功",
		},
		map[string]interface{}{
			"id": uint(3),
			"time": "2025-05-06T12:40:00Z",
			"type": "消息",
			"content": "发送了测试消息模板到50个目标",
			"operator": "admin",
			"status": "进行中",
		},
	}
	
	// 设置模拟行为
	mockRepo.On("GetRecentActivities", mock.Anything).Return(mockActivities, 3, nil)
	
	// 测试获取最近活动
	activities, total, err := dashboardService.GetRecentActivities(query)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 3, len(activities))
	assert.Equal(t, int64(3), total)
	mockRepo.AssertExpectations(t)
}

// 测试获取任务趋势
func TestGetTaskTrend(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	dashboardService := service.NewDashboardService()
	
	// 模拟参数
	params := map[string]interface{}{
		"days": 7,
	}
	
	// 模拟返回数据
	mockTrend := []interface{}{
		map[string]interface{}{
			"date": "2025-04-30",
			"created": 5,
			"completed": 3,
			"failed": 1,
		},
		map[string]interface{}{
			"date": "2025-05-01",
			"created": 8,
			"completed": 6,
			"failed": 2,
		},
		map[string]interface{}{
			"date": "2025-05-02",
			"created": 6,
			"completed": 5,
			"failed": 0,
		},
		map[string]interface{}{
			"date": "2025-05-03",
			"created": 10,
			"completed": 8,
			"failed": 1,
		},
		map[string]interface{}{
			"date": "2025-05-04",
			"created": 7,
			"completed": 6,
			"failed": 1,
		},
		map[string]interface{}{
			"date": "2025-05-05",
			"created": 9,
			"completed": 7,
			"failed": 2,
		},
		map[string]interface{}{
			"date": "2025-05-06",
			"created": 3,
			"completed": 0,
			"failed": 1,
		},
	}
	
	// 设置模拟行为
	mockRepo.On("GetTaskTrend", mock.Anything).Return(mockTrend, nil)
	
	// 测试获取任务趋势
	trend, err := dashboardService.GetTaskTrend(params)
	
	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, 7, len(trend))
	mockRepo.AssertExpectations(t)
}
