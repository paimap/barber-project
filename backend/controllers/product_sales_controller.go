package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"
	"strconv"
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateProductSales(c *gin.Context) {
	// Ambil barber
	barberIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}
	barberID := barberIDValue.(uint)

	var barber models.Barber
	if err := config.DB.Where("user_id = ?", barberID).First(&barber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "baber not found"})
	}

	// Bind input
	var input struct {
		PaymentType string `json:"PaymentType" binding:"required,oneof=CASH QRIS"`
		Products    []struct {
			ProductID int64 `json:"ID" binding:"required"`
			Quantity  int64 `json:"Quantity" binding:"required,gt=0"`
		} `json:"Products" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil semua product data
	productIDs := make([]int64, 0, len(input.Products))
	for _, p := range input.Products {
		productIDs = append(productIDs, p.ProductID)
	}

	var products []models.Product
	if err := config.DB.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(products) != len(productIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "one or more products not found"})
		return
	}

	// Map productID -> price
	priceMap := make(map[uint]int64)
	for _, p := range products {
		priceMap[p.ID] = p.Price
	}

	// Hitung total price
	var totalPrice int64
	for _, p := range input.Products {
		totalPrice += priceMap[uint(p.ProductID)] * p.Quantity
	}

	// Mulai transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ambil inventory barber (satu outlet)
	var inventory models.OutletInventory
	if err := tx.Where("id = ?", barber.OutletID).First(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Buat ProductSales
	productSales := models.ProductSales{
		PaymentType: input.PaymentType,
		PriceAtSale: totalPrice,
		OutletID:    barber.OutletID,
	}
	if err := tx.Create(&productSales).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Loop tiap product input
	var product_productSales []models.Product_ProductSales
	for _, p := range input.Products {
		var poi models.Product_OutletInventory

		// Lock row stock product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("outlet_inventory_id = ? AND product_id = ?", inventory.ID, p.ProductID).
			First(&poi).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %d not in inventory", p.ProductID)})
			return
		}

		// Cek stock cukup
		if poi.Quantity < p.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("stock not enough for product %d", p.ProductID)})
			return
		}

		// Kurangi stock
		if err := tx.Model(&poi).Update("quantity", gorm.Expr("quantity - ?", p.Quantity)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Siapkan pivot table
		product_productSales = append(product_productSales, models.Product_ProductSales{
			ProductSalesID: productSales.ID,
			ProductID:      uint(p.ProductID),
			Quantity:       p.Quantity,
		})
	}

	// Simpan pivot table
	if err := tx.Create(&product_productSales).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "product sales created successfully",
		"data":    productSales.ID,
	})
}


type BarberSalesResponse struct {
    SalesID    uint           `json:"ServiceId"`
    PriceAtSale  int64      `json:"PriceAtSale"`
    PaymentType string          `json:"PaymentType"`
    ProductName pq.StringArray `json:"ProductName" gorm:"type:text[]"`
	CreatedAt string `json:"CreatedAt"`
	Quantity int64          `json:"Quantity"`
}
func GetBarberSales(c *gin.Context) {
    userIDValue, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    userID := userIDValue.(uint)

    var barber models.Barber
    if err := config.DB.Where("user_id = ?", userID).First(&barber).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "barber not found"})
        return
    }

    var sales []BarberSalesResponse

    // Implementasi Query SQL ke GORM
    err := config.DB.Table("product_product_sales").
        Select(`
            product_sales.id as sales_id, 
            product_sales.price_at_sale, 
            product_sales.payment_type, 
            array_agg(products.name) as product_name, 
            product_sales.created_at, 
            SUM(product_product_sales.quantity) as quantity
        `).
        Joins("JOIN products ON products.id = product_product_sales.product_id").
        Joins("JOIN product_sales ON product_sales.id = product_product_sales.product_sales_id").
        Where("product_sales.outlet_id = ?", barber.OutletID).
		Where("product_sales.deleted_at IS NULL").
        Where("products.deleted_at IS NULL").
        Where("product_product_sales.deleted_at IS NULL").
        Group("product_sales.id").
        Scan(&sales).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"sales": sales})
}

func UpdateSalesPaymentMethod(c *gin.Context){
	salesIdParam := c.Param("id")
	salesId, err := strconv.ParseUint(salesIdParam, 10, 64)

	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sales id"})
	return
}

	var input struct{
		PaymentType string `json:"PaymentType" binding:"required,oneof=CASH QRIS"`
	}

	if err := c.ShouldBindJSON(&input);  err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map [string]interface{}{}

	if input.PaymentType != ""{
		updates["payment_type"] = input.PaymentType
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}


	var sales models.ProductSales
	if err := config.DB.First(&sales, salesId).Error; err!=nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&sales).Updates(updates).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to update sales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sales updated successfully"})
}