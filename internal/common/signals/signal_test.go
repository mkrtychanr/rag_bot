package signals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignals(t *testing.T) {
	assert.NotPanics(t, func() { SetupStackDump() })
	assert.NotPanics(t, func() { PrintStack() })
}
