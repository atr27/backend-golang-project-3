package services

import (
	"errors"
	"hr-backend/internal/models"
	"hr-backend/internal/repositories"
	"hr-backend/internal/utils"
	"time"

	"gorm.io/gorm"
)

type PayrollService struct {
	payrollRepo  *repositories.PayrollRepository
	employeeRepo *repositories.EmployeeRepository
	db           *gorm.DB
}

func NewPayrollService(payrollRepo *repositories.PayrollRepository, employeeRepo *repositories.EmployeeRepository, db *gorm.DB) *PayrollService {
	return &PayrollService{
		payrollRepo:  payrollRepo,
		employeeRepo: employeeRepo,
		db:           db,
	}
}

func (s *PayrollService) GeneratePayroll(req *models.GeneratePayrollRequest) ([]models.Payroll, error) {
	// Get all active employees
	employees, _, err := s.employeeRepo.FindAll(1, 1000, nil, "active", "")
	if err != nil {
		return nil, err
	}

	var payrolls []models.Payroll

	for _, employee := range employees {
		// Check if payroll already exists
		existing, err := s.payrollRepo.FindByEmployeeAndPeriod(employee.ID, req.Month, req.Year)
		if err == nil && existing.ID > 0 {
			continue // Skip if already exists
		}

		// Calculate tax (simple flat 10% for example)
		tax := employee.Salary * 0.10

		payroll := models.Payroll{
			EmployeeID:  employee.ID,
			Month:       req.Month,
			Year:        req.Year,
			BasicSalary: employee.Salary,
			Allowances:  0,
			Deductions:  0,
			Tax:         tax,
			NetSalary:   employee.Salary - tax,
			Status:      "pending",
		}

		if err := s.payrollRepo.Create(&payroll); err != nil {
			return nil, err
		}

		payrolls = append(payrolls, payroll)
	}

	return payrolls, nil
}

func (s *PayrollService) GetPayrolls(month, year, page, limit int) ([]models.Payroll, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.payrollRepo.FindAll(month, year, page, limit)
}

func (s *PayrollService) GetPayrollByID(id uint) (*models.Payroll, error) {
	return s.payrollRepo.FindByID(id)
}

func (s *PayrollService) UpdatePayroll(id uint, payroll *models.Payroll) (*models.Payroll, error) {
	existing, err := s.payrollRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	existing.BasicSalary = payroll.BasicSalary
	existing.Allowances = payroll.Allowances
	existing.Deductions = payroll.Deductions
	existing.Tax = payroll.Tax
	existing.NetSalary = payroll.BasicSalary + payroll.Allowances - payroll.Deductions - payroll.Tax

	if err := s.payrollRepo.Update(existing); err != nil {
		return nil, err
	}

	return s.payrollRepo.FindByID(id)
}

func (s *PayrollService) ProcessPayment(id uint) (*models.Payroll, error) {
	payroll, err := s.payrollRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if payroll.Status == "paid" {
		return nil, errors.New("payroll already paid")
	}

	now := time.Now()
	payroll.PaymentDate = &now
	payroll.Status = "paid"

	if err := s.payrollRepo.Update(payroll); err != nil {
		return nil, err
	}

	return payroll, nil
}

func (s *PayrollService) GetPayrollSummary(month, year int) (*models.PayrollSummary, error) {
	return s.payrollRepo.GetPayrollSummary(month, year)
}

func (s *PayrollService) GeneratePayrollPDF(id uint) ([]byte, error) {
	payroll, err := s.payrollRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if payroll.Status != "paid" {
		return nil, errors.New("payroll must be paid before downloading")
	}

	return utils.GeneratePayrollPDF(payroll)
}
