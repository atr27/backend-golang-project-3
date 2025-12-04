package services

import (
	"hr-backend/internal/models"
	"hr-backend/internal/repositories"
)

type DepartmentService struct {
	deptRepo *repositories.DepartmentRepository
}

func NewDepartmentService(deptRepo *repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		deptRepo: deptRepo,
	}
}

func (s *DepartmentService) CreateDepartment(dept *models.Department) (*models.Department, error) {
	if err := s.deptRepo.Create(dept); err != nil {
		return nil, err
	}
	return s.deptRepo.FindByID(dept.ID)
}

func (s *DepartmentService) GetDepartments() ([]models.Department, error) {
	return s.deptRepo.FindAll()
}

func (s *DepartmentService) GetDepartmentByID(id uint) (*models.Department, error) {
	return s.deptRepo.FindByID(id)
}

func (s *DepartmentService) UpdateDepartment(id uint, dept *models.Department) (*models.Department, error) {
	existing, err := s.deptRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	existing.Name = dept.Name
	existing.Description = dept.Description
	existing.ManagerID = dept.ManagerID

	if err := s.deptRepo.Update(existing); err != nil {
		return nil, err
	}

	return s.deptRepo.FindByID(id)
}

func (s *DepartmentService) DeleteDepartment(id uint) error {
	return s.deptRepo.Delete(id)
}
