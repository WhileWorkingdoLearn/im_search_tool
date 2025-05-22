package queries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSorting(t *testing.T) {
	dist := StringDistance("Hello World", "Hello World")
	assert.Equal(t, 0, int(dist))
	dist2 := StringDistance("Hello World", "Hello")
	assert.Equal(t, 1, int(dist2))
}
