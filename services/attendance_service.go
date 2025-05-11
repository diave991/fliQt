package services

import (
	"fliQt/models/dto"
	"time"

	"fliQt/models"
	"fliQt/repositories"
)

// AttendanceService 封裝打卡相關業務邏輯
type AttendanceService struct {
	repo *repositories.AttendanceRepository
}

// NewAttendanceService 建立新的 AttendanceService
func NewAttendanceService(repo *repositories.AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

// Create 新增一筆打卡記錄
func (s *AttendanceService) Create(att *models.Attendance) error {
	return s.repo.Create(att)
}

// GetAll 取得所有打卡記錄
func (s *AttendanceService) GetAll() ([]models.Attendance, error) {
	return s.repo.GetAll()
}

// GetByID 根據 ID 取得單筆打卡記錄
func (s *AttendanceService) GetByID(id uint) (*models.Attendance, error) {
	return s.repo.GetByID(id)
}

// GetByEmployee 取得特定員工所有打卡記錄
func (s *AttendanceService) GetByEmployee(empID uint) ([]models.Attendance, error) {
	return s.repo.GetByEmployee(empID)
}

// GetByEmployeeAndDate 取得特定員工在某日期的打卡紀錄
func (s *AttendanceService) GetByEmployeeAndDate(empID uint, date time.Time) ([]dto.Attendance, error) {
	// 篩選同一天的紀錄
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)
	return s.repo.GetByEmployeeAndPeriod(empID, start, end)
}
