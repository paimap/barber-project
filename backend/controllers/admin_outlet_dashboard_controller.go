package controllers

import (
    "backend/config"
    "backend/models"
    "net/http"
    "time"
	"strconv"
	"math"

    "github.com/gin-gonic/gin"
)

func GetOutletSummaryByID(c *gin.Context) {
    outletID := c.Param("id")
    db := config.DB

    // Menentukan rentang waktu hari ini (00:00:00 sampai sekarang)
    now := time.Now()
    beginningOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

    type RevenueResult struct {
        Address        string `json:"address"`
        ProductRevenue int64  `json:"product_revenue"`
        ServiceRevenue int64  `json:"service_revenue"`
        TotalRevenue   int64  `json:"total_revenue"`
        Profit         int64  `json:"profit"`
    }

    var result RevenueResult

    // 1. Ambil Data Outlet (Alamat/Identitas)
    var outlet models.Outlet
    if err := db.Select("address").First(&outlet, outletID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
        return
    }
    result.Address = outlet.Address

    // 2. Hitung Product Revenue untuk Outlet ini (Hanya Hari Ini)
    var prodRev int64
    err := db.Model(&models.ProductSales{}).
        Where("outlet_id = ? AND created_at >= ?", outletID, beginningOfDay).
        Select("COALESCE(SUM(price_at_sale), 0)").
        Scan(&prodRev).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung product revenue"})
        return
    }

    // 3. Hitung Service Revenue via Barbers di Outlet ini (Hanya Hari Ini)
    var servRev int64
    err = db.Model(&models.Service{}).
        Joins("JOIN barbers ON barbers.id = services.barber_id").
        Where("barbers.outlet_id = ? AND services.created_at >= ?", outletID, beginningOfDay).
        Select("COALESCE(SUM(services.price_at_sale), 0)").
        Scan(&servRev).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung service revenue"})
        return
    }

    // 4. Kalkulasi data akhir
    result.ProductRevenue = prodRev
    result.ServiceRevenue = servRev
    result.TotalRevenue = prodRev + servRev
    result.Profit = int64(float64(result.TotalRevenue) * 0.20)

    c.JSON(http.StatusOK, gin.H{
        "status":    "success",
        "outlet_id": outletID,
        "period":    "today",
        "data":      result,
    })
}

func GetOutletRevenueChart(c *gin.Context) {
    // 1. Ambil OutletID dari URL
    outletID := c.Param("id")
    
    // Tentukan rentang 7 hari terakhir
    now := time.Now()
    startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    sevenDaysAgo := startOfToday.AddDate(0, 0, -6)

    // Struct untuk menampung hasil scan query
    type DailyResult struct {
        Date  string
        Total int64
    }

    // 2. Query Revenue Product spesifik Outlet ini
    var productResults []DailyResult
    config.DB.Raw(`
        SELECT TO_CHAR(ps.created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date, 
               SUM(ps.price_at_sale) as total
        FROM product_sales ps
        WHERE ps.outlet_id = ? AND ps.created_at >= ?
        GROUP BY 1
    `, outletID, sevenDaysAgo).Scan(&productResults)

    // 3. Query Revenue Service spesifik Outlet ini via Barbers
    var serviceResults []DailyResult
    config.DB.Raw(`
        SELECT TO_CHAR(s.created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date, 
               SUM(s.price_at_sale) as total
        FROM services s
        JOIN barbers b ON b.id = s.barber_id
        WHERE b.outlet_id = ? AND s.created_at >= ?
        GROUP BY 1
    `, outletID, sevenDaysAgo).Scan(&serviceResults)

    // 4. Mapping hasil query ke map untuk akses cepat
    mapP := make(map[string]int64)
    for _, r := range productResults {
        mapP[r.Date] = r.Total
    }
    mapM := make(map[string]int64)
    for _, r := range serviceResults {
        mapM[r.Date] = r.Total
    }

    // 5. Susun data final untuk 7 hari terakhir
    type ChartData struct {
        Day string `json:"day"`
        P   int64  `json:"p"` // Product
        M   int64  `json:"m"` // Service (M-nya mungkin dari "Massage" atau "Main service")
    }
    
    var finalData []ChartData
    days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

    for i := 0; i < 7; i++ {
        t := sevenDaysAgo.AddDate(0, 0, i)
        dateStr := t.Format("2006-01-02")
        
        finalData = append(finalData, ChartData{
            Day: days[int(t.Weekday())],
            P:   mapP[dateStr],
            M:   mapM[dateStr],
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    "success",
        "outlet_id": outletID,
        "data":      finalData,
    })
}

func GetOutletDistributions(c *gin.Context) {
    // 1. Ambil ID outlet dari parameter URL
    outletID := c.Param("id")
    
    // Tentukan rentang waktu hari ini
    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    endOfDay := startOfDay.Add(24 * time.Hour)

    type DistributionData struct {
        Name  string `json:"name"`
        Value int64  `json:"value"`
    }

    // 2. Total Produk Terjual (Per Produk) - Khusus Outlet Ini
    var productDist []DistributionData
    config.DB.Raw(`
        SELECT p.name, SUM(pps.quantity) as value
        FROM product_product_sales pps
        JOIN products p ON p.id = pps.product_id
        JOIN product_sales ps ON ps.id = pps.product_sales_id
        WHERE ps.outlet_id = ? 
          AND ps.created_at >= ? AND ps.created_at < ?
          AND ps.deleted_at IS NULL
          AND pps.deleted_at IS NULL
        GROUP BY p.name
        ORDER BY value DESC
    `, outletID, startOfDay, endOfDay).Scan(&productDist)

    // 3. Total Service Dilakukan (Per Tipe Service) - Khusus Outlet Ini
    var serviceDist []DistributionData
    config.DB.Raw(`
        SELECT st.name, COUNT(sst.id) as value
        FROM service_service_types sst
        JOIN service_types st ON st.id = sst.service_type_id
        JOIN services s ON s.id = sst.service_id
        JOIN barbers b ON b.id = s.barber_id
        WHERE b.outlet_id = ? 
          AND s.created_at >= ? AND s.created_at < ?
          AND s.deleted_at IS NULL
          AND sst.deleted_at IS NULL
        GROUP BY st.name
        ORDER BY value DESC
    `, outletID, startOfDay, endOfDay).Scan(&serviceDist)

    c.JSON(http.StatusOK, gin.H{
        "status":    "success",
        "outlet_id": outletID,
        "data": gin.H{
            "products": productDist,
            "services": serviceDist,
        },
    })
}

func GetOutletTransactions(c *gin.Context) {
	outletID := c.Param("id")
	db := config.DB

	// 1. Ambil parameter page & limit dari Query String (Contoh: ?page=1&limit=10)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	// Validasi dasar agar tidak error
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	offset := (page - 1) * limit

	// 2. Hitung Total Seluruh Data (untuk menghitung Total Pages)
	var totalRows int64
	countQuery := `
		SELECT (
			(SELECT COUNT(*) FROM product_sales WHERE outlet_id = ? AND deleted_at IS NULL) +
			(SELECT COUNT(*) FROM services s JOIN barbers b ON b.id = s.barber_id 
			 WHERE b.outlet_id = ? AND s.deleted_at IS NULL)
		) as total_count
	`
	db.Raw(countQuery, outletID, outletID).Scan(&totalRows)

	// 3. Kalkulasi Metadata Pagination
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	// 4. Struct untuk menampung hasil Union
	type TransactionResponse struct {
		ID          uint      `json:"id"`
		Type        string    `json:"type"`         // "Product" atau "Service"
		PaymentType string    `json:"payment_type"` // "CASH" atau "QRIS"
		PriceAtSale int64     `json:"price_at_sale"`
		CreatedAt   time.Time `json:"created_at"`
		PerformedBy string    `json:"performed_by"` // Nama Barber (null jika produk)
	}

	var transactions []TransactionResponse

	// 5. Query Utama dengan Union, Order, dan Limit Offset
	dataQuery := `
		SELECT * FROM (
			SELECT ps.id, 'Product' as type, ps.payment_type, ps.price_at_sale, ps.created_at, NULL as performed_by
			FROM product_sales ps 
			WHERE ps.outlet_id = ? AND ps.deleted_at IS NULL
			
			UNION ALL
			
			SELECT s.id, 'Service' as type, s.payment_type, s.price_at_sale, s.created_at, b.name as performed_by
			FROM services s 
			JOIN barbers b ON b.id = s.barber_id
			WHERE b.outlet_id = ? AND s.deleted_at IS NULL
		) AS combined_data
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := db.Raw(dataQuery, outletID, outletID, limit, offset).Scan(&transactions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}

	// 6. Response JSON ke Frontend
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactions,
		"meta": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_rows":   totalRows,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_prev":     page > 1,
		},
	})
}

func GetBarberServiceByID(c *gin.Context) {
    // 1. Ambil ID Barber dari parameter URL
    barberID := c.Param("id")

    // 2. Cek apakah barber tersebut ada (Opsional, untuk validasi)
    var barber models.Barber
    if err := config.DB.First(&barber, barberID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "barber not found"})
        return
    }

    var services []BarberServiceResponse

    // 3. Query menggunakan Join untuk mengambil data service beserta tipe servisnya (array_agg)
    err := config.DB.Table("services").
        Select(`
            services.id AS service_id,
            services.payment_type,
            services.price_at_sale,
            TO_CHAR(services.created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD HH24:MI:SS') AS created_at,
            array_agg(service_types.name) AS service_types 
        `).
        Joins("JOIN service_service_types sst ON sst.service_id = services.id").
        Joins("JOIN service_types ON service_types.id = sst.service_type_id").
        Where("services.barber_id = ?", barber.ID).
        Where("services.deleted_at IS NULL").
        Where("service_types.deleted_at IS NULL").
        Where("sst.deleted_at IS NULL").
        Group("services.id, services.payment_type, services.price_at_sale, services.created_at").
        Order("services.created_at DESC"). // Urutkan dari yang terbaru
        Scan(&services).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat servis barber"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":   "success",
        "barber_id": barber.ID,
        "services":  services,
    })
}

func GetBarberStatsTodayByID(c *gin.Context) {
    db := config.DB
    // 1. Ambil ID Barber dari parameter URL
    barberID := c.Param("id")

    // 2. Tentukan rentang waktu "Hari Ini" (00:00:00 s/d sekarang)
    now := time.Now()
    todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

    var stats BarberStatsResponse

    // 3. Query Agregasi: Hitung total servis dan revenue hari ini
    err := db.Model(&models.Service{}).
        Select(`
            services.barber_id, 
            barbers.name as barber_name, 
            COUNT(services.id) as total_service, 
            COALESCE(SUM(services.price_at_sale), 0) as total_revenue
        `).
        Joins("JOIN barbers ON barbers.id = services.barber_id").
        Where("services.barber_id = ? AND services.created_at >= ?", barberID, todayStart).
        Group("services.barber_id, barbers.name").
        Scan(&stats).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data statistik barber"})
        return
    }

    // 4. Penanganan jika barber belum ada servis sama sekali hari ini
    if stats.BarberID == 0 {
        var b models.Barber
        if err := db.Select("id, name").First(&b, barberID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Barber tidak ditemukan"})
            return
        }
        stats.BarberID = b.ID
        stats.BarberName = b.Name
        stats.TotalService = 0
        stats.TotalRevenue = 0
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "period": "today",
        "data":   stats,
    })
}