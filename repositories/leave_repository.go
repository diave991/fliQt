package repositories

import (
	"fliQt/models"
	"fliQt/models/dto"
	"gorm.io/gorm"
)

type LeaveRepository struct {
	DB *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{DB: db}
}

func (r *LeaveRepository) Create(leave *models.Leave) error {
	return r.DB.Create(leave).Error
}

func (r *LeaveRepository) GetAll() ([]models.Leave, error) {
	var leaves []models.Leave
	err := r.DB.Find(&leaves).Error
	return leaves, err
}

func (r *LeaveRepository) GetAllWithStaff() ([]dto.Leave, error) {
	var list []dto.Leave

	// 注意表名要与模型 TableName 一致，默认 plural 小写
	err := r.DB.
		Table("leaves").
		Select(`leaves.id,
                leaves.employee_id,
                employees.name   AS employee_name,
                employees.position AS employee_position,
                leaves.start_date,
                leaves.end_date,
                leaves.reason,
                leaves.created_at`).
		Joins("LEFT JOIN employees ON leaves.employee_id = employees.id").
		Scan(&list).Error

	return list, err
}
func (r *LeaveRepository) GetByID(id uint) (*models.Leave, error) {
	var leave models.Leave
	err := r.DB.First(&leave, id).Error
	return &leave, err
}

func (r *LeaveRepository) Update(leave *models.Leave) error {
	return r.DB.Save(leave).Error
}

func (r *LeaveRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Leave{}, id).Error
}
