package services

import "fliQt/models"
import "fliQt/repositories"

type EmployeeService struct {
	Repo *repositories.EmployeeRepository
}

func NewEmployeeService(repo *repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{Repo: repo}
}

func (s *EmployeeService) Create(employee *models.Employee) error {
	return s.Repo.Create(employee)
}

func (s *EmployeeService) GetAll() ([]models.Employee, error) {
	return s.Repo.GetAll()
}

func (s *EmployeeService) GetByID(id uint) (*models.Employee, error) {
	return s.Repo.GetByID(id)
}

func (s *EmployeeService) Update(employee *models.Employee) error {
	return s.Repo.Update(employee)
}

func (s *EmployeeService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
