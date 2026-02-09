package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	log.Println("Running migrations...")

	if err := db.AutoMigrate(
		&User{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

func DropMigration(db *gorm.DB) error {
	log.Println("Dropping all tables...")

	if err := db.Migrator().DropTable(
		&User{},
	); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("All tables dropped")
	return nil
}
