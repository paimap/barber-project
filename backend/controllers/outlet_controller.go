package controllers

import (
	"net/http"
	"strconv"

	"backend/models"
	"backend/config"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func CreateOutlet(c *gin.Context) {
	userID:= c.GetUint("user_id")

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var input struct {
		Address string `json:"address" binding:"required"`
		Phone   string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outlet := models.Outlet{
		Address:       input.Address,
		PhoneNumber:   input.Phone,
		MitraID:       mitra.ID, // pakai Mitra.ID, bukan User.ID
	}

	if err := config.DB.Create(&outlet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"outlet": outlet})
}

func GetOutlets(c *gin.Context){
	userID := c.GetUint("user_id")

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var outlets []models.Outlet
	if err := config.DB.Where("mitra_id = ?", mitra.ID).Find(&outlets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"outlets": outlets})
}

func UpdateOutlet(c *gin.Context) {
	userID := c.GetUint("user_id")
	outletIDParam := c.Param("id")
	outletID, _ := strconv.ParseUint(outletIDParam, 10, 64)

	owns, err := middleware.EnsureMitraOwnsOutlet(userID, uint(outletID))
	if err != nil || !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	var input struct {
		Address string `json:"address"`
		Phone   string `json:"phone"`
	}
	c.ShouldBindJSON(&input)

	var outlet models.Outlet
	if err := config.DB.First(&outlet, outletID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "outlet not found"})
		return
	}

	config.DB.Model(&outlet).Updates(models.Outlet{
		Address:     input.Address,
		PhoneNumber: input.Phone,
	})

	c.JSON(http.StatusOK, gin.H{"outlet": outlet})
}

func DeleteOutlet(c *gin.Context) {
	userID := c.GetUint("user_id")
	outletIDParam := c.Param("id")
	outletID, _ := strconv.ParseUint(outletIDParam, 10, 64)

	owns, err := middleware.EnsureMitraOwnsOutlet(userID, uint(outletID))
	if err != nil || !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
		return
	}

	if err := config.DB.Delete(&models.Outlet{}, outletID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}