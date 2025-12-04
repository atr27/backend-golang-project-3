package services

import (
	"errors"
	"hr-backend/internal/models"
	"hr-backend/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type LeaveService struct {
	leaveRepo    *repositories.LeaveRepository
	employeeRepo *repositories.EmployeeRepository
}

func NewLeaveService(leaveRepo *repositories.LeaveRepository, employeeRepo *repositories.EmployeeRepository) *LeaveService {
	return &LeaveService{
		leaveRepo:    leaveRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *LeaveService) CreateLeave(req *models.CreateLeaveRequest) (*models.Leave, error) {
	// Check if employee exists
	_, err := s.employeeRepo.FindByID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	// Validate dates
	if req.EndDate.Time.Before(req.StartDate.Time) {
		return nil, errors.New("end date must be after start date")
	}

	// Calculate total days
	totalDays := int(req.EndDate.Time.Sub(req.StartDate.Time).Hours()/24) + 1

	// Check leave balance
	year := req.StartDate.Time.Year()
	balance, err := s.leaveRepo.FindLeaveBalance(req.EmployeeID, req.LeaveType, year)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if balance != nil && balance.RemainingDays < totalDays {
		return nil, errors.New("insufficient leave balance")
	}

	leave := &models.Leave{
		EmployeeID: req.EmployeeID,
		LeaveType:  req.LeaveType,
		StartDate:  req.StartDate.Time,
		EndDate:    req.EndDate.Time,
		TotalDays:  totalDays,
		Reason:     req.Reason,
		Status:     "pending",
	}

	if err := s.leaveRepo.Create(leave); err != nil {
		return nil, err
	}

	return s.leaveRepo.FindByID(leave.ID)
}

func (s *LeaveService) GetLeaves(employeeID *uint, status string, page, limit int) ([]models.Leave, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.leaveRepo.FindAll(employeeID, status, page, limit)
}

func (s *LeaveService) GetLeaveByID(id uint) (*models.Leave, error) {
	return s.leaveRepo.FindByID(id)
}

func (s *LeaveService) ApproveLeave(id uint, approverID uint, req *models.ApproveLeaveRequest) (*models.Leave, error) {
	leave, err := s.leaveRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if leave.Status != "pending" {
		return nil, errors.New("leave request is not pending")
	}

	leave.Status = req.Status
	leave.ApprovedBy = &approverID
	now := time.Now()
	leave.ApprovedAt = &now

	// Update leave balance if approved
	if req.Status == "approved" {
		year := leave.StartDate.Year()
		balance, err := s.leaveRepo.FindLeaveBalance(leave.EmployeeID, leave.LeaveType, year)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if balance != nil {
			balance.UsedDays += leave.TotalDays
			balance.RemainingDays = balance.TotalDays - balance.UsedDays
			if err := s.leaveRepo.UpdateLeaveBalance(balance); err != nil {
				return nil, err
			}
		}
	}

	if err := s.leaveRepo.Update(leave); err != nil {
		return nil, err
	}

	return s.leaveRepo.FindByID(id)
}

func (s *LeaveService) GetLeaveBalance(employeeID uint) ([]models.LeaveBalance, error) {
	year := time.Now().Year()
	balances, err := s.leaveRepo.FindAllLeaveBalances(employeeID, year)

	// Initialize default balances if not found
	if len(balances) == 0 {
		defaultLeaveTypes := []struct {
			Type string
			Days int
		}{
			{"annual", 15},
			{"sick", 10},
			{"casual", 7},
		}

		for _, lt := range defaultLeaveTypes {
			balance := &models.LeaveBalance{
				EmployeeID:    employeeID,
				LeaveType:     lt.Type,
				TotalDays:     lt.Days,
				UsedDays:      0,
				RemainingDays: lt.Days,
				Year:          year,
			}
			s.leaveRepo.CreateLeaveBalance(balance)
		}

		balances, err = s.leaveRepo.FindAllLeaveBalances(employeeID, year)
	}

	return balances, err
}
