package account_group

import "tg_manager_api/model"

// Response 通用响应格式
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// ListAccountGroupsResponse 获取账号分组列表响应
type ListAccountGroupsResponse struct {
	Total         int64               `json:"total"`          // 总数
	AccountGroups []model.AccountGroup `json:"account_groups"` // 账号分组列表
}
