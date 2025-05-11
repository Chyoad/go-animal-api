package db

import (
	"fmt"
	"go-animal-api/internal/domain"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dsn string, ginMode string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DSN database tidak boleh kosong")
	}

	logLevel := logger.Info
	if ginMode == "release" {
		logLevel = logger.Warn // Kurangi log di mode release
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false, // Log error record not found jika diinginkan
			Colorful:                  true,
		},
	)

	var gormDB *gorm.DB
	var err error
	// Retry koneksi
	for i := 0; i < 5; i++ {
		gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err == nil {
			log.Println("Berhasil terhubung ke database.")
			log.Println("Memigrasikan skema database...")
			err = gormDB.AutoMigrate(&domain.Animal{}) // Hanya migrasi entitas Animal
			if err != nil {
				log.Printf("Gagal memigrasikan skema database: %v. Menutup koneksi.", err)
				sqlDB, dbErr := gormDB.DB()
				if dbErr == nil {
					sqlDB.Close()
				}
				return nil, fmt.Errorf("gagal memigrasikan skema database: %w", err)
			}
			log.Println("Migrasi database selesai.")
			return gormDB, nil
		}
		log.Printf("Gagal terhubung ke database (percobaan %d/5 menggunakan DSN): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("gagal terhubung ke database setelah beberapa percobaan (DSN: %s): %w", dsn, err)
}