package queries

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	ngrma := GenerateNGrams(InputNormalizer.normalizeString("Hallo Da Draußen"), 4)
	fmt.Println(ngrma)
	assert.NotNil(t, ngrma)
	assert.Equal(t, "Ha,al,ll,lo,oD,Da,aD,Dr,ra,au,us,ss,se,en", ngrma)
}

func TestProcessor(t *testing.T) {
	result := ProcessString("Hallo Da Draußen")
	assert.Len(t, result, 17)
	assert.Equal(t, result, "h$ll!d$dr$!ss=+")
}
