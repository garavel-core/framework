package cache

type Repository interface {
	// Fetches a value from the cache.
	Get(key string, defaultValue ...any) any

	// Persists data in the cache, uniquely referenced by a key with an optional expiration TTL time.
	Set(key string, value any, ttl ...any) bool

	// Delete an item from the cache by its unique key.
	Delete(key string) bool

	// Wipes clean the entire cache's keys.
	Clear() bool

	// Obtains multiple cache items by their unique keys.
	GetMultiple(keys []string, defaultValue ...any) map[string]any

	// Persists a set of key => value pairs in the cache, with an optional TTL.
	SetMultiple(values map[string]any, ttl ...any) bool

	// Deletes multiple cache items in a single operation.
	DeleteMultiple(keys []string) bool

	// Determines whether an item is present in the cache.
	//
	//  NOTE: It is recommended that has() is only to be used for cache warming type purposes
	//  and not to be used within your live applications operations for get/set, as this method
	//  is subject to a race condition where your has() will return true and immediately after,
	//  another script can remove it making the state of your app out of date.
	Has(key string) bool

	// Retrieve an item from the cache and delete it.
	Pull(key string, defaultValue ...any) any

	// Store an item in the cache.
	Put(key string, value any, ttl ...any) bool

	// Store an item in the cache if the key does not exist.
	Add(key string, value any, ttl ...any) bool

	// Increment the value of an item in the cache.
	Increment(key string, value ...any) any

	// Decrement the value of an item in the cache.
	Decrement(key string, value ...any) any

	// Store an item in the cache indefinitely.
	Forever(key string, value any) bool

	// Get an item from the cache, or execute the given Closure and store the result.
	Remember(key string, ttl any, callback func() any) any

	// Get an item from the cache, or execute the given Closure and store the result forever.
	Sear(key string, callback func() any) any

	// Get an item from the cache, or execute the given Closure and store the result forever.
	RememberForever(key string, callback func() any) any

	// Remove an item from the cache.
	Forget(key string) bool

	// Get the cache store implementation.
	GetStore() Store
}
