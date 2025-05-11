package repositories

import (
	"fliQt/models"
	"fliQt/models/dto"
	"gorm.io/gorm"
	"time"
)

// AttendanceRepository 負責打卡記錄的資料庫存取
// 提供 CRUD 及依條件查詢功能

type AttendanceRepository struct {
	DB *gorm.DB
}

// NewAttendanceRepository 建立新的 AttendanceRepository
func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

// Create 新增一筆打卡紀錄
func (r *AttendanceRepository) Create(att *models.Attendance) error {
	return r.DB.Create(att).Error
}

// GetAll 取得所有打卡紀錄
func (r *AttendanceRepository) GetAll() ([]models.Attendance, error) {
	var records []models.Attendance
	err := r.DB.Find(&records).Error
	return records, err
}

// GetByID 根據 ID 取得單筆打卡紀錄
func (r *AttendanceRepository) GetByID(id uint) (*models.Attendance, error) {
	var att models.Attendance
	err := r.DB.First(&att, id).Error
	return &att, err
}

// GetByEmployee 取得特定員工所有打卡紀錄
func (r *AttendanceRepository) GetByEmployee(empID uint) ([]models.Attendance, error) {
	var records []models.Attendance
	err := r.DB.Where("employee_id = ?", empID).Order("timestamp").Find(&records).Error
	return records, err
}

//// GetByEmployeeAndPeriod 取得特定員工於指定期間的打卡紀錄
//func (r *AttendanceRepository) GetByEmployeeAndPeriod(empID uint, start, end time.Time) ([]models.Attendance, error) {
//	var records []models.Attendance
//	err := r.DB.Debug().Where("employee_id = ? AND timestamp >= ? AND timestamp < ?", empID, start, end).
//		Order("timestamp").Find(&records).Error
//
//	return records, err
//}

func (r *AttendanceRepository) GetByEmployeeAndPeriod(empID uint, start, end time.Time) ([]dto.Attendance, error) {
	var list []dto.Attendance
	err := r.DB.
		Table("attendances").
		Select(`attendances.id,
                attendances.employee_id,
                employees.name   AS employee_name,
                employees.position AS employee_position,
                attendances.type,
                attendances.timestamp`).
		Joins("LEFT JOIN employees ON attendances.employee_id = employees.id").
		Where("attendances.employee_id = ? AND attendances.timestamp BETWEEN ? AND ?", empID, start, end).
		Scan(&list).
		Error
	return list, err
}
