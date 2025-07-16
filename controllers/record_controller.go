package controllers

import (
	"library/database" // Sesuaikan dengan nama proyekmu
	"library/helpers"  // Sesuaikan dengan nama proyekmu
	"library/models"   // Sesuaikan dengan nama proyekmu
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateRecords membuat pengguna baru
func CreateRecord(c *fiber.Ctx) error {
	record := new(models.Lending_records)

	if err := c.BodyParser(record); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if result := database.DBClient.Create(&record); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	// Memuat data terkait 'Book' dan 'User'
	if err := database.DBClient.Preload("Book").Preload("User").First(&record, record.ID).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to preload related data: "+err.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusCreated, "Record created successfully", record)
}

// GetAllUsers mendapatkan semua pengguna
func GetAllRecord(c *fiber.Ctx) error {
	var Records []models.Lending_records
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

	if result := database.DBClient.Model(&models.Lending_records{}).Count(&total); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}
	if result := database.DBClient.Preload("Book").
		Preload("User").Limit(limit).Offset(offset).Find(&Records); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}
	if len(Records) == 0 && page > 1 { // Jika halaman lebih dari 1 dan tidak ada buku, berarti halaman kosong
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "No Records found on this page")
	} else if len(Records) == 0 { // Jika di halaman 1 pun tidak ada buku sama sekali
		return helpers.SuccessResponse(c, fiber.StatusOK, "No Records found", []models.Lending_records{})
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	return helpers.SuccessResponse(c, fiber.StatusOK, "Records retrieved successfully", fiber.Map{
		"data":         Records,
		"total_items":  total,
		"current_page": page,
		"per_page":     limit,
		"total_pages":  totalPages,
	})
}

// GetRecordsByID mendapatkan pengguna berdasarkan ID
func GetRecordByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	RecordID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Record ID format")
	}
	Records := new(models.Lending_records)

	if result := database.DBClient.First(&Records, RecordID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Records not found")
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Records retrieved successfully", Records)
}

// UpdateRecords memperbarui pengguna
func UpdateRecords(c *fiber.Ctx) error {
	idStr := c.Params("id")

	RecordsID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Records ID format")
	}

	Records := new(models.Lending_records)

	if result := database.DBClient.First(&Records, RecordsID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Records not found")
	}

	updates := new(models.Lending_records)
	if err := c.BodyParser(updates); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if updates.Book_id != "" {
		Records.Book_id = updates.Book_id
	}
	if updates.User_id != "" {
		Records.User_id = updates.User_id
	}
	if updates.Borrow_date.IsZero() {
		Records.Borrow_date = updates.Borrow_date
	}
	if updates.ReturnDate.IsZero() {
		Records.ReturnDate = updates.ReturnDate
	}

	if result := database.DBClient.Save(&Records); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}
	if err := database.DBClient.Model(&Records).Association("Book").Find(&Records.Book); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to load Book data")
	}
	if err := database.DBClient.Model(&Records).Association("User").Find(&Records.User); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to load User data")
	}
	return helpers.SuccessResponse(c, fiber.StatusOK, "Records updated successfully", Records)
}

// DeleteRecords menghapus pengguna
func DeleteRecords(c *fiber.Ctx) error {
	idStr := c.Params("id")
	RecordsID, err := uuid.Parse(idStr)
	if err != nil {
		// Jika ID dari URL bukan UUID yang valid
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Records ID format")
	}

	Records := new(models.Lending_records)

	if result := database.DBClient.First(&Records, RecordsID); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Records not found")
	}

	if result := database.DBClient.Delete(&Records); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, result.Error.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Records deleted successfully", nil)
}
