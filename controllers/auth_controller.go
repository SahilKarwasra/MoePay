package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sahilkarwasra/moepay/config"
	"github.com/sahilkarwasra/moepay/models"
	"github.com/sahilkarwasra/moepay/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SendOtpRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOtpRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Otp         string `json:"otp" binding:"required"`
}

func SendOTP(ctx *gin.Context) {
	var req SendOtpRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	// Baisc Validations
	if len(req.PhoneNumber) < 8 {
		utils.BadRequest(ctx, "Invalid phone number")
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// delete old otp of this number
	otpCol := config.MongoDB.Collection(config.OtpCollections)
	otpCol.DeleteMany(c, bson.M{"phone_number": req.PhoneNumber})

	// generate otp
	otp := utils.GenerateOTP()
	otpDoc := models.Otp{
		ID:          primitive.NewObjectID(),
		PhoneNumber: req.PhoneNumber,
		Otp:         otp,
		IsUsed:      false,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	_, err := otpCol.InsertOne(c, otpDoc)
	if err != nil {
		utils.InternalServerError(ctx, "Failed to generate OTP")
		return
	}

	utils.LogOTP(req.PhoneNumber, otp)

	utils.Success(ctx, "OTP Sent Successfully", http.StatusOK, gin.H{
		"phone_number": req.PhoneNumber,
		"expires_in":   "10 Minutes",
	})

}

func VerifyOTP(ctx *gin.Context) {
	var req VerifyOtpRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	// Basic Validations
	if len(req.PhoneNumber) < 8 {
		utils.BadRequest(ctx, "Invalid phone number")
		return
	}

	if len(req.Otp) != 6 {
		utils.BadRequest(ctx, "Invalid OTP")
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	otpCol := config.MongoDB.Collection(config.OtpCollections)
	var otpData models.Otp

	err := otpCol.FindOne(c, bson.M{"phone_number": req.PhoneNumber, "otp": req.Otp, "is_used": false}).Decode(&otpData)

	if err == mongo.ErrNoDocuments {
		utils.BadRequest(ctx, "Invalid OTP")
		return
	}

	if err != nil {
		utils.InternalServerError(ctx, "Failed to verify OTP")
		return
	}

	// checking if otp is expired
	if time.Now().After(otpData.ExpiresAt) {
		utils.BadRequest(ctx, "OTP has expired")
		return
	}

	// mark otp as used
	otpCol.UpdateOne(c, bson.M{"_id": otpData.ID}, bson.M{"$set": bson.M{"is_used": true}})

	userCol := config.MongoDB.Collection(config.UsersCollections)

	// find user and check if user is exist and its data
	var user models.Users
	var isUserNew = false
	err = userCol.FindOne(c, bson.M{"phone_number": req.PhoneNumber}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		// create new user
		user = models.Users{
			ID:          primitive.NewObjectID(),
			PhoneNumber: req.PhoneNumber,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserName:    "",
			FirstName:   "",
			LastName:    "",
			ProfilePic:  "",
			WalletId:    "",
			KycStatus:   "",
		}
		_, err = userCol.InsertOne(c, user)
		if err != nil {
			utils.InternalServerError(ctx, "Failed to create user")
			return
		}
	} else if err != nil {
		utils.InternalServerError(ctx, "Failed to find user")
		return
	}

	if user.WalletId == "" {
		isUserNew = true
	}

	// generate token
	tokenPair, err := utils.GenerateTokenPair(user.ID.Hex(), user.PhoneNumber, "")
	if err != nil {
		utils.InternalServerError(ctx, "Failed to generate tokens")
		return
	}
	// insert the refresh token into db
	tokenCol := config.MongoDB.Collection(config.TokenCollection)
	tokenCol.InsertOne(c, models.RefreshToken{
		ID:        primitive.NewObjectID(),
		UserId:    user.ID,
		Token:     tokenPair.RefreshToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	utils.Success(ctx, "OTP Verified Successfully", http.StatusOK, gin.H{
		"accessToken":  tokenPair.AccessToken,
		"refreshToken": tokenPair.RefreshToken,
		"isUserNew":    isUserNew,
	})
}
