package model

import (
	"time"

	"gorm.io/gorm"
)

type Resource struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Description string         `json:"description"`
	Unit        string         `json:"unit"` // кг, л и т.п.
	Quantity    int            `json:"quantity"`
}
