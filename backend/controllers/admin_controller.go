package controllers

import (
	"net/http"
	"gorm.io/gorm"
	"strconv"

	"backend/models"
	"backend/config"
	
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllMitra(c *gin.Context){
	var mitra []models.Mitra
	if err := config.DB.Find(&mitra).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mitra": mitra})
}

func CreateMitra(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Email:    input.Email,
			Password: string(hashedPassword),
			Role:     models.RoleMitra,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		mitra := models.Mitra{
			Name:        input.Name,
			UserID:      user.ID,
			PhoneNumber: input.Phone,
		}

		if err := tx.Create(&mitra).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "mitra created successfully"})
}

func UpdateMitra(c *gin.Context) {
	mitraIdParam := c.Param("id")

	mitraId, err := strconv.ParseUint(mitraIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mitra id"})
		return
	}

	var input struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mitra models.Mitra
	if err := config.DB.First(&mitra, mitraId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mitra not found"})
		return
	}

	updates := map[string]interface{}{}

	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Phone != "" {
		updates["phone_number"] = input.Phone
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no fields to update",
		})
		return
	}

	if err := config.DB.Model(&mitra).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ambil data terbaru
	config.DB.First(&mitra, mitra.ID)

	c.JSON(http.StatusOK, gin.H{"mitra": mitra})
}

func DeleteMitra(c *gin.Context){
	mitraIdParam := c.Param("id")
	mitraId, err := strconv.ParseUint(mitraIdParam, 10, 64); if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mitra id"})
		return
	}

	var mitra models.Mitra
	if err := config.DB.First(&mitra, mitraId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&mitra).Error; err != nil{
			return err
		}

		if err := tx.Delete(&models.User{}, mitra.UserID).Error; err != nil{
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Deleted"})
}

