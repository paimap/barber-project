package controllers

import (
	"backend/config"
	"backend/models"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

/* ============================================================
   STRUCT & TYPES
============================================================ */

type ChatController struct {
	Client *genai.Client
}

type ContextProvider func() string

/* ============================================================
   INIT
============================================================ */

func NewChatController() *ChatController {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("GEMINI_API_KEY tidak ditemukan")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	return &ChatController{Client: client}
}

/* ============================================================
   PROMPT BUILDER (CENTRALIZED)
============================================================ */

func buildSystemPrompt(contextType string, data string) string {
	return fmt.Sprintf( 
		`PERAN: Anda adalah AI Business Analyst internal untuk sistem manajemen barbershop. 

		KONTEKS HALAMAN AKTIF: %s 

		DATA TERSEDIA: %s 

		ATURAN WAJIB: 
		- Gunakan HANYA data yang diberikan di atas. 
		- Jangan mengarang angka, asumsi, atau informasi tambahan. 
		- Jika data tidak tersedia atau tidak relevan dengan pertanyaan, katakan dengan jujur. 
		- Jangan membahas topik di luar bisnis barbershop. 
		
		FORMAT JAWABAN: 
		- Ringkas, jelas, dan profesional. 
		- Gunakan bullet point jika lebih dari satu poin. 
		- Berikan insight bisnis dan operasional pada setiap jawaban
		- Hindari penjelasan teknis sistem atau database. 
		
		TUJUAN: Membantu Superadmin / Owner mengambil keputusan berbasis data.` , contextType, data) }

/* ============================================================ 
   CONTEXT PROVIDERS 
   ============================================================ */

func (cc *ChatController) getDashboardContext() string {
    now := time.Now()
    // Pastikan lokasi Asia/Jakarta agar sinkron dengan SQL
    loc, _ := time.LoadLocation("Asia/Jakarta") 
    startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
    sevenDaysAgo := startOfToday.AddDate(0, 0, -6)
    endOfToday := startOfToday.Add(24 * time.Hour)

    /* ---------- HARI INI (Logika Perkalian pps.quantity) ---------- */
    var productRev, serviceRev int64
    config.DB.Raw(`
        SELECT COALESCE(SUM(pps.quantity * ps.price_at_sale),0) 
        FROM product_product_sales pps 
        JOIN product_sales ps ON ps.id = pps.product_sales_id 
        WHERE ps.created_at >= ? AND ps.created_at < ?`, 
        startOfToday, endOfToday).Scan(&productRev)

    config.DB.Raw(`
        SELECT COALESCE(SUM(price_at_sale),0) 
        FROM services 
        WHERE created_at >= ? AND created_at < ?`, 
        startOfToday, endOfToday).Scan(&serviceRev)

    /* ---------- TREN 7 HARI (Sinkron dengan getRevenueChart) ---------- */
    type DailyResult struct {
        DateKey string // YYYY-MM-DD untuk mapping
        DayName string // Mon, Tue, dll
        Total   int64
    }
    
    var productTrend, serviceTrend []DailyResult
    
    // Query Product dengan logic Price * Quantity + Timezone
    config.DB.Raw(`
        SELECT 
            TO_CHAR(ps.created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date_key,
            TO_CHAR(ps.created_at AT TIME ZONE 'Asia/Jakarta', 'Dy') as day_name,
            SUM(pps.quantity * ps.price_at_sale) as total
        FROM product_sales ps
        JOIN product_product_sales pps ON ps.id = pps.product_sales_id
        WHERE ps.created_at >= ? 
        GROUP BY 1, 2 ORDER BY 1 ASC`, sevenDaysAgo).Scan(&productTrend)

    // Query Service dengan Timezone
    config.DB.Raw(`
        SELECT 
            TO_CHAR(created_at AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM-DD') as date_key,
            TO_CHAR(created_at AT TIME ZONE 'Asia/Jakarta', 'Dy') as day_name,
            SUM(price_at_sale) as total
        FROM services 
        WHERE created_at >= ? 
        GROUP BY 1, 2 ORDER BY 1 ASC`, sevenDaysAgo).Scan(&serviceTrend)

    // Mapping menggunakan DateKey (YYYY-MM-DD) agar tidak tertukar antar minggu
    mapP := map[string]int64{}
    for _, v := range productTrend { mapP[v.DateKey] = v.Total }
    mapS := map[string]int64{}
    for _, v := range serviceTrend { mapS[v.DateKey] = v.Total }

    // Susun Tren secara dinamis (6 hari lalu sampai hari ini)
    trendText := ""
    for i := 0; i < 7; i++ {
        t := sevenDaysAgo.AddDate(0, 0, i)
        dKey := t.Format("2006-01-02")
        dName := t.Format("Mon")
        trendText += fmt.Sprintf("%s(P:%d,S:%d) ", dName, mapP[dKey], mapS[dKey])
    }

    // ... (Bagian Ranking & Produk Terlaris tetap sama atau sesuaikan startOfToday) ...
    
    totalRevenue := productRev + serviceRev
    estimatedProfit := float64(totalRevenue) * 0.2

    return fmt.Sprintf(
        "=== DASHBOARD SUMMARY (%s) ===\n"+
            "[HARI INI]\n- Total Revenue: Rp%d\n- Estimasi Profit: Rp%.0f\n- Produk: Rp%d\n- Servis: Rp%d\n\n"+
            "[TREN 7 HARI]\n%s",
        now.In(loc).Format("02 Jan 2006"), totalRevenue, estimatedProfit, productRev, serviceRev, trendText,
    )
}

func (cc *ChatController) getMitraContext() string {
	now := time.Now()
	loc := now.Location()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalMitras, totalOutlets int64
	var productRev, serviceRev int64

	// 1. Hitung Total Mitra & Outlet (Kumulatif)
	config.DB.Model(&models.Mitra{}).Count(&totalMitras)
	config.DB.Model(&models.Outlet{}).Count(&totalOutlets)

	// 2. Revenue Product Hari Ini
	config.DB.Raw(`
		SELECT COALESCE(SUM(pps.quantity * ps.price_at_sale), 0) 
		FROM product_product_sales pps 
		JOIN product_sales ps ON ps.id = pps.product_sales_id 
		WHERE ps.created_at >= ? AND ps.created_at < ? AND ps.deleted_at IS NULL`, 
		startOfDay, endOfDay).Scan(&productRev)

	// 3. Revenue Service Hari Ini
	config.DB.Raw(`
		SELECT COALESCE(SUM(price_at_sale), 0) 
		FROM services 
		WHERE created_at >= ? AND created_at < ? AND deleted_at IS NULL`, 
		startOfDay, endOfDay).Scan(&serviceRev)

	todayRevenue := productRev + serviceRev

	/* ---------- FINAL CONTEXT ---------- */
	return fmt.Sprintf(
		"=== MITRA & OUTLET SUMMARY (%s) ===\n"+
			"[DATA KUMULATIF]\n- Total Mitra Aktif: %d Mitra\n- Total Cabang/Outlet: %d Lokasi\n\n"+
			"[PERFORMA HARI INI]\n- Total Pendapatan Gabungan: Rp%d\n- Kontribusi Produk: Rp%d\n- Kontribusi Layanan/Service: Rp%d\n\n"+
			"Catatan: Data ini mencakup seluruh ekosistem mitra dan cabang yang terdaftar.",
		now.Format("02 Jan 2006"), totalMitras, totalOutlets, todayRevenue, productRev, serviceRev,
	)
}

func (cc *ChatController) getServiceContext() string {
	now := time.Now()
	loc := now.Location()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var metrics struct {
		ProductRevenue int64
		ProductSold    int64
		ServiceRevenue int64
		ServiceCount   int64
	}

	// 1. Ambil data Produk
	config.DB.Raw(`
		SELECT COALESCE(SUM(ps.price_at_sale * pps.quantity), 0) as product_revenue, 
		COALESCE(SUM(pps.quantity), 0) as product_sold 
		FROM product_sales ps 
		LEFT JOIN product_product_sales pps ON pps.product_sales_id = ps.id 
		WHERE ps.created_at >= ? AND ps.created_at < ? AND ps.deleted_at IS NULL`, 
		startOfDay, endOfDay).Scan(&metrics)

	// 2. Ambil data Service
	var serviceRes struct {
		Rev   int64
		Count int64
	}
	config.DB.Raw(`
		SELECT COALESCE(SUM(price_at_sale), 0), COUNT(id) 
		FROM services 
		WHERE created_at >= ? AND created_at < ? AND deleted_at IS NULL`, 
		startOfDay, endOfDay).Scan(&serviceRes)

	metrics.ServiceRevenue = serviceRes.Rev
	metrics.ServiceCount = serviceRes.Count

	/* ---------- FINAL CONTEXT ---------- */
	return fmt.Sprintf(
		"=== SERVICE & PRODUCT SUMMARY (%s) ===\n"+
			"[METRIK PRODUK]\n- Total Penjualan: Rp%d\n- Item Terjual: %d pcs\n\n"+
			"[METRIK LAYANAN/SERVICE]\n- Total Pendapatan Jasa: Rp%d\n- Layanan Dilakukan: %d tindakan\n\n"+
			"[ANALISIS SINGKAT]\n- Rasio Jasa vs Produk: %d transaksi layanan banding %d produk terjual.\n- Total Aktivitas Operasional: %d (Layanan + Produk).",
		now.Format("02 Jan 2006"), metrics.ProductRevenue, metrics.ProductSold, metrics.ServiceRevenue, metrics.ServiceCount, metrics.ServiceCount, metrics.ProductSold, metrics.ServiceCount+metrics.ProductSold,
	)
}

func (cc *ChatController) getInventoryContext() string {
	type InventoryItem struct {
		ProductName string
		Price       int64
		Quantity    int64
	}
	var items []InventoryItem

	err := config.DB.
		Table("product_main_inventories").
		Select("products.name as product_name, products.price, product_main_inventories.quantity").
		Joins("JOIN products ON products.id = product_main_inventories.product_id").
		Where("product_main_inventories.deleted_at IS NULL").
		Where("products.deleted_at IS NULL").
		Scan(&items).Error

	if err != nil {
		return "Gagal mengambil data inventaris."
	}

	var totalItems int64
	lowStockItems := ""
	inventoryValue := int64(0)

	for _, item := range items {
		totalItems += item.Quantity
		inventoryValue += (item.Quantity * item.Price)
		if item.Quantity < 10 {
			lowStockItems += fmt.Sprintf("- %s (Sisa: %d)\n", item.ProductName, item.Quantity)
		}
	}

	if lowStockItems == "" {
		lowStockItems = "Semua stok dalam kondisi aman."
	}

	/* ---------- FINAL CONTEXT ---------- */
	return fmt.Sprintf(
		"=== INVENTORY/STOCK SUMMARY ===\n"+
			"[IKHTISAR STOK]\n- Total Unit di Gudang: %d item\n- Estimasi Nilai Aset Stok: Rp%d\n\n"+
			"[PERINGATAN STOK MENIPIS]\n%s\n"+
			"[DETAIL SEMUA PRODUK]\n%v",
		totalItems, inventoryValue, lowStockItems, items,
	)
}

func (cc *ChatController) getMitraDetailContext(mitraID uint) string {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	sevenDaysAgo := startOfDay.AddDate(0, 0, -6)

	// 1. NAMA MITRA & REVENUE HARI INI
	var mitraName string
	config.DB.Model(&models.Mitra{}).Select("name").Where("id = ?", mitraID).Scan(&mitraName)

	var prodRev, servRev int64
	config.DB.Model(&models.ProductSales{}).
		Joins("JOIN outlets ON outlets.id = product_sales.outlet_id").
		Where("outlets.mitra_id = ? AND product_sales.created_at >= ? AND product_sales.created_at < ?", mitraID, startOfDay, endOfDay).
		Select("COALESCE(SUM(product_sales.price_at_sale), 0)").Scan(&prodRev)

	config.DB.Model(&models.Service{}).
		Joins("JOIN barbers ON barbers.id = services.barber_id").
		Where("barbers.mitra_id = ? AND services.created_at >= ? AND services.created_at < ?", mitraID, startOfDay, endOfDay).
		Select("COALESCE(SUM(services.price_at_sale), 0)").Scan(&servRev)

	// 2. TREN 7 HARI
	var pTrend []struct {
		Date  string
		Total int64
	}
	config.DB.Raw(`
		SELECT TO_CHAR(ps.created_at, 'Dy') as date, SUM(ps.price_at_sale) as total 
		FROM product_sales ps 
		JOIN outlets o ON o.id = ps.outlet_id 
		WHERE o.mitra_id = ? AND ps.created_at >= ? 
		GROUP BY 1`, 
		mitraID, sevenDaysAgo).Scan(&pTrend)

	trendMap := "Tren 7 Hari: "
	for _, v := range pTrend {
		trendMap += fmt.Sprintf("[%s: P:Rp%d] ", v.Date, v.Total)
	}

	// 3. DISTRIBUSI & OUTLET
	var topProd string
	config.DB.Raw(`
		SELECT p.name 
		FROM product_product_sales pps 
		JOIN products p ON p.id = pps.product_id 
		JOIN product_sales ps ON ps.id = pps.product_sales_id 
		JOIN outlets o ON o.id = ps.outlet_id 
		WHERE o.mitra_id = ? AND ps.created_at >= ? 
		GROUP BY p.name ORDER BY SUM(pps.quantity) DESC LIMIT 1`, 
		mitraID, startOfDay).Scan(&topProd)

	var outletCount int64
	config.DB.Model(&models.Outlet{}).Where("mitra_id = ?", mitraID).Count(&outletCount)

	return fmt.Sprintf(
		"=== DETAIL MITRA: %s ===\n"+
			"- Status Hari Ini: Revenue Produk Rp%d, Servis Rp%d (Profit 20%%: Rp%.0f)\n"+
			"- Skala: Memiliki %d Outlet aktif.\n"+
			"- Produk Terlaris Hari Ini: %s\n- %s",
		mitraName, prodRev, servRev, float64(prodRev+servRev)*0.2, outletCount, topProd, trendMap,
	)
}

func (cc *ChatController) getOutletDetailContext(outletID uint) string {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	sevenDaysAgo := startOfDay.AddDate(0, 0, -6)

	var outlet models.Outlet
	config.DB.Select("address").First(&outlet, outletID)

	var prodRev, servRev int64
	config.DB.Model(&models.ProductSales{}).
		Where("outlet_id = ? AND created_at >= ? AND created_at < ?", outletID, startOfDay, endOfDay).
		Select("COALESCE(SUM(price_at_sale), 0)").Scan(&prodRev)

	config.DB.Model(&models.Service{}).
		Joins("JOIN barbers ON barbers.id = services.barber_id").
		Where("barbers.outlet_id = ? AND services.created_at >= ? AND services.created_at < ?", outletID, startOfDay, endOfDay).
		Select("COALESCE(SUM(services.price_at_sale), 0)").Scan(&servRev)

	type LogTx struct {
		Type        string
		Price       int64
		Time        time.Time
		PerformedBy string
	}
	var todayTxs []LogTx
	config.DB.Raw(`
		SELECT * FROM (
			SELECT 'Product' as type, price_at_sale as price, created_at, 'Staff Toko' as performed_by 
			FROM product_sales 
			WHERE outlet_id = ? AND created_at >= ? AND created_at < ?
			UNION ALL
			SELECT 'Service' as type, s.price_at_sale as price, s.created_at, b.name as performed_by 
			FROM services s 
			JOIN barbers b ON b.id = s.barber_id 
			WHERE b.outlet_id = ? AND s.created_at >= ? AND s.created_at < ?
		) AS daily_tx ORDER BY created_at ASC`, 
		outletID, startOfDay, endOfDay, outletID, startOfDay, endOfDay).Scan(&todayTxs)

	txText := ""
	for _, t := range todayTxs {
		txText += fmt.Sprintf("- [%s] %s: Rp%d (Oleh: %s)\n", t.Time.Format("15:04"), t.Type, t.Price, t.PerformedBy)
	}
	if txText == "" {
		txText = "Belum ada transaksi hari ini."
	}

	var pTrend []struct {
		Date  string
		Total int64
	}
	config.DB.Raw(`
		SELECT TO_CHAR(created_at, 'Dy') as date, SUM(price_at_sale) as total 
		FROM product_sales 
		WHERE outlet_id = ? AND created_at >= ? GROUP BY 1`, 
		outletID, sevenDaysAgo).Scan(&pTrend)

	trendMap := "Tren: "
	for _, v := range pTrend {
		trendMap += fmt.Sprintf("%s:Rp%d ", v.Date, v.Total)
	}

	return fmt.Sprintf(
		"=== DETAIL OUTLET: %s ===\n"+
			"[SUMMARY HARI INI]\n- Revenue: Produk Rp%d, Servis Rp%d\n- Total: Rp%d (Profit: Rp%.0f)\n\n"+
			"[LOG TRANSAKSI HARI INI]\n%s\n"+
			"[INFO TAMBAHAN]\n- %s",
		outlet.Address, prodRev, servRev, prodRev+servRev, float64(prodRev+servRev)*0.2, txText, trendMap,
	)
}

/* ============================================================
   HANDLER
============================================================ */

func (cc *ChatController) HandleChat(c *gin.Context) {
	/* 1. AUTH */
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := val.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user_id type"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User tidak ditemukan"})
		return
	}

	now := time.Now()

	/* 2. RESET KUOTA HARIAN */
	if user.LastChatDate.IsZero() || user.LastChatDate.Format("2006-01-02") != now.Format("2006-01-02") {
		user.ChatLimit = 5
		user.LastChatDate = now
		config.DB.Where("user_id = ?", userID).Delete(&models.ChatHistory{})
		config.DB.Save(&user)
	}

	if user.ChatLimit <= 0 {
		c.JSON(429, gin.H{
			"error": "Kuota chat harian sudah habis. Silakan coba lagi besok.",
		})
		return
	}

	/* 3. INPUT VALIDATION */
	var input struct {
		Message     string `json:"message" binding:"required"`
		PageContext string `json:"page_context"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Message == "" {
		c.JSON(400, gin.H{"error": "Pesan tidak boleh kosong"})
		return
	}

	/* 4. CONTEXT RESOLVER */
	var (
		dataText    string
		contextType = "general"
		id          uint
	)

	if n, _ := fmt.Sscanf(input.PageContext, "outlet_%d", &id); n == 1 {
		dataText = cc.getOutletDetailContext(id)
		contextType = "outlet_detail"
	} else if n, _ := fmt.Sscanf(input.PageContext, "mitra_%d", &id); n == 1 {
		dataText = cc.getMitraDetailContext(id)
		contextType = "mitra_detail"
	} else {
		staticContexts := map[string]ContextProvider{
			"dashboard": cc.getDashboardContext,
			"stock":     cc.getInventoryContext,
			"mitra":     cc.getMitraContext,
			"services":  cc.getServiceContext,
		}

		if fn, ok := staticContexts[input.PageContext]; ok {
			dataText = fn()
			contextType = input.PageContext
		} else {
			dataText = "Tidak ada data spesifik pada halaman ini."
			contextType = "unknown"
		}
	}

	/* 5. CHAT HISTORY (LIMIT 5) */
	const MaxHistory = 5
	var historyDB []models.ChatHistory
	config.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(MaxHistory).Find(&historyDB)

	var genaiHistory []*genai.Content
	for i := len(historyDB) - 1; i >= 0; i-- {
		role := historyDB[i].Role
		if role != "user" && role != "model" {
			role = "user"
		}
		genaiHistory = append(genaiHistory, &genai.Content{
			Role: role,
			Parts: []genai.Part{
				genai.Text(historyDB[i].Content),
			},
		})
	}

	/* 6. GEMINI EXECUTION (TIMEOUT 30 DETIK, RETRY 2X, LOG ERROR) */
	model := cc.Client.GenerativeModel("gemini-2.5-flash")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(buildSystemPrompt(contextType, dataText)),
		},
	}

	cs := model.StartChat()
	cs.History = genaiHistory

	ctxTimeout, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var resp *genai.GenerateContentResponse
	var err error

	for i := 0; i < 2; i++ { // retry 2x
		resp, err = cs.SendMessage(ctxTimeout, genai.Text(input.Message))
		if err != nil {
			fmt.Println("SendMessage error:", err)
		} else if resp != nil && len(resp.Candidates) > 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	aiReply := "Maaf, AI sedang sibuk. Silakan coba beberapa saat lagi."
	if resp != nil && len(resp.Candidates) > 0 {
		aiReply = ""
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					aiReply += fmt.Sprintf("%v", part)
				}
			}
		}
		if aiReply == "" {
			aiReply = "Maaf, saya belum bisa memberikan insight dari data ini."
		}
	}

	/* 7. DB TRANSACTION */
	tx := config.DB.Begin()
	if err := tx.Create(&models.ChatHistory{
		UserID:  userID,
		Role:    "user",
		Content: input.Message,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Gagal menyimpan chat user"})
		return
	}

	if err := tx.Create(&models.ChatHistory{
		UserID:  userID,
		Role:    "model",
		Content: aiReply,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Gagal menyimpan chat AI"})
		return
	}

	user.ChatLimit--
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Gagal update kuota"})
		return
	}

	tx.Commit()

	/* 8. RESPONSE */
	c.JSON(200, gin.H{
		"status":           "success",
		"reply":            aiReply,
		"remaining_limit":  user.ChatLimit,
		"context_detected": contextType,
	})
}


/* ============================================================
   GET CHAT HISTORY
============================================================ */

func (cc *ChatController) GetChatHistory(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	history := []models.ChatHistory{}
	err := config.DB.Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startOfDay, endOfDay).
		Order("created_at asc").
		Limit(20).
		Find(&history).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengambil riwayat chat hari ini"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"history": history,
	})
}
