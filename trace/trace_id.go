package trace

import "context"

const (
	IdKey = "tid" // The key string of trace id.
)

type ctxKey struct{}

// ContextWithId store the trace id in context.Value.
func ContextWithId(ctx context.Context, id string) context.Context {
	// Do not store duplicate.
	if val := ctx.Value(ctxKey{}); val != nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, id)
}

// IdFromContext get the trace id from context.Value.
func IdFromContext(ctx context.Context) (id string) {
	var ok bool
	id, ok = ctx.Value(ctxKey{}).(string)
	if !ok {
		id = "none"
	}
	return
}
