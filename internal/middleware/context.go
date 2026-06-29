package middleware

import "context"

type ctxKey string

const (
	ctxRequestID ctxKey = "request_id"
	ctxUserID    ctxKey = "user_id"
	ctxUserEmail ctxKey = "user_email"
)

func newContextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxRequestID, id)
}

func GetRequestID(ctx context.Context) string {
	v, _ := ctx.Value(ctxRequestID).(string)
	return v
}

func NewContextWithUser(ctx context.Context, userID, email string) context.Context {
	ctx = context.WithValue(ctx, ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxUserEmail, email)
	return ctx
}

func GetUserID(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserID).(string)
	return v
}

func GetUserEmail(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserEmail).(string)
	return v
}
