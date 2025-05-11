package worker

import (
	"github.com/gin-gonic/gin"
	
	"tg_manager_api/model"
	"tg_manager_api/model/response"
	"tg_manager_api/services/worker"
	"tg_manager_api/utils"
)

// RegisterWorkerRequest 注册工作节点请求
type RegisterWorkerRequest struct {
	Hostname string `json:"hostname" binding:"required"` // 主机名
	IP       string `json:"ip" binding:"required"`       // IP地址
	MaxTasks int    `json:"max_tasks"`                   // 最大可同时执行的任务数
	Tags     string `json:"tags"`                        // 标签，用逗号分隔
	Version  string `json:"version"`                     // Worker版本
}

// HeartbeatRequest 心跳请求
type HeartbeatRequest struct {
	WorkerID string `json:"worker_id" binding:"required"` // 工作节点ID
}

// WorkerController 工作节点控制器
type WorkerController struct{}

// RegisterWorker 注册工作节点
// @Summary 注册工作节点
// @Description 注册一个新的工作节点或重新激活已存在的节点
// @Tags Worker
// @Accept json
// @Produce json
// @Param data body RegisterWorkerRequest true "工作节点注册数据"
// @Success 200 {object} response.Response{data=string} "注册成功，返回worker_id"
// @Router /api/v1/worker/register [post]
func (ctrl *WorkerController) RegisterWorker(c *gin.Context) {
	var req RegisterWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}
	
	// 设置默认值
	if req.MaxTasks <= 0 {
		req.MaxTasks = 10 // 默认最大任务数
	}
	
	// 获取工作节点服务
	workerService := worker.GetWorkerServiceFromContext(c)
	
	// 注册工作节点
	workerID, err := workerService.RegisterWorker(c, req.Hostname, req.IP, req.MaxTasks, req.Tags)
	if err != nil {
		response.FailWithMessage("注册工作节点失败: "+err.Error(), c)
		return
	}
	
	response.OkWithData(workerID, c)
}

// Heartbeat 工作节点心跳
// @Summary 工作节点心跳
// @Description 更新工作节点的心跳时间，保持节点在线状态
// @Tags Worker
// @Accept json
// @Produce json
// @Param data body HeartbeatRequest true "心跳请求数据"
// @Success 200 {object} response.Response "心跳成功"
// @Router /api/v1/worker/heartbeat [post]
func (ctrl *WorkerController) Heartbeat(c *gin.Context) {
	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}
	
	// 获取工作节点服务
	workerService := worker.GetWorkerServiceFromContext(c)
	
	// 更新心跳
	if err := workerService.UpdateHeartbeat(c, req.WorkerID); err != nil {
		response.FailWithMessage("更新心跳失败: "+err.Error(), c)
		return
	}
	
	response.OkWithMessage("心跳更新成功", c)
}

// GetWorkerList 获取工作节点列表
// @Summary 获取工作节点列表
// @Description 分页获取所有工作节点
// @Tags Worker
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Worker}} "获取成功"
// @Router /api/v1/workers [get]
func (ctrl *WorkerController) GetWorkerList(c *gin.Context) {
	// 获取分页参数
	page, pageSize := utils.GetPage(c)
	
	// 获取工作节点服务
	workerService := worker.GetWorkerServiceFromContext(c)
	
	// 获取工作节点列表
	workers, total, err := workerService.GetAllWorkers(c, page, pageSize)
	if err != nil {
		response.FailWithMessage("获取工作节点列表失败: "+err.Error(), c)
		return
	}
	
	response.OkWithDetailed(response.PageResult{
		List:     workers,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// GetWorkerDetail 获取工作节点详情
// @Summary 获取工作节点详情
// @Description 根据ID获取工作节点的详细信息
// @Tags Worker
// @Accept json
// @Produce json
// @Param id path string true "工作节点ID"
// @Success 200 {object} response.Response{data=model.Worker} "获取成功"
// @Router /api/v1/worker/{id} [get]
func (ctrl *WorkerController) GetWorkerDetail(c *gin.Context) {
	// 获取工作节点ID
	workerID := c.Param("id")
	
	// 获取工作节点服务
	workerService := worker.GetWorkerServiceFromContext(c)
	
	// 获取工作节点状态
	worker, err := workerService.GetWorkerStatus(c, workerID)
	if err != nil {
		response.FailWithMessage("获取工作节点详情失败: "+err.Error(), c)
		return
	}
	
	response.OkWithData(worker, c)
}

// GetWorkerTasks 获取工作节点的任务列表
// @Summary 获取工作节点的任务列表
// @Description 获取指定工作节点正在执行或已执行的任务
// @Tags Worker
// @Accept json
// @Produce json
// @Param id path string true "工作节点ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Task}} "获取成功"
// @Router /api/v1/worker/{id}/tasks [get]
func (ctrl *WorkerController) GetWorkerTasks(c *gin.Context) {
	// 获取工作节点ID
	workerID := c.Param("id")
	
	// 获取分页参数
	page, pageSize := utils.GetPage(c)
	
	// 获取工作节点服务
	workerService := worker.GetWorkerServiceFromContext(c)
	
	// 获取工作节点任务列表
	tasks, total, err := workerService.GetWorkerTasks(c, workerID, page, pageSize)
	if err != nil {
		response.FailWithMessage("获取工作节点任务列表失败: "+err.Error(), c)
		return
	}
	
	response.OkWithDetailed(response.PageResult{
		List:     tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}
