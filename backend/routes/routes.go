package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"catatan-api/config"
	"catatan-api/controllers"
	"catatan-api/middleware"
	"catatan-api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.Origin},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	auth := &controllers.AuthController{DB: db, Secret: cfg.JWTSecret}
	notes := &controllers.NoteController{DB: db}

	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) { utils.OK(c, 200, "ok", gin.H{"time": time.Now()}) })
	api.POST("/auth/register", auth.Register)
	api.POST("/auth/login", auth.Login)

	secured := api.Group("", middleware.Auth(cfg.JWTSecret))
	{
		secured.GET("/me", auth.Me)
		secured.GET("/stats", notes.Stats)
		secured.GET("/notes", notes.List)
		secured.POST("/notes", notes.Create)
		secured.GET("/notes/:id", notes.Show)
		secured.PUT("/notes/:id", notes.Update)
		secured.PATCH("/notes/:id/pin", notes.TogglePin)
		secured.DELETE("/notes/:id", notes.Delete)
	}

	// Sajikan hasil build React bila folder dist ada (mode produksi satu biner).
	if _, err := os.Stat(filepath.Join(cfg.StaticDir, "index.html")); err == nil {
		r.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				utils.Fail(c, 404, "endpoint tidak ditemukan")
				return
			}
			file := filepath.Join(cfg.StaticDir, filepath.Clean("/"+c.Request.URL.Path))
			if st, err := os.Stat(file); err == nil && !st.IsDir() {
				c.File(file)
				return
			}
			c.File(filepath.Join(cfg.StaticDir, "index.html"))
		})
	} else {
		r.NoRoute(func(c *gin.Context) { c.JSON(http.StatusNotFound, utils.Response{Message: "tidak ditemukan"}) })
	}
	return r
}
