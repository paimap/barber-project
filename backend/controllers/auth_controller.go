package controllers

import (
	"net/http"
	"gorm.io/gorm"

	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role     string `json:"role" binding:"required,oneof=SUPERADMIN MITRA BARBER"`
    Name     string `json:"name" binding:"required"`
    Phone    string `json:"phone" binding:"required"`
}


type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword(
        []byte(input.Password),
        bcrypt.DefaultCost,
    )

    err := config.DB.Transaction(func(tx *gorm.DB) error {
        user := models.User{
            Email:    input.Email,
            Password: string(hashedPassword),
            Role:     input.Role,
        }

        if err := tx.Create(&user).Error; err != nil {
            return err
        }

        switch input.Role {
        case models.RoleSuperAdmin:
            return tx.Create(&models.SuperAdmin{
                UserID:      user.ID,
                Name:        input.Name,
                PhoneNumber: input.Phone,
            }).Error

        case models.RoleMitra:
            return tx.Create(&models.Mitra{
                UserID:      user.ID,
                Name:        input.Name,
                PhoneNumber: input.Phone,
            }).Error

        case models.RoleBarber:
            return tx.Create(&models.Barber{
                UserID:      user.ID,
                Name:        input.Name,
                PhoneNumber: input.Phone,
            }).Error
        }

        return nil
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "register success"})
}


func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}



