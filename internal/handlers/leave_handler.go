package handlers

import (
	"hr-backend/internal/models"
	"hr-backend/internal/services"
	"hr-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	leaveService *services.LeaveService
}

func NewLeaveHandler(leaveService *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{leaveService: leaveService}
}

func (h *LeaveHandler) CreateLeave(c *gin.Context) {
	var req models.CreateLeaveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	leave, err := h.leaveService.CreateLeave(&req)
	if err != nil {
		utils.ErrorResponse(c, 400, "CREATE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 201, "Leave request created successfully", leave)
}

func (h *LeaveHandler) GetLeaves(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	var employeeID *uint
	if empIDStr := c.Query("employee_id"); empIDStr != "" {
		if id, err := strconv.ParseUint(empIDStr, 10, 32); err == nil {
			uid := uint(id)
			employeeID = &uid
		}
	}

	leaves, total, err := h.leaveService.GetLeaves(employeeID, status, page, limit)
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, leaves, total, page, limit)
}

func (h *LeaveHandler) GetLeaveByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid leave ID")
		return
	}

	leave, err := h.leaveService.GetLeaveByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 404, "NOT_FOUND", "Leave not found")
		return
	}

	utils.SuccessResponse(c, 200, "Leave retrieved successfully", leave)
}

func (h *LeaveHandler) ApproveLeave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid leave ID")
		return
	}

	var req models.ApproveLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	leave, err := h.leaveService.ApproveLeave(uint(id), userID.(uint), &req)
	if err != nil {
		utils.ErrorResponse(c, 400, "APPROVAL_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Leave request updated successfully", leave)
}

func (h *LeaveHandler) GetLeaveBalance(c *gin.Context) {
	employeeIDStr := c.Param("employee_id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid employee ID")
		return
	}

	balances, err := h.leaveService.GetLeaveBalance(uint(employeeID))
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Leave balance retrieved successfully", balances)
}
