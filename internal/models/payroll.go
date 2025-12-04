package models

import (
	"time"
)

type Payroll struct {
	BaseModel
	EmployeeID  uint       `gorm:"not null" json:"employee_id" binding:"required"`
	Employee    *Employee  `gorm:"constraint:OnDelete:CASCADE;" json:"employee,omitempty"`
	Month       int        `gorm:"not null" json:"month" binding:"required,min=1,max=12"`
	Year        int        `gorm:"not null" json:"year" binding:"required"`
	BasicSalary float64    `gorm:"not null" json:"basic_salary" binding:"required"`
	Allowances  float64    `gorm:"default:0" json:"allowances"`
	Deductions  float64    `gorm:"default:0" json:"deductions"`
	Tax         float64    `gorm:"default:0" json:"tax"`
	NetSalary   float64    `gorm:"not null" json:"net_salary"`
	PaymentDate *time.Time `json:"payment_date"`
	Status      string     `gorm:"default:'pending'" json:"status"`
}

type GeneratePayrollRequest struct {
	Month int `json:"month" binding:"required,min=1,max=12"`
	Year  int `json:"year" binding:"required"`
}

type PayrollSummary struct {
	TotalEmployees  int     `json:"total_employees"`
	TotalBasicPay   float64 `json:"total_basic_pay"`
	TotalAllowances float64 `json:"total_allowances"`
	TotalDeductions float64 `json:"total_deductions"`
	TotalTax        float64 `json:"total_tax"`
	TotalNetPay     float64 `json:"total_net_pay"`
}
