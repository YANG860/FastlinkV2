package db

import "gorm.io/gorm"

// Admin 记录有特权的用户ID
type Admin struct {
	gorm.Model
	UserID      uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Description string `gorm:"size:256" json:"description"`
}



