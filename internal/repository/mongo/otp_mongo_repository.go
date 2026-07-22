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

func (r *OtpMongoRepository) FindOtpByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Otp, error) {
	var otpDoc models.Otp
	err := r.collection.FindOne(ctx, bson.M{
		"phone_number": phoneNumber,
	}).Decode(&otpDoc)

	return &otpDoc, err
}

func (r *OtpMongoRepository) UpdateOne(ctx context.Context, otp *models.Otp) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": otp.ID}, bson.M{"$set": otp})
	return err
}
