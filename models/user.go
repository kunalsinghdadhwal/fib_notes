package models

import (
	"github.com/google/uuid"
	"github.com/kunalsinghdadhwal/fib_notes/utils"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Notes     []Note    `gorm:"foreignKey:UserID" json:"notes,omitempty"`
	CreatedAt uint      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt uint      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	if u.Password != "" {
		hashedPassword, err := utils.Hash(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}

	return
}

func (u *User) CheckPassword(password string) error {
	return utils.CheckPassword(u.Password, password)
}
