package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Task 任务模型
type Task struct {
	BaseModel
	TaskID      string          `gorm:"uniqueIndex;column:task_id;comment:任务ID" json:"task_id"`        // 任务ID
	TaskType    string          `gorm:"column:task_type;comment:任务类型" json:"task_type"`              // 任务类型: send_message, join_group, add_contact等
	AccountID   uint            `gorm:"index;column:account_id;comment:账号ID" json:"account_id"`       // 关联的账号ID
	Params      TaskParams      `gorm:"type:json;column:params;comment:任务参数" json:"params"`           // 任务参数，JSON格式
	Status      string          `gorm:"column:status;comment:任务状态" json:"status"`                    // 状态: pending, assigned, processing, completed, failed
	Priority    int             `gorm:"column:priority;default:0;comment:任务优先级" json:"priority"`      // 优先级，数字越大优先级越高
	ErrorMessage string         `gorm:"column:error_message;comment:错误信息" json:"error_message"`      // 错误信息
	TimeoutSec  int             `gorm:"column:timeout_sec;default:300;comment:超时时间(秒)" json:"timeout_sec"` // 执行超时时间，单位秒
	StartedAt   *time.Time      `gorm:"column:started_at;comment:开始时间" json:"started_at"`            // 开始执行时间
	CompletedAt *time.Time      `gorm:"column:completed_at;comment:完成时间" json:"completed_at"`        // 完成时间
	
	// 外键关系
	Account    *Account         `json:"account,omitempty" gorm:"foreignKey:AccountID"` // 关联的账号
	TaskRecords []TaskRecord    `json:"task_records,omitempty" gorm:"foreignKey:TaskID;references:TaskID"` // 任务执行记录
}

// TableName 设置表名
func (Task) TableName() string {
	return "tasks"
}

// TaskParams 任务参数，使用JSON存储
type TaskParams map[string]interface{}

// Value 实现driver.Valuer接口
func (p TaskParams) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan 实现sql.Scanner接口
func (p *TaskParams) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, &p)
}
