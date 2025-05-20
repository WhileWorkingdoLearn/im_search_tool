package queries

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQL_Queries(t *testing.T) {
	query, args := BuildQuery(GenerateNgrams("HelloworldinTown", 4))
	fmt.Println(query)
	fmt.Println(args...)
	assert.NotEqual(t, "", query)
	assert.NotNil(t, "", args...)
}
