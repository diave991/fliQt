// services/employee_service_test.go
package services

import (
	"testing"

	"fliQt/models"
	"fliQt/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test db: %v", err)
	}
	// 同步最新模型结构
	db.AutoMigrate(&models.Employee{}, &models.Leave{})
	return db
}

func TestEmployeeService_CRUD_SoftDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewEmployeeRepository(db)
	svc := NewEmployeeService(repo)

	// 1. Create
	emp := &models.Employee{
		Name:     "王小明",
		Position: "工程师",
		Contact:  "xiaoming@example.com",
		Salary:   1200000,
	}
	if err := svc.Create(emp); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 2. 初始 GetAll，应该有 1 条且 status=1
	list, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(list))
	}
	if list[0].Status != 1 {
		t.Errorf("expected status=1 after create; got %d", list[0].Status)
	}

	// 3. Soft Delete
	if err := svc.Delete(emp.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// 4. GetByID 应报错（因为 ID 条件里有 status=1 筛选） :contentReference[oaicite:0]{index=0}:contentReference[oaicite:1]{index=1}
	if _, err := svc.GetByID(emp.ID); err == nil {
		t.Error("expected error from GetByID after soft delete, got none")
	}

	// 5. 删除后 GetAll 仍会返回记录，但其 Status 应为 0
	list2, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll after delete failed: %v", err)
	}
	if len(list2) != 1 {
		t.Fatalf("expected 1 employee after delete, got %d", len(list2))
	}
	if list2[0].Status != 0 {
		t.Errorf("expected status=0 after soft delete; got %d", list2[0].Status)
	}
}
