package handlers

import (
	"hr-backend/internal/models"
	"hr-backend/internal/services"
	"hr-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	deptService *services.DepartmentService
}

func NewDepartmentHandler(deptService *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{deptService: deptService}
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var dept models.Department

	if err := c.ShouldBindJSON(&dept); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	created, err := h.deptService.CreateDepartment(&dept)
	if err != nil {
		utils.ErrorResponse(c, 400, "CREATE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 201, "Department created successfully", created)
}

func (h *DepartmentHandler) GetDepartments(c *gin.Context) {
	departments, err := h.deptService.GetDepartments()
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Departments retrieved successfully", departments)
}

func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid department ID")
		return
	}

	department, err := h.deptService.GetDepartmentByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 404, "NOT_FOUND", "Department not found")
		return
	}

	utils.SuccessResponse(c, 200, "Department retrieved successfully", department)
}

func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid department ID")
		return
	}

	var dept models.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	updated, err := h.deptService.UpdateDepartment(uint(id), &dept)
	if err != nil {
		utils.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Department updated successfully", updated)
}

func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid department ID")
		return
	}

	if err := h.deptService.DeleteDepartment(uint(id)); err != nil {
		utils.ErrorResponse(c, 400, "DELETE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Department deleted successfully", nil)
}
