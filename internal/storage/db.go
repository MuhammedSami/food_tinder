package storage

import (
	"fmt"
	"foodtinder/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func NewDb(cfg config.DBConn) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port,
	)

	var db *gorm.DB
	var err error

	timeout := time.After(10 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			log.Fatalf("failed to connect to database within 10 seconds: %v", err)
		case <-tick:
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				SkipDefaultTransaction: true,
				PrepareStmt:            true,
				Logger:                 logger.Default.LogMode(logger.Warn),
			})
			if err == nil {
				sqlDB, err := db.DB()
				if err != nil {
					log.Fatalf("failed to get sql.DB: %v", err)
				}

				sqlDB.SetMaxOpenConns(25)
				sqlDB.SetMaxIdleConns(25)
				sqlDB.SetConnMaxLifetime(5 * time.Minute)
				sqlDB.SetConnMaxIdleTime(2 * time.Minute)

				log.Println("PostgreSQL connection succeeded via GORM")
				return db
			}

			log.Println("DB not ready yet, retrying...")
		}
	}
}
