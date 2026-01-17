package main

import (
	"log"
	"backend/database" // sesuaikan path
	"backend/config"   // sesuaikan path koneksi DB Anda
    "backend/models"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Koneksi ke Database
	err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file tidak ditemukan, pastikan env vars sudah diset")
    }

	config.ConnectDatabase()
    db := config.DB // Ambil variabel global dari package config

	log.Println("Cleaning database...")
    db.Migrator().DropTable(
        &models.Service_ServiceType{}, &models.Product_ProductSales{},
		&models.Product_OutletInventory{}, &models.Product_MainInventory{},
		&models.ProductSales{}, &models.Service{}, &models.ServiceType{},
		&models.OutletInventory{}, &models.Barber{}, &models.Outlet{},
		&models.Mitra{}, &models.SuperAdmin{}, &models.Product{},
		&models.MainInventory{}, &models.User{},
    )

    db.AutoMigrate(&models.Service_ServiceType{}, &models.Product_ProductSales{},
		&models.Product_OutletInventory{}, &models.Product_MainInventory{},
		&models.ProductSales{}, &models.Service{}, &models.ServiceType{},
		&models.OutletInventory{}, &models.Barber{}, &models.Outlet{},
		&models.Mitra{}, &models.SuperAdmin{}, &models.Product{},
		&models.MainInventory{}, &models.User{},) 

	// 4. Jalankan Seeder
    log.Println("Seeding...")
    err = database.SeedAll(db)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Success!")
}