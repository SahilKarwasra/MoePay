package repository

import "errors"

var (
	ErrInvalidPhoneNumber = errors.New("Invalid Phone Number")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrOtpNotFound        = errors.New("otp not found")
)
