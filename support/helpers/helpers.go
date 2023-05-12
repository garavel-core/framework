package helpers

// Call the given Closure with the given value then return the value.
func Tap[T any](value T, callback ...func(T)) T {
	if callback == nil || callback[0] == nil {
		// TODO HigherOrderTapProxy
		return value
	}

	callback[0](value)

	return value
}

// Return the default value of the given value.
func Value(value any, args ...any) any {
	if fn, ok := value.(func(args ...any) any); ok {
		return fn(args...)
	}

	return value
}
