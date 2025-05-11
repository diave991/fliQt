package dto

// DailyStatus represents one day attendance status
type DailyStatus struct {
	Date   string `json:"date"`   // YYYY-MM-DD
	Status string `json:"status"` // "present", "leave", "absent"
}

type AttendanceReport []DailyStatus

// EmployeeReport 包含員工基本資訊及其 7 天出缺勤報表
type EmployeeReport struct {
	EmployeeID       uint             `json:"employee_id"`
	EmployeeName     string           `json:"employee_name"`
	EmployeePosition string           `json:"employee_position"`
	Report           AttendanceReport `json:"report"`
}
