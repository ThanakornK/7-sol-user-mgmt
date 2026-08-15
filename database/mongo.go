package database

import (
	"context"
	"errors"
	"strings"
	"time"
	"user-mgmt/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDBClient(cfg config.MongoDBConfig) (*mongo.Client, error) {

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
