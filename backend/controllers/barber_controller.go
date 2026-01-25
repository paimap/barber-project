package controllers

import (
	"net/http"
	"time"

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

type BarberStatsResponse struct {
    BarberID     uint   `json:"barber_id"`
    BarberName   string `json:"barber_name"`
    TotalService int64  `json:"total_service"`
    TotalRevenue int64  `json:"total_revenue"`
}
func GetBarberStatsToday(c *gin.Context) {
    db := config.DB
	userID := c.GetUint("user_id")

	var barber models.Barber
	if err := db.Where("user_id = ?", userID).First(&barber).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error" : "barber not found"})
	}

    // 1. Tentukan rentang waktu "Hari Ini" (00:00:00 sampai sekarang)
    now := time.Now()
    todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

    var stats BarberStatsResponse

    // 2. Query Agregasi dengan Filter Waktu
    err := db.Model(&models.Service{}).
        Select("services.barber_id, barbers.name as barber_name, COUNT(services.id) as total_service, SUM(services.price_at_sale) as total_revenue").
        Joins("JOIN barbers ON barbers.id = services.barber_id").
        Where("services.barber_id = ? AND services.created_at >= ?", barber.ID, todayStart).
        Group("services.barber_id, barbers.name").
        Scan(&stats).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data hari ini"})
        return
    }

    // 3. Penanganan jika belum ada service hari ini
    if stats.BarberID == 0 {
        var barberSelect models.Barber
        db.First(&barberSelect, barber.ID)
        stats.BarberID = barberSelect.ID
        stats.BarberName = barberSelect.Name
        stats.TotalService = 0
        stats.TotalRevenue = 0
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "period": "today",
        "data":   stats,
    })
}
