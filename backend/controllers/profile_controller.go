package controllers

import (
	"net/http"
	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {
    db := config.DB
    
    // 1. Ambil userID yang sudah diset oleh middleware Auth
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. Cari user di database beserta rolenya
    var user models.User
    if err := db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
        return
    }

    // 3. Kembalikan data penting saja
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data": gin.H{
            "id":    user.ID,
            "email": user.Email,
            "role":  user.Role,
        },
    })
}
