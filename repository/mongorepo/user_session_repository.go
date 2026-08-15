package mongorepo

import (
	"context"
	"time"
	"user-mgmt/domain"
	"user-mgmt/repository"
	"user-mgmt/repository/mongorepo/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const userSessionCollection = "user_sessions"

type userSessionRepository struct{ db *mongo.Database }

func NewUserSessionRepository(db *mongo.Database) repository.UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r *userSessionRepository) Create(ctx context.Context, session *domain.UserSession) error {
	_, err := r.db.Collection(userSessionCollection).InsertOne(ctx, model.FromUserSessionDomain(session))
	return err
}

func (r *userSessionRepository) GetByID(ctx context.Context, id string) (*domain.UserSession, error) {
	var session model.UserSession
	if err := r.db.Collection(userSessionCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&session); err != nil {
		return nil, err
	}
	return session.ToDomain(), nil
}

func (r *userSessionRepository) Rotate(ctx context.Context, id, currentHash, nextHash string, usedAt time.Time) error {
	result, err := r.db.Collection(userSessionCollection).UpdateOne(ctx, bson.M{
		"_id": id, "tokenHash": currentHash, "revokedAt": bson.M{"$exists": false}, "expiresAt": bson.M{"$gt": usedAt},
	}, bson.M{"$set": bson.M{"tokenHash": nextHash, "lastUsedAt": usedAt.UTC()}})
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return domain.ErrInvalidRefreshToken
	}
	return nil
}

func (r *userSessionRepository) Revoke(ctx context.Context, id string, revokedAt time.Time) error {
	_, err := r.db.Collection(userSessionCollection).UpdateOne(ctx, bson.M{"_id": id, "revokedAt": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"revokedAt": revokedAt.UTC()}})
	return err
}

func (r *userSessionRepository) RevokeAllByUserID(ctx context.Context, userID string, revokedAt time.Time) error {
	_, err := r.db.Collection(userSessionCollection).UpdateMany(ctx, bson.M{"userId": userID, "revokedAt": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"revokedAt": revokedAt.UTC()}})
	return err
}
