package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lib/pq"
)


func CreateService(c *gin.Context){
	var input struct{
		ServiceTypeID []uint `json:"ServiceTypeID" binding:"required"`
		PaymentType string `json:"PaymentType" binding:"required,oneof=CASH QRIS"`
	}
	if err:= c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	barberIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	barberID := barberIDValue.(uint)

	var barber models.Barber
	if err := config.DB.Where("user_id = ?", barberID).First(&barber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "baber not found"})
	}

	var servicesType []models.ServiceType
	if err := config.DB.Where("id IN ?", input.ServiceTypeID).Find(&servicesType).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error" : err.Error()})
		return
	}

	if len(servicesType) != len(input.ServiceTypeID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "one or more service types not found",
		})
		return
	}

	var totalPrice int64 = 0
	for _,service := range servicesType {
		totalPrice += service.Price
	}

	tx := config.DB.Begin()

	service := models.Service{
		PaymentType: input.PaymentType,
		PriceAtSale: totalPrice,
		BarberID: barber.ID,
	}
	
	if err := tx.Create(&service).Error; err != nil{
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	var links []models.Service_ServiceType
	for _,serviceType := range servicesType {
		links = append(links, models.Service_ServiceType{
			ServiceID: service.ID,
			ServiceTypeID: serviceType.ID,
		}) 
	}

	if err := tx.Create(&links).Error; err != nil{
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message" : "service created successfully"})
}

type BarberServiceResponse struct {
    ServiceID    uint           `json:"ServiceId"`
    PaymentType  string         `json:"PaymentType"`
    PriceAtSale  int64          `json:"PriceAtSale"`
    ServiceTypes pq.StringArray `json:"ServiceTypes" gorm:"type:text[]"`
	CreatedAt string `json:"CreatedAt"`
}
func GetBarberService(c *gin.Context) {
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

	var services []BarberServiceResponse

    err := config.DB.Table("services").
        Select(`
            services.id AS service_id,
            services.payment_type,
            services.price_at_sale,
			services.created_at,
            array_agg(service_types.name) AS service_types 
        `).
        Joins("JOIN service_service_types sst ON sst.service_id = services.id").
        Joins("JOIN service_types ON service_types.id = sst.service_type_id").
        Where("services.barber_id = ?", barber.ID).
        Where("services.deleted_at IS NULL").
        Where("service_types.deleted_at IS NULL").
        Where("sst.deleted_at IS NULL").
        Group("services.id, services.payment_type, services.price_at_sale"). 
        Scan(&services).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}


func UpdateServicePaymentMethod(c *gin.Context){
	serviceIdParam := c.Param("id")
	serviceId, err := strconv.ParseUint(serviceIdParam, 10, 64)

	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
	return
}

	var input struct{
		PaymentType string `json:"PaymentType" binding:"required"`
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


	var service models.Service
	if err := config.DB.First(&service, serviceId).Error; err!=nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&service).Updates(updates).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service updated successfully"})
}