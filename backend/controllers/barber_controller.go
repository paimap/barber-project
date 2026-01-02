package controllers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func CreateBarber(c *gin.Context) {
	userID := c.GetUint("user_id")

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

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

	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
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
			Role:     models.RoleBarber,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		barber := models.Barber{
			Name:        input.Name,
			UserID:      user.ID,
			PhoneNumber: input.Phone,
			MitraID:     mitra.ID,
		}

		if err := tx.Create(&barber).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "barber created successfully"})
}

func GetBarberByMitra(c *gin.Context) {
	userID := c.GetUint("user_id")

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var barbers []models.Barber
	if err := config.DB.Preload("Outlet").Where("mitra_id = ?", mitra.ID).Find(&barbers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response with outlet address
	response := make([]map[string]interface{}, 0, len(barbers))
	for _, b := range barbers {
		barberMap := map[string]interface{}{
			"ID":          b.ID,
			"Name":        b.Name,
			"PhoneNumber": b.PhoneNumber,
			"OutletID":    b.OutletID,
			"CreatedAt":   b.CreatedAt,
		}
		if b.OutletID != 0 && b.Outlet.Address != "" {
			barberMap["OutletAddress"] = b.Outlet.Address
		} else {
			barberMap["OutletAddress"] = nil
		}
		response = append(response, barberMap)
	}

	c.JSON(http.StatusOK, gin.H{"barbers": response})
}

func UpdateBarber(c *gin.Context) {
	userID := c.GetUint("user_id")
	barberID := c.Param("id")

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var barber models.Barber
	if err := config.DB.Where("id = ? AND mitra_id = ?", barberID, mitra.ID).First(&barber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "barber does not belong to this mitra"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Phone != "" {
		updates["phone_number"] = input.Phone
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := config.DB.Model(&barber).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "barber updated successfully"})
}

func DeleteBarber(c *gin.Context) {
	userID := c.GetUint("user_id")
	barberID := c.Param("id")

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var barber models.Barber
	if err := config.DB.Where("id = ? AND mitra_id = ?", barberID, mitra.ID).First(&barber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "barber does not belong to this mitra"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, barber.UserID).Error; err != nil {
			return err
		}

		if err := tx.Delete(&barber).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "barber deleted successfully"})
}
func AssignBarberToOutlet(c *gin.Context) {
	userID := c.GetUint("user_id")
	barberID := c.Param("id")

	var mitra models.Mitra
	if err := config.DB.Where("user_id = ?", userID).First(&mitra).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a mitra"})
		return
	}

	var input struct {
		OutletID uint `json:"outlet_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi barber milik mitra
	var barber models.Barber
	if err := config.DB.Where("id = ? AND mitra_id = ?", barberID, mitra.ID).First(&barber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "barber does not belong to this mitra"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var outlet models.Outlet
	if err := config.DB.Where("id = ? AND mitra_id = ?", input.OutletID, mitra.ID).First(&outlet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "outlet does not belong to this mitra"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&barber).Update("outlet_id", input.OutletID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "barber assigned to outlet successfully"})
}
