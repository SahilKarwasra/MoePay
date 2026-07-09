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
	userRepo  repository.UserRepository
	otpRepo   repository.OtpRepository
	tokenRepo repository.TokenRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	otpRepo repository.OtpRepository,
	tokenRepo repository.TokenRepository,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		tokenRepo: tokenRepo,
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

func (s *UserService) VerifyOtpRequest(ctx context.Context, phoneNumber string, otp string) (models.VerifyOTPResponse, error) {
	var isUserNew = false

	// basic validation
	if len(phoneNumber) < 8 {
		return models.VerifyOTPResponse{}, repository.ErrInvalidPhoneNumber
	}
	if len(otp) != 6 {
		return models.VerifyOTPResponse{}, repository.ErrInvalidOtp
	}

	// find otp
	otpDoc, err := s.otpRepo.FindOtpByPhoneNumberAndOtp(ctx, phoneNumber, otp)

	// checking for error and not found
	if err != nil {
		return models.VerifyOTPResponse{}, repository.ErrOtpNotFound
	}

	// checking if otp is used or not
	if otpDoc.IsUsed {
		return models.VerifyOTPResponse{}, repository.ErrOtpUsed
	}

	// checking if otp is expired or not
	if otpDoc.ExpiresAt.Before(time.Now()) {
		return models.VerifyOTPResponse{}, repository.ErrOtpExpired
	}

	// mark otp as used
	otpDoc.IsUsed = true

	err = s.otpRepo.UpdateOne(ctx, otpDoc)
	if err != nil {
		return models.VerifyOTPResponse{}, err
	}

	var user *models.User
	user, err = s.userRepo.GetUserByPhoneNumber(ctx, phoneNumber)

	// check if user exists or not
	if err == repository.ErrUserNotFound {
		user = &models.User{
			ID:          primitive.NewObjectID(),
			PhoneNumber: phoneNumber,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserName:    "",
			FirstName:   "",
			LastName:    "",
			ProfilePic:  "",
			WalletId:    "",
			KycStatus:   "",
		}

		err = s.userRepo.CreateUser(ctx, user)
		if err != nil {
			return models.VerifyOTPResponse{}, err
		}
	} else if err != nil {
		return models.VerifyOTPResponse{}, err
	}

	if user.WalletId == "" {
		isUserNew = true
	}

	// generate token
	tokenPair, err := utils.GenerateTokenPair(user.ID.Hex(), user.PhoneNumber, "")
	if err != nil {
		return models.VerifyOTPResponse{}, err
	}

	// save token to db
	err = s.tokenRepo.CreateToken(ctx, &models.RefreshToken{
		ID:        primitive.NewObjectID(),
		UserId:    user.ID,
		Token:     tokenPair.RefreshToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return models.VerifyOTPResponse{}, err
	}

	return models.VerifyOTPResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		IsUserNew:    isUserNew,
	}, nil

}
