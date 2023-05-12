package cache

type AtomicStore interface {
	// Store an item in the cache if the key does not exist.
	Add(key string, value any, ttl ...any) bool
}
