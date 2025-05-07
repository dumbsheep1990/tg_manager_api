package account_group

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"tg_manager_api/model"
	"tg_manager_api/services/account"
)

// CreateAccountGroup 创建Telegram账号分组
// @Summary 创建Telegram账号分组
// @Description 创建新的Telegram账号分组
// @Tags AccountGroup
// @Accept json
// @Produce json
// @Param data body CreateAccountGroupRequest true "分组信息"
// @Success 200 {object} Response{data=model.AccountGroup} "成功"
// @Router /api/v1/account-groups [post]
func CreateAccountGroup(c *gin.Context) {
	var req CreateAccountGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 准备分组数据
	accountGroup := &model.AccountGroup{
		Name:        req.Name,
		Description: req.Description,
		Status:      "ACTIVE",
	}

	// 调用服务层创建分组
	id, err := accountService.CreateAccountGroup(context.Background(), accountGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "创建账号分组失败: " + err.Error(),
		})
		return
	}

	// 获取创建后的完整对象
	group, err := accountService.GetAccountGroup(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取创建后的账号分组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "创建账号分组成功",
		Data:    group,
	})
}

// GetAccountGroup 获取账号分组详情
// @Summary 获取指定账号分组详情
// @Description 根据ID获取Telegram账号分组详情
// @Tags AccountGroup
// @Accept json
// @Produce json
// @Param id path int true "分组ID"
// @Success 200 {object} Response{data=model.AccountGroup} "成功"
// @Router /api/v1/account-groups/{id} [get]
func GetAccountGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的分组ID",
		})
		return
	}

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 调用服务层获取分组详情
	group, err := accountService.GetAccountGroup(context.Background(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    http.StatusNotFound,
				Message: "账号分组不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号分组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "获取账号分组成功",
		Data:    group,
	})
}

// UpdateAccountGroup 更新账号分组
// @Summary 更新账号分组信息
// @Description 更新指定账号分组的信息
// @Tags AccountGroup
// @Accept json
// @Produce json
// @Param id path int true "分组ID"
// @Param data body UpdateAccountGroupRequest true "更新的分组信息"
// @Success 200 {object} Response{data=model.AccountGroup} "成功"
// @Router /api/v1/account-groups/{id} [put]
func UpdateAccountGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的分组ID",
		})
		return
	}

	var req UpdateAccountGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 首先获取分组
	group, err := accountService.GetAccountGroup(context.Background(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    http.StatusNotFound,
				Message: "账号分组不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号分组失败: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Description != "" {
		group.Description = req.Description
	}
	if req.Status != "" {
		group.Status = req.Status
	}

	// 调用服务层更新分组
	if err := accountService.UpdateAccountGroup(context.Background(), group); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "更新账号分组失败: " + err.Error(),
		})
		return
	}

	// 重新获取更新后的分组
	group, err = accountService.GetAccountGroup(context.Background(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取更新后的账号分组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "更新账号分组成功",
		Data:    group,
	})
}

// DeleteAccountGroup 删除账号分组
// @Summary 删除账号分组
// @Description 删除指定ID的账号分组
// @Tags AccountGroup
// @Accept json
// @Produce json
// @Param id path int true "分组ID"
// @Success 200 {object} Response "成功"
// @Router /api/v1/account-groups/{id} [delete]
func DeleteAccountGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的分组ID",
		})
		return
	}

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 首先验证分组是否存在
	_, err = accountService.GetAccountGroup(context.Background(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    http.StatusNotFound,
				Message: "账号分组不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号分组失败: " + err.Error(),
		})
		return
	}

	// 调用服务层删除分组
	if err = accountService.DeleteAccountGroup(context.Background(), uint(id)); err != nil {
		// 判断是否是特定的错误类型
		if err.Error() == "分组下存在账号，无法删除" {
			c.JSON(http.StatusBadRequest, Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "删除账号分组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "删除账号分组成功",
	})
}

// ListAccountGroups 获取账号分组列表
// @Summary 获取账号分组列表
// @Description 获取Telegram账号分组列表，支持分页和搜索
// @Tags AccountGroup
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} Response{data=ListAccountGroupsResponse} "成功"
// @Router /api/v1/account-groups [get]
func ListAccountGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	// keyword := c.Query("keyword") // TODO: 实现服务层的关键词搜索

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 调用服务层获取分组列表
	groups, total, err := accountService.GetAccountGroups(context.Background(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号分组列表失败: " + err.Error(),
		})
		return
	}

	// 转换指针列表为对象列表
	var accountGroups []model.AccountGroup
	for _, g := range groups {
		// 确保每个分组的Accounts字段为空，避免数据过大
		g.Accounts = nil 
		accountGroups = append(accountGroups, *g)
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "获取账号分组列表成功",
		Data: ListAccountGroupsResponse{
			Total:        total,
			AccountGroups: accountGroups,
		},
	})
}
