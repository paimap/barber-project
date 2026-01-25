package routes

import (
	"backend/controllers"
	"backend/middleware"

	"backend/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	protected := r.Group("/api")
	protected.Use(middleware.JwtAuthMiddleware())
	{
		protected.GET("/me", controllers.GetCurrentUser)
		protected.GET("/products", controllers.GetAllProduct)
		protected.GET("/service-type", controllers.GetAllServiceType)

	}

	// routes/routes.go
	mitra := protected.Group("/mitra")
	mitra.Use(middleware.RequireRole(models.RoleMitra))
	{
		mitra.POST("/outlets", controllers.CreateOutlet)
		mitra.GET("/outlets", controllers.GetOutlets)
		mitra.PUT("/outlets/:id", controllers.UpdateOutlet)
		mitra.DELETE("/outlets/:id", controllers.DeleteOutlet)

		mitra.POST("/barber", controllers.CreateBarber)
		mitra.GET("/barber", controllers.GetBarberByMitra)
		mitra.PUT("/barber/:id", controllers.UpdateBarber)
		mitra.DELETE("/barber/:id", controllers.DeleteBarber)
		mitra.POST("/barber/:id/assign", controllers.AssignBarberToOutlet)

		mitra.PUT("/inventory", controllers.UpsertOutletInventoryItem)
		mitra.GET("/inventory", controllers.GetOutletInventory)
		mitra.PUT("/inventory/:product_id", controllers.UpdateOutletInventoryStock)
		mitra.DELETE("/inventory/:product_id", controllers.DeleteProductFromOutletInventory)

		mitra.GET("dashboard/:id/summary", controllers.GetMitraSummaryByID)
		mitra.GET("dashboard/:id/chart", controllers.GetMitraRevenueChart)
		mitra.GET("dashboard/:id/distribution", controllers.GetMitraDistributions)
		mitra.GET("dashboard/:id/leaderboard", controllers.GetMitraLeaderboard)

		mitra.GET("outlet/:id/summary", controllers.GetOutletSummaryByID)
		mitra.GET("outlet/:id/chart", controllers.GetOutletRevenueChart)
		mitra.GET("outlet/:id/distribution", controllers.GetOutletDistributions)
		mitra.GET("outlet/:id/transactions", controllers.GetOutletTransactions)
	}

	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(models.RoleSuperAdmin))
	{
		admin.GET("/mitra", controllers.GetAllMitra)
		admin.POST("/mitra", controllers.CreateMitra)
		admin.PUT("/mitra/:id", controllers.UpdateMitra)
		admin.DELETE("/mitra/:id", controllers.DeleteMitra)
		admin.GET("mitra/summary", controllers.GetMitraSummary)

		admin.POST("/products", controllers.CreateProduct)
		admin.PUT("/products/:id", controllers.UpdateProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)

		admin.POST("/service-type", controllers.CreateServiceType)
		admin.PUT("/service-type/:id", controllers.UpdateServiceType)
		admin.DELETE("/service-type/:id", controllers.DeleteServiceType)

		admin.GET("/product-service/summary", controllers.GetProductServiceSummary)

		admin.PUT("/inventory", controllers.UpsertMainInventoryItem)
		admin.GET("/inventory", controllers.GetMainInventory)
		admin.PUT("/inventory/:product_id", controllers.UpdateMainInventoryStock)
		admin.DELETE("/inventory/:product_id", controllers.DeleteProductFromMainInventory)

		admin.GET("dashboard/summary", controllers.GetDailySummary)
		admin.GET("dashboard/chart", controllers.GetRevenueChart)
		admin.GET("dashboard/distribution", controllers.GetDistributions)
		admin.GET("dashboard/leaderboard", controllers.GetRankings)

		admin.GET("mitra/:id/summary", controllers.GetMitraSummaryByID)
		admin.GET("mitra/:id/chart", controllers.GetMitraRevenueChart)
		admin.GET("mitra/:id/distribution", controllers.GetMitraDistributions)
		admin.GET("mitra/:id/outlets", controllers.GetMitraOutlets)
		admin.POST("mitra/:id/outlets", controllers.CreateOutletByMitraID)
		admin.DELETE("mitra/outlets/:id", controllers.DeleteOutletByAdmin)
		admin.PUT("mitra/outlets/:id", controllers.UpdateOutletByAdmin)
		
		admin.GET("outlet/:id/summary", controllers.GetOutletSummaryByID)
		admin.GET("outlet/:id/chart", controllers.GetOutletRevenueChart)
		admin.GET("outlet/:id/distribution", controllers.GetOutletDistributions)
		admin.GET("outlet/:id/transactions", controllers.GetOutletTransactions)
	}

	barber := protected.Group("/barber")
	barber.Use(middleware.RequireRole(models.RoleBarber))
	{
		barber.GET("/service", controllers.GetBarberService)
		barber.POST("/service", controllers.CreateService)
		barber.PUT("/service/:id", controllers.UpdateServicePaymentMethod)

		barber.GET("/sales", controllers.GetBarberSales)
		barber.POST("/sales", controllers.CreateProductSales)
		barber.PUT("/sales/:id", controllers.UpdateSalesPaymentMethod)

		barber.GET("/service/summary", controllers.GetBarberStatsToday)
	}
}
