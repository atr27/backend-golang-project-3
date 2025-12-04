package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func ValidationErrorResponse(c *gin.Context, errors []FieldError) {
	c.JSON(400, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid input data",
			Details: errors,
		},
	})
}

func PaginatedSuccessResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	c.JSON(200, PaginatedResponse{
		Success: true,
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}
