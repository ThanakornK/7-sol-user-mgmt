package utils

// ContextKey type for context key to store in context.
type ContextKey string

const (
	UserIDContextKey  ContextKey = "user_id"
	TraceIDContextKey ContextKey = "trace_id"
)
