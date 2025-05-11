package models

import "time"

// Employee 員工資料模型
// 包含員工編號、姓名、職位、聯絡方式與薪資
type Employee struct {
	ID        uint      `gorm:"primaryKey;comment:員工編號" json:"id"`
	Name      string    `gorm:"size:100;comment:員工姓名" json:"name"`
	Position  string    `gorm:"size:50;comment:職位" json:"position"`
	Contact   string    `gorm:"size:100;comment:聯絡資訊" json:"contact"`
	Salary    float64   `gorm:"comment:薪資" json:"salary"`
	Status    int       `gorm:"type:tinyint;default:1;comment:'1=有效,0=已刪除'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime;<-:create;comment:創建時間" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新時間" json:"updated_at"`
}
