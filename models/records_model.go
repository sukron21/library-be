package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan model pengguna
type Lending_records struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Book_id     string     `json:"book_id"`
	User_id     string     `json:"user_id"`
	Borrow_date time.Time  `json:"borrow_date"`
	ReturnDate  *time.Time `json:"return_date"`

	Book Book `gorm:"foreignKey:Book_id;references:ID"`
	User User `gorm:"foreignKey:User_id;references:ID"`
}

func (u *Lending_records) BeforeCreate(tx *gorm.DB) (err error) {
	// Ensure ID is unique if not already set (e.g., by DB default or manually)
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
