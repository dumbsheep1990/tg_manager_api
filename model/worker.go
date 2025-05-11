package model

import "time"

// Worker 工作节点模型
type Worker struct {
	BaseModel
	WorkerID      string    `gorm:"uniqueIndex;column:worker_id;comment:工作节点ID" json:"worker_id"`  // 工作节点唯一标识
	Hostname      string    `gorm:"column:hostname;comment:主机名" json:"hostname"`                   // 主机名
	IP            string    `gorm:"column:ip;comment:IP地址" json:"ip"`                             // IP地址
	Status        string    `gorm:"column:status;comment:状态" json:"status"`                       // 状态: online, offline, busy
	LastHeartbeat time.Time `gorm:"column:last_heartbeat;comment:最后心跳时间" json:"last_heartbeat"`   // 最后心跳时间
	MaxTasks      int       `gorm:"column:max_tasks;default:10;comment:最大任务数" json:"max_tasks"`    // 最大可同时执行的任务数
	CurrentTasks  int       `gorm:"column:current_tasks;default:0;comment:当前任务数" json:"current_tasks"` // 当前正在执行的任务数
	Tags          string    `gorm:"column:tags;comment:标签(逗号分隔)" json:"tags"`                     // 标签，用于任务分配策略
	Version       string    `gorm:"column:version;comment:Worker版本" json:"version"`                // Worker版本号
	
	// 外键关系
	TaskRecords []TaskRecord `json:"task_records,omitempty" gorm:"foreignKey:WorkerID;references:WorkerID"` // 执行的任务记录
	TaskAssignments []TaskAssignment `json:"task_assignments,omitempty" gorm:"foreignKey:WorkerID;references:WorkerID"` // 分配的任务
}

// TableName 设置表名
func (Worker) TableName() string {
	return "workers"
}
