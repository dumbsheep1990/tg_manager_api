package account_group

// CreateAccountGroupRequest 创建账号分组请求
type CreateAccountGroupRequest struct {
	Name        string `json:"name" binding:"required"`        // 分组名称
	Description string `json:"description" binding:"required"` // 分组描述
}

// UpdateAccountGroupRequest 更新账号分组请求
type UpdateAccountGroupRequest struct {
	Name        string `json:"name"`        // 分组名称
	Description string `json:"description"` // 分组描述
	Status      string `json:"status"`      // 状态: ACTIVE(活跃), INACTIVE(停用)
}
