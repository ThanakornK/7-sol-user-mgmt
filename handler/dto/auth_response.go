package dto

import "user-mgmt/domain"

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int64  `json:"expiresIn"`
}

func NewTokenResponse(pair *domain.TokenPair) TokenResponse {
	return TokenResponse{
		AccessToken: pair.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(pair.AccessExpiresIn.Seconds()),
	}
}
