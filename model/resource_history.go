package model

import (
	"time"

	"gorm.io/gorm"
)

// ResourceHistory tracks all changes made to resources
type ResourceHistory struct {
	gorm.Model
	ResourceID  uint      `gorm:"not null" json:"resource_id"`                         // ID of the resource being changed
	Action      string    `gorm:"not null" json:"action"`                              // CREATE, UPDATE, DELETE
	UserID      uint      `gorm:"not null" json:"user_id"`                             // User who made the change
	OldData     string    `gorm:"type:text" json:"old_data,omitempty"`                 // JSON of old data (for UPDATE/DELETE)
	NewData     string    `gorm:"type:text" json:"new_data,omitempty"`                 // JSON of new data (for CREATE/UPDATE)
	Timestamp   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"timestamp"` // When the change happened
	Description string    `json:"description,omitempty"`                               // Optional description of the change

	// Relations
	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource,omitempty"`
	User     User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}
