package mongo

import (
	"context"

	"github.com/sahilkarwasra/moepay/internal/models"
	"github.com/sahilkarwasra/moepay/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type OtpMongoRepository struct {
	collection *mongoDriver.Collection
}

func NewOtpMongoRepository(
	collection *mongoDriver.Collection,
) repository.OtpRepository {
	return &OtpMongoRepository{
		collection: collection,
	}
}

func (r *OtpMongoRepository) DeleteManyByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) error {
	_, err := r.collection.DeleteMany(
		ctx,
		bson.M{"phone_number": phoneNumber},
	)

	return err
}

func (r *OtpMongoRepository) InsertOne(
	ctx context.Context,
	otp *models.Otp,
) error {
	_, err := r.collection.InsertOne(ctx, otp)
	return err
}
