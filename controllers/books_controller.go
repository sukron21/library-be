package controllers

import (
	"library/database" // Sesuaikan dengan nama proyekmu
	"library/helpers"  // Sesuaikan dengan nama proyekmu
	"library/models"   // Sesuaikan dengan nama proyekmu
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateBooks membuat pengguna baru
func CreateBook(c *fiber.Ctx) error {
	Books := new(models.Book)

	if err := c.BodyParser(Books); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if result := database.DBClient.Create(&Books); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusCreated, "Books created successfully", Books)
}

// GetAllBook mendapatkan semua pengguna
func GetAllBooks(c *fiber.Ctx) error {
	var books []models.Book
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
	if result := database.DBClient.Limit(limit).Offset(offset).Find(&books); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}
	if len(books) == 0 && page > 1 { // Jika halaman lebih dari 1 dan tidak ada buku, berarti halaman kosong
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "No books found on this page")
	} else if len(books) == 0 { // Jika di halaman 1 pun tidak ada buku sama sekali
		return helpers.SuccessResponse(c, fiber.StatusOK, "No books found", []models.Book{})
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	return helpers.SuccessResponse(c, fiber.StatusOK, "Books retrieved successfully", fiber.Map{
		"data":         books,
		"total_items":  total,
		"current_page": page,
		"per_page":     limit,
		"total_pages":  totalPages,
	})
}

// GetBooksByID mendapatkan pengguna berdasarkan ID
func GetBooksByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	bookID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid book ID format")
	}
	Books := new(models.Book)

	if result := database.DBClient.First(&Books, bookID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Books not found")
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Books retrieved successfully", Books)
}

// UpdateBooks memperbarui pengguna
func UpdateBooks(c *fiber.Ctx) error {
	idStr := c.Params("id")

	BooksID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Books ID format")
	}

	Books := new(models.Book)

	if result := database.DBClient.First(&Books, BooksID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Books not found")
	}

	updates := new(models.Book)
	if err := c.BodyParser(updates); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if updates.Title != "" {
		Books.Title = updates.Title
	}
	if updates.Author != "" {
		Books.Author = updates.Author
	}
	if updates.Isbn != "" {
		Books.Isbn = updates.Isbn
	}
	if updates.Quantity != "" {
		Books.Quantity = updates.Quantity
	}
	if updates.Category != "" {
		Books.Category = updates.Category
	}
	if result := database.DBClient.Save(&Books); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Books updated successfully", Books)
}

// DeleteBooks menghapus pengguna
func DeleteBooks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	BooksID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Books ID format")
	}

	Books := new(models.Book)

	if result := database.DBClient.First(&Books, BooksID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Books not found")
	}

	if result := database.DBClient.Delete(&Books); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Books deleted successfully", nil)
}
