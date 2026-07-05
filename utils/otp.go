package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

func GenerateOTP() string {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "123456"
	}
	otp := n.Int64() + 100000
	return fmt.Sprintf("%d", otp)
}

func LogOTP(phone, otp string) {
	log.Printf("Phone: %s | OTP: %s | (This would be sent via SMS in production)", phone, otp)
}
