package middleware

import (
	"fmt"
	"library/config"  // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"library/helpers" // SESUAIKAN DENGAN NAMA MODUL GO ANDA
	"strings"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthRequired adalah middleware untuk memverifikasi token JWT (Access Token)
func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization header required")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == "" {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Bearer token not found")
	}

	cfg := config.LoadConfig() // Muat konfigurasi untuk secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode penandatanganan adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Mengembalikan secret key
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		// Log error lebih detail untuk debugging
		fmt.Printf("JWT parsing error: %v\n", err)
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Set user ID atau informasi lainnya ke dalam context Fiber jika diperlukan
		// Ini memungkinkan controller mengakses user_id
		c.Locals("userID", claims["user_id"])
		return c.Next() // Lanjutkan ke handler berikutnya
	} else {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token claims")
	}
}

// GenerateAccessToken menghasilkan Access Token JWT
func GenerateAccessToken(userID uuid.UUID, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Minute * 30)), // Access Token berlaku 30 menit
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateRefreshToken menghasilkan Refresh Token JWT
func GenerateRefreshToken(userID uuid.UUID, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Refresh Token berlaku 7 hari
	})

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
