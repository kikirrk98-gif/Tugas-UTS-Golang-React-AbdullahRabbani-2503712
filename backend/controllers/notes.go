package controllers

import (
	"errors"
	"strconv"
	"strings"

	"catatan-api/middleware"
	"catatan-api/models"
	"catatan-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoteController struct {
	DB *gorm.DB
}

type noteInput struct {
	Title    string `json:"title" binding:"required,min=1,max=200"`
	Content  string `json:"content"`
	Category string `json:"category" binding:"max=50"`
	Pinned   bool   `json:"pinned"`
}

// scope membatasi query ke catatan milik user yang sedang login.
func (n *NoteController) mine(c *gin.Context) *gorm.DB {
	return n.DB.Where("user_id = ?", middleware.UserID(c))
}

func (n *NoteController) find(c *gin.Context) (models.Note, bool) {
	var note models.Note
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, 400, "id tidak valid")
		return note, false
	}
	if err := n.mine(c).First(&note, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Fail(c, 404, "catatan tidak ditemukan")
		} else {
			utils.Fail(c, 500, "kesalahan database")
		}
		return note, false
	}
	return note, true
}

func (n *NoteController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "9"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 9
	}
	// Session membuat tx aman dipakai ulang setelah Count.
	tx := n.mine(c).Model(&models.Note{}).Session(&gorm.Session{})
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		like := "%" + q + "%"
		tx = tx.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	if cat := c.Query("category"); cat != "" {
		tx = tx.Where("category = ?", cat)
	}
	var total int64
	tx.Count(&total)

	notes := []models.Note{}
	tx.Order("pinned DESC, updated_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&notes)

	pages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(200, utils.Response{Success: true, Message: "daftar catatan", Data: notes,
		Meta: gin.H{"page": page, "limit": limit, "total": total, "pages": pages}})
}

func (n *NoteController) Create(c *gin.Context) {
	var in noteInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Fail(c, 422, "judul wajib diisi (maks 200 karakter)")
		return
	}
	cat := strings.TrimSpace(in.Category)
	if cat == "" {
		cat = "Umum"
	}
	note := models.Note{UserID: middleware.UserID(c), Title: strings.TrimSpace(in.Title),
		Content: in.Content, Category: cat, Pinned: in.Pinned}
	if err := n.DB.Create(&note).Error; err != nil {
		utils.Fail(c, 500, "gagal menyimpan catatan")
		return
	}
	utils.OK(c, 201, "catatan dibuat", note)
}

func (n *NoteController) Show(c *gin.Context) {
	if note, ok := n.find(c); ok {
		utils.OK(c, 200, "detail catatan", note)
	}
}

func (n *NoteController) Update(c *gin.Context) {
	note, ok := n.find(c)
	if !ok {
		return
	}
	var in noteInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.Fail(c, 422, "judul wajib diisi (maks 200 karakter)")
		return
	}
	cat := strings.TrimSpace(in.Category)
	if cat == "" {
		cat = "Umum"
	}
	// map dipakai agar nilai false/kosong tetap ikut tersimpan
	n.DB.Model(&note).Updates(map[string]any{
		"title": strings.TrimSpace(in.Title), "content": in.Content, "category": cat, "pinned": in.Pinned,
	})
	n.DB.First(&note, note.ID)
	utils.OK(c, 200, "catatan diperbarui", note)
}

func (n *NoteController) TogglePin(c *gin.Context) {
	note, ok := n.find(c)
	if !ok {
		return
	}
	n.DB.Model(&note).Update("pinned", !note.Pinned)
	n.DB.First(&note, note.ID)
	utils.OK(c, 200, "status pin diubah", note)
}

func (n *NoteController) Delete(c *gin.Context) {
	note, ok := n.find(c)
	if !ok {
		return
	}
	n.DB.Delete(&note)
	utils.OK(c, 200, "catatan dihapus", nil)
}

func (n *NoteController) Stats(c *gin.Context) {
	type row struct {
		Category string `json:"category"`
		Total    int64  `json:"total"`
	}
	var rows []row
	n.mine(c).Model(&models.Note{}).Select("category, count(*) as total").
		Group("category").Order("total DESC").Scan(&rows)
	var total, pinned int64
	n.mine(c).Model(&models.Note{}).Count(&total)
	n.mine(c).Model(&models.Note{}).Where("pinned = ?", true).Count(&pinned)
	if rows == nil {
		rows = []row{}
	}
	utils.OK(c, 200, "statistik", gin.H{"total": total, "pinned": pinned, "categories": rows})
}
