package slices

func In[T any](needle any, haystack []T) bool {
	var value any

	for _, value = range haystack {
		if needle == value {
			return true
		}
	}

	return false
}

func Get[T any](slice []T, index int, defaultValue ...any) (value any) {
	if defaultValue != nil {
		value = defaultValue[0]
	}

	if index >= 0 && index < len(slice) {
		value = slice[index]
	}

	return
}
