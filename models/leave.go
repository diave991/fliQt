package models

import "time"

// Leave 請假資料模型
// 包含請假編號、員工編號、開始與結束日期、請假原因

// 欄位透過 `comment` 標註欄位說明

type Leave struct {
	ID         uint      `gorm:"primaryKey;comment:請假編號" json:"id"`
	EmployeeID uint      `gorm:"comment:員工編號" json:"employee_id"`
	StartDate  time.Time `gorm:"comment:開始日期" json:"start_date"`
	EndDate    time.Time `gorm:"comment:結束日期" json:"end_date"`
	Reason     string    `gorm:"type:text;comment:請假原因" json:"reason"`
	CreatedAt  time.Time `gorm:"autoCreateTime;comment:建立時間" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;comment:更新時間" json:"updated_at"`
}
