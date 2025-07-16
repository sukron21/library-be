package controllers

import (
	"library/database" // Sesuaikan dengan nama proyekmu
	"library/helpers"  // Sesuaikan dengan nama proyekmu
	"library/models"   // Sesuaikan dengan nama proyekmu
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// GetAllUsers mendapatkan semua pengguna
func GetAllUsers(c *fiber.Ctx) error {
	var user []models.User
	var total int64
	page, err := strconv.Atoi(c.Query("page", "1")) // Ambil "page" dari URL, default "1"
	if err != nil || page < 1 {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit", "10")) // Ambil "limit" dari URL, default "10"
	if err != nil || limit < 1 {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid limit number")
	}
	offset := (page - 1) * limit

	if result := database.DBClient.Model(&models.Book{}).Count(&total); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}
	if result := database.DBClient.Limit(limit).Offset(offset).Find(&user); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}
	if len(user) == 0 && page > 1 { // Jika halaman lebih dari 1 dan tidak ada buku, berarti halaman kosong
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "No books found on this page")
	} else if len(user) == 0 { // Jika di halaman 1 pun tidak ada buku sama sekali
		return helpers.SuccessResponse(c, fiber.StatusOK, "No books found", []models.Book{})
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	return helpers.SuccessResponse(c, fiber.StatusOK, "Books retrieved successfully", fiber.Map{
		"data":         user,
		"total_items":  total,
		"current_page": page,
		"per_page":     limit,
		"total_pages":  totalPages,
	})
}

// GetUserByID mendapatkan pengguna berdasarkan ID
func GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	bookID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format")
	}
	user := new(models.User)

	if result := database.DBClient.First(&user, bookID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}

// UpdateUser memperbarui pengguna
func UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")

	userID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format")
	}

	user := new(models.User)

	if result := database.DBClient.First(&user, userID); result.Error != nil {
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
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format")
	}

	user := new(models.User)

	if result := database.DBClient.First(&user, userID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	if result := database.DBClient.Delete(&user); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}

func GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in token")
	}
	// Konversi ke UUID
	id, err := uuid.Parse(userID.(string))
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format")
	}

	user := new(models.User)
	if result := database.DBClient.First(&user, id); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	user.Password = ""

	return helpers.SuccessResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}

func GetAllUsersNoPagination(c *fiber.Ctx) error {
	var users []models.User

	// Ambil semua user
	if result := database.DBClient.Find(&users); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	// Hapus password sebelum return
	for i := range users {
		users[i].Password = ""
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "All users retrieved successfully", users)
}
