package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMitraSummaryByID(c *gin.Context) {
    mitraID := c.Param("id")
    db := config.DB

    // Menentukan rentang waktu hari ini (00:00:00 sampai sekarang)
    now := time.Now()
    beginningOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

    type RevenueResult struct {
        Name           string `json:"name"`
        ProductRevenue int64  `json:"product_revenue"`
        ServiceRevenue int64  `json:"service_revenue"`
        TotalRevenue   int64  `json:"total_revenue"`
        Profit         int64  `json:"profit"`
    }

    var result RevenueResult

    // 1. Ambil Nama Mitra
    var mitra models.Mitra
    if err := db.Select("name").First(&mitra, mitraID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Mitra tidak ditemukan"})
        return
    }
    result.Name = mitra.Name

    // 2. Hitung Product Revenue via Outlets (Hanya Hari Ini)
    var prodRev int64
    err := db.Model(&models.ProductSales{}).
        Joins("JOIN outlets ON outlets.id = product_sales.outlet_id").
        Where("outlets.mitra_id = ? AND product_sales.created_at >= ?", mitraID, beginningOfDay).
        Select("COALESCE(SUM(product_sales.price_at_sale), 0)").
        Scan(&prodRev).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung product revenue"})
        return
    }

    // 3. Hitung Service Revenue via Barbers (Hanya Hari Ini)
    var servRev int64
    err = db.Model(&models.Service{}).
        Joins("JOIN barbers ON barbers.id = services.barber_id").
        Where("barbers.mitra_id = ? AND services.created_at >= ?", mitraID, beginningOfDay).
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
        "status":   "success",
        "mitra_id": mitraID,
        "period":   "today",
        "data":     result,
    })
}

func GetMitraRevenueChart(c *gin.Context) {
    // 1. Ambil MitraID dari URL
    mitraID := c.Param("id")
    
    // Tentukan rentang 7 hari terakhir
    now := time.Now()
    startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    sevenDaysAgo := startOfToday.AddDate(0, 0, -6)

    // 2. Query Revenue Product (P) spesifik Mitra via Outlets
    var productResults []struct {
        Date  string
        Total int64
    }
    config.DB.Raw(`
        SELECT TO_CHAR(ps.created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date, 
               SUM(ps.price_at_sale) as total
        FROM product_sales ps
        JOIN outlets o ON o.id = ps.outlet_id
        WHERE o.mitra_id = ? AND ps.created_at >= ?
        GROUP BY 1
    `, mitraID, sevenDaysAgo).Scan(&productResults)

    // 3. Query Revenue Service (M) spesifik Mitra via Barbers
    var serviceResults []struct {
        Date  string
        Total int64
    }
    config.DB.Raw(`
        SELECT TO_CHAR(s.created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date, 
               SUM(s.price_at_sale) as total
        FROM services s
        JOIN barbers b ON b.id = s.barber_id
        WHERE b.mitra_id = ? AND s.created_at >= ?
        GROUP BY 1
    `, mitraID, sevenDaysAgo).Scan(&serviceResults)

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
        "status":   "success",
        "mitra_id": mitraID,
        "data":     finalData,
    })
}

func GetMitraDistributions(c *gin.Context) {
    // 1. Ambil ID mitra dari parameter URL
    mitraID := c.Param("id")
    
    // Tentukan rentang waktu hari ini
    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    endOfDay := startOfDay.Add(24 * time.Hour)

    // 2. Total Produk Terjual (Per Produk) - Khusus Mitra Ini
    var productDist []DistributionData
    config.DB.Raw(`
        SELECT p.name, SUM(pps.quantity) as value
        FROM product_product_sales pps
        JOIN products p ON p.id = pps.product_id
        JOIN product_sales ps ON ps.id = pps.product_sales_id
        JOIN outlets o ON o.id = ps.outlet_id
        WHERE o.mitra_id = ? 
          AND ps.created_at >= ? AND ps.created_at < ?
          AND ps.deleted_at IS NULL
          AND pps.deleted_at IS NULL
        GROUP BY p.name
        ORDER BY value DESC
    `, mitraID, startOfDay, endOfDay).Scan(&productDist)

    // 3. Total Service Dilakukan (Per Tipe Service) - Khusus Mitra Ini
    var serviceDist []DistributionData
    config.DB.Raw(`
        SELECT st.name, COUNT(sst.id) as value
        FROM service_service_types sst
        JOIN service_types st ON st.id = sst.service_type_id
        JOIN services s ON s.id = sst.service_id
        JOIN barbers b ON b.id = s.barber_id
        WHERE b.mitra_id = ? 
          AND s.created_at >= ? AND s.created_at < ?
          AND s.deleted_at IS NULL
          AND sst.deleted_at IS NULL
        GROUP BY st.name
        ORDER BY value DESC
    `, mitraID, startOfDay, endOfDay).Scan(&serviceDist)

    c.JSON(http.StatusOK, gin.H{
        "status":   "success",
        "mitra_id": mitraID,
        "data": gin.H{
            "products": productDist,
            "services": serviceDist,
        },
    })
}

func GetMitraOutlets(c *gin.Context) {
    // 1. Ambil ID mitra dari parameter URL /:id
    mitraID := c.Param("id")
    db := config.DB

    // 2. Struct untuk menampung data yang spesifik diminta
    type OutletResponse struct {
        Address     string    `json:"address"`
        PhoneNumber string    `json:"phone_number"`
        CreatedAt   time.Time `json:"created_at"`
		ID          uint      `json:"id"`
    }

    var outlets []OutletResponse

    // 3. Query ke tabel outlets milik mitra tersebut
    // Menggunakan Model models.Outlet untuk memastikan table name benar sesuai GORM
    err := db.Model(&models.Outlet{}).
        Select("address, phone_number, created_at, id").
        Where("mitra_id = ?", mitraID).
        Find(&outlets).Error

    // 4. Handle jika terjadi error database
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Gagal mengambil data outlet",
        })
        return
    }

    // 5. Response sukses
    c.JSON(http.StatusOK, gin.H{
        "status":   "success",
        "mitra_id": mitraID,
        "data":     outlets,
    })
}