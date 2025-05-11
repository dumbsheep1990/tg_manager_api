package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// TaskRecord 任务执行记录模型
type TaskRecord struct {
	BaseModel
	TaskID       string       `gorm:"index;column:task_id;comment:任务ID" json:"task_id"`            // 关联的任务ID
	WorkerID     string       `gorm:"index;column:worker_id;comment:工作节点ID" json:"worker_id"`     // 执行任务的工作节点ID
	Status       string       `gorm:"column:status;comment:执行状态" json:"status"`                   // 状态: processing, completed, failed
	Result       TaskResult   `gorm:"type:json;column:result;comment:执行结果" json:"result"`          // 执行结果，JSON格式
	ErrorMessage string       `gorm:"column:error_message;comment:错误信息" json:"error_message"`     // 错误信息
	StartedAt    time.Time    `gorm:"column:started_at;comment:开始时间" json:"started_at"`           // 开始执行时间
	CompletedAt  *time.Time   `gorm:"column:completed_at;comment:完成时间" json:"completed_at"`       // 完成时间
	ExecutionTime int         `gorm:"column:execution_time;comment:执行耗时(毫秒)" json:"execution_time"` // 执行耗时，单位毫秒
	
	// 外键关系
	Task         *Task        `json:"task,omitempty" gorm:"foreignKey:TaskID;references:TaskID"` // 关联的任务
	Worker       *Worker      `json:"worker,omitempty" gorm:"foreignKey:WorkerID;references:WorkerID"` // 关联的工作节点
}

// TableName 设置表名
func (TaskRecord) TableName() string {
	return "task_records"
}

// TaskResult 任务执行结果，使用JSON存储
type TaskResult map[string]interface{}

// Value 实现driver.Valuer接口
func (r TaskResult) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

// Scan 实现sql.Scanner接口
func (r *TaskResult) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, &r)
}
