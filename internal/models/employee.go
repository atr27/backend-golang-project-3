package models

import (
	"time"
)

type Employee struct {
	BaseModel
	UserID           *uint       `gorm:"uniqueIndex" json:"user_id"`
	User             *User       `gorm:"constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	EmployeeCode     string      `gorm:"uniqueIndex;not null" json:"employee_code" binding:"required"`
	FirstName        string      `gorm:"not null" json:"first_name" binding:"required"`
	LastName         string      `gorm:"not null" json:"last_name" binding:"required"`
	DateOfBirth      *time.Time  `json:"date_of_birth"`
	Gender           string      `json:"gender"`
	Phone            string      `json:"phone"`
	Address          string      `json:"address"`
	DepartmentID     *uint       `json:"department_id"`
	Department       *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Position         string      `json:"position"`
	HireDate         time.Time   `gorm:"not null" json:"hire_date" binding:"required"`
	EmploymentStatus string      `gorm:"default:'active'" json:"employment_status"`
	Salary           float64     `json:"salary"`
	ProfilePicture   string      `json:"profile_picture"`
}

type CreateEmployeeRequest struct {
	Email            string     `json:"email" binding:"required,email"`
	Password         string     `json:"password" binding:"required,min=6"`
	EmployeeCode     string     `json:"employee_code" binding:"required"`
	FirstName        string     `json:"first_name" binding:"required"`
	LastName         string     `json:"last_name" binding:"required"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Gender           string     `json:"gender"`
	Phone            string     `json:"phone"`
	Address          string     `json:"address"`
	DepartmentID     *uint      `json:"department_id"`
	Position         string     `json:"position"`
	HireDate         time.Time  `json:"hire_date" binding:"required"`
	EmploymentStatus string     `json:"employment_status"`
	Salary           float64    `json:"salary"`
	Role             string     `json:"role" binding:"required,oneof=admin hr_manager department_manager employee"`
}

type UpdateEmployeeRequest struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Gender           string     `json:"gender"`
	Phone            string     `json:"phone"`
	Address          string     `json:"address"`
	DepartmentID     *uint      `json:"department_id"`
	Position         string     `json:"position"`
	EmploymentStatus string     `json:"employment_status"`
	Salary           float64    `json:"salary"`
}
