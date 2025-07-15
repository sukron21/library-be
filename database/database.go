package database

import (
	"fmt"
	"library/config" // Sesuaikan dengan nama proyekmu
	"library/models" // Sesuaikan dengan nama proyekmu
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBClient *gorm.DB

// InitDatabase menginisialisasi koneksi database
func InitDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	var err error
	DBClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	// Migrasi skema database
	DBClient.AutoMigrate(&models.User{})
	log.Println("Database migration complete!")
}
