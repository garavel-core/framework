// From:
//  - https://www.php.net/manual/en/book.array.php
//  - https://github.com/laravel/framework/blob/9.x/src/Illuminate/Collections/Arr.php

package arr

import (
	"reflect"
	"strings"

	"github.com/garavel-core/framework/support/ctype"
	"github.com/garavel-core/framework/support/str"
)

type Array[K comparable, V any] map[K]V

const (
	CASE_LOWER = iota
	CASE_UPPER
)

// Changes the case of all keys in an array
func (a Array[K, V]) ChangeKeyCase(cased ...int) Array[K, V] {
	// 设置参数的默认值
	if len(cased) == 0 {
		cased = append(cased, CASE_LOWER)
	}

	result := make(Array[K, V], a.Count())

	for key, value := range a {
		// 如果是字符串键就进行转换
		// 否则不转换键
		if k, ok := any(key).(string); ok {
			if cased[0] == CASE_UPPER {
				key = any(strings.ToUpper(k)).(K)
			} else {
				key = any(strings.ToLower(k)).(K)
			}
		}

		result[key] = value
	}

	return result
}

// Split an array into chunks
func (a Array[K, V]) Chunk(length int, preserveKeys ...bool) []Array[K, V] {
	// 设置参数的默认值
	if len(preserveKeys) == 0 {
		preserveKeys = append(preserveKeys, false)
	}

	if length < 1 {
		panic("Array.Chunk(): Argument #1 (length) must be greater than 0")
	}

	result := make([]Array[K, V], a.Count()/length)
	i := 0
	j := 1

	for key, value := range a {
		result[i][key] = value

		if j%length == 0 {
			i++
		}

		j++
	}

	return result
}

// Checks if the given key or index exists in the array.
func (a Array[K, V]) Exists(key K) bool {
	_, exists := a[key]
	return exists
}

// Return the values from a single column in the input array
func (a Array[K, V]) Column(columnKey any, indexKey ...K) Array[any, any] {
	var value any

	result := make(Array[any, any])
	i := 0
	for _, v := range a {
		if columnKey == nil {
			value = v
		} else {
			value, _ = Get(v, columnKey)
		}

		if key, exists := Get(v, indexKey); exists {
			result[key] = value
		} else {
			result[i] = value
			i++
		}
	}

	return result
}

func (a Array[K, V]) Count() int {
	return len(a)
}

func Keys[K comparable, V any](array map[K]V) []K {
	keys := make([]K, len(array))
	i := 0

	for key := range array {
		keys[i] = key
		i++
	}

	return keys
}

func Map[K comparable, V any](array map[K]V, callback func(V, K) V) map[K]V {
	for key, value := range array {
		array[key] = callback(value, key)
	}

	return array
}

func In[K comparable, V any](needle any, haystack map[K]V) bool {
	var value any

	for _, value = range haystack {
		if needle == value {
			return true
		}
	}

	return false
}

// Get the value of the given key in a collection using reflection. Returns nil if the array is not a collection.
func Get(array any, key any) (any, bool) {
	if key == nil || array == nil {
		return nil, false
	}

	t := reflect.ValueOf(array)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Map:
		if v := t.MapIndex(reflect.ValueOf(key)); v.IsValid() {
			return v.Interface(), true
		}
	case reflect.Array, reflect.Slice:
		if k, ok := key.(int); ok && k > -1 && k < t.Len() {
			return t.Index(k).Interface(), true
		}
	case reflect.Struct:

		if k, ok := key.(string); ok {
			// 如果全是小写则转成大写驼峰形式，因为我们仅获取导出的字段
			if ctype.Lower(k) {
				k = str.Studly(k)
			}

			if v := t.FieldByName(k); v.IsValid() {
				return v.Interface(), true
			}
		}
	}

	return nil, false
}
