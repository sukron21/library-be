package routes

import (
	"library/controllers" // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"library/middleware"  // SESUAIKAN DENGAN NAMA MODUL GO ANDA

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mengatur semua rute API
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Rute Autentikasi (Publik)
	api.Post("/auth/login", controllers.Login)
	api.Post("/auth/refresh", controllers.RefreshAccessToken)

	// Rute Publik lainnya
	api.Post("/users", controllers.CreateUser)

	authenticated := api.Group("/protected")
	authenticated.Use(middleware.AuthRequired)
	authenticated.Get("/users", controllers.GetAllUsers)
	authenticated.Get("/users/:id", controllers.GetUserByID)
	authenticated.Put("/users/:id", controllers.UpdateUser)
	authenticated.Delete("/users/:id", controllers.DeleteUser)
	//books
	authenticated.Post("/books", controllers.CreateBook)
	authenticated.Get("/books", controllers.GetAllBooks)
	authenticated.Get("/books/:id", controllers.GetBooksByID)
	authenticated.Put("/books/:id", controllers.UpdateBooks)
	authenticated.Delete("/books/:id", controllers.DeleteBooks)
}
