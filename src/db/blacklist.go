package db

import "gorm.io/gorm"

// Blacklist 记录被封禁的网址
type Blacklist struct {
	gorm.Model
	URL         string `gorm:"uniqueIndex;not null" json:"url"`
	Description string `gorm:"size:256" json:"description"`
}