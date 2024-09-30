package internal_test

import (
	"testing"

	"github.com/quickfixgo/quickfix/internal"
	"github.com/stretchr/testify/assert"
)

func TestOnceReset(t *testing.T) {
	var calls int
	var once internal.Once
	once.Do(func() {
		calls++
	})
	once.Do(func() {
		calls++
	})
	assert.EqualValues(t, 1, calls)
	once.Reset()
	once.Do(func() {
		calls++
	})
	once.Do(func() {
		calls++
	})
	once.Do(func() {
		calls++
	})
	assert.EqualValues(t, 2, calls)
}
