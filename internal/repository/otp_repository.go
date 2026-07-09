package repository

import (
	"context"

	"github.com/sahilkarwasra/moepay/internal/models"
)

type OtpRepository interface {
	DeleteManyByPhoneNumber(ctx context.Context, phoneNumber string) error
	InsertOne(ctx context.Context, otp *models.Otp) error
	FindOtpByPhoneNumberAndOtp(ctx context.Context, phoneNumber string, otp string) (*models.Otp, error)
	UpdateOne(ctx context.Context, otp *models.Otp) error
}
