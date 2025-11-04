package models

import "gorm.io/gorm"

func RunMigrations(db *gorm.DB) error {
	if err := MigrateUser(db); err != nil {
		return err
	}
	if err := MigrateBooks(db); err != nil {
		return err
	}
	return nil
}
