package repositories

import (
	"fliQt/models"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	DB *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.DB.Create(employee).Error
}

func (r *EmployeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := r.DB.Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepository) GetByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.DB.Where("id = ? AND status = 1", id).First(&employee).Error
	return &employee, err
}

func (r *EmployeeRepository) Update(e *models.Employee) error {
	return r.DB.
		Model(&models.Employee{}).
		Where("id = ? AND status = 1", e.ID).
		Updates(map[string]interface{}{
			"name":     e.Name,
			"position": e.Position,
			"contact":  e.Contact,
			"salary":   e.Salary,
		}).Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	return r.DB.Model(&models.Employee{}).Where("id = ?", id).Update("status", 0).Error
}
