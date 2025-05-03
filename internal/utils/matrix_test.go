package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrix(t *testing.T) {
	m := NewMatrix[int](3)

	for i := range 9 {
		m.Add(i)
	}

	expected := [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}

	assert.Equal(t, expected, m.GetMatrix())
}
