package controllers

import (
	"library/database" // Sesuaikan dengan nama proyekmu
	"library/helpers"  // Sesuaikan dengan nama proyekmu
	"library/models"   // Sesuaikan dengan nama proyekmu

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser membuat pengguna baru
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Could not hash password")
	}
	user.Password = string(hashedPassword)

	if result := database.DBClient.Create(&user); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusCreated, "User created successfully", user)
}

// GetUserByID mendapatkan pengguna berdasarkan ID
// GetAllUsers mendapatkan semua pengguna
func GetAllUsers(c *fiber.Ctx) error {
	var user []models.User // Gunakan slice (array dinamis) untuk menampung banyak user

	// Menggunakan Find() tanpa kondisi untuk mengambil semua data
	if result := database.DBClient.Find(&user); result.Error != nil {
		// Jika terjadi error saat query (selain tidak ditemukan), kembalikan Internal Server Error
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	// Sembunyikan password untuk setiap user sebelum dikirim sebagai respons
	for i := range user {
		user[i].Password = ""
	}

	// Jika tidak ada user ditemukan, Find() tidak mengembalikan error,
	// tetapi slice `users` akan kosong. Kita bisa mengirimkan pesan yang sesuai.
	if len(user) == 0 {
		return helpers.SuccessResponse(c, fiber.StatusOK, "No users found", []models.User{}) // Mengembalikan array kosong
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Users retrieved successfully", user)
}

// GetUserByID mendapatkan pengguna berdasarkan ID
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.User)

	if result := database.DBClient.First(&user, id); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}

// UpdateUser memperbarui pengguna
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.User)

	if result := database.DBClient.First(&user, id); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	updates := new(models.User)
	if err := c.BodyParser(updates); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if updates.Name != "" {
		user.Name = updates.Name
	}
	if updates.Email != "" {
		user.Email = updates.Email
	}
	// Perbarui password jika ada
	if updates.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updates.Password), bcrypt.DefaultCost)
		if err != nil {
			return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Could not hash password")
		}
		user.Password = string(hashedPassword)
	}

	if result := database.DBClient.Save(&user); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "User updated successfully", user)
}

// DeleteUser menghapus pengguna
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.User)

	if result := database.DBClient.First(&user, id); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	if result := database.DBClient.Delete(&user); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}
