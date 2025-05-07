package tdata_account

// UpdateTdataAccountRequest 更新tdata账号请求
type UpdateTdataAccountRequest struct {
	AccountGroupID uint   `json:"account_group_id"` // 账号分组ID
	Status         string `json:"status"`           // 状态
	AccountLevel   int    `json:"account_level"`    // 账号等级：1-普通，2-中级，3-高级
}
