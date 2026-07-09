package mongo

import (
	"context"
	"errors"

	"github.com/sahilkarwasra/moepay/internal/models"
	"github.com/sahilkarwasra/moepay/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepository struct {
	collection *mongoDriver.Collection
}

func NewUserMongoRepository(
	collection *mongoDriver.Collection,
) repository.UserRepository {
	return &UserMongoRepository{
		collection: collection,
	}
}

func (r *UserMongoRepository) CreateUser(
	ctx context.Context,
	user *models.User,
) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongoDriver.IsDuplicateKeyError(err) {
			return repository.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserMongoRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	var user models.User

	err := r.collection.
		FindOne(ctx, bson.M{"email": email}).
		Decode(&user)

	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserMongoRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var user models.User

	err = r.collection.
		FindOne(ctx, bson.M{"_id": objectID}).
		Decode(&user)

	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserMongoRepository) UpdateUser(
	ctx context.Context,
	user *models.User,
) error {
	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"user_name":   user.UserName,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"profile_pic": user.ProfilePic,
			"wallet_id":   user.WalletId,
			"kyc_status":  user.KycStatus,
			"updated_at":  user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongoDriver.IsDuplicateKeyError(err) {
			return repository.ErrUserAlreadyExists
		}
		return err
	}

	if result.MatchedCount == 0 {
		return repository.ErrUserNotFound
	}

	return nil
}

func (r *UserMongoRepository) DeleteUser(
	ctx context.Context,
	userID string,
) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(
		ctx,
		bson.M{"_id": objectID},
	)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return repository.ErrUserNotFound
	}

	return nil
}

func (r *UserMongoRepository) GetUserByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) (*models.User, error) {
	var user models.User

	err := r.collection.
		FindOne(ctx, bson.M{"phone_number": phoneNumber}).
		Decode(&user)

	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
