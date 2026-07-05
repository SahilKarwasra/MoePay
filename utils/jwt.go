package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserId      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func getAccessSecret() []byte {
	secret := os.Getenv("access_secret")
	if secret == "" {
		secret = "dummy_access_secret"
	}
	return []byte(secret)
}

func getRefreshSecret() []byte {
	secret := os.Getenv("refresh_secret")
	if secret == "" {
		secret = "dummy_refresh_secret"
	}
	return []byte(secret)
}

func GenerateTokenPair(userId string, phone string, name string) (*TokenPair, error) {
	// access token for 15 minutes
	accessExpiry := time.Now().Add(15 * time.Minute)
	accessClaims := TokenClaims{
		UserId:      userId,
		PhoneNumber: phone,
		Name:        name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "moepay",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			Subject:   userId,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(getAccessSecret())
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %v", err)
	}

	// refresh token for 7 days
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := TokenClaims{
		UserId:      userId,
		PhoneNumber: phone,
		Name:        name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "moepay",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			Subject:   userId,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(getRefreshSecret())
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %v", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessExpiry.Unix()),
	}, nil
}

func ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return getAccessSecret(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return getRefreshSecret(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
