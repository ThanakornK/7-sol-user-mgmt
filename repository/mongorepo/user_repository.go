// package mongorepo for MongoDB repository implementations
package mongorepo

import (
	"context"
	"time"
	"user-mgmt/domain"
	"user-mgmt/repository"
	"user-mgmt/repository/mongorepo/model"
	"user-mgmt/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	userModel := model.FromUserDomain(user)

	_, err := r.db.Collection("users").InsertOne(ctx, userModel)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	userModel := &model.User{}
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(userModel)
	if err != nil {
		return nil, err
	}

	return userModel.ToDomain(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	userModel := &model.User{}
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(userModel)
	if err != nil {
		return nil, err
	}

	return userModel.ToDomain(), nil
}

func (r *userRepository) GetUserList(ctx context.Context, page, pageSize int) ([]*domain.User, utils.Pagination, error) {
	var users []*model.User
	cursor, err := r.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, utils.Pagination{}, err
	}
	defer cursor.Close(ctx)

	var total int64
	total, err = r.db.Collection("users").CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	if err := cursor.All(ctx, &users); err != nil {
		return nil, utils.Pagination{}, err
	}

	// Convert to domain.User
	domainUsers := make([]*domain.User, 0, len(users))
	for _, user := range users {
		domainUsers = append(domainUsers, user.ToDomain())
	}

	return domainUsers, utils.Pagination{
		Total: total,
	}, nil
}

func (r *userRepository) Update(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	now := time.Now().UTC()

	filter := bson.M{
		"_id": user.ID,
	}

	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"email":     user.Email,
			"updatedAt": now,
		},
	}

	// Return updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedModel model.User

	err := r.db.Collection("users").
		FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updatedModel)
	if err != nil {
		return nil, err
	}

	return updatedModel.ToDomain(), nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Collection("users").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
