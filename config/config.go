package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig
	JWT           JWTConfig
	MongoDBConfig MongoDBConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type MongoDBConfig struct {
	URI            string
	DatabaseName   string
	CollectionName string
}

type JWTConfig struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshExpiresIn time.Duration
}

func LoadEnv() (*Config, error) {

	_ = godotenv.Load()

	jwtExpiresIn, _ := time.ParseDuration(getEnvWithDefault("JWT_EXPIRES_IN", "15m"))
	jwtRefreshExpiresIn, _ := time.ParseDuration(getEnvWithDefault("JWT_REFRESH_EXPIRES_IN", "1h"))

	return &Config{
		Server: ServerConfig{
			Port:    getEnvWithDefault("PORT", "8080"),
			GinMode: getEnvWithDefault("GIN_MODE", "release"),
		},
		MongoDBConfig: MongoDBConfig{
			URI:            getEnvWithDefault("MONGO_URI", "mongodb://admin:admin12345@localhost:27017"),
			DatabaseName:   getEnvWithDefault("MONGO_DATABASE_NAME", "user_mgmt"),
			CollectionName: getEnvWithDefault("MONGO_COLLECTION_NAME", "users"),
		},
		JWT: JWTConfig{
			Secret:           getEnvWithDefault("JWT_SECRET", "default_jwt_secret"),
			ExpiresIn:        jwtExpiresIn,
			RefreshExpiresIn: jwtRefreshExpiresIn,
		},
	}, nil

}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
