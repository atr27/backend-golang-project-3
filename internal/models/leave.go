package models

import (
	"time"
)

type Leave struct {
	BaseModel
	EmployeeID uint       `gorm:"not null" json:"employee_id" binding:"required"`
	Employee   *Employee  `gorm:"constraint:OnDelete:CASCADE;" json:"employee,omitempty"`
	LeaveType  string     `gorm:"not null" json:"leave_type" binding:"required,oneof=annual sick casual unpaid"`
	StartDate  time.Time  `gorm:"type:date;not null" json:"start_date" binding:"required"`
	EndDate    time.Time  `gorm:"type:date;not null" json:"end_date" binding:"required"`
	TotalDays  int        `gorm:"not null" json:"total_days"`
	Reason     string     `json:"reason"`
	Status     string     `gorm:"default:'pending'" json:"status"`
	ApprovedBy *uint      `json:"approved_by"`
	Approver   *User      `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	ApprovedAt *time.Time `json:"approved_at"`
}

type LeaveBalance struct {
	BaseModel
	EmployeeID    uint      `gorm:"not null" json:"employee_id"`
	Employee      *Employee `gorm:"constraint:OnDelete:CASCADE;" json:"employee,omitempty"`
	LeaveType     string    `gorm:"not null" json:"leave_type"`
	TotalDays     int       `gorm:"not null" json:"total_days"`
	UsedDays      int       `gorm:"default:0" json:"used_days"`
	RemainingDays int       `gorm:"not null" json:"remaining_days"`
	Year          int       `gorm:"not null" json:"year"`
}

type CreateLeaveRequest struct {
	EmployeeID uint         `json:"employee_id" binding:"required"`
	LeaveType  string       `json:"leave_type" binding:"required,oneof=annual sick casual unpaid"`
	StartDate  FlexibleDate `json:"start_date" binding:"required"`
	EndDate    FlexibleDate `json:"end_date" binding:"required"`
	Reason     string       `json:"reason"`
}

type ApproveLeaveRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}
