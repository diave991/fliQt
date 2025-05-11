package services

import (
	"context"
	"encoding/json"
	"fliQt/models"
	"fmt"
	"time"

	"fliQt/repositories"

	"fliQt/models/dto"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ReportService handles report generation and storage
type ReportService struct {
	egRepo *repositories.EmployeeRepository
	lvRepo *repositories.LeaveRepository
	atRepo *repositories.AttendanceRepository
	rdb    *redis.Client
	db     *gorm.DB
}

// NewReportService constructs a ReportService
func NewReportService(db *gorm.DB, rdb *redis.Client,
	egRepo *repositories.EmployeeRepository,
	lvRepo *repositories.LeaveRepository,
	atRepo *repositories.AttendanceRepository) *ReportService {
	return &ReportService{egRepo, lvRepo, atRepo, rdb, db}
}

// GenerateReportForEmployee computes 7-day report for one employee
func (s *ReportService) GenerateReportForEmployee(ctx context.Context, empID uint) (dto.AttendanceReport, error) {
	today := time.Now().Truncate(24 * time.Hour)
	report := make(dto.AttendanceReport, 7)

	for i := 0; i < 7; i++ {
		day := today.AddDate(0, 0, -i)
		start := day
		end := day.Add(24 * time.Hour)

		// 改用 GetAll 並自行篩選同一員工的請假
		allLeaves, _ := s.lvRepo.GetAll()
		var leaves []models.Leave
		for _, l := range allLeaves {
			if l.EmployeeID == empID {
				leaves = append(leaves, l)
			}
		}
		status := "absent"
		for _, l := range leaves {
			if !l.StartDate.After(end) && !l.EndDate.Before(start) {
				status = "leave"
				break
			}
		}
		if status != "leave" {
			// check attendance
			atts, _ := s.atRepo.GetByEmployeeAndPeriod(empID, start, end)
			if len(atts) > 0 {
				status = "present"
			}
		}

		report[i] = dto.DailyStatus{
			Date:   day.Format("2006-01-02"),
			Status: status,
		}
	}

	// reverse so oldest first
	for i, j := 0, len(report)-1; i < j; i, j = i+1, j-1 {
		report[i], report[j] = report[j], report[i]
	}

	// store to Redis
	key := fmt.Sprintf("attendance:report:%d", empID)
	data, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	if err := s.rdb.Set(ctx, key, data, 0).Err(); err != nil {
		return nil, err
	}
	return report, nil
}

// RunScheduler schedules every ten minutes with 5 worker concurrency
func (s *ReportService) RunScheduler(ctx context.Context) {
	// 1. 先拿到所有員工
	empList, err := s.egRepo.GetAll()
	if err != nil {
		return
	}

	jobs := make(chan uint, len(empList))
	// feed jobs
	for _, emp := range empList {
		jobs <- emp.ID
	}
	close(jobs)

	// start workers
	for w := 0; w < 5; w++ {
		go func() {
			for empID := range jobs {
				s.GenerateReportForEmployee(ctx, empID)
			}
		}()
	}
}

func (s *ReportService) GetReport(ctx context.Context, empID uint) (dto.AttendanceReport, error) {
	key := fmt.Sprintf("attendance:report:%d", empID)
	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var report dto.AttendanceReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}
	return report, nil
}

// GetAllReports 分頁取得所有員工的出缺勤報表
// page 從 1 開始，每頁 10 筆
func (s *ReportService) GetAllReports(ctx context.Context, page int) ([]dto.EmployeeReport, error) {
	emps, err := s.egRepo.GetAll() // 取得所有員工
	if err != nil {
		return nil, err
	}

	// 計算分頁範圍
	const pageSize = 10
	start := (page - 1) * pageSize
	if start >= len(emps) {
		return []dto.EmployeeReport{}, nil
	}
	end := start + pageSize
	if end > len(emps) {
		end = len(emps)
	}

	var results []dto.EmployeeReport
	for _, emp := range emps[start:end] {
		// 從 Redis 讀取已生成報表；若無則自動生成
		report, err := s.GetReport(ctx, emp.ID)
		if err != nil {
			// fallback：現場生成並存入 Redis
			report, _ = s.GenerateReportForEmployee(ctx, emp.ID)
		}
		results = append(results, dto.EmployeeReport{
			EmployeeID:       emp.ID,
			EmployeeName:     emp.Name,
			EmployeePosition: emp.Position,
			Report:           report,
		})
	}
	return results, nil
}
