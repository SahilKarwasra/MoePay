package service

import (
	"context"
	"time"

	"github.com/sahilkarwasra/moepay/internal/models"
	"github.com/sahilkarwasra/moepay/internal/repository"
	"github.com/sahilkarwasra/moepay/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo repository.UserRepository
	otpRepo  repository.OtpRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	otpRepo repository.OtpRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		otpRepo:  otpRepo,
	}
}

func (s *UserService) SendOTP(ctx context.Context, phoneNumber string) (string, error) {
	if len(phoneNumber) < 8 {
		return "", repository.ErrInvalidPhoneNumber
	}

	// delete old otp of this number
	err := s.otpRepo.DeleteManyByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return "", err
	}

	// generate otp
	otp := utils.GenerateOTP()
	otpDoc := models.Otp{
		ID:          primitive.NewObjectID(),
		PhoneNumber: phoneNumber,
		Otp:         otp,
		IsUsed:      false,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	err = s.otpRepo.InsertOne(ctx, &otpDoc)
	if err != nil {
		return "", err
	}

	utils.LogOTP(phoneNumber, otp)

	return otp, nil
}
