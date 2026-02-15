package authctx

import (
	"context"
	"errors"
)

type contextKey string

const UserIdKey contextKey = "user_id"

func GeUserIdFromContext(ctx context.Context) (int64, error) {
    userId, ok := ctx.Value(UserIdKey).(int64)
    if !ok || userId == 0 {
        return 0, errors.New("user id not found in context")
    }

    return userId, nil
}