package migrations

import (
	"fliQt/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// 1. 建表的遷移
		{

			ID: "20250511_create_employee_and_leave_and_Attendance",
			Migrate: func(tx *gorm.DB) error {
				log.Printf("進行 建表的遷移 ~~~~~~~~~…  \n")
				return tx.AutoMigrate(&models.Employee{}, &models.Leave{}, &models.Attendance{})
			},
			Rollback: func(tx *gorm.DB) error {
				// ... rollback code ...
				return nil
			},
		},
		// 2. 插入初始資料的遷移
		{
			ID: "20250511_seed_initial_data",
			Migrate: func(tx *gorm.DB) error {
				log.Printf("進行 插入初始資料的遷移 ~~~~~~~~~…  \n")
				// 準備要插入的員工資料
				emps := []models.Employee{
					{Name: "王小明", Position: "資深後端工程師", Contact: "xiaoming.wang@example.com", Salary: 1200000, Status: 1},
					{Name: "李麗華", Position: "前端工程師", Contact: "li.lihua@example.com", Salary: 1000000, Status: 1},
				}
				// 逐筆寫入
				for _, e := range emps {
					log.Printf("進行 Migrations ~~~~~~~~~…  \n")
					//if err := tx.Create(&e).Error; err != nil {
					//	return err
					//}
					if err := tx.Where("name = ?", e.Name).FirstOrCreate(&e).Error; err != nil {
						return err
					}
				}
				log.Printf("進行 Migrations ERROR~~~~~~~~~…  \n")
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				log.Printf("進行 Rollback ERROR~~~~~~~~~…  \n")
				// 回滾時可刪掉這些資料
				return tx.Exec("DELETE FROM employees WHERE name IN (?, ?)", "王小明", "李麗華").Error
			},
		},
		// 可以繼續新增更多 migration…
	})

	return m.Migrate()
}
