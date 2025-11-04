package models

import "gorm.io/gorm"

type Book struct {
	Id        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	Author    string `json:"author"`
	Title     string `gorm:"not null" json:"title"`
	Publisher string `gorm:"not null" json:"publisher"`
	Year      int    `gorm:"not null" json:"year"`
	User      User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
