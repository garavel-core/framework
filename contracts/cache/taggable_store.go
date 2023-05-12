package cache

type TaggableStore interface {
	// Begin executing a new tags operation.
	Tags(names ...string)
}
