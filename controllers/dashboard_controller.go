// controllers/dashboard_controller.go
package controllers

import (
	"fmt"
	"library/database" // Sesuaikan dengan nama modul Go Anda
	"library/helpers"
	"library/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// DashboardSummaryResponse adalah struktur untuk respons ringkasan dashboard
type DashboardSummaryResponse struct {
	TotalBooks          int64   `json:"total_books"`
	TotalMembers        int64   `json:"total_members"` // Menggunakan TotalMembers, asumsikan semua user adalah anggota
	BorrowingsThisMonth int64   `json:"borrowings_this_month"`
	AvgDailyBorrowings  float64 `json:"avg_daily_borrowings"`
}

// GetDashboardSummary mengambil metrik ringkasan dashboard
func GetDashboardSummary(c *fiber.Ctx) error {
	var totalBooks int64
	if err := database.DBClient.Model(&models.Book{}).Count(&totalBooks).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil total buku: "+err.Error())
	}

	var totalMembers int64
	// Menghitung total semua user sebagai anggota
	if err := database.DBClient.Model(&models.User{}).Count(&totalMembers).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil total anggota: "+err.Error())
	}

	// Hitung peminjaman bulan ini
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond) // Hari terakhir bulan ini

	var borrowingsThisMonth int64
	if err := database.DBClient.Model(&models.Lending_records{}).
		Where("borrow_date BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Count(&borrowingsThisMonth).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil peminjaman bulan ini: "+err.Error())
	}

	// Hitung rata-rata peminjaman harian untuk bulan ini
	var distinctBorrowingDays int64
	// Penting: Pastikan nama tabel 'records' dan kolom 'borrow_date' sesuai di database
	sqlQuery := `SELECT COUNT(DISTINCT DATE(borrow_date)) FROM lending_records WHERE borrow_date BETWEEN ? AND ?`
	if err := database.DBClient.Raw(sqlQuery, startOfMonth, endOfMonth).Scan(&distinctBorrowingDays).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghitung hari peminjaman unik: "+err.Error())
	}

	var avgDailyBorrowings float64
	if distinctBorrowingDays > 0 {
		avgDailyBorrowings = float64(borrowingsThisMonth) / float64(distinctBorrowingDays)
	} else {
		avgDailyBorrowings = 0
	}

	response := DashboardSummaryResponse{
		TotalBooks:          totalBooks,
		TotalMembers:        totalMembers,
		BorrowingsThisMonth: borrowingsThisMonth,
		AvgDailyBorrowings:  avgDailyBorrowings,
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Ringkasan dashboard berhasil diambil", response)
}

// MonthlyTrendResponse adalah struktur untuk respons tren bulanan
type MonthlyTrendResponse struct {
	MonthName  string `json:"month"`
	MonthNum   int    `json:"month_num"`
	Borrowings int64  `json:"borrowings"`
	Returns    int64  `json:"returns"`
}

// GetMonthlyBorrowingTrend mengambil data tren peminjaman dan pengembalian bulanan
func GetMonthlyBorrowingTrend(c *fiber.Ctx) error {
	yearStr := c.Query("year", fmt.Sprintf("%d", time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Format tahun tidak valid")
	}

	var results []MonthlyTrendResponse
	// Inisialisasi data untuk 12 bulan
	for i := 1; i <= 12; i++ {
		results = append(results, MonthlyTrendResponse{
			MonthName:  time.Month(i).String()[:3], // Jan, Feb, Mar, dst.
			MonthNum:   i,
			Borrowings: 0,
			Returns:    0,
		})
	}

	// Kueri untuk peminjaman
	var borrowedData []struct {
		Month int   `gorm:"column:month"`
		Count int64 `gorm:"column:count"`
	}
	// Menggunakan EXTRACT(MONTH FROM borrow_date) untuk mendapatkan bulan
	// PASTIKAN database Anda mendukung fungsi ini (PostgreSQL, MySQL).
	if err := database.DBClient.Model(&models.Lending_records{}).
		Select("EXTRACT(MONTH FROM borrow_date) as month, COUNT(*) as count").
		Where("EXTRACT(YEAR FROM borrow_date) = ?", year).
		Group("month").
		Order("month asc").
		Scan(&borrowedData).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data peminjaman bulanan: "+err.Error())
	}

	for _, data := range borrowedData {
		for i := range results {
			if results[i].MonthNum == data.Month {
				results[i].Borrowings = data.Count
				break
			}
		}
	}

	// Kueri untuk pengembalian (hanya jika return_date tidak NULL)
	var returnedData []struct {
		Month int   `gorm:"column:month"`
		Count int64 `gorm:"column:count"`
	}
	if err := database.DBClient.Model(&models.Lending_records{}).
		Select("EXTRACT(MONTH FROM return_date) as month, COUNT(*) as count").
		Where("EXTRACT(YEAR FROM return_date) = ? AND return_date IS NOT NULL", year).
		Group("month").
		Order("month asc").
		Scan(&returnedData).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data pengembalian bulanan: "+err.Error())
	}

	for _, data := range returnedData {
		for i := range results {
			if results[i].MonthNum == data.Month {
				results[i].Returns = data.Count
				break
			}
		}
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Tren peminjaman bulanan berhasil diambil", results)
}

type LatestActivityResponse struct {
	Type      string    `json:"type"`
	UserName  string    `json:"user_name"`
	BookTitle string    `json:"book_title"`
	Date      time.Time `json:"date"`
}

// controllers/dashboard_controller.go
// ... (Your existing code)

// GetLatestActivity mengambil daftar aktivitas peminjaman dan pengembalian terbaru
func GetLatestActivity(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Limit tidak valid")
	}

	var records []models.Lending_records
	// Assuming you've correctly updated models/record_model.go for GORM relations (Option 1)
	if err := database.DBClient.
		Preload("User").
		Preload("Book").
		Order("borrow_date DESC"). // Ordering by created_at of the record itself
		Limit(limit).
		Find(&records).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil aktivitas terbaru: "+err.Error())
	}

	var activities []LatestActivityResponse
	for _, record := range records {
		// Add borrow activity
		activities = append(activities, LatestActivityResponse{
			Type:      "borrow",
			UserName:  record.User.Name,
			BookTitle: record.Book.Title,
			Date:      record.Borrow_date,
		})

		// Add return activity ONLY if ReturnDate is not nil and not a zero time value
		// Use record.ReturnDate (camelCase)
		if record.ReturnDate != nil && !record.ReturnDate.IsZero() {
			activities = append(activities, LatestActivityResponse{
				Type:      "return",
				UserName:  record.User.Name,
				BookTitle: record.Book.Title,
				Date:      *record.ReturnDate, // Dereference the pointer
			})
		}
	}

	// Sort all activities (borrow and return) by Date in descending order
	// This uses a simple bubble sort; for larger datasets, consider `sort.Slice` for efficiency.
	for i := 0; i < len(activities); i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[j].Date.After(activities[i].Date) { // Sort Descending (latest first)
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}

	// Ensure the final slice doesn't exceed the requested limit after adding both borrow/return activities
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Aktivitas terbaru berhasil diambil", activities)
}

//

// GetTopBorrowedBooks: Buku Paling Banyak Dipinjam 🏆

//Ini akan mengembalikan top buku berdasarkan jumlah peminjaman.

// TopBorrowedBookResponse adalah struktur untuk respons buku paling banyak dipinjam
type TopBorrowedBookResponse struct {
	BookTitle   string `json:"title"`
	BorrowCount int64  `json:"borrow_count"`
}

// GetTopBorrowedBooks mengambil daftar buku yang paling banyak dipinjam
func GetTopBorrowedBooks(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "7")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Limit tidak valid")
	}

	// Ambil bulan dan tahun dari query, jika tidak ada, gunakan bulan dan tahun sekarang
	monthStr := c.Query("month", fmt.Sprintf("%d", time.Now().Month()))
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Bulan tidak valid")
	}

	yearStr := c.Query("year", fmt.Sprintf("%d", time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Tahun tidak valid")
	}

	// Hitung start dan end date untuk periode yang diminta
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var topBooks []TopBorrowedBookResponse

	// Join Records dengan Books, GROUP BY BookID dan hitung jumlah peminjaman
	if err := database.DBClient.Model(&models.Lending_records{}).
		Select("books.title, COUNT(lending_records.id) as borrow_count").
		Joins("JOIN books ON lending_records.book_id = books.id"). // Pastikan nama tabel 'books'
		Where("lending_records.borrow_date BETWEEN ? AND ?", startDate, endDate).
		Group("books.id, books.title"). // Penting: sertakan books.title di GROUP BY jika di SELECT
		Order("borrow_count DESC").
		Limit(limit).
		Scan(&topBooks).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil buku paling banyak dipinjam: "+err.Error())
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Buku paling banyak dipinjam berhasil diambil", topBooks)
}

type CategoryDistributionResponse struct {
	Category  string `json:"category"`
	BookCount int64  `json:"book_count"`
}

// GetBookCategoriesDistribution mengambil distribusi buku per kategori
func GetBookCategoriesDistribution(c *fiber.Ctx) error {
	var categories []CategoryDistributionResponse

	// Group Books berdasarkan Category dan hitung jumlahnya
	if err := database.DBClient.Model(&models.Book{}).
		Select("category, COUNT(id) as book_count").
		Group("category").
		Order("book_count DESC").
		Scan(&categories).Error; err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil distribusi kategori buku: "+err.Error())
	}

	// Jika ingin menyertakan kategori yang tidak memiliki buku, perlu logika tambahan
	// Misalnya, ambil semua kategori unik lalu join dengan hasil count

	return helpers.SuccessResponse(c, fiber.StatusOK, "Distribusi kategori buku berhasil diambil", categories)
}
