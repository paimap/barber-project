package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/models"
)

func GetAllProduct(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func CreateProduct(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Price int64  `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:  input.Name,
		Price: input.Price,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product Created Successfully"})
}

func UpdateProduct(c *gin.Context) {
	productIdParam := c.Param("id")
	productId, err := strconv.ParseUint(productIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Name  string `json:"name"`
		Price int64  `json:"price"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ProductNotFound"})
		return
	}

	updates := map[string]interface{}{}

	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Price != 0 {
		updates["price"] = input.Price
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := config.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product Updated"})
}

func DeleteProduct(c *gin.Context) {
	productIdParam := c.Param("id")
	productId, err := strconv.ParseUint(productIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product Not Found"})
		return
	}

	tx := config.DB.Begin()

	// Soft delete Product_MainInventory yang terkait
	if err := tx.Delete(&models.Product_MainInventory{}, "product_id = ?", product.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Soft delete Product_OutletInventory yang terkait
	if err := tx.Delete(&models.Product_OutletInventory{}, "product_id = ?", product.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Soft delete Product
	if err := tx.Delete(&models.Product{}, product.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Product Deleted Successfully"})
}
