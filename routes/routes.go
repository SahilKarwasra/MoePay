package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sahilkarwasra/moepay/controllers"
	"github.com/sahilkarwasra/moepay/utils"
)

func RegisterRoutes(r *gin.Engine) {
	// Health Check
	r.GET("/health", func(ctx *gin.Context) {
		utils.Success(ctx, "OK", http.StatusOK, gin.H{"message": "OK"})
	})

	api := r.Group("/api/v1")

	// public auth routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/send-otp", controllers.SendOTP)
		authRoutes.POST("/verify-otp", controllers.VerifyOTP)
	}
}
