package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan model pengguna
type User struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"password"` // Jangan sertakan password saat marshal JSON
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Ensure ID is unique if not already set (e.g., by DB default or manually)
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
