package handlers

import (
	"context"
	"errors"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
)

// Context keys for storing values in request context
type contextKey string

const (
	// UserContextKey is the key for storing the authenticated user in context
	UserContextKey contextKey = "user"
)

// ErrNoUserInContext is returned when no user is found in the request context
var ErrNoUserInContext = errors.New("no user found in context")

// GetUserFromContext retrieves the authenticated user from the request context
func GetUserFromContext(ctx context.Context) (*models.User, error) {
	logger.Debug().Msg("Attempting to retrieve user from context")
	user, ok := ctx.Value(UserContextKey).(*models.User)
	if !ok || user == nil {
		logger.Error().Bool("ok", ok).Msg("No user found in context - authentication middleware may not have run")
		return nil, ErrNoUserInContext
	}
	logger.Debug().Str("user_id", user.ID).Msg("User found in context")
	return user, nil
}

// SetUserInContext sets the authenticated user in the request context
func SetUserInContext(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}
