package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// refreshSecretBytes is the number of bytes in the refresh token secret.
const refreshSecretBytes = 32

// RefreshCredential struct represents the refresh token credential.
type RefreshCredential struct {
	SessionID string
	Secret    string
	Raw       string
	Hash      string
}

// NewRefreshCredential creates a new refresh token credential.
func NewRefreshCredential() (RefreshCredential, error) {
	return newRefreshCredential(uuid.NewString())
}

// NewRotatedRefreshCredential creates a new rotated refresh token credential.
func NewRotatedRefreshCredential(sessionID string) (RefreshCredential, error) {
	if _, err := uuid.Parse(sessionID); err != nil {
		return RefreshCredential{}, errors.New("invalid user session ID")
	}
	return newRefreshCredential(sessionID)
}

// newRefreshCredential creates a new refresh token credential.
func newRefreshCredential(sessionID string) (RefreshCredential, error) {
	secretBytes := make([]byte, refreshSecretBytes)
	if _, err := rand.Read(secretBytes); err != nil {
		return RefreshCredential{}, err
	}
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)
	return RefreshCredential{SessionID: sessionID, Secret: secret, Raw: sessionID + "." + secret, Hash: HashRefreshSecret(secret)}, nil
}

// ParseRefreshCredential parses the refresh token credential.
func ParseRefreshCredential(raw string) (string, string, error) {
	if raw != strings.TrimSpace(raw) {
		return "", "", errors.New("invalid refresh token")
	}
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return "", "", errors.New("invalid refresh token")
	}
	if _, err := uuid.Parse(parts[0]); err != nil {
		return "", "", errors.New("invalid refresh token")
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || len(decoded) != refreshSecretBytes {
		return "", "", errors.New("invalid refresh token")
	}
	return parts[0], parts[1], nil
}

// HashRefreshSecret hashes the refresh token secret.
func HashRefreshSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}
