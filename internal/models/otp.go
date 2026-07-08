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
