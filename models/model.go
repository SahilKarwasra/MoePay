package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Otp struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Otp         string             `bson:"otp" json:"otp"`
	IsUsed      bool               `bson:"is_used" json:"is_used"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
}

type Users struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	UserName    string             `bson:"user_name" json:"user_name"`
	FirstName   string             `bson:"first_name" json:"first_name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	ProfilePic  string             `bson:"profile_pic" json:"profile_pic"`
	WalletId    string             `bson:"wallet_id" json:"wallet_id"`
	KycStatus   string             `bson:"kyc_status" json:"kyc_status"`
}

type RefreshToken struct {
	ID        primitive.ObjectID `bson:"_id,omnitempty" json:"id"`
	UserId    primitive.ObjectID `bson:"_id,omnitempty" json:"user_id"`
	Token     string             `bson:"token" json:"token"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
