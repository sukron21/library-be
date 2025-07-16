package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan model pengguna
type Book struct {
	gorm.Model
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Isbn     string    `json:"isbn" gorm:"unique"`
	Quantity string    `json:"quantity"`
	Category string    `json:"category"`
}

func (u *Book) BeforeCreate(tx *gorm.DB) (err error) {
	// Ensure ID is unique if not already set (e.g., by DB default or manually)
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
