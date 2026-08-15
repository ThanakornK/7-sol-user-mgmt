// package config provides configuration for the application.
package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config is the configuration struct for the application.
type Config struct {
	Server        ServerConfig
	JWT           JWTConfig
	Refresh       RefreshConfig
	MongoDBConfig MongoDBConfig
}

// ServerConfig is the server configuration struct.
type ServerConfig struct {
	Port          string
	GinMode       string
	AllowedOrigin string
}

// MongoDBConfig is the MongoDB configuration struct.
type MongoDBConfig struct {
	URI            string
	DatabaseName   string
	CollectionName string
}

// JWTConfig is the JWT configuration struct.
type JWTConfig struct {
	Secret    string
	Issuer    string
	Audience  string
	ExpiresIn time.Duration
}

// RefreshConfig is the refresh configuration struct.
type RefreshConfig struct {
	ExpiresIn    time.Duration
	CookieName   string
	CookiePath   string
	CookieSecure bool
}

// LoadEnv loads the environment variables into the Config struct.
func LoadEnv() (*Config, error) {

	_ = godotenv.Load()

	// Get environment variables that need to parse before using.
	jwtExpiresIn, err := time.ParseDuration(getEnvWithDefault("JWT_EXPIRES_IN", "15m"))
	if err != nil {
		return nil, err
	}
	refreshExpiresIn, err := time.ParseDuration(getEnvWithDefault("REFRESH_EXPIRES_IN", "168h"))
	if err != nil {
		return nil, err
	}
	cookieSecure, err := strconv.ParseBool(getEnvWithDefault("COOKIE_SECURE", "true"))
	if err != nil {
		return nil, err
	}

	// Get environment variables that required.
	mongoDbURI := os.Getenv("MONGO_URI")
	if mongoDbURI == "" || mongoDbURI == "default_mongo_uri" {
		return nil, errors.New("MONGO_URI must be configured")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" || jwtSecret == "default_jwt_secret" {
		return nil, errors.New("JWT_SECRET must be configured")
	}

	return &Config{
		Server: ServerConfig{
			Port:          getEnvWithDefault("PORT", "8080"),
			GinMode:       getEnvWithDefault("GIN_MODE", "debug"),
			AllowedOrigin: getEnvWithDefault("ALLOWED_ORIGIN", "http://localhost:3000"),
		},
		MongoDBConfig: MongoDBConfig{
			URI:            mongoDbURI,
			DatabaseName:   getEnvWithDefault("MONGO_DATABASE_NAME", "user_mgmt"),
			CollectionName: getEnvWithDefault("MONGO_COLLECTION_NAME", "users"),
		},
		JWT: JWTConfig{
			Secret:    jwtSecret,
			Issuer:    getEnvWithDefault("JWT_ISSUER", "user-mgmt"),
			Audience:  getEnvWithDefault("JWT_AUDIENCE", "user-mgmt-api"),
			ExpiresIn: jwtExpiresIn,
		},
		Refresh: RefreshConfig{
			ExpiresIn:    refreshExpiresIn,
			CookieName:   getEnvWithDefault("REFRESH_COOKIE_NAME", "refresh_token"),
			CookiePath:   getEnvWithDefault("REFRESH_COOKIE_PATH", "/api/v1/auth"),
			CookieSecure: cookieSecure,
		},
	}, nil

}

// getEnvWithDefault returns the environment variable value or the default value if it's not set.
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
