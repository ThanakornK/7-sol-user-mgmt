package handler

import "user-mgmt/service"

type authHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *authHandler {
	return &authHandler{
		userService: userService,
	}
}
