package config

import "go.mongodb.org/mongo-driver/mongo"

var MongoDB *mongo.Database

const (
	UsersCollections = "users"
	OtpCollections   = "otp"
	TokenCollection  = "refresh_tokens"
)
