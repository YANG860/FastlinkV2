package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;size:64;not null"`
	PasswordHash string `gorm:"size:128;not null"`

	Banned        bool   `gorm:"not null;default:false"`
	AccessTokenID string `gorm:"size:128;not null"`
	CreatedLinks  []Link `gorm:"foreignKey:CreatorID;references:ID"`
}
