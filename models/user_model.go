package models

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan model pengguna
type User struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"` // Jangan sertakan password saat marshal JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
