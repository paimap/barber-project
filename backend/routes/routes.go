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
	}

	// routes/routes.go
mitra := protected.Group("api/mitra")
mitra.Use(middleware.RequireRole(models.RoleMitra))
{
	mitra.POST("/outlets", controllers.CreateOutlet)
	mitra.GET("/outlets", controllers.GetOutlets)
	mitra.PUT("/outlets/:id", controllers.UpdateOutlet)
	mitra.DELETE("/outlets/:id", controllers.DeleteOutlet)
}

	
}

