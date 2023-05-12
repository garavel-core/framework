package cache

type Store interface {
	// Retrieve an item from the cache by key.
	Get(key string) any

	// Retrieve multiple items from the cache by key.
	//
	// Items not found in the cache will have a null value.
	Many(keys []string) map[string]any

	// Store an item in the cache for a given number of seconds.
	Put(key string, value any, seconds int) bool

	// Store multiple items in the cache for a given number of seconds.
	PutMany(values map[string]any, seconds int) bool

	// Increment the value of an item in the cache.
	Increment(key string, value ...any) any

	// Decrement the value of an item in the cache.
	Decrement(key string, value ...any) any

	// Store an item in the cache indefinitely.
	Forever(key string, value any) bool

	// Remove an item from the cache.
	Forget(key string) bool

	// Remove all items from the cache.
	Flush() bool

	// Get the cache key prefix.
	GetPrefix() string
}
