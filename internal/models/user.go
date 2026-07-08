package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
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

type PublicUser struct {
	ID    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	// Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (u *User) ToPublic() PublicUser {
	return PublicUser{
		ID:   u.ID.Hex(),
		Name: u.UserName,
		// Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
