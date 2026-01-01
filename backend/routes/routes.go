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
		protected.GET("/profile", controllers.Profile)
		protected.GET("/products", controllers.GetAllProduct)
		
	}

	// routes/routes.go
	mitra := protected.Group("/mitra")
	mitra.Use(middleware.RequireRole(models.RoleMitra))
	{
		mitra.POST("/outlets", controllers.CreateOutlet)
		mitra.GET("/outlets", controllers.GetOutlets)
		mitra.PUT("/outlets/:id", controllers.UpdateOutlet)
		mitra.DELETE("/outlets/:id", controllers.DeleteOutlet)
	}

	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(models.RoleSuperAdmin))
	{
		admin.GET("/mitra", controllers.GetAllMitra)
		admin.POST("/mitra", controllers.CreateMitra)
		admin.PUT("/mitra/:id", controllers.UpdateMitra)
		admin.DELETE("/mitra/:id", controllers.DeleteMitra)

		admin.POST("/products", controllers.CreateProduct)
		admin.PUT("/products/:id", controllers.UpdateProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)
	}
	
}

