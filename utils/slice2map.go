package utils

// SliceToMap converts a slice to a map using a key function to extract the key from each element.
func SliceToMap[T any, V comparable](src []T, key func(T) V) map[V]T {
	var result = make(map[V]T)
	for _, v := range src {
		result[key(v)] = v
	}
	return result
}
