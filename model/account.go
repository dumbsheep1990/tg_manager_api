package model

// Account Telegram账号模型
type Account struct {
	BaseModel
	AccountGroupID  uint   `json:"account_group_id"`  // 账号分组ID
	Phone           string `json:"phone"`             // 电话号码
	Username        string `json:"username"`          // 用户名
	FirstName       string `json:"first_name"`        // 名字
	LastName        string `json:"last_name"`         // 姓氏
	Status          string `json:"status"`            // 状态: PENDING_IMPORT(导入中), ACTIVE(活跃), BANNED(被封禁), RESTRICTED(受限), IMPORT_FAILED(导入失败)
	TaskID          string `json:"task_id"`           // 任务ID
	TdataPath       string `json:"tdata_path"`        // Tdata文件路径
	TdataFilename   string `json:"tdata_filename"`    // Tdata文件名
	ErrorMessage    string `json:"error_message"`     // 错误信息
	LastLoginAt     int64  `json:"last_login_at"`     // 最后登录时间
	LastCheckAt     int64  `json:"last_check_at"`     // 最后检测时间
	CheckResult     string `json:"check_result"`      // 检测结果
	AccountLevel    int    `json:"account_level"`     // 账号等级：1-普通，2-中级，3-高级
	CreatedByUserID uint   `json:"created_by_user_id"` // 创建用户ID
	
	// 外键关系
	AccountGroup AccountGroup `json:"account_group" gorm:"foreignKey:AccountGroupID"` // 关联的账号分组
}
