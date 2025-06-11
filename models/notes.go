package models

import (
	"github.com/google/uuid"
)

type Note struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt uint      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt uint      `gorm:"autoUpdateTime" json:"updated_at"`
}
