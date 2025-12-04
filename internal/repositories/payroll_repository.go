package repositories

import (
	"hr-backend/internal/models"

	"gorm.io/gorm"
)

type PayrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) *PayrollRepository {
	return &PayrollRepository{db: db}
}

func (r *PayrollRepository) Create(payroll *models.Payroll) error {
	return r.db.Create(payroll).Error
}

func (r *PayrollRepository) FindAll(month, year int, page, limit int) ([]models.Payroll, int64, error) {
	var payrolls []models.Payroll
	var total int64

	query := r.db.Model(&models.Payroll{}).Preload("Employee.Department").Preload("Employee")

	if month > 0 {
		query = query.Where("month = ?", month)
	}

	if year > 0 {
		query = query.Where("year = ?", year)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&payrolls).Error

	return payrolls, total, err
}

func (r *PayrollRepository) FindByID(id uint) (*models.Payroll, error) {
	var payroll models.Payroll
	err := r.db.Preload("Employee.Department").Preload("Employee").First(&payroll, id).Error
	return &payroll, err
}

func (r *PayrollRepository) FindByEmployeeAndPeriod(employeeID uint, month, year int) (*models.Payroll, error) {
	var payroll models.Payroll
	err := r.db.Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).
		First(&payroll).Error
	return &payroll, err
}

func (r *PayrollRepository) Update(payroll *models.Payroll) error {
	return r.db.Save(payroll).Error
}

func (r *PayrollRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payroll{}, id).Error
}

func (r *PayrollRepository) GetPayrollSummary(month, year int) (*models.PayrollSummary, error) {
	var summary models.PayrollSummary

	err := r.db.Model(&models.Payroll{}).
		Where("month = ? AND year = ?", month, year).
		Select(`
			COUNT(*) as total_employees,
			COALESCE(SUM(basic_salary), 0) as total_basic_pay,
			COALESCE(SUM(allowances), 0) as total_allowances,
			COALESCE(SUM(deductions), 0) as total_deductions,
			COALESCE(SUM(tax), 0) as total_tax,
			COALESCE(SUM(net_salary), 0) as total_net_pay
		`).
		Scan(&summary).Error

	return &summary, err
}
