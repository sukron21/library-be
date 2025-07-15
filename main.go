package main

import (
	"library/config"   // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"library/database" // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"library/routes"   // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet" // Opsional: untuk security headers
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Muat konfigurasi aplikasi
	cfg := config.LoadConfig()

	// Inisialisasi koneksi database
	database.InitDatabase(cfg)

	// Buat instance aplikasi Fiber
	app := fiber.New()

	// Middleware Global
	app.Use(logger.New()) // Logging setiap permintaan ke konsol
	app.Use(cors.New())   // Mengizinkan Cross-Origin Resource Sharing (CORS)
	app.Use(helmet.New()) // Opsional: Menambahkan berbagai security HTTP headers

	// Setup semua rute API
	routes.SetupRoutes(app)

	// Jalankan server Fiber di port yang ditentukan dalam konfigurasi
	log.Fatal(app.Listen(":" + cfg.Port))
}
