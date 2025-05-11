package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string // Host database : localhost, mysql_db (Docker)
	DBPort     string // Port database : 3306
	DBUser     string // User database
	DBPassword string // Password database
	DBName     string // Nama database
	DBSchema   string // Parameter tambahan DSN : charset=utf8mb4&parseTime=True&loc=Local
	AppPort    string // Port aplikasi
	GinMode    string // Mode Gin: debug, release, test
	DBDSN      string // DSN lengkap, bisa dibuat otomatis atau di-override
}

func LoadConfig(envPath ...string) (*Config, error) {
	
	if len(envPath) > 0 {
		err := godotenv.Load(envPath[0])
		if err != nil {
			log.Printf("Peringatan: Tidak dapat memuat file %s: %v. Menggunakan env variabel secara langsung.", envPath[0], err)
		} else {
			log.Printf("Konfigurasi dimuat dari %s", envPath[0])
		}
	} else {
		err := godotenv.Load()
		if err == nil {
			log.Println("Konfigurasi dimuat dari file .env default.")
		}
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "animal_db"),
		DBSchema:   getEnv("DB_SCHEMA", "charset=utf8mb4&parseTime=True&loc=Local"),
		AppPort:    getEnv("APP_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"), // default "debug" untuk lokal
	}

	// Format: USER:PASSWORD@tcp(HOST:PORT)/DBNAME?PARAMETER
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSchema,
	)
	// DB_DSN di-override jika diset
	cfg.DBDSN = getEnv("DB_DSN", dsn)

	// Validasi
	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN tidak diset")
	}
	if _, err := strconv.Atoi(cfg.AppPort); err != nil {
		return nil, fmt.Errorf("APP_PORT tidak valid: %s, error: %w", cfg.AppPort, err)
	}
	validGinModes := map[string]bool{"debug": true, "release": true, "test": true}
	if !validGinModes[cfg.GinMode] {
		log.Printf("Peringatan: GIN_MODE tidak valid ('%s'), akan menggunakan 'debug'. Mode yang valid: debug, release, test.", cfg.GinMode)
		cfg.GinMode = "debug"
	}

	return cfg, nil
}

// getEnv membaca env atau return nilai default.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Env variabel '%s' tidak ditemukan, menggunakan nilai default: '%s'", key, defaultValue)
	return defaultValue
}