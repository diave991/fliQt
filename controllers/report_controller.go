package controllers

import (
	"context"
	"net/http"
	"strconv"

	"fliQt/services"
	"github.com/gin-gonic/gin"
)

// ReportController 提供員工出缺勤報表查詢 API
type ReportController struct {
	service *services.ReportService
}

// NewReportController 建構函式
func NewReportController(s *services.ReportService) *ReportController {
	return &ReportController{service: s}
}

// RegisterRoutes 註冊路由到 gin Engine
func (c *ReportController) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/v1/reports", c.ListReports) // 全體員工分頁報表
	r.GET("/reports/:employee_id", c.GetReport)
}

// ListReports 回傳指定頁數的員工出缺勤報表，每頁 10 筆
func (c *ReportController) ListReports(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	data, err := c.service.GetAllReports(context.Background(), page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"page": page,
		"data": data,
	})
}

// GetReport 回傳指定員工的 7 天出缺勤報表
func (c *ReportController) GetReport(ctx *gin.Context) {
	empID, err := strconv.Atoi(ctx.Param("employee_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	// 從 Redis 讀取
	report, err := c.service.GetReport(context.Background(), uint(empID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, report)
}
