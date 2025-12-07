package db

// Admin 记录有特权的用户ID
type Admin struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	Description string `gorm:"size:256" json:"description"`
}
