package service

import (
	"context"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
	"user-mgmt/config"
	"user-mgmt/domain"
	"user-mgmt/repository"
	"user-mgmt/utils"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// authService implements the AuthService interface.
type authService struct {
	users      repository.UserRepository
	sessions   repository.UserSessionRepository
	jwtConfig  *config.JWTConfig
	refreshCfg config.RefreshConfig
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(users repository.UserRepository, sessions repository.UserSessionRepository, jwtConfig *config.JWTConfig, refreshConfig config.RefreshConfig) AuthService {
	return &authService{users: users, sessions: sessions, jwtConfig: jwtConfig, refreshCfg: refreshConfig}
}

// Login logs in a user and returns a token pair.
func (s *authService) Login(ctx context.Context, email, password string) (*domain.TokenPair, error) {
	// Get user by email.
	user, err := s.users.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}
	// Check password matches.
	if !utils.CheckPassword(password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate access token.
	now := time.Now().UTC()
	accessToken, err := utils.GenerateAccessToken(s.jwtConfig, user.ID.String(), now)
	if err != nil {
		return nil, err
	}
	credential, err := utils.NewRefreshCredential()
	if err != nil {
		return nil, err
	}
	// Create user session.
	session := domain.NewUserSession(credential.SessionID, user.ID.String(), credential.Hash, s.refreshCfg.ExpiresIn)
	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, err
	}
	return &domain.TokenPair{AccessToken: accessToken, RefreshToken: credential.Raw, AccessExpiresIn: s.jwtConfig.ExpiresIn}, nil
}

// Refresh refreshes a token pair.
func (s *authService) Refresh(ctx context.Context, rawRefreshToken string) (*domain.TokenPair, error) {
	// Parse refresh token. to get session ID and secret.
	sessionID, secret, err := utils.ParseRefreshCredential(rawRefreshToken)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	// Get user session by ID.
	session, err := s.sessions.GetByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrInvalidRefreshToken
		}
		return nil, err
	}
	// Check session is valid.
	now := time.Now().UTC()
	if session.RevokedAt != nil {
		return nil, domain.ErrInvalidRefreshToken
	}
	if !now.Before(session.ExpiresAt) {
		return nil, domain.ErrRefreshTokenExpired
	}
	// Check session token hash matches.
	expected, expectedErr := hex.DecodeString(session.TokenHash)
	presented, presentedErr := hex.DecodeString(utils.HashRefreshSecret(secret))
	if expectedErr != nil || presentedErr != nil || subtle.ConstantTimeCompare(expected, presented) != 1 {
		if err := s.sessions.Revoke(ctx, session.ID, now); err != nil {
			return nil, fmt.Errorf("revoke user session: %w", err)
		}
		return nil, domain.ErrInvalidRefreshToken
	}

	// Get user by ID.
	user, err := s.users.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	// Rotate session token.
	next, err := utils.NewRotatedRefreshCredential(session.ID)
	if err != nil {
		return nil, err
	}
	// Update session token hash.
	if err := s.sessions.Rotate(ctx, session.ID, session.TokenHash, next.Hash, now); err != nil {
		if errors.Is(err, domain.ErrInvalidRefreshToken) {
			if revokeErr := s.sessions.Revoke(ctx, session.ID, now); revokeErr != nil {
				return nil, fmt.Errorf("revoke user session: %w", revokeErr)
			}
		}
		return nil, err
	}
	// Generate new access token.
	now = time.Now().UTC()
	accessToken, err := utils.GenerateAccessToken(s.jwtConfig, user.ID.String(), now)
	if err != nil {
		return nil, err
	}
	return &domain.TokenPair{AccessToken: accessToken, RefreshToken: next.Raw, AccessExpiresIn: s.jwtConfig.ExpiresIn}, nil
}

// Logout logs out a user and revokes their session.
func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	// Parse refresh token. to get session ID.
	sessionID, _, err := utils.ParseRefreshCredential(rawRefreshToken)
	if err != nil {
		return nil
	}
	// Revoke session.
	return s.sessions.Revoke(ctx, sessionID, time.Now().UTC())
}
