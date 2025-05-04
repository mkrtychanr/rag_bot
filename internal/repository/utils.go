package repository

func Int64SliceToPGXArray(data []int64) []interface{} {
	result := make([]interface{}, 0, len(data))

	for _, v := range data {
		result = append(result, v)
	}

	return result
}
