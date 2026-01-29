package controllers

import (
	"backend/config"
	"net/http"
	"time"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func GetDailySummary(c *gin.Context) {
	// 1. Tentukan rentang waktu hari ini (00:00:00 - 23:59:59)
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var productRevenue int64
	var serviceRevenue int64

	// 2. Hitung Revenue Product (Hanya yang transaksi hari ini)
	err := config.DB.Model(&models.Product_ProductSales{}).
		Select("COALESCE(SUM(product_sales.price_at_sale * product_product_sales.quantity), 0)").
		Joins("JOIN product_sales ON product_sales.id = product_product_sales.product_sales_id").
		Where("product_sales.created_at >= ? AND product_sales.created_at < ?", startOfDay, endOfDay).
		Scan(&productRevenue).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	// 3. Hitung Revenue Service (Hanya yang transaksi hari ini)
	err = config.DB.Model(&models.Service{}).
		Select("COALESCE(SUM(price_at_sale), 0)").
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Scan(&serviceRevenue).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data service"})
		return
	}

	// 4. Kalkulasi Total
	totalRevenue := productRevenue + serviceRevenue
	profit := float64(totalRevenue) * 0.2

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"date":   now.Format("2006-01-02"), // Informasi tanggal hari ini
		"data": gin.H{
			"total_revenue":   totalRevenue,
			"profit_20":       int64(profit),
			"product_revenue": productRevenue,
			"service_revenue": serviceRevenue,
		},
	})
}

type ChartData struct {
	Day string `json:"day"`
	P   int64  `json:"p"` // Product Revenue
	M   int64  `json:"m"` // Mitra/Service Revenue
}
func GetRevenueChart(c *gin.Context) {
    // 1. Tentukan rentang 7 hari terakhir (Gunakan Local agar sinkron dengan Asia/Jakarta)
    now := time.Now()
    startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    sevenDaysAgo := startOfToday.AddDate(0, 0, -6)

    // 2. Query untuk Revenue Product (p)
    var productResults []struct {
        Date  string
        Total int64
    }
    // Menggunakan GROUP BY 1 (urutan kolom pertama) agar lebih aman
    config.DB.Raw(`
		SELECT TO_CHAR(ps.created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date, 
			SUM(ps.price_at_sale * pps.quantity) as total
		FROM product_sales ps
		JOIN product_product_sales pps ON ps.id = pps.product_sales_id
		WHERE ps.created_at >= ?
		GROUP BY 1
`	, sevenDaysAgo).Scan(&productResults)

    // 3. Query untuk Revenue Service (m)
    var serviceResults []struct {
        Date  string
        Total int64
    }
    // Perbaikan: Ganti ps.created_at menjadi created_at karena tabelnya services
    config.DB.Raw(`
        SELECT TO_CHAR(created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date, 
               SUM(price_at_sale) as total
        FROM services
        WHERE created_at >= ?
        GROUP BY 1
    `, sevenDaysAgo).Scan(&serviceResults)

    // 4. Mapping
    mapP := make(map[string]int64)
    for _, r := range productResults {
        mapP[r.Date] = r.Total
    }
    mapM := make(map[string]int64)
    for _, r := range serviceResults {
        mapM[r.Date] = r.Total
    }

    // 5. Susun data
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
        "status": "success",
        "data":   finalData,
    })
}

type DistributionData struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}
func GetDistributions(c *gin.Context) {
    // 1. Tentukan rentang waktu hari ini
    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    endOfDay := startOfDay.Add(24 * time.Hour)

    // 2. Total Produk Terjual (Per Produk) - Hari Ini
    var productDist []DistributionData
    config.DB.Raw(`
        SELECT p.name, SUM(pps.quantity) as value
        FROM product_product_sales pps
        JOIN products p ON p.id = pps.product_id
        JOIN product_sales ps ON ps.id = pps.product_sales_id
        WHERE ps.created_at >= ? AND ps.created_at < ?
          AND ps.deleted_at IS NULL
          AND pps.deleted_at IS NULL
        GROUP BY p.name
        ORDER BY value DESC
    `, startOfDay, endOfDay).Scan(&productDist)

    // 3. Total Service Dilakukan (Per Tipe Service) - Hari Ini
    var serviceDist []DistributionData
    config.DB.Raw(`
        SELECT st.name, COUNT(sst.id) as value
        FROM service_service_types sst
        JOIN service_types st ON st.id = sst.service_type_id
        JOIN services s ON s.id = sst.service_id
        WHERE s.created_at >= ? AND s.created_at < ?
          AND s.deleted_at IS NULL
          AND sst.deleted_at IS NULL
        GROUP BY st.name
        ORDER BY value DESC
    `, startOfDay, endOfDay).Scan(&serviceDist)

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data": gin.H{
            "products": productDist,
            "services": serviceDist,
        },
    })
}

type RankingData struct {
	Name    string `json:"name"`
	Revenue int64  `json:"revenue"`
}
func GetRankings(c *gin.Context) {
	// 1. Rentang waktu hari ini
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 2. Top 3 Cabang (Outlet) Berdasarkan Revenue Produk
	var topOutlets []RankingData
	config.DB.Raw(`
		SELECT o.address as name, SUM(ps.price_at_sale) as revenue
		FROM product_sales ps
		JOIN outlets o ON o.id = ps.outlet_id
		WHERE ps.created_at >= ? AND ps.created_at < ?
		  AND ps.deleted_at IS NULL
		GROUP BY o.id, o.address
		ORDER BY revenue DESC
		LIMIT 3
	`, startOfDay, endOfDay).Scan(&topOutlets)

	// 3. Top 3 Mitra Berdasarkan Revenue Service
	var topMitras []RankingData
	config.DB.Raw(`
		SELECT m.name, SUM(s.price_at_sale) as revenue
		FROM services s
		JOIN barbers b ON b.id = s.barber_id
		JOIN mitras m ON m.id = b.mitra_id
		WHERE s.created_at >= ? AND s.created_at < ?
		  AND s.deleted_at IS NULL
		GROUP BY m.id, m.name
		ORDER BY revenue DESC
		LIMIT 3
	`, startOfDay, endOfDay).Scan(&topMitras)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"top_outlets": topOutlets,
			"top_mitras":  topMitras,
		},
	})
}

type MitraSummary struct {
	TotalMitras   int64 `json:"total_mitras"`
	TotalOutlets  int64 `json:"total_outlets"`
	TodayRevenue  int64 `json:"today_revenue"`
}
func GetMitraSummary(c *gin.Context) {
	// 1. Tentukan rentang waktu hari ini
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var summary MitraSummary

	// 2. Hitung Total Mitra (Kumulatif/Semua yang aktif)
	config.DB.Model(&models.Mitra{}).Count(&summary.TotalMitras)

	// 3. Hitung Total Cabang/Outlet (Kumulatif/Semua yang aktif)
	config.DB.Model(&models.Outlet{}).Count(&summary.TotalOutlets)

	// 4. Hitung Total Revenue (Gabungan Product + Service) Khusus Hari Ini
	var productRev int64
	var serviceRev int64

	// Revenue Product Hari Ini
	err := config.DB.Model(&models.Product_ProductSales{}).
		Select("COALESCE(SUM(product_sales.price_at_sale * product_product_sales.quantity), 0)").
		Joins("JOIN product_sales ON product_sales.id = product_product_sales.product_sales_id").
		Where("product_sales.created_at >= ? AND product_sales.created_at < ?", startOfDay, endOfDay).
		Scan(&productRev).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	// Revenue Service Hari Ini
	err = config.DB.Model(&models.Service{}).
		Select("COALESCE(SUM(price_at_sale), 0)").
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Scan(&serviceRev).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data service"})
		return
	}

	summary.TodayRevenue = productRev + serviceRev

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   summary,
	})
}

type MetricsSummary struct {
	ProductRevenue  int64 `json:"product_revenue"`
	ServiceRevenue  int64 `json:"service_revenue"`
	ProductSold     int64 `json:"product_sold"`      // Jumlah item terjual
	ServiceCount    int64 `json:"service_performed"` // Jumlah tindakan servis
}
func GetProductServiceSummary(c *gin.Context) {
    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    endOfDay := startOfDay.Add(24 * time.Hour)

    var metrics MetricsSummary

    // 1. Ambil data Produk ke struct langsung
    config.DB.Raw(`
        SELECT 
            COALESCE(SUM(ps.price_at_sale * pps.quantity), 0) as product_revenue,
            COALESCE(SUM(pps.quantity), 0) as product_sold
        FROM product_sales ps
        LEFT JOIN product_product_sales pps ON pps.product_sales_id = ps.id
        WHERE ps.created_at >= ? AND ps.created_at < ? 
          AND ps.deleted_at IS NULL
    `, startOfDay, endOfDay).Scan(&metrics)

    // 2. Ambil data Service ke variabel temporary dulu supaya tidak menimpa product_revenue
    var serviceResult struct {
        ServiceRevenue   int64
        ServicePerformed int64
    }

    config.DB.Raw(`
        SELECT 
            COALESCE(SUM(price_at_sale), 0) as service_revenue,
            COUNT(id) as service_performed
        FROM services
        WHERE created_at >= ? AND created_at < ? 
          AND deleted_at IS NULL
    `, startOfDay, endOfDay).Scan(&serviceResult)

    // 3. Masukkan hasil service ke struct utama
    metrics.ServiceRevenue = serviceResult.ServiceRevenue
    metrics.ServiceCount = serviceResult.ServicePerformed

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   metrics,
    })
}