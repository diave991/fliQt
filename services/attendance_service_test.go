// services/attendance_service_test.go
package services

import (
	"testing"
	"time"

	"fliQt/models"
	"fliQt/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupAttendanceDB 初始化 SQLite in-memory DB 并自动迁移
func setupAttendanceDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.Employee{}, &models.Attendance{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	return db
}

func TestAttendanceService_CRUD(t *testing.T) {
	db := setupAttendanceDB(t)

	// 创建一个员工，供外键使用
	emp := &models.Employee{
		Name:     "TestUser",
		Position: "Developer",
		Contact:  "test@example.com",
		Salary:   1000,
		Status:   1,
	}
	if err := db.Create(emp).Error; err != nil {
		t.Fatalf("create employee failed: %v", err)
	}

	// 使用真实仓储与服务
	repo := repositories.NewAttendanceRepository(db) // :contentReference[oaicite:0]{index=0}:contentReference[oaicite:1]{index=1}
	svc := NewAttendanceService(repo)                // :contentReference[oaicite:2]{index=2}:contentReference[oaicite:3]{index=3}

	// 1. 初始 GetAll 应返回空切片
	all0, err := svc.GetAll()
	if err != nil {
		t.Fatalf("initial GetAll failed: %v", err)
	}
	if len(all0) != 0 {
		t.Fatalf("expected 0 records initially, got %d", len(all0))
	}

	// 2. Create 一条打卡记录
	now := time.Now().UTC()
	att := &models.Attendance{
		EmployeeID: emp.ID,
		Type:       "IN",
		Timestamp:  now,
	}
	if err := svc.Create(att); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if att.ID == 0 {
		t.Fatal("expected att.ID to be set after Create")
	}

	// 3. GetAll 应返回一条记录
	all1, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll after create failed: %v", err)
	}
	if len(all1) != 1 {
		t.Fatalf("expected 1 record, got %d", len(all1))
	}

	// 4. GetByID 应返回同一条记录
	got, err := svc.GetByID(att.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.ID != att.ID || got.EmployeeID != emp.ID || got.Type != "IN" || !got.Timestamp.Equal(att.Timestamp) {
		t.Errorf("GetByID returned %+v; want %+v", got, att)
	}

	// 5. GetByEmployee 应返回同一条记录
	byEmp, err := svc.GetByEmployee(emp.ID)
	if err != nil {
		t.Fatalf("GetByEmployee failed: %v", err)
	}
	if len(byEmp) != 1 || byEmp[0].ID != att.ID {
		t.Errorf("GetByEmployee returned %v; want ID %d", byEmp, att.ID)
	}

	// 6. GetByEmployeeAndDate 在当天应返回记录
	sameDay, err := svc.GetByEmployeeAndDate(emp.ID, now)
	if err != nil {
		t.Fatalf("GetByEmployeeAndDate failed: %v", err)
	}
	if len(sameDay) != 1 {
		t.Errorf("GetByEmployeeAndDate returned %d records; want 1", len(sameDay))
	}

	// 7. GetByEmployeeAndDate 在非当天（隔天）应返回空
	outOfDay, err := svc.GetByEmployeeAndDate(emp.ID, now.AddDate(0, 0, -1))
	if err != nil {
		t.Fatalf("GetByEmployeeAndDate (outside) failed: %v", err)
	}
	if len(outOfDay) != 0 {
		t.Errorf("GetByEmployeeAndDate (outside) returned %d; want 0", len(outOfDay))
	}
}
