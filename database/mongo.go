package database

import (
	"context"
	"errors"
	"strings"
	"time"
	"user-mgmt/config"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// IndexModels returns the index models for users and user sessions. it need to define when initialize the database.
func IndexModels() (users []mongo.IndexModel, userSessions []mongo.IndexModel) {
	users = []mongo.IndexModel{{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)}}
	userSessions = []mongo.IndexModel{
		{Keys: bson.D{{Key: "expiresAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
		{Keys: bson.D{{Key: "userId", Value: 1}}},
	}
	return users, userSessions
}

// EnsureIndexes ensures the index models are created in the database.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	users, userSessions := IndexModels()
	if _, err := db.Collection("users").Indexes().CreateMany(ctx, users); err != nil {
		return err
	}
	_, err := db.Collection("user_sessions").Indexes().CreateMany(ctx, userSessions)
	return err
}

func NewMongoDBClient(cfg config.MongoDBConfig) (*mongo.Client, error) {

	// Init MongoDB
	if strings.TrimSpace(cfg.URI) == "" {
		return nil, errors.New("MONGO_URI must be set")
	}
	if strings.TrimSpace(cfg.DatabaseName) == "" {
		return nil, errors.New("MONGO_DATABASE_NAME must be set")
	}

	clientOption := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(clientOption)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	return client, nil
}

// GetDatabase returns a database instance
func GetDatabase(client *mongo.Client, cfg config.MongoDBConfig) *mongo.Database {
	return client.Database(cfg.DatabaseName)
}
