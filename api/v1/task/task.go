package task

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"tg_manager_api/global"
	"tg_manager_api/model/response"
	"tg_manager_api/services/task"
	"tg_manager_api/utils"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TaskType  string                 `json:"task_type" binding:"required"`  // 任务类型: send_message, join_group, add_contact等
	AccountID uint                   `json:"account_id" binding:"required"` // 关联的账号ID
	Params    map[string]interface{} `json:"params" binding:"required"`     // 任务参数，JSON格式
	Priority  *int                   `json:"priority"`                      // 优先级，可选
	TimeoutSec *int                  `json:"timeout_sec"`                   // 超时时间，可选
}

// TaskController 任务控制器
type TaskController struct{}

// CreateTask 创建任务
// @Summary 创建新任务
// @Description 创建一个新的Telegram任务
// @Tags Task
// @Accept json
// @Produce json
// @Param data body CreateTaskRequest true "创建任务的数据"
// @Success 200 {object} response.Response{data=model.Task} "创建成功"
// @Router /api/v1/task [post]
func (ctrl *TaskController) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}
	
	// 获取任务服务
	taskService := task.GetTaskServiceFromContext(c)
	
	// 创建任务
	newTask, err := taskService.CreateTask(c, req.TaskType, req.AccountID, req.Params)
	if err != nil {
		response.FailWithMessage("创建任务失败: "+err.Error(), c)
		return
	}
	
	// 如果提供了可选参数，更新任务
	if req.Priority != nil || req.TimeoutSec != nil {
		updates := map[string]interface{}{}
		
		if req.Priority != nil {
			updates["priority"] = *req.Priority
		}
		
		if req.TimeoutSec != nil {
			updates["timeout_sec"] = *req.TimeoutSec
		}
		
		// 如果有额外参数需要更新
		if len(updates) > 0 {
			if err := global.DB.Model(&newTask).Updates(updates).Error; err != nil {
				// 继续流程，仅记录更新失败，不影响任务创建
				global.LOG.Error("更新任务优先级或超时时间失败: " + err.Error())
			}
		}
	}
	
	response.OkWithData(newTask, c)
}

// GetTaskList 获取任务列表
// @Summary 获取任务列表
// @Description 分页获取任务列表
// @Tags Task
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Task}} "获取成功"
// @Router /api/v1/tasks [get]
func (ctrl *TaskController) GetTaskList(c *gin.Context) {
	// 获取分页参数
	page, pageSize := utils.GetPage(c)
	
	// 获取任务服务
	taskService := task.GetTaskServiceFromContext(c)
	
	// 获取任务列表
	tasks, total, err := taskService.GetTasks(c, page, pageSize)
	if err != nil {
		response.FailWithMessage("获取任务列表失败: "+err.Error(), c)
		return
	}
	
	response.OkWithDetailed(response.PageResult{
		List:     tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// GetTaskDetail 获取任务详情
// @Summary 获取任务详情
// @Description 根据任务ID获取任务详情
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} response.Response{data=model.Task} "获取成功"
// @Router /api/v1/task/{id} [get]
func (ctrl *TaskController) GetTaskDetail(c *gin.Context) {
	// 获取任务ID
	taskID := c.Param("id")
	
	// 获取任务服务
	taskService := task.GetTaskServiceFromContext(c)
	
	// 获取任务详情
	task, err := taskService.GetTask(c, taskID)
	if err != nil {
		response.FailWithMessage("获取任务详情失败: "+err.Error(), c)
		return
	}
	
	response.OkWithData(task, c)
}

// CancelTask 取消任务
// @Summary 取消任务
// @Description 取消指定的任务
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} response.Response "取消成功"
// @Router /api/v1/task/{id}/cancel [post]
func (ctrl *TaskController) CancelTask(c *gin.Context) {
	// 获取任务ID
	taskID := c.Param("id")
	
	// 获取任务服务
	taskService := task.GetTaskServiceFromContext(c)
	
	// 取消任务
	if err := taskService.CancelTask(c, taskID); err != nil {
		response.FailWithMessage("取消任务失败: "+err.Error(), c)
		return
	}
	
	response.OkWithMessage("取消任务成功", c)
}

// GetTasksByAccount 获取账号关联的任务
// @Summary 获取账号关联的任务
// @Description 获取指定账号关联的所有任务
// @Tags Task
// @Accept json
// @Produce json
// @Param account_id path uint true "账号ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Task}} "获取成功"
// @Router /api/v1/account/{account_id}/tasks [get]
func (ctrl *TaskController) GetTasksByAccount(c *gin.Context) {
	// 获取账号ID
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的账号ID", c)
		return
	}
	
	// 获取分页参数
	page, pageSize := utils.GetPage(c)
	
	// 查询账号关联的任务
	var tasks []*model.Task
	var total int64
	
	// 获取总记录数
	if err := global.DB.Model(&model.Task{}).
		Where("account_id = ?", accountID).
		Count(&total).Error; err != nil {
		response.FailWithMessage("获取任务总数失败: "+err.Error(), c)
		return
	}
	
	// 查询任务列表
	offset := (page - 1) * pageSize
	if err := global.DB.Where("account_id = ?", accountID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		response.FailWithMessage("获取任务列表失败: "+err.Error(), c)
		return
	}
	
	response.OkWithDetailed(response.PageResult{
		List:     tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}
