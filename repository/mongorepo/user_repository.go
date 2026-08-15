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

// userRepository struct implements the UserRepository interface.
type userRepository struct {
	db         *mongo.Database
	collection userCollection
}

type userCollection interface {
	InsertOne(context.Context, any, ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	FindOne(context.Context, any, ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	Find(context.Context, any, ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	CountDocuments(context.Context, any, ...options.Lister[options.CountOptions]) (int64, error)
	FindOneAndUpdate(context.Context, any, any, ...options.Lister[options.FindOneAndUpdateOptions]) *mongo.SingleResult
	DeleteOne(context.Context, any, ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &userRepository{db: db}
}

func newUserRepository(collection userCollection) *userRepository {
	return &userRepository{collection: collection}
}

func (r *userRepository) users() userCollection {
	if r.collection != nil {
		return r.collection
	}
	return r.db.Collection("users")
}

// Create creates a new user in the database.
func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	userModel := model.FromUserDomain(user)

	_, err := r.users().InsertOne(ctx, userModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrEmailExists
		}
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID from the database.
func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	userModel := &model.User{}
	err := r.users().FindOne(ctx, bson.M{"_id": id}).Decode(userModel)
	if err != nil {
		return nil, err
	}

	return userModel.ToDomain(), nil
}

// GetByEmail retrieves a user by email from the database.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	userModel := &model.User{}
	err := r.users().FindOne(ctx, bson.M{"email": email}).Decode(userModel)
	if err != nil {
		return nil, err
	}

	return userModel.ToDomain(), nil
}

// GetUserList retrieves a list of users from the database with pagination.
func (r *userRepository) GetUserList(ctx context.Context, page, pageSize int) ([]*domain.User, utils.Pagination, error) {
	var users []*model.User
	findOptions := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))
	cursor, err := r.users().Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, utils.Pagination{}, err
	}
	defer cursor.Close(ctx)

	var total int64
	total, err = r.users().CountDocuments(ctx, bson.M{})
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
		Page:     int64(page),
		PageSize: int64(pageSize),
		Total:    total,
	}, nil
}

// Update updates a user in the database.
func (r *userRepository) Update(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	now := time.Now().UTC()

	filter := bson.M{
		"_id": user.ID.String(),
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

	err := r.users().FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updatedModel)
	if err != nil {
		return nil, err
	}

	return updatedModel.ToDomain(), nil
}

// Delete deletes a user from the database.
func (r *userRepository) Delete(ctx context.Context, id string) error {
	_, err := r.users().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
