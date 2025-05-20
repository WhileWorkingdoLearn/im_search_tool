package queries

import (
	"strings"
)

func GenerateNgrams(text string, n int) string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	if len(words) < n {
		return text // Falls zu wenige Wörter, gib einfach den Text zurück.
	}
	var ngrams []string
	for i := 0; i <= len(words)-n; i++ {
		gram := strings.Join(words[i:i+n], " ")
		ngrams = append(ngrams, gram)
	}
	return strings.Join(ngrams, ", ")
}
