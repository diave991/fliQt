package services

import (
	"fliQt/models"
	"fliQt/models/dto"
)
import "fliQt/repositories"

type LeaveService struct {
	Repo *repositories.LeaveRepository
}

func NewLeaveService(repo *repositories.LeaveRepository) *LeaveService {
	return &LeaveService{Repo: repo}
}

func (s *LeaveService) Create(leave *models.Leave) error {
	return s.Repo.Create(leave)
}

func (s *LeaveService) GetAll() ([]dto.Leave, error) {
	return s.Repo.GetAllWithStaff()
}

func (s *LeaveService) GetByID(id uint) (*models.Leave, error) {
	return s.Repo.GetByID(id)
}

func (s *LeaveService) Update(leave *models.Leave) error {
	return s.Repo.Update(leave)
}

func (s *LeaveService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
