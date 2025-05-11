package db

import (
	"fmt"
	"log"

	"fliQt/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Printf(">>> DSN=%s", dsn)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
