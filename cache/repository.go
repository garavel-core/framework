package cache

import (
	contracts "github.com/garavel-core/framework/contracts/cache"
	"github.com/garavel-core/framework/support/arr"
	"github.com/garavel-core/framework/support/helpers"
)

type Repository struct {
	// The cache store implementation.
	store contracts.Store
	// The default number of seconds to store items.
	defaultCacheTime int
	// The event dispatcher implementation.
	// events contracts.Dispatcher
}

// Create a new cache repository instance.
func NewRepository(store contracts.Store) *Repository {
	return &Repository{store: store, defaultCacheTime: 3600}
}

// Determine if an item exists in the cache.
func (r *Repository) Has(key string) bool {
	return r.store.Get(key) != nil
}

// Determine if an item doesn't exist in the cache.
func (r *Repository) Missing(key string) bool {
	return !r.Has(key)
}

// Retrieve an item from the cache by key.
func (r *Repository) Get(key string, defaultValue ...any) any {
	value := r.store.Get(r.itemKey(key))

	// If we could not find the cache value, we will fire the missed event and get
	// the default value for this cache value. This default could be a callback
	// so we will execute the value function which will resolve it if needed.
	if value == nil {
		// r.event(NewCacheMissed(key))

		if defaultValue != nil && defaultValue[0] != nil {
			value = helpers.Value(defaultValue)
		}

	} else {
		// r.event(NewCacheHit(key, value))
	}

	return value
}

// Retrieve multiple items from the cache by key.
//
// Items not found in the cache will have a null value.
func (r *Repository) Many(keys any) map[string]any {
	if defaults, ok := keys.(map[string]any); ok {
		return arr.Map(r.store.Many(arr.Keys(defaults)), func(value any, key string) any {
			if value == nil {
				return defaults[key]
			}

			return value
		})
	}

	return r.store.Many(keys.([]string))
}

// Obtains multiple cache items by their unique keys.
func (r *Repository) GetMultiple(keys []string, defaultValue ...any) map[string]any {
	var value any

	if len(defaultValue) != 0 {
		value = defaultValue[0]
	}

	defaults := make(map[string]any, len(keys))

	for _, key := range keys {
		defaults[key] = value
	}

	return r.Many(defaults)
}

// Retrieve an item from the cache and delete it.
func (r *Repository) Pull(key string, defaultValue ...any) any {
	return helpers.Tap(r.Get(key, defaultValue...), func(_ any) {
		r.Forget(key)
	})
}

// Store an item in the cache.
func (r *Repository) Put(key string, value any, ttl ...any) bool {
	if ttl == nil || ttl[0] == nil {
		return r.Forever(key, value)
	}

	seconds := r.getSeconds(ttl)

	if seconds <= 0 {
		return r.Forget(key)
	}

	result := r.store.Put(r.itemKey(key), value, seconds)

	// if result {
	//     r.event(NewKeyWritten(key, value, seconds))
	// }

	return result
}

// Persists data in the cache, uniquely referenced by a key with an optional expiration TTL time.
func (r *Repository) Set(key string, value any, ttl ...any) bool {
	return r.Put(key, value, ttl...)
}

// Store multiple items in the cache for a given number of seconds.
func (r *Repository) PutMany(values map[string]any, ttl ...any) bool {
	if ttl == nil || ttl[0] == nil {
		return r.putManyForever(values)
	}

	seconds := r.getSeconds(ttl)

	if seconds <= 0 {
		return r.DeleteMultiple(arr.Keys(values))
	}

	result := r.store.PutMany(values, seconds)

	// if result {
	//     for key, value := range values {
	//         r.event(NewKeyWritten(key, value, seconds))
	//     }
	// }

	return result
}

// Store multiple items in the cache indefinitely.
func (r *Repository) putManyForever(values map[string]any) bool {
	result := true

	for key, value := range values {
		if !r.Forever(key, value) {
			result = false
		}
	}

	return result
}

// Persists a set of key => value pairs in the cache, with an optional TTL.
func (r *Repository) SetMultiple(values map[string]any, ttl ...any) bool {
	return r.PutMany(values, ttl...)
}

// Store an item in the cache if the key does not exist.
func (r *Repository) Add(key string, value any, ttl ...any) bool {
	var seconds any

	if ttl != nil && ttl[0] != nil {
		seconds = r.getSeconds(ttl)

		if seconds.(int) <= 0 {
			return false
		}

		// If the store has an "add" method we will call the method on the store so it
		// has a chance to override this logic. Some drivers better support the way
		// this operation should work with a total "atomic" implementation of it.
		if store, ok := any(r.store).(contracts.AtomicStore); ok {
			return store.Add(r.itemKey(key), value, seconds)
		}
	}

	// If the value did not exist in the cache, we will put the value in the cache
	// so it exists for subsequent requests. Then, we will return true so it is
	// easy to know if the value gets added. Otherwise, we will return false.
	if r.Get(key) == nil {
		return r.Put(key, value, seconds)
	}

	return false
}

// Increment the value of an item in the cache.
func (r *Repository) Increment(key string, value ...any) any {
	return r.store.Increment(key, value...)
}

// Decrement the value of an item in the cache.
func (r *Repository) Decrement(key string, value ...any) any {
	return r.store.Decrement(key, value...)
}

// Store an item in the cache indefinitely.
func (r *Repository) Forever(key string, value any) bool {
	result := r.store.Forever(r.itemKey(key), value)

	// if result {
	//     r.event(NewKeyWritten(key, value))
	// }

	return result
}

// Get an item from the cache, or execute the given Closure and store the result.
func (r *Repository) Remember(key string, ttl any, callback func() any) any {
	value := r.Get(key)

	// If the item exists in the cache we will just return this immediately and if
	// not we will execute the given Closure and cache the result of that for a
	// given number of seconds so it's available for all subsequent requests.
	if value != nil {
		return value
	}

	value = callback()

	r.Put(key, value, helpers.Value(ttl, value))

	return value
}

// Get an item from the cache, or execute the given Closure and store the result forever.
func (r *Repository) Sear(key string, callback func() any) any {
	return r.RememberForever(key, callback)
}

// Get an item from the cache, or execute the given Closure and store the result forever.
func (r *Repository) RememberForever(key string, callback func() any) any {
	value := r.Get(key)

	// If the item exists in the cache we will just return this immediately
	// and if not we will execute the given Closure and cache the result
	// of that forever so it is available for all subsequent requests.
	if value != nil {
		return value
	}

	return helpers.Tap(callback(), func(value any) {
		r.Forever(key, value)
	})
}

// Remove an item from the cache.
func (r *Repository) Forget(key string) bool {
	return helpers.Tap(r.store.Forget(r.itemKey(key)), func(result bool) {
		// if result {
		//     r.event(NewKeyForgotten(key))
		// }
	})
}

// Delete an item from the cache by its unique key.
func (r *Repository) Delete(key string) bool {
	return r.Forget(key)
}

// Deletes multiple cache items in a single operation.
func (r *Repository) DeleteMultiple(keys []string) bool {
	result := true

	for _, key := range keys {
		if !r.Forget(key) {
			result = false
		}
	}

	return result
}

// Wipes clean the entire cache's keys.
func (r *Repository) Clear() bool {
	return r.store.Flush()
}

// Begin executing a new tags operation if the store supports it.
// func (r *Repository) Tags(names ...string) (contracts.TaggableStore, error) {
//     store, ok := any(r.store).(contracts.TaggableStore);

//     if !ok {
//         return nil, errors.New("This cache store does not support tagging.")
//     }

//     cache := store.Tags(names...)

//     if r.events != nil {
//         cache.SetEventDispatcher(r.events)
//     }

//     return cache.SetDefaultCacheTime(r.defaultCacheTime), nil
// }

// Format the key for a cache item.
func (r *Repository) itemKey(key string) string {
	return key
}

// Calculate the number of seconds for the given TTL.
func (r *Repository) getSeconds(ttl any) int {
	// TODO
	return 0
}

// Get the default cache time.
func (r *Repository) GetDefaultCacheTime() int {
	return r.defaultCacheTime
}

// Set the default cache time in seconds.
func (r *Repository) SetDefaultCacheTime(seconds int) *Repository {
	r.defaultCacheTime = seconds

	return r
}

// Get the cache store implementation.
func (r *Repository) GetStore() contracts.Store {
	return r.store
}

// Fire an event for this cache instance.
// func (r *Repository)	 event(event any) {
//     if r.events != nil {
//         r.events.dispatch(event)
//     }
// }

// Get the event dispatcher instance.
// func (r *Repository) GetEventDispatcher() contracts.Dispatcher {
//     return r.events;
// }

// Set the event dispatcher instance.
// func (r *Repository)	SetEventDispatcher(events contracts.Dispatcher)  {
//     r.events = events;
// }
