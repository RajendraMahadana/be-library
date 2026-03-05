package database

import (
	"log"

	"gorm.io/gorm"
)

type BookModel struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"size:255;not null"`
	Author string `gorm:"size:255;not null"`
	ISBN   string `gorm:"size:100;unique"`
	Stock  int
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&BookModel{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("✅ Database migrated")
}