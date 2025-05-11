package response

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code int         `json:"code"`    // 状态码
	Data interface{} `json:"data"`    // 数据
	Msg  string      `json:"message"` // 消息
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总数
	Page     int         `json:"page"`      // 页码
	PageSize int         `json:"page_size"` // 每页大小
}

const (
	SUCCESS = 0   // 成功
	ERROR   = 7   // 失败
)

// Result 返回结果
func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// Ok 成功返回
func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

// OkWithMessage 成功返回带消息
func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

// OkWithData 成功返回带数据
func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

// OkWithDetailed 成功返回带数据和消息
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

// Fail 失败返回
func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

// FailWithMessage 失败返回带消息
func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

// FailWithDetailed 失败返回带数据和消息
func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
