package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sahilkarwasra/moepay/internal/service"
	"github.com/sahilkarwasra/moepay/internal/utils"
)

type SendOtpRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOtpRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Otp         string `json:"otp" binding:"required"`
}

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) SendOTP(c *gin.Context) {

	var req SendOtpRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	ctx := c.Request.Context()

	_, err := h.userService.SendOTP(ctx, req.PhoneNumber)

	if err != nil {
		utils.InternalServerError(c, "Failed to send OTP")
		return
	}

	utils.Success(c, "OTP Sent Successfully", http.StatusOK, gin.H{
		"phone_number": req.PhoneNumber,
		"expires_in":   "10 Minutes",
	})

}
