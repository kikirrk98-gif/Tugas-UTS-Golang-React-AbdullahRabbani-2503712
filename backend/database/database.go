package database

import (
	"log"

	"catatan-api/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("gagal membuka database: ", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Note{}); err != nil {
		log.Fatal("gagal migrasi: ", err)
	}
	return db
}
