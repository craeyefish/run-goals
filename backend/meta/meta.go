package meta

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyUserID = contextKey("userID")
)

// GetUserIDFromContext gets the userID value from the context.
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(int64)
	return userID, ok
}
