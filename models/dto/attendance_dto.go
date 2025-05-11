package dto

import "time"

type Attendance struct {
	ID               uint      `json:"id"`
	EmployeeID       uint      `json:"employee_id"`
	EmployeeName     string    `json:"employee_name"`
	EmployeePosition string    `json:"employee_position"`
	Type             string    `json:"type"`
	Timestamp        time.Time `json:"timestamp"`
}
