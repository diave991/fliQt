// services/leave_service_test.go
package services

import (
	"errors"
	"testing"
	"time"

	"fliQt/models"
	"fliQt/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupLeaveTestDB 在内存 SQLite 中自动迁移 Employee 与 Leave 模型
func setupLeaveTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.Employee{}, &models.Leave{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	return db
}

func TestLeaveService_CRUD(t *testing.T) {
	db := setupLeaveTestDB(t)

	// 1. 建立員工，供外鍵使用
	emp := &models.Employee{
		Name:     "張三",
		Position: "助理",
		Contact:  "zhangsan@example.com",
		Salary:   800000,
		Status:   1,
	}
	if err := db.Create(emp).Error; err != nil {
		t.Fatalf("create employee failed: %v", err)
	}

	// 2. 構造 Service
	leaveRepo := repositories.NewLeaveRepository(db)
	leaveSvc := NewLeaveService(leaveRepo)

	// 3. 初始 GetAll 應回空
	all0, err := leaveSvc.GetAll()
	if err != nil {
		t.Fatalf("initial GetAll failed: %v", err)
	}
	if len(all0) != 0 {
		t.Fatalf("expected 0 leaves initially, got %d", len(all0))
	}

	// 4. Create
	start := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 6, 3, 0, 0, 0, 0, time.UTC)
	leave := &models.Leave{
		EmployeeID: emp.ID,
		StartDate:  start,
		EndDate:    end,
		Reason:     "年假",
	}
	if err := leaveSvc.Create(leave); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if leave.ID == 0 {
		t.Fatal("expected leave.ID to be set after Create")
	}

	// 5. GetAll 應回一筆
	all1, err := leaveSvc.GetAll()
	if err != nil {
		t.Fatalf("GetAll after create failed: %v", err)
	}
	if len(all1) != 1 {
		t.Fatalf("expected 1 leave, got %d", len(all1))
	}

	// 6. GetByID 應回該筆
	got, err := leaveSvc.GetByID(leave.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got == nil {
		t.Fatal("GetByID returned nil")
	}
	if !got.StartDate.Equal(start) || !got.EndDate.Equal(end) || got.Reason != "年假" {
		t.Errorf("GetByID returned %+v; want start=%v end=%v reason=年假", got, start, end)
	}

	// 7. Update
	got.Reason = "病假"
	if err := leaveSvc.Update(got); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	updated, err := leaveSvc.GetByID(leave.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if updated.Reason != "病假" {
		t.Errorf("expected Reason=病假 after update, got %q", updated.Reason)
	}

	// 8. Delete
	if err := leaveSvc.Delete(leave.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	// Delete 后 GetByID 应返回 ErrRecordNotFound
	_, err = leaveSvc.GetByID(leave.ID)
	if err == nil {
		t.Error("expected error from GetByID after delete, got nil")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}

	// 9. Delete 后 GetAll 应回空
	all2, err := leaveSvc.GetAll()
	if err != nil {
		t.Fatalf("GetAll after delete failed: %v", err)
	}
	if len(all2) != 0 {
		t.Errorf("expected 0 leaves after delete, got %d", len(all2))
	}
}
