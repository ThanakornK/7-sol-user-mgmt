package handler

import (
	"net/http"
	"time"
	"user-mgmt/config"
	"user-mgmt/domain"
	"user-mgmt/handler/dto"
	"user-mgmt/service"
	"user-mgmt/utils"

	"github.com/gin-gonic/gin"
)

// authHandler struct implements the AuthHandler interface.
type authHandler struct {
	authService   service.AuthService
	refreshConfig config.RefreshConfig
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authService service.AuthService, refreshConfig config.RefreshConfig) *authHandler {
	return &authHandler{authService: authService, refreshConfig: refreshConfig}
}

// Login handles the login request.
func (h *authHandler) Login(c *gin.Context) {
	h.noStore(c)
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponseStruct("invalid request body", err.Error()))
		return
	}
	if errs := req.Validate(); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponseStruct("validation failed", utils.ErrorMessages(errs)))
		return
	}
	pair, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(c, err)
		return
	}
	h.setRefreshCookie(c, pair.RefreshToken)
	c.JSON(http.StatusOK, utils.NewSuccessResponseStruct("login successful", dto.NewTokenResponse(pair)))
}

// RefreshToken handles the refresh token request.
func (h *authHandler) RefreshToken(c *gin.Context) {
	h.noStore(c)
	raw, err := c.Cookie(h.refreshConfig.CookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponseStruct("authentication failed", "invalid refresh token"))
		return
	}
	pair, err := h.authService.Refresh(c.Request.Context(), raw)
	if err != nil {
		writeAuthError(c, err)
		return
	}
	h.setRefreshCookie(c, pair.RefreshToken)
	c.JSON(http.StatusOK, utils.NewSuccessResponseStruct("token refreshed successfully", dto.NewTokenResponse(pair)))
}

// Logout handles the logout request.
func (h *authHandler) Logout(c *gin.Context) {
	h.noStore(c)
	raw, _ := c.Cookie(h.refreshConfig.CookieName)
	err := h.authService.Logout(c.Request.Context(), raw)
	h.clearRefreshCookie(c)
	if err != nil {
		writeAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResponseStruct("logout successful", nil))
}

// writeAuthError writes the authentication error response.
func writeAuthError(c *gin.Context, err error) {
	status, response := utils.MapErrorResponse(
		err,
		http.StatusUnauthorized,
		"authentication failed",
		"invalid credentials",
		domain.ErrInvalidCredentials,
		domain.ErrInvalidRefreshToken,
		domain.ErrRefreshTokenExpired,
	)
	c.JSON(status, response)
}

// setRefreshCookie sets the refresh token cookie.
func (h *authHandler) setRefreshCookie(c *gin.Context, value string) {
	h.writeCookie(c, value, int(h.refreshConfig.ExpiresIn.Seconds()))
}

// clearRefreshCookie clears the refresh token cookie.
func (h *authHandler) clearRefreshCookie(c *gin.Context) { h.writeCookie(c, "", -1) }

// writeCookie writes the cookie.
func (h *authHandler) writeCookie(c *gin.Context, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{Name: h.refreshConfig.CookieName, Value: value, Path: h.refreshConfig.CookiePath, MaxAge: maxAge, Expires: time.Now().Add(time.Duration(maxAge) * time.Second), HttpOnly: true, Secure: h.refreshConfig.CookieSecure, SameSite: http.SameSiteStrictMode})
}

// noStore disables caching for the response.
func (h *authHandler) noStore(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
}
