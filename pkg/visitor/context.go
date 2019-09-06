package visitor

// ContextKey key used for "Visitor" context in http.Request
type ContextKey struct {
	Key string
}

// GetContextKey return context key for "Visitor"
func GetContextKey() ContextKey {
	return ContextKey{Key: "Visitor"}
}
