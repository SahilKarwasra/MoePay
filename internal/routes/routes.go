package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sahilkarwasra/moepay/internal/controllers"
	"github.com/sahilkarwasra/moepay/internal/service"
)

// func RegisterRoutes(r *gin.Engine) {
// 	// Health Check
// 	r.GET("/health", func(ctx *gin.Context) {
// 		utils.Success(ctx, "OK", http.StatusOK, gin.H{"message": "OK"})
// 	})

// 	api := r.Group("/api/v1")

// 	// public auth routes
// 	authRoutes := api.Group("/auth")
// 	{
// 		authRoutes.POST("/send-otp", controllers.SendOTP)
// 		authRoutes.POST("/verify-otp", controllers.VerifyOTP)
// 	}
// }

func SetupRouter(
	userService *service.UserService,
) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	authHandler := controllers.NewAuthHandler(userService)

	api := engine.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/send-otp", authHandler.SendOTP)
	}

	return engine
}
