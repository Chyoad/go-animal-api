package main

import (
	"fmt"
	animalHandler "go-animal-api/internal/handler"
	animalRepo "go-animal-api/internal/repository"
	animalUseCase "go-animal-api/internal/usecase"
	"go-animal-api/pkg/config"
	"go-animal-api/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
)


func main() {
	cfg, err := config.LoadConfig() // config.LoadConfig(".env.production"). Jika path tidak diberikan, akan mencoba memuat ".env".
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// Set Mode Gin dari konfigurasi
	gin.SetMode(cfg.GinMode)
	log.Printf("Mode Gin diatur ke: %s", cfg.GinMode)

	// Inisialisasi Database
	gormDB, err := db.InitDB(cfg.DBDSN, cfg.GinMode)
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}

	// --- Setup Layers ---
	animalRepository := animalRepo.NewMysqlAnimalRepository(gormDB)
	animalUC := animalUseCase.NewAnimalUseCase(animalRepository)
	animalHdlr := animalHandler.NewAnimalHandler(animalUC)

	// --- Router Setup ---
	router := gin.Default()

	v1 := router.Group("/v1")
	animalHandler.RegisterAnimalRoutes(v1, animalHdlr)

	// --- Start Server ---
	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server berjalan di http://localhost%s (Port Aplikasi: %s)", serverAddr, cfg.AppPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}