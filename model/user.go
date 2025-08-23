package model

import (
	"time"

	"gorm.io/gorm"
)

// User struct
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"uniqueIndex;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email"`
	Password  string         `gorm:"not null;" validate:"required,min=6,max=50" json:"password"`
	Names     string         `json:"names"`
}
