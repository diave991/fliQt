package main

import (
	"context"
	"fliQt/config"
	"fliQt/controllers"
	"fliQt/db"
	"fliQt/migrations"
	"fliQt/repositories"
	"fliQt/routes"
	"fliQt/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	var mysqlDB *gorm.DB
	var err error
	// 最多重试 10 次，每次间隔 3 秒
	for i := 1; i <= 10; i++ {
		mysqlDB, err = db.ConnectMySQL(cfg)
		if err == nil {
			break
		}
		log.Printf("等待 MySQL 就绪，第 %d 次重试…  (%v)\n", i, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("无法连接 MySQL，退出: %v", err)
	}

	// 2. 执行 gormigrate 迁移
	log.Println(">>> 开始执行数据库迁移")
	if err := migrations.Migrate(mysqlDB); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	log.Println(">>> 数据库迁移完成")
	redisClient := db.ConnectRedis(cfg)
	if err := db.PingRedis(context.Background(), redisClient); err != nil {
		log.Fatalf("failed to connect Redis: %v", err)
	}

	empRepo := repositories.NewEmployeeRepository(mysqlDB)
	empSvc := services.NewEmployeeService(empRepo)
	empCtrl := controllers.NewEmployeeController(empSvc)

	leaveRepo := repositories.NewLeaveRepository(mysqlDB)
	leaveSvc := services.NewLeaveService(leaveRepo)
	leaveCtrl := controllers.NewLeaveController(leaveSvc)

	// Attendance wiring
	attRepo := repositories.NewAttendanceRepository(mysqlDB)
	attSvc := services.NewAttendanceService(attRepo)
	attCtrl := controllers.NewAttendanceController(attSvc)

	r := gin.Default()
	routes.Setup(r, empCtrl, leaveCtrl, attCtrl)

	r.Run() // default :8080
}
