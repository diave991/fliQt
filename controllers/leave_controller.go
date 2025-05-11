package controllers

import (
	"fliQt/models"
	"fliQt/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LeaveController struct {
	Service *services.LeaveService
}

func NewLeaveController(service *services.LeaveService) *LeaveController {
	return &LeaveController{Service: service}
}

func (c *LeaveController) Create(ctx *gin.Context) {
	var leave models.Leave
	if err := ctx.ShouldBindJSON(&leave); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.Service.Create(&leave); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, leave)
}

func (c *LeaveController) GetAll(ctx *gin.Context) {
	leaves, err := c.Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, leaves)
}
