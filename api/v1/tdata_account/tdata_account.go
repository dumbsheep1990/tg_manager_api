package tdata_account

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"tg_manager_api/model"
	"tg_manager_api/services/account"
)

// ImportTdataAccount 导入tdata账号
// @Summary 导入tdata账号
// @Description 上传tdata文件并导入到系统中
// @Tags TdataAccount
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "tdata文件"
// @Param account_group_id formData int true "账号分组ID"
// @Param account_level formData int true "账号等级：1-普通，2-中级，3-高级" default(1)
// @Param status formData string true "账号状态" default("ACTIVE")
// @Success 200 {object} Response{data=model.Account} "成功"
// @Router /api/v1/tdata-accounts/import [post]
func ImportTdataAccount(c *gin.Context) {
	// 获取表单数据
	accountGroupID, err := strconv.Atoi(c.PostForm("account_group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的账号分组ID",
		})
		return
	}

	accountLevel, err := strconv.Atoi(c.PostForm("account_level"))
	if err != nil || accountLevel < 1 || accountLevel > 3 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的账号等级，应为1-3之间的整数",
		})
		return
	}

	status := c.PostForm("status")
	if status == "" {
		status = "ACTIVE"
	}

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 验证分组是否存在
	_, err = accountService.GetAccountGroup(context.Background(), uint(accountGroupID))
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

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "获取上传文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 创建存储目录
	uploadDir := "./storage/tdata"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "创建存储目录失败: " + err.Error(),
		})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "创建文件失败: " + err.Error(),
		})
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "保存文件失败: " + err.Error(),
		})
		return
	}

	// 准备账号数据
	account := &model.Account{
		AccountGroupID:  uint(accountGroupID),
		Status:          "PENDING_IMPORT", // 导入中状态
		TdataPath:       filePath,
		TdataFilename:   filename,
		AccountLevel:    accountLevel,
		CreatedByUserID: 1, // 假设当前用户ID为1，实际项目中应从认证中获取
	}

	// 调用服务层创建账号
	accountID, err := accountService.CreateAccount(context.Background(), account)
	if err != nil {
		// 如果创建记录失败，删除已上传的文件
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "创建账号记录失败: " + err.Error(),
		})
		return
	}

	// 获取创建后的账号完整对象
	createdAccount, err := accountService.GetAccount(context.Background(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取创建后的账号失败: " + err.Error(),
		})
		return
	}

	// 在实际项目中，这里应该发送一个消息到RabbitMQ，让Python Worker去处理tdata导入
	// 例如：sendImportTaskToQueue(accountID)

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "tdata账号文件上传成功，正在处理中",
		Data:    createdAccount,
	})
}

// GetTdataAccount 获取tdata账号详情
// @Summary 获取指定tdata账号详情
// @Description 根据ID获取tdata账号详情
// @Tags TdataAccount
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Success 200 {object} Response{data=model.Account} "成功"
// @Router /api/v1/tdata-accounts/{id} [get]
func GetTdataAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的账号ID",
		})
		return
	}

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 调用服务层获取账号详情
	tdataAccount, err := accountService.GetAccount(context.Background(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    http.StatusNotFound,
				Message: "账号不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "获取账号详情成功",
		Data:    tdataAccount,
	})
}

// ListTdataAccounts 获取tdata账号列表
// @Summary 获取tdata账号列表
// @Description 获取tdata账号列表，支持分页、搜索和分组筛选
// @Tags TdataAccount
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param account_group_id query int false "账号分组ID"
// @Param status query string false "账号状态"
// @Success 200 {object} Response{data=ListTdataAccountsResponse} "成功"
// @Router /api/v1/tdata-accounts [get]
func ListTdataAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	// keyword := c.Query("keyword") // TODO: 实现服务层的关键词搜索
	accountGroupID, _ := strconv.Atoi(c.Query("account_group_id"))
	// status := c.Query("status") // TODO: 实现服务层的状态筛选

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 调用服务层获取账号列表
	// 注意: 这里目前只提供了按分组查询的功能，其他筛选条件需要在服务层扩展
	accounts, total, err := accountService.GetAccounts(context.Background(), uint(accountGroupID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号列表失败: " + err.Error(),
		})
		return
	}

	// 转换指针列表为对象列表
	var accountList []model.Account
	for _, a := range accounts {
		accountList = append(accountList, *a)
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "获取账号列表成功",
		Data: ListTdataAccountsResponse{
			Total:    total,
			Accounts: accountList,
		},
	})
}

// UpdateTdataAccount 更新tdata账号
// @Summary 更新tdata账号信息
// @Description 更新指定tdata账号的信息
// @Tags TdataAccount
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param data body UpdateTdataAccountRequest true "更新的账号信息"
// @Success 200 {object} Response{data=model.Account} "成功"
// @Router /api/v1/tdata-accounts/{id} [put]
func UpdateTdataAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的账号ID",
		})
		return
	}

	var req UpdateTdataAccountRequest
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

	// 首先获取当前账号
	tdataAccount, err := accountService.GetAccount(context.Background(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    http.StatusNotFound,
				Message: "账号不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号失败: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.AccountGroupID > 0 {
		// 检查账号分组是否存在
		_, err := accountService.GetAccountGroup(context.Background(), req.AccountGroupID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, Response{
					Code:    http.StatusBadRequest,
					Message: "账号分组不存在",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, Response{
				Code:    http.StatusInternalServerError,
				Message: "检查分组失败: " + err.Error(),
			})
			return
		}
		tdataAccount.AccountGroupID = req.AccountGroupID
	}
	
	if req.Status != "" {
		tdataAccount.Status = req.Status
	}
	
	if req.AccountLevel > 0 && req.AccountLevel <= 3 {
		tdataAccount.AccountLevel = req.AccountLevel
	}

	// 调用服务层更新账号
	if err := accountService.UpdateAccount(context.Background(), tdataAccount); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "更新账号失败: " + err.Error(),
		})
		return
	}

	// 重新获取更新后的账号信息
	updatedAccount, err := accountService.GetAccount(context.Background(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取更新后的账号失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "更新账号成功",
		Data:    updatedAccount,
	})
}

// DeleteTdataAccount 删除tdata账号
// @Summary 删除tdata账号
// @Description 删除指定ID的tdata账号
// @Tags TdataAccount
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Success 200 {object} Response "成功"
// @Router /api/v1/tdata-accounts/{id} [delete]
func DeleteTdataAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "无效的账号ID",
		})
		return
	}

	// 创建服务实例
	serviceFactory := account.NewServiceFactory()
	accountService := serviceFactory.AccountService()

	// 验证账号是否存在
	_, err = accountService.GetAccount(context.Background(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    http.StatusNotFound,
				Message: "账号不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "获取账号失败: " + err.Error(),
		})
		return
	}

	// 调用服务层删除账号
	if err := accountService.DeleteAccount(context.Background(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "删除账号失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "删除账号成功",
	})
}
