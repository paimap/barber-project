package controllers

import (
	"net/http"
	"strconv"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func CreateOutlet(c *gin.Context) {
	userID := c.GetUint("user_id")

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
		Address:     input.Address,
		PhoneNumber: input.Phone,
		MitraID:     mitra.ID, // pakai Mitra.ID, bukan User.ID
	}

	if err := config.DB.Create(&outlet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"outlet": outlet})
}

func GetOutlets(c *gin.Context) {
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
	outletID, err := strconv.ParseUint(outletIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid outlet id"})
		return
	}

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var outlet models.Outlet
	if err := config.DB.Where("id = ? AND mitra_id = ?", outletID, mitra.ID).First(&outlet).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "outlet does not belong to this mitra"})
		return
	}

	var input struct {
		Address string `json:"address"`
		Phone   string `json:"phone"`
	}
	c.ShouldBindJSON(&input)

	updates := map[string]interface{}{}
	if input.Address != "" {
		updates["address"] = input.Address
	}
	if input.Phone != "" {
		updates["phone_number"] = input.Phone
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := config.DB.Model(&outlet).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"outlet": outlet})
}

func DeleteOutlet(c *gin.Context) {
	userID := c.GetUint("user_id")
	outletIDParam := c.Param("id")
	outletID, err := strconv.ParseUint(outletIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid outlet id"})
		return
	}

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var outlet models.Outlet
	if err := config.DB.Where("id = ? AND mitra_id = ?", outletID, mitra.ID).First(&outlet).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "outlet does not belong to this mitra"})
		return
	}

	// Delete outlet (cascade delete will handle barbers and outlet_inventories)
	if err := config.DB.Delete(&outlet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "outlet deleted successfully"})
}
