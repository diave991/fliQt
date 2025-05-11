package controllers

import (
	"net/http"
	"strconv"
	"time"

	"fliQt/models"
	"fliQt/services"
	"github.com/gin-gonic/gin"
)

// AttendanceController 處理打卡相關 HTTP 請求
// 提供新增、查詢等多種操作

type AttendanceController struct {
	Service *services.AttendanceService
}

// NewAttendanceController 建立新的 AttendanceController 實例
func NewAttendanceController(service *services.AttendanceService) *AttendanceController {
	return &AttendanceController{Service: service}
}

// Create 新增一筆打卡記錄
func (c *AttendanceController) Create(ctx *gin.Context) {
	var att models.Attendance
	if err := ctx.ShouldBindJSON(&att); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 若未指定 Timestamp，預設為當前時間
	if att.Timestamp.IsZero() {
		att.Timestamp = time.Now()
	}
	if err := c.Service.Create(&att); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, att)
}

// GetAll 取得所有打卡記錄
func (c *AttendanceController) GetAll(ctx *gin.Context) {
	records, err := c.Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, records)
}

// GetByID 根據記錄 ID 取得單筆打卡記錄
func (c *AttendanceController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	record, err := c.Service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, record)
}

// GetByEmployee 依員工編號及可選日期查詢打卡記錄
// 若未帶 date 參數，回傳員工所有記錄；
// 若帶 date，格式應為 YYYY-MM-DD，僅回傳該日期的記錄
func (c *AttendanceController) GetByEmployee(ctx *gin.Context) {
	empIDParam := ctx.Query("employee_id")
	empIDUint64, err := strconv.ParseUint(empIDParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}
	dateParam := ctx.Query("date")
	if dateParam != "" {
		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		records, err := c.Service.GetByEmployeeAndDate(uint(empIDUint64), date)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, records)
		return
	}
	records, err := c.Service.GetByEmployee(uint(empIDUint64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, records)
}
