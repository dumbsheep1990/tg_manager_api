package model

// AccountGroup Telegram账号分组模型
type AccountGroup struct {
	BaseModel
	Name        string `json:"name"`        // 分组名称
	Description string `json:"description"` // 分组描述
	Status      string `json:"status"`      // 状态: ACTIVE(活跃), INACTIVE(停用)
	
	// 关联关系
	Accounts []Account `json:"accounts" gorm:"foreignKey:AccountGroupID"` // 分组下的账号列表
}
