package middleware

import (
	"fmt"
	"library/config"  // Sesuaikan dengan nama modulmu
	"library/helpers" // Sesuaikan dengan nama modulmu
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

	// Hilangkan "Bearer " di depan
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == "" {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Bearer token not found")
	}

	cfg := config.LoadConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode penandatanganan adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		fmt.Printf("JWT parsing error: %v\n", err)
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Ambil user_id dari claims
		userIDClaim := claims["user_id"]
		if userIDClaim == nil {
			return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in token claims")
		}

		// Pastikan dalam bentuk string
		userIDStr := fmt.Sprintf("%v", userIDClaim)

		// Validasi jika benar UUID
		if _, err := uuid.Parse(userIDStr); err != nil {
			return helpers.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format in token")
		}

		// Simpan ke context
		c.Locals("userID", userIDStr)
		return c.Next()
	} else {
		return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token claims")
	}
}

// GenerateAccessToken menghasilkan Access Token JWT
func GenerateAccessToken(userID uuid.UUID, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),                                      // Pastikan UUID dikonversi string
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Minute * 30)), // Expired 30 menit
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
		"user_id": userID.String(),                                        // Pastikan UUID dikonversi string
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Expired 7 hari
	})
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
