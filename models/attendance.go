package models

import (
	"time"
)

// Attendance 打卡記錄模型
// 包含打卡編號、員工編號、打卡類型（上班或下班）、打卡時間
// 及建立/更新時間

type Attendance struct {
	ID         uint      `gorm:"primaryKey;comment:打卡編號" json:"id"`
	EmployeeID uint      `gorm:"index;comment:員工編號" json:"employee_id"`
	Type       string    `gorm:"size:20;comment:打卡類型，上班:IN，下班:OUT" json:"type"`
	Timestamp  time.Time `gorm:"comment:打卡時間" json:"timestamp"`
	CreatedAt  time.Time `gorm:"autoCreateTime;comment:記錄建立時間" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;comment:記錄更新時間" json:"updated_at"`
}
