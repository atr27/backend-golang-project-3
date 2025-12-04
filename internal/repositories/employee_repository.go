package repositories

import (
	"hr-backend/internal/models"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *EmployeeRepository) FindAll(page, limit int, departmentID *uint, status, search string) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.Model(&models.Employee{}).Preload("User").Preload("Department")

	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	if status != "" {
		query = query.Where("employment_status = ?", status)
	}

	if search != "" {
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR employee_code ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&employees).Error

	return employees, total, err
}

func (r *EmployeeRepository) FindByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Preload("User").Preload("Department").First(&employee, id).Error
	return &employee, err
}

func (r *EmployeeRepository) FindByEmployeeCode(code string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("employee_code = ?", code).First(&employee).Error
	return &employee, err
}

func (r *EmployeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

func (r *EmployeeRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Employee{}).Where("employment_status = ?", status).Count(&count).Error
	return count, err
}

func (r *EmployeeRepository) CountByDepartment(departmentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Employee{}).Where("department_id = ?", departmentID).Count(&count).Error
	return count, err
}
