package main

import (
	"log"

	"catatan-api/config"
	"catatan-api/database"
	"catatan-api/routes"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // opsional: baca file .env
	cfg := config.Load()
	db := database.Connect(cfg.DBPath)
	r := routes.Setup(db, cfg)
	log.Printf("server berjalan di http://localhost:%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
