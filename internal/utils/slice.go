package utils

func GetUniqueSlice[K comparable](d []K) []K {
	return MapKeysToSlice(SliceToMap(d))
}
