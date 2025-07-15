package controllers

import (
	"fmt"
	"library/config"     // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"library/database"   // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"library/helpers"    // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"library/middleware" // Untuk generate token
	"library/models"     // SESUAIKAN DENGAN NAMA MODUL GO ANDA

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest struct untuk parsing body permintaan login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshTokenRequest struct untuk parsing body permintaan refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Login mengautentikasi pengguna dan mengembalikan Access Token & Refresh Token
func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	user := new(models.User)
	// Cari pengguna berdasarkan email
	if result := database.DBClient.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Bandingkan password yang di-hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	cfg := config.LoadConfig()

	// Generate Access Token
	accessToken, err := middleware.GenerateAccessToken(user.ID, cfg)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Could not generate access token")
	}

	// Generate Refresh Token
	refreshToken, err := middleware.GenerateRefreshToken(user.ID, cfg)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Could not generate refresh token")
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Login successful", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshAccessToken memperbarui Access Token menggunakan Refresh Token
func RefreshAccessToken(c *fiber.Ctx) error {
	req := new(RefreshTokenRequest)
	if err := c.BodyParser(req); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	cfg := config.LoadConfig()

	// Parse dan verifikasi Refresh Token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		// Log error lebih detail untuk debugging
		fmt.Printf("Refresh token parsing error: %v\n", err)
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token claims")
	}

	// Ambil userID dari claims
	userID := uint(claims["user_id"].(float64)) // claims["user_id"] biasanya di-parse sebagai float64

	// Generate Access Token baru
	newAccessToken, err := middleware.GenerateAccessToken(userID, cfg)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Could not generate new access token")
	}

	// Opsional: Generate Refresh Token baru juga (rotasi refresh token)
	// Jika Anda ingin mengimplementasikan rotasi refresh token,
	// tambahkan logika untuk membuat newRefreshToken dan update di database (jika disimpan).
	newRefreshToken, err := middleware.GenerateRefreshToken(userID, cfg)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, "Could not generate new refresh token")
	}

	return helpers.SuccessResponse(c, fiber.StatusOK, "Access token refreshed successfully", fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken, // Kirim refresh token baru jika rotasi diaktifkan
	})
}
