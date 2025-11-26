package db

import "gorm.io/gorm"

const (
	LinkTypeGeneral = "general"
	LinkTypeOneShot = "one_shot"
	LinkTypePrivate = "private"
	LinkTypeCustom  = "custom"
)


type Link struct {
	gorm.Model

	Type      string `gorm:"size:32;not null" json:"type"`
	CreatorID uint   `gorm:"index;not null" json:"creator_id"`
	SourceURL string `gorm:"size:2048;not null" json:"source_url"`
	ShortCode string `gorm:"uniqueIndex;not null;size:32" json:"short_code"`
	Clicks    uint   `gorm:"default:0" json:"clicks"`
}


