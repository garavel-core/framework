// 解决代码之间的环形依赖

package maps

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, len(m))
	i := 0

	for _, v := range m {
		values[i] = v
		i++
	}

	return values
}
