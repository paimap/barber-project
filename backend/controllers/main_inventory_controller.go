package controllers

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func EnsureMainInventory() (uint, error) {
	var inventory models.MainInventory

	err := config.DB.First(&inventory).Error
	if err == nil {
		return inventory.ID, nil
	}

	if err := config.DB.Create(&inventory).Error; err != nil {
		return 0, err
	}

	return inventory.ID, nil
}

func UpsertMainInventoryItem(c *gin.Context) {
	var input struct {
		ProductID uint  `json:"product_id" binding:"required"`
		Quantity  int64 `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product not found or already deleted",
		})
		return
	}

	mainInventoryID, err := EnsureMainInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	var item models.Product_MainInventory
	err = tx.
		Where("main_inventory_id = ? AND product_id = ?",
			mainInventoryID, input.ProductID).
		First(&item).Error

	if err != nil {
		item = models.Product_MainInventory{
			MainInventoryID: mainInventoryID,
			ProductID:       input.ProductID,
			Quantity:        input.Quantity,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := tx.Model(&item).
			Update("quantity", gorm.Expr("quantity + ?", input.Quantity)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Main inventory updated",
	})
}

func UpdateMainInventoryStock(c *gin.Context) {
	productIDParam := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var input struct {
		Quantity int64 `json:"quantity" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product not found or already deleted",
		})
		return
	}

	mainInventoryID, err := EnsureMainInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	var item models.Product_MainInventory
	if err := tx.
		Where("main_inventory_id = ? AND product_id = ?",
			mainInventoryID, productID).
		First(&item).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found in main inventory",
		})
		return
	}

	if err := tx.Model(&item).
		Update("quantity", input.Quantity).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":  "Main inventory stock updated",
		"quantity": input.Quantity,
	})
}

type MainInventoryItemResponse struct {
	ProductID   uint
	ProductName string
	Price       int64
	Quantity    int64
}

func GetMainInventory(c *gin.Context) {
	result := []MainInventoryItemResponse{}

	err := config.DB.
		Table("product_main_inventories").
		Select(`
			products.id as product_id,
			products.name as product_name,
			products.price,
			product_main_inventories.quantity
		`).
		Joins("JOIN products ON products.id = product_main_inventories.product_id").
		Where("product_main_inventories.deleted_at IS NULL").
		Where("products.deleted_at IS NULL").
		Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result,
	})
}

func DeleteProductFromMainInventory(c *gin.Context) {
	productIDParam := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	mainInventoryID, err := EnsureMainInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	if err := tx.
		Where("main_inventory_id = ? AND product_id = ?",
			mainInventoryID, productID).
		Delete(&models.Product_MainInventory{}).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Product removed from main inventory",
	})
}
