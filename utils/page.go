package utils

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
)

// GetPage 获取分页参数
// 如果请求中没有提供分页参数，则使用默认值
func GetPage(c *gin.Context) (int, int) {
	// 默认值
	defaultPage := 1
	defaultPageSize := 10
	
	// 从请求中获取页码
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}
	
	// 从请求中获取每页大小
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = defaultPageSize
	}
	
	return page, pageSize
}

// GetOffset 获取偏移量
func GetOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 10
	}
	
	return (page - 1) * pageSize
}
