package rid

import "context"

type ridCtxKey struct{}

// WithContext save the request id with context.Value.
func WithContext(ctx context.Context, reqId string) context.Context {
	// Do not store duplicate.
	if _id, ok := ctx.Value(ridCtxKey{}).(string); ok && _id == reqId {
		return ctx
	}
	return context.WithValue(ctx, ridCtxKey{}, reqId)
}

// FromContext get the request id from context.Value.
func FromContext(ctx context.Context) (reqId string) {
	var ok bool
	reqId, ok = ctx.Value(ridCtxKey{}).(string)
	if !ok {
		reqId = "none"
	}
	return
}
