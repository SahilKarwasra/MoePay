package repository

import (
	"context"

	"github.com/sahilkarwasra/moepay/internal/models"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, token *models.RefreshToken) error
	GetTokenByUserID(ctx context.Context, userID string) (*models.RefreshToken, error)
	UpdateToken(ctx context.Context, token *models.RefreshToken) error
	DeleteToken(ctx context.Context, userID string) error
}
