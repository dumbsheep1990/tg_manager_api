package tdata_account

import "tg_manager_api/model"

// Response 通用响应格式
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// ListTdataAccountsResponse 获取tdata账号列表响应
type ListTdataAccountsResponse struct {
	Total    int64           `json:"total"`    // 总数
	Accounts []model.Account `json:"accounts"` // 账号列表
}
