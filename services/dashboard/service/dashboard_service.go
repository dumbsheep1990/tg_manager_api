package service

import (
	"context"
	"time"
)

// StatsPeriod 统计周期
type StatsPeriod string

const (
	StatsPeriodDay   StatsPeriod = "DAY"   // 日
	StatsPeriodWeek  StatsPeriod = "WEEK"  // 周
	StatsPeriodMonth StatsPeriod = "MONTH" // 月
)

// Statistics 统计数据
type Statistics struct {
	AccountCount         int            `json:"account_count"`          // 账号总数
	ActiveAccountCount   int            `json:"active_account_count"`   // 活跃账号数
	BannedAccountCount   int            `json:"banned_account_count"`   // 被封禁账号数
	TaskCount            int            `json:"task_count"`             // 任务总数
	PendingTaskCount     int            `json:"pending_task_count"`     // 待处理任务数
	SuccessTaskCount     int            `json:"success_task_count"`     // 成功任务数
	FailedTaskCount      int            `json:"failed_task_count"`      // 失败任务数
	MessageSentCount     int            `json:"message_sent_count"`     // 已发送消息数
	MessageSentCountByDay map[string]int `json:"message_sent_count_by_day"` // 按天统计的发送消息数
	UpdatedAt           time.Time       `json:"updated_at"`             // 更新时间
}

// DashboardService 仪表盘服务接口
type DashboardService interface {
	// 获取系统概览统计数据
	GetStatistics(ctx context.Context) (*Statistics, error)
	
	// 获取指定周期的消息发送统计
	GetMessageSentStats(ctx context.Context, period StatsPeriod) (map[string]int, error)
	
	// 获取账号状态分布
	GetAccountStatusDistribution(ctx context.Context) (map[string]int, error)
	
	// 获取任务状态分布 
	GetTaskStatusDistribution(ctx context.Context) (map[string]int, error)
	
	// 获取最近活动日志
	GetRecentActivityLogs(ctx context.Context, limit int) ([]map[string]interface{}, error)
}

// NewDashboardService 创建仪表盘服务实例
func NewDashboardService() DashboardService {
	return &dashboardService{}
}

// dashboardService 仪表盘服务实现
type dashboardService struct{}

// GetStatistics 获取系统概览统计数据
func (s *dashboardService) GetStatistics(ctx context.Context) (*Statistics, error) {
	// 实现获取系统概览统计数据逻辑
	return nil, nil
}

// GetMessageSentStats 获取指定周期的消息发送统计
func (s *dashboardService) GetMessageSentStats(ctx context.Context, period StatsPeriod) (map[string]int, error) {
	// 实现获取消息发送统计逻辑
	return nil, nil
}

// GetAccountStatusDistribution 获取账号状态分布
func (s *dashboardService) GetAccountStatusDistribution(ctx context.Context) (map[string]int, error) {
	// 实现获取账号状态分布逻辑
	return nil, nil
}

// GetTaskStatusDistribution 获取任务状态分布
func (s *dashboardService) GetTaskStatusDistribution(ctx context.Context) (map[string]int, error) {
	// 实现获取任务状态分布逻辑
	return nil, nil
}

// GetRecentActivityLogs 获取最近活动日志
func (s *dashboardService) GetRecentActivityLogs(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	// 实现获取最近活动日志逻辑
	return nil, nil
}
