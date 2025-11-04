package models

import "gorm.io/gorm"

type Book struct {
	Id        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Author    string `json:"author"`
	Title     string `gorm:"not null" json:"title"`
	Publisher string `gorm:"not null" json:"publisher"`
	Year      int    `gorm:"not null" json:"year"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
