package database

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"backend/models"
)

// Helper untuk mendapatkan waktu acak dalam 7 hari terakhir
func randomTimeInLast7Days() time.Time {
	now := time.Now()
	// Kurangi antara 0-6 hari, 0-23 jam, 0-59 menit
	daysAgo := rand.Intn(7)
	hoursAgo := rand.Intn(24)
	minsAgo := rand.Intn(60)

	return now.AddDate(0, 0, -daysAgo).
		Add(-time.Duration(hoursAgo) * time.Hour).
		Add(-time.Duration(minsAgo) * time.Minute)
}

func SeedAll(db *gorm.DB) error {
	rand.Seed(time.Now().UnixNano())

	return db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("--> Memulai Seeding Masif (500+ Transaksi)...")

		hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		pass := string(hashed)

		// 1. SEED MASTER DATA (ADMIN & PRODUCTS)
		superUser := models.User{Email: "super@mail.com", Password: pass, Role: models.RoleSuperAdmin}
		tx.Create(&superUser)
		tx.Create(&models.SuperAdmin{Name: "Owner Pusat", UserID: superUser.ID, PhoneNumber: "08111"})

		products := []models.Product{
			{Name: "Pomade High Shine", Price: 75000},
			{Name: "Matte Clay", Price: 85000},
			{Name: "Beard Oil Deluxe", Price: 60000},
			{Name: "Hair Serum", Price: 120000},
		}
		for i := range products {
			tx.Create(&products[i])
		}

		serviceTypes := []models.ServiceType{
			{Name: "Gentleman Cut", Price: 45000},
			{Name: "Hot Towel Shave", Price: 35000},
			{Name: "Hair Wash & Massage", Price: 25000},
			{Name: "Complete Grooming", Price: 100000},
		}
		for i := range serviceTypes {
			tx.Create(&serviceTypes[i])
		}

		// 2. SEED MITRA, OUTLET, BARBER
		var allBarbers []models.Barber
		var allOutlets []models.Outlet

		for i := 1; i <= 5; i++ { // 5 Mitra
			uMitra := models.User{Email: fmt.Sprintf("mitra%d@mail.com", i), Password: pass, Role: models.RoleMitra}
			tx.Create(&uMitra)
			mitra := models.Mitra{Name: fmt.Sprintf("Mitra %d", i), UserID: uMitra.ID}
			tx.Create(&mitra)

			for j := 1; j <= 2; j++ { // 2 Outlet per Mitra (Total 10 Outlet)
				outlet := models.Outlet{Address: fmt.Sprintf("Lokasi %d-%d", i, j), MitraID: mitra.ID}
				tx.Create(&outlet)
				allOutlets = append(allOutlets, outlet)

				// Tambah stok awal outlet
				outInv := models.OutletInventory{OutletID: outlet.ID}
				tx.Create(&outInv)
				for _, p := range products {
					tx.Create(&models.Product_OutletInventory{ProductID: p.ID, OutletInventoryID: outInv.ID, Quantity: 200})
				}

				for k := 1; k <= 3; k++ { // 3 Barber per Outlet (Total 30 Barber)
					uBarber := models.User{Email: fmt.Sprintf("barber-o%dm%dk%d@mail.com", j, i, k), Password: pass, Role: models.RoleBarber}
					tx.Create(&uBarber)
					barber := models.Barber{Name: fmt.Sprintf("Barber %d", len(allBarbers)+1), UserID: uBarber.ID, MitraID: mitra.ID, OutletID: outlet.ID}
					tx.Create(&barber)
					allBarbers = append(allBarbers, barber)
				}
			}
		}

		// 3. SEEDING TRANSAKSI MASIF (Target 500+)
		fmt.Println("--> Generating 400 Service Transactions...")
		for i := 0; i < 400; i++ {
			randomBarber := allBarbers[rand.Intn(len(allBarbers))]
			randomSType := serviceTypes[rand.Intn(len(serviceTypes))]
			createdDate := randomTimeInLast7Days()

			svc := models.Service{
				PaymentType: func() string { if rand.Intn(2) == 0 { return models.Cash }; return models.QRIS }(),
				PriceAtSale: randomSType.Price,
				BarberID:    randomBarber.ID,
			}
			svc.CreatedAt = createdDate
			tx.Create(&svc)

			tx.Create(&models.Service_ServiceType{
				ServiceID: svc.ID, ServiceTypeID: randomSType.ID,
			})
		}

		fmt.Println("--> Generating 150 Product Sales Transactions...")
		for i := 0; i < 150; i++ {
			randomOutlet := allOutlets[rand.Intn(len(allOutlets))]
			randomProd := products[rand.Intn(len(products))]
			qty := int64(rand.Intn(2) + 1)
			createdDate := randomTimeInLast7Days()

			sale := models.ProductSales{
				PaymentType: func() string { if rand.Intn(2) == 0 { return models.Cash }; return models.QRIS }(),
				PriceAtSale: randomProd.Price * qty,
				OutletID:    randomOutlet.ID,
			}
			sale.CreatedAt = createdDate
			tx.Create(&sale)

			tx.Create(&models.Product_ProductSales{
				ProductSalesID: sale.ID,
				ProductID:      randomProd.ID,
				Quantity:       qty,
			})
		}

		fmt.Printf("--> Berhasil menambahkan %d Transaksi Jasa dan %d Penjualan Produk.\total", 400, 150)
		return nil
	})
}