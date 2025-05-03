package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver := ""

	assert.NotPanics(t, func() { ver = BuildVersionString("test") })
	assert.NotEmpty(t, ver)

	Version = "test"

	assert.NotPanics(t, func() { ver = BuildVersionString("test") })
	assert.NotEmpty(t, ver)
}
