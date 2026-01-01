package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/config"
	"backend/models"
)

func GetAllServiceType(c *gin.Context){
	var serviceType []models.ServiceType
	if err := config.DB.Find(&serviceType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"servicesType" : serviceType})
}

func CreateServiceType(c *gin.Context){
	var input struct{
		Name string `json:"name" binding:"required"`
		Price int64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	serviceType := models.ServiceType{
		Name: input.Name,
		Price: input.Price,
	}

	if err := config.DB.Create(&serviceType).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service Type Created Successfully"})
}

func UpdateServiceType(c *gin.Context){
	serviceTypeIdParam := c.Param("id")
	servicesTypeId, err := strconv.ParseUint(serviceTypeIdParam, 10, 64)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	var input struct{
		Name string `json:"name"`
		Price int64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var servicesType models.ServiceType
	if err := config.DB.First(&servicesType, servicesTypeId).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error" : "Service Type Not Found"})
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

	if err := config.DB.Model(&servicesType).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service Type Updated"})	
}

func DeleteServiceType(c *gin.Context){
	servicesTypeIdParam := c.Param("id")
	servicesTypeId, err := strconv.ParseUint(servicesTypeIdParam, 10, 64)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var serviceType models.ServiceType
	if err := config.DB.First(&serviceType, servicesTypeId).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Service Type Not Found"})
		return
	}

	if err := config.DB.Delete(&models.ServiceType{}, serviceType.ID).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message" : "Service Type Deleted Successfully"})
}