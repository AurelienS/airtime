package handler

import (
	"github.com/AurelienS/cigare/internal/user"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) UserHandler {
	return UserHandler{
		userService: userService,
	}
}
