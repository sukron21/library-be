package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct untuk menyimpan konfigurasi aplikasi
type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

// LoadConfig memuat konfigurasi dari variabel lingkungan atau file .env
func LoadConfig() *Config {
	err := godotenv.Load() // Memuat variabel dari .env
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:      getEnv("PORT", "3000"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "password"),
		DBName:    getEnv("DB_NAME", "mydb"),
		JWTSecret: getEnv("JWT_SECRET", "supersecretjwtkey"),
	}
}

// getEnv membantu mendapatkan nilai variabel lingkungan dengan nilai default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
