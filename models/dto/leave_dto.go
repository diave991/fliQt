package dto

import "time"

type Leave struct {
	ID               uint      `json:"id"`
	EmployeeID       uint      `json:"employee_id"`
	EmployeeName     string    `json:"employee_name"`
	EmployeePosition string    `json:"employee_position"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Reason           string    `json:"reason"`
	CreatedAt        time.Time `json:"created_at"`
}
