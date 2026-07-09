package mongo

import (
	"context"

	"github.com/sahilkarwasra/moepay/internal/models"
	"github.com/sahilkarwasra/moepay/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenMongoRepository struct {
	collection *mongo.Collection
}

func NewTokenMongoRepository(collection *mongo.Collection) repository.TokenRepository {
	return &TokenMongoRepository{collection: collection}
}

func (r *TokenMongoRepository) CreateToken(ctx context.Context, token *models.RefreshToken) error {
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *TokenMongoRepository) GetTokenByUserID(ctx context.Context, userID string) (*models.RefreshToken, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var token models.RefreshToken
	err = r.collection.FindOne(ctx, bson.M{"user_id": objectID}).Decode(&token)
	return &token, err
}

func (r *TokenMongoRepository) UpdateToken(ctx context.Context, token *models.RefreshToken) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": token.UserId}, bson.M{"$set": token})
	return err
}

func (r *TokenMongoRepository) DeleteToken(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"user_id": objectID})
	return err
}
