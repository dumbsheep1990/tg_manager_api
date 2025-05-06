package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 所有表的公共字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`       // 主键自增
	CreatedAt time.Time      `json:"created_at"`                  // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                  // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`             // 删除时间（软删除）
}
