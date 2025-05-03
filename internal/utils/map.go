package utils

func MapKeysToSlice[K comparable, V any](d map[K]V) []K {
	r := make([]K, 0, len(d))

	for k := range d {
		r = append(r, k)
	}

	return r
}

func MapValuesToSlice[K comparable, V any](d map[K]V) []V {
	r := make([]V, 0, len(d))

	for _, v := range d {
		r = append(r, v)
	}

	return r
}

func SliceToMap[K comparable](d []K) map[K]struct{} {
	r := make(map[K]struct{}, len(d))

	for _, v := range d {
		r[v] = struct{}{}
	}

	return r
}
