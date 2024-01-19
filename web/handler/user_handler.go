package handler

import (
	"github.com/AurelienS/cigare/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}
