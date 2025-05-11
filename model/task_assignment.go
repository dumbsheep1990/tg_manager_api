package model

import "time"

// TaskAssignment 任务分配模型
type TaskAssignment struct {
	BaseModel
	TaskID          string     `gorm:"index;column:task_id;comment:任务ID" json:"task_id"`                  // 关联的任务ID
	WorkerID        string     `gorm:"index;column:worker_id;comment:工作节点ID" json:"worker_id"`           // 分配给的工作节点ID
	Status          string     `gorm:"column:status;comment:分配状态" json:"status"`                         // 状态: pending, accepted, rejected, completed
	AssignedAt      time.Time  `gorm:"column:assigned_at;comment:分配时间" json:"assigned_at"`               // 分配时间
	AcceptedAt      *time.Time `gorm:"column:accepted_at;comment:接受时间" json:"accepted_at"`               // 接受时间
	CompletedAt     *time.Time `gorm:"column:completed_at;comment:完成时间" json:"completed_at"`             // 完成时间
	RejectionReason string     `gorm:"column:rejection_reason;comment:拒绝原因" json:"rejection_reason"`     // 拒绝原因
	Priority        int        `gorm:"column:priority;default:0;comment:优先级" json:"priority"`             // 优先级，数字越大优先级越高
	
	// 外键关系
	Task            *Task      `json:"task,omitempty" gorm:"foreignKey:TaskID;references:TaskID"`       // 关联的任务
	Worker          *Worker    `json:"worker,omitempty" gorm:"foreignKey:WorkerID;references:WorkerID"` // 关联的工作节点
}

// TableName 设置表名
func (TaskAssignment) TableName() string {
	return "task_assignments"
}
