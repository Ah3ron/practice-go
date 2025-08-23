package model

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"` // кг, л и т.п.
	Quantity    int    `json:"quantity"`
}
