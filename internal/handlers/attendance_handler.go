package handlers

import (
	"hr-backend/internal/models"
	"hr-backend/internal/services"
	"hr-backend/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	attendanceService *services.AttendanceService
}

func NewAttendanceHandler(attendanceService *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attendanceService: attendanceService}
}

func (h *AttendanceHandler) ClockIn(c *gin.Context) {
	var req models.ClockInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	attendance, err := h.attendanceService.ClockIn(&req)
	if err != nil {
		utils.ErrorResponse(c, 400, "CLOCK_IN_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 201, "Clocked in successfully", attendance)
}

func (h *AttendanceHandler) ClockOut(c *gin.Context) {
	var req models.ClockOutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	attendance, err := h.attendanceService.ClockOut(&req)
	if err != nil {
		utils.ErrorResponse(c, 400, "CLOCK_OUT_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Clocked out successfully", attendance)
}

func (h *AttendanceHandler) GetAttendance(c *gin.Context) {
	employeeIDStr := c.Query("employee_id")
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_EMPLOYEE_ID", "Invalid employee ID")
		return
	}

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	attendances, err := h.attendanceService.GetAttendanceByEmployee(uint(employeeID), month, year)
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Attendance records retrieved successfully", attendances)
}

func (h *AttendanceHandler) GetAttendanceReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_DATE", "Invalid start date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_DATE", "Invalid end date format")
		return
	}

	attendances, err := h.attendanceService.GetAttendanceReport(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Attendance report retrieved successfully", attendances)
}

func (h *AttendanceHandler) CreateManualAttendance(c *gin.Context) {
	var attendance models.Attendance

	if err := c.ShouldBindJSON(&attendance); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := h.attendanceService.CreateManualAttendance(&attendance)
	if err != nil {
		utils.ErrorResponse(c, 400, "CREATE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 201, "Manual attendance created successfully", result)
}
