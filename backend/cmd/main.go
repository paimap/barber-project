package main

import (
	"log"

	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDatabase()

	// Auto migrate
	config.DB.AutoMigrate(
		&models.User{},
		&models.SuperAdmin{},
		&models.Mitra{},
		&models.Barber{},

		&models.MainInventory{},
		&models.Product_MainInventory{},

		&models.Outlet{},
		&models.OutletInventory{},
		&models.Product_OutletInventory{},

		&models.Product{},
		&models.ProductSales{},
		&models.Product_ProductSales{},

		&models.Service{},
		&models.ServiceType{},          // ← INI YANG KURANG
		&models.Service_ServiceType{},  // ← junction table
	)

	

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}
