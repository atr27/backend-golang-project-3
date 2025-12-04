package handlers

import (
	"hr-backend/internal/models"
	"hr-backend/internal/services"
	"hr-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	employeeService *services.EmployeeService
}

func NewEmployeeHandler(employeeService *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.CreateEmployeeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	employee, err := h.employeeService.CreateEmployee(&req)
	if err != nil {
		utils.ErrorResponse(c, 400, "CREATE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 201, "Employee created successfully", employee)
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	search := c.Query("search")

	var departmentID *uint
	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		if id, err := strconv.ParseUint(deptIDStr, 10, 32); err == nil {
			uid := uint(id)
			departmentID = &uid
		}
	}

	employees, total, err := h.employeeService.GetEmployees(page, limit, departmentID, status, search)
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, employees, total, page, limit)
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid employee ID")
		return
	}

	employee, err := h.employeeService.GetEmployeeByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 404, "NOT_FOUND", "Employee not found")
		return
	}

	utils.SuccessResponse(c, 200, "Employee retrieved successfully", employee)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid employee ID")
		return
	}

	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	employee, err := h.employeeService.UpdateEmployee(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Employee updated successfully", employee)
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid employee ID")
		return
	}

	if err := h.employeeService.DeleteEmployee(uint(id)); err != nil {
		utils.ErrorResponse(c, 400, "DELETE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Employee deleted successfully", nil)
}

func (h *EmployeeHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.employeeService.GetDashboardStats()
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Dashboard stats retrieved successfully", stats)
}

func (h *EmployeeHandler) GenerateEmployeeCode(c *gin.Context) {
	code, err := h.employeeService.GenerateEmployeeCode()
	if err != nil {
		utils.ErrorResponse(c, 500, "GENERATION_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Employee code generated", gin.H{"employee_code": code})
}
