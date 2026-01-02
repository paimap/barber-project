package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`

	Role string `gorm:"type:varchar(20);not null"`

	SuperAdmin *SuperAdmin `gorm:"foreignKey:UserID"`
	Mitra *Mitra `gorm:"foreignKey:UserID"`
	Barber *Barber `gorm:"foreignKey:UserID"`
}

const (
    RoleSuperAdmin = "SUPERADMIN"
    RoleMitra      = "MITRA"
    RoleBarber     = "BARBER"
)

type SuperAdmin struct {
	gorm.Model
	Name string
	UserID uint `gorm:"uniqueIndex;not null"`
	PhoneNumber string `gorm:"type:varchar(20);index"`
}

type Mitra struct {
	gorm.Model
	Name string
	UserID uint `gorm:"uniqueIndex;not null"`
	PhoneNumber string `gorm:"type:varchar(20);index"`

	Outlets []Outlet
	Barbers []Barber
}

type Barber struct {
	gorm.Model
	Name string
	UserID uint `gorm:"uniqueIndex;not null"`
	PhoneNumber string `gorm:"type:varchar(20);index"`

	MitraID uint `gorm:"not null"`
	Mitra Mitra

	OutletID uint  `gorm:"default:null"`
	Outlet Outlet `gorm:"foreignKey:OutletID;references:ID"`

	Services []Service
}

type MainInventory struct {
	gorm.Model
	
	Products []Product_MainInventory
}

type Product_MainInventory struct {
	gorm.Model

	ProductID uint `gorm:"not null;constraint:OnDelete:CASCADE"`
	Product Product

	MainInventoryID uint `gorm:"not null"`
	MainInventory MainInventory

	Quantity int64
}

type Outlet struct {
	gorm.Model
	Address string
	PhoneNumber string `gorm:"type:varchar(20)"`

	MitraID uint `gorm:"not null"`
	Mitra Mitra

	OutletInventories []OutletInventory
	ProductSales []ProductSales
}

type OutletInventory struct {
	gorm.Model 

	OutletID uint `gorm:"not null"`
	Outlet Outlet
	
	Products []Product_OutletInventory
}

type Product_OutletInventory struct {
	gorm.Model
	Quantity int64

	ProductID uint `gorm:"not null;constraint:OnDelete:CASCADE"`
	Product Product

	OutletInventoryID uint `gorm:"not null"`
	OutletInventory OutletInventory
}

type ProductSales struct{
	gorm.Model
	PaymentType string `gorm:"type:varchar(20)"`
	PriceAtSale int64

	OutletID uint `gorm:"not null"`
	Outlet Outlet

	Products []Product_ProductSales
}

const (
    Cash = "CASH"
    QRIS = "QRIS"
)

type Product_ProductSales struct{
	gorm.Model
	Quantity int64
	
	ProductSalesID uint `gorm:"not null"`
	ProductSales ProductSales

	ProductID uint `gorm:"not null"`
	Product Product
}

type Service struct{
	gorm.Model
	PaymentType string `gorm:"type:varchar(20)"`
	PriceAtSale int64

	BarberID uint `gorm:"not null"`
	Barber Barber

	ServiceTypes []Service_ServiceType
}

type ServiceType struct{
	gorm.Model
	Name string
	Price int64

	Services []Service_ServiceType
}

type Service_ServiceType struct{
	gorm.Model

	ServiceID uint `gorm:"not null"`
	Service Service

	ServiceTypeID uint `gorm:"not null"`
	ServiceType ServiceType
}

type Product struct{
	gorm.Model
	Name string
	Price int64

	MainInventories []Product_MainInventory
	OutletInventories []Product_OutletInventory
	ProductSales []Product_ProductSales
}
