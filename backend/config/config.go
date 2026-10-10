package config

import (
	"log"
	"os"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret []byte
	Origin    string
	StaticDir string
}

func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() Config {
	secret := get("JWT_SECRET", "")
	if secret == "" {
		log.Println("PERINGATAN: JWT_SECRET kosong, memakai secret dev. Jangan dipakai di produksi.")
		secret = "dev-secret-ganti-di-produksi"
	}
	return Config{
		Port:      get("APP_PORT", "8080"),
		DBPath:    get("DB_PATH", "catatan.db"),
		JWTSecret: []byte(secret),
		Origin:    get("CORS_ORIGIN", "http://localhost:5173"),
		StaticDir: get("STATIC_DIR", "../frontend/dist"),
	}
}
