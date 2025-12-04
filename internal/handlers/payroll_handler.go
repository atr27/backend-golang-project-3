package handlers

import (
	"hr-backend/internal/models"
	"hr-backend/internal/services"
	"hr-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	payrollService *services.PayrollService
}

func NewPayrollHandler(payrollService *services.PayrollService) *PayrollHandler {
	return &PayrollHandler{payrollService: payrollService}
}

func (h *PayrollHandler) GeneratePayroll(c *gin.Context) {
	var req models.GeneratePayrollRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	payrolls, err := h.payrollService.GeneratePayroll(&req)
	if err != nil {
		utils.ErrorResponse(c, 400, "GENERATION_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 201, "Payroll generated successfully", payrolls)
}

func (h *PayrollHandler) GetPayrolls(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	payrolls, total, err := h.payrollService.GetPayrolls(month, year, page, limit)
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, payrolls, total, page, limit)
}

func (h *PayrollHandler) GetPayrollByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid payroll ID")
		return
	}

	payroll, err := h.payrollService.GetPayrollByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 404, "NOT_FOUND", "Payroll not found")
		return
	}

	utils.SuccessResponse(c, 200, "Payroll retrieved successfully", payroll)
}

func (h *PayrollHandler) UpdatePayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid payroll ID")
		return
	}

	var payroll models.Payroll
	if err := c.ShouldBindJSON(&payroll); err != nil {
		utils.ErrorResponse(c, 400, "VALIDATION_ERROR", err.Error())
		return
	}

	updated, err := h.payrollService.UpdatePayroll(uint(id), &payroll)
	if err != nil {
		utils.ErrorResponse(c, 400, "UPDATE_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Payroll updated successfully", updated)
}

func (h *PayrollHandler) ProcessPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid payroll ID")
		return
	}

	payroll, err := h.payrollService.ProcessPayment(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 400, "PAYMENT_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Payment processed successfully", payroll)
}

func (h *PayrollHandler) GetPayrollSummary(c *gin.Context) {
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	if month == 0 || year == 0 {
		utils.ErrorResponse(c, 400, "INVALID_PARAMETERS", "Month and year are required")
		return
	}

	summary, err := h.payrollService.GetPayrollSummary(month, year)
	if err != nil {
		utils.ErrorResponse(c, 500, "FETCH_FAILED", err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "Payroll summary retrieved successfully", summary)
}

func (h *PayrollHandler) DownloadPayrollSlip(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "INVALID_ID", "Invalid payroll ID")
		return
	}

	pdfBytes, err := h.payrollService.GeneratePayrollPDF(uint(id))
	if err != nil {
		utils.ErrorResponse(c, 400, "PDF_GENERATION_FAILED", err.Error())
		return
	}

	// Set headers for PDF download
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=slip_gaji.pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	c.Data(200, "application/pdf", pdfBytes)
}
