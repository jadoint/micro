package contextkey

// ContextKey key used for "Visitor" context in http.Request
type ContextKey struct {
	Key string
}

// GetVisitorKey return context key for "Visitor"
func GetVisitorKey() ContextKey {
	return ContextKey{Key: "Visitor"}
}
