package db

import "gorm.io/gorm"

// UrlBlacklist 记录被封禁的网址
type UrlBlacklist struct {
	gorm.Model
	URL         string `gorm:"uniqueIndex;not null" json:"url"`
	Description string `gorm:"size:256" json:"description"`
}