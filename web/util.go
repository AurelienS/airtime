package web

import (
	"context"
	"fmt"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/webserver/middleware"
)

// GetUserFromRequestContext retrieves the user from the request context.
func GetUserFromRequestContext(ctx context.Context) (model.User, bool) {
	user, ok := ctx.Value(middleware.UserContextKey).(model.User)
	fmt.Println("file: util.go ~ line 13 ~ ok : ", ok)
	fmt.Println("file: util.go ~ line 13 ~ user : ", user)
	return user, ok
}
