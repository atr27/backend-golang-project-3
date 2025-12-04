package repositories

import (
	"hr-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepository) FindByEmployeeAndDate(employeeID uint, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.Where("employee_id = ? AND date = ?", employeeID, date).First(&attendance).Error
	return &attendance, err
}

func (r *AttendanceRepository) FindByEmployee(employeeID uint, startDate, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := r.db.Where("employee_id = ? AND date BETWEEN ? AND ?", employeeID, startDate, endDate).
		Order("date DESC").
		Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) FindByDateRange(startDate, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := r.db.Preload("Employee").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date DESC").
		Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) Update(attendance *models.Attendance) error {
	return r.db.Save(attendance).Error
}

func (r *AttendanceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Attendance{}, id).Error
}

func (r *AttendanceRepository) CountByStatus(status string, date time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Attendance{}).
		Where("status = ? AND date = ?", status, date).
		Count(&count).Error
	return count, err
}

func (r *AttendanceRepository) CountTodayPresent() (int64, error) {
	today := time.Now().Format("2006-01-02")
	var count int64
	err := r.db.Model(&models.Attendance{}).
		Where("date = ? AND status = ?", today, "present").
		Count(&count).Error
	return count, err
}
