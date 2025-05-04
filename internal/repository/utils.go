package repository

import (
	"strconv"
	"strings"
)

func Int64SliceToString(data []int64) string {
	if len(data) == 0 {
		return ""
	}

	b := strings.Builder{}

	b.WriteString(strconv.Itoa(int(data[0])))

	for i := 1; i < len(data); i++ {
		b.WriteString(", " + strconv.Itoa(int(data[i])))
	}

	return b.String()
}
