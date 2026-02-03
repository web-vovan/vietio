package authctx

import "context"

type contextKey string

const UserIdKey contextKey = "user_id"

func GeUserIdFromContext(ctx context.Context) int64 {
    return ctx.Value(UserIdKey).(int64)
}