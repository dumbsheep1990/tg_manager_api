package service

import (
	"context"
	"fmt"
	"os"
	"tg_manager_api/global"
	"tg_manager_api/model"
)

// AccountService 账号服务接口
type AccountService interface {
	// 创建账号分组
	CreateAccountGroup(ctx context.Context, group *model.AccountGroup) (uint, error)
	
	// 获取账号分组列表
	GetAccountGroups(ctx context.Context, page, pageSize int) ([]*model.AccountGroup, int64, error)
	
	// 获取账号分组详情
	GetAccountGroup(ctx context.Context, id uint) (*model.AccountGroup, error)
	
	// 更新账号分组
	UpdateAccountGroup(ctx context.Context, group *model.AccountGroup) error
	
	// 删除账号分组
	DeleteAccountGroup(ctx context.Context, id uint) error
	
	// 创建账号
	CreateAccount(ctx context.Context, account *model.Account) (uint, error)
	
	// 获取账号列表
	GetAccounts(ctx context.Context, groupID uint, page, pageSize int) ([]*model.Account, int64, error)
	
	// 获取账号详情
	GetAccount(ctx context.Context, id uint) (*model.Account, error)
	
	// 更新账号
	UpdateAccount(ctx context.Context, account *model.Account) error
	
	// 删除账号
	DeleteAccount(ctx context.Context, id uint) error
}

// NewAccountService 创建账号服务实例
func NewAccountService() AccountService {
	return &accountService{}
}

// accountService 账号服务实现
type accountService struct{}

// CreateAccountGroup 创建账号分组
func (s *accountService) CreateAccountGroup(ctx context.Context, group *model.AccountGroup) (uint, error) {
	result := global.DB.Create(group)
	if result.Error != nil {
		return 0, result.Error
	}
	return group.ID, nil
}

// GetAccountGroups 获取账号分组列表
func (s *accountService) GetAccountGroups(ctx context.Context, page, pageSize int) ([]*model.AccountGroup, int64, error) {
	var groups []*model.AccountGroup
	var total int64

	// 获取总数
	if err := global.DB.Model(&model.AccountGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := global.DB.Limit(pageSize).Offset(offset).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// GetAccountGroup 获取账号分组详情
func (s *accountService) GetAccountGroup(ctx context.Context, id uint) (*model.AccountGroup, error) {
	var group model.AccountGroup
	if err := global.DB.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateAccountGroup 更新账号分组
func (s *accountService) UpdateAccountGroup(ctx context.Context, group *model.AccountGroup) error {
	return global.DB.Save(group).Error
}

// DeleteAccountGroup 删除账号分组
func (s *accountService) DeleteAccountGroup(ctx context.Context, id uint) error {
	// 检查分组下是否有账号
	var count int64
	if err := global.DB.Model(&model.Account{}).Where("account_group_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("分组下存在账号，无法删除")
	}

	// 执行删除操作
	return global.DB.Delete(&model.AccountGroup{}, id).Error
}

// CreateAccount 创建账号
func (s *accountService) CreateAccount(ctx context.Context, account *model.Account) (uint, error) {
	result := global.DB.Create(account)
	if result.Error != nil {
		return 0, result.Error
	}
	return account.ID, nil
}

// GetAccounts 获取账号列表
func (s *accountService) GetAccounts(ctx context.Context, groupID uint, page, pageSize int) ([]*model.Account, int64, error) {
	var accounts []*model.Account
	var total int64

	// 构建查询
	query := global.DB.Model(&model.Account{}).Preload("AccountGroup")
	if groupID > 0 {
		query = query.Where("account_group_id = ?", groupID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// GetAccount 获取账号详情
func (s *accountService) GetAccount(ctx context.Context, id uint) (*model.Account, error) {
	var account model.Account
	if err := global.DB.Preload("AccountGroup").First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateAccount 更新账号
func (s *accountService) UpdateAccount(ctx context.Context, account *model.Account) error {
	return global.DB.Save(account).Error
}

// DeleteAccount 删除账号
func (s *accountService) DeleteAccount(ctx context.Context, id uint) error {
	// 获取账号信息以便删除文件
	var account model.Account
	if err := global.DB.First(&account, id).Error; err != nil {
		return err
	}

	// 开启事务
	tx := global.DB.Begin()
	
	// 删除账号记录
	if err := tx.Delete(&account).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 如果存在tdata文件，删除文件
	if account.TdataPath != "" {
		if err := os.Remove(account.TdataPath); err != nil && !os.IsNotExist(err) {
			tx.Rollback()
			return err
		}
	}
	
	// 提交事务
	return tx.Commit().Error
}
