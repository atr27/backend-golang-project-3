package repositories

import (
	"hr-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type LeaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

func (r *LeaveRepository) Create(leave *models.Leave) error {
	return r.db.Create(leave).Error
}

func (r *LeaveRepository) FindAll(employeeID *uint, status string, page, limit int) ([]models.Leave, int64, error) {
	var leaves []models.Leave
	var total int64

	query := r.db.Model(&models.Leave{}).Preload("Employee").Preload("Approver")

	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&leaves).Error

	return leaves, total, err
}

func (r *LeaveRepository) FindByID(id uint) (*models.Leave, error) {
	var leave models.Leave
	err := r.db.Preload("Employee").Preload("Approver").First(&leave, id).Error
	return &leave, err
}

func (r *LeaveRepository) Update(leave *models.Leave) error {
	return r.db.Save(leave).Error
}

func (r *LeaveRepository) Delete(id uint) error {
	return r.db.Delete(&models.Leave{}, id).Error
}

func (r *LeaveRepository) FindLeaveBalance(employeeID uint, leaveType string, year int) (*models.LeaveBalance, error) {
	var balance models.LeaveBalance
	err := r.db.Where("employee_id = ? AND leave_type = ? AND year = ?", employeeID, leaveType, year).
		First(&balance).Error
	return &balance, err
}

func (r *LeaveRepository) CreateLeaveBalance(balance *models.LeaveBalance) error {
	return r.db.Create(balance).Error
}

func (r *LeaveRepository) UpdateLeaveBalance(balance *models.LeaveBalance) error {
	return r.db.Save(balance).Error
}

func (r *LeaveRepository) FindAllLeaveBalances(employeeID uint, year int) ([]models.LeaveBalance, error) {
	var balances []models.LeaveBalance
	err := r.db.Where("employee_id = ? AND year = ?", employeeID, year).Find(&balances).Error
	return balances, err
}

func (r *LeaveRepository) CountPendingLeaves() (int64, error) {
	var count int64
	err := r.db.Model(&models.Leave{}).Where("status = ?", "pending").Count(&count).Error
	return count, err
}

func (r *LeaveRepository) CountTodayOnLeave() (int64, error) {
	today := time.Now()
	var count int64
	err := r.db.Model(&models.Leave{}).
		Where("status = ? AND ? BETWEEN start_date AND end_date", "approved", today).
		Count(&count).Error
	return count, err
}
