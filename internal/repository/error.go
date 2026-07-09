package repository

import "errors"

var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrOtpNotFound        = errors.New("otp not found")
	ErrInvalidOtp         = errors.New("invalid otp")
	ErrOtpUsed            = errors.New("otp already used")
	ErrOtpExpired         = errors.New("otp expired")
)
