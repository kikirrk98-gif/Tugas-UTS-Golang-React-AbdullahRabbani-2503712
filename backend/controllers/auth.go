package controllers

import (
	"strings"
	"time"

	"catatan-api/middleware"
	"catatan-api/models"
	"catatan-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB     *gorm.DB
	Secret []byte
}

type registerInput struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthController) token(u models.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return t.SignedString(a.Secret)
}

func (a *AuthController) Register(c *gin.Context) {
	var in registerInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Fail(c, 422, "data tidak valid: nama min 2 huruf, email benar, password min 6 karakter")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(c, 500, "gagal memproses password")
		return
	}
	user := models.User{Name: strings.TrimSpace(in.Name), Email: strings.ToLower(in.Email), Password: string(hash)}
	if err := a.DB.Create(&user).Error; err != nil {
		utils.Fail(c, 409, "email sudah terdaftar")
		return
	}
	tok, _ := a.token(user)
	utils.OK(c, 201, "registrasi berhasil", gin.H{"token": tok, "user": user})
}

func (a *AuthController) Login(c *gin.Context) {
	var in loginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Fail(c, 422, "email dan password wajib diisi")
		return
	}
	var user models.User
	if a.DB.Where("email = ?", strings.ToLower(in.Email)).First(&user).Error != nil ||
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
		utils.Fail(c, 401, "email atau password salah")
		return
	}
	tok, _ := a.token(user)
	utils.OK(c, 200, "login berhasil", gin.H{"token": tok, "user": user})
}

func (a *AuthController) Me(c *gin.Context) {
	var user models.User
	if err := a.DB.First(&user, middleware.UserID(c)).Error; err != nil {
		utils.Fail(c, 404, "user tidak ditemukan")
		return
	}
	utils.OK(c, 200, "profil", user)
}
