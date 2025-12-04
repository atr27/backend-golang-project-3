package services

import (
	"errors"
	"fmt"
	"hr-backend/internal/models"
	"hr-backend/internal/repositories"
	"hr-backend/internal/utils"

	"gorm.io/gorm"
)

type EmployeeService struct {
	employeeRepo *repositories.EmployeeRepository
	userRepo     *repositories.UserRepository
	db           *gorm.DB
}

func NewEmployeeService(employeeRepo *repositories.EmployeeRepository, userRepo *repositories.UserRepository, db *gorm.DB) *EmployeeService {
	return &EmployeeService{
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
		db:           db,
	}
}

func (s *EmployeeService) CreateEmployee(req *models.CreateEmployeeRequest) (*models.Employee, error) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Check if employee code already exists
	_, err = s.employeeRepo.FindByEmployeeCode(req.EmployeeCode)
	if err == nil {
		return nil, errors.New("employee code already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user and employee in transaction
	var employee *models.Employee
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create user
		user := &models.User{
			Email:        req.Email,
			PasswordHash: hashedPassword,
			Role:         req.Role,
			IsActive:     true,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Create employee
		employee = &models.Employee{
			UserID:           &user.ID,
			EmployeeCode:     req.EmployeeCode,
			FirstName:        req.FirstName,
			LastName:         req.LastName,
			DateOfBirth:      req.DateOfBirth,
			Gender:           req.Gender,
			Phone:            req.Phone,
			Address:          req.Address,
			DepartmentID:     req.DepartmentID,
			Position:         req.Position,
			HireDate:         req.HireDate,
			EmploymentStatus: req.EmploymentStatus,
			Salary:           req.Salary,
		}

		if employee.EmploymentStatus == "" {
			employee.EmploymentStatus = "active"
		}

		return tx.Create(employee).Error
	})

	if err != nil {
		return nil, err
	}

	return s.employeeRepo.FindByID(employee.ID)
}

func (s *EmployeeService) GetEmployees(page, limit int, departmentID *uint, status, search string) ([]models.Employee, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.employeeRepo.FindAll(page, limit, departmentID, status, search)
}

func (s *EmployeeService) GetEmployeeByID(id uint) (*models.Employee, error) {
	return s.employeeRepo.FindByID(id)
}

func (s *EmployeeService) UpdateEmployee(id uint, req *models.UpdateEmployeeRequest) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		employee.FirstName = req.FirstName
	}
	if req.LastName != "" {
		employee.LastName = req.LastName
	}
	if req.DateOfBirth != nil {
		employee.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		employee.Gender = req.Gender
	}
	if req.Phone != "" {
		employee.Phone = req.Phone
	}
	if req.Address != "" {
		employee.Address = req.Address
	}
	if req.DepartmentID != nil {
		employee.DepartmentID = req.DepartmentID
	}
	if req.Position != "" {
		employee.Position = req.Position
	}
	if req.EmploymentStatus != "" {
		employee.EmploymentStatus = req.EmploymentStatus
	}
	if req.Salary > 0 {
		employee.Salary = req.Salary
	}

	err = s.employeeRepo.Update(employee)
	if err != nil {
		return nil, err
	}

	return s.employeeRepo.FindByID(id)
}

func (s *EmployeeService) DeleteEmployee(id uint) error {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete employee
		if err := tx.Delete(employee).Error; err != nil {
			return err
		}

		// Delete associated user
		if employee.UserID != nil {
			if err := tx.Delete(&models.User{}, *employee.UserID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *EmployeeService) GetDashboardStats() (map[string]interface{}, error) {
	totalEmployees, err := s.employeeRepo.CountByStatus("active")
	if err != nil {
		return nil, err
	}

	// Get today's attendance count
	var presentToday int64
	if err := s.db.Model(&models.Attendance{}).
		Where("date = CURRENT_DATE AND status = ?", "present").
		Count(&presentToday).Error; err != nil {
		return nil, err
	}

	// Get today's leave count
	var onLeaveToday int64
	if err := s.db.Model(&models.Leave{}).
		Where("status = ? AND CURRENT_DATE BETWEEN start_date AND end_date", "approved").
		Count(&onLeaveToday).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_employees": totalEmployees,
		"present_today":   presentToday,
		"on_leave_today":  onLeaveToday,
		"absent_today":    totalEmployees - presentToday - onLeaveToday,
	}

	return stats, nil
}

func (s *EmployeeService) GenerateEmployeeCode() (string, error) {
	var count int64
	s.db.Model(&models.Employee{}).Count(&count)
	return fmt.Sprintf("EMP%04d", count+1), nil
}
