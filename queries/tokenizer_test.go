package queries

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	ngrma := GenerateNgrams("Hallo Da Draußen", 2)
	fmt.Println(ngrma)
	assert.NotNil(t, ngrma)
	assert.Equal(t, ngrma, "hallo da, da draußen")
}
