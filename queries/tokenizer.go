package queries

import (
	"strings"
)

type Normalizer map[string]string

func (norm Normalizer) normalizeString(input string) string {
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)
	for k, v := range norm {
		input = strings.ReplaceAll(input, k, v)
	}
	return input
}

var InputNormalizer = Normalizer{
	"'": "",
	"ä": "ae",
	"ü": "ue",
	"ß": "ss",
}

var InputTranspiler = Normalizer{
	"a": "$",
	"e": "=",
	"i": "=",
	"o": "!",
	"u": "!",
	"p": "§",
	"b": "§",
}

func (norm Normalizer) transpileString(input string) string {
	input = strings.ToLower(input)
	for k, v := range norm {
		input = strings.ReplaceAll(input, k, v)
	}
	return input
}

func CountVowels(s string) int {
	vowels := "aeiouAEIOU"
	count := 0
	for _, char := range s {
		if strings.ContainsRune(vowels, char) {
			count++
		}
	}
	return count
}

// countConsonants zählt die Anzahl der Konsonanten im String
func CountConsonants(s string) int {
	consonants := "bcdfghjklmnpqrstvwxyz"
	count := 0
	for _, char := range s {
		if strings.ContainsRune(consonants, char) {
			count++
		}
	}
	return count
}

func GenerateNGrams(input string, n int) string {
	// Remove spaces and normalize the string
	input = strings.ReplaceAll(input, " ", "")
	length := len(input)

	if n > length || n <= 0 {
		return input
	}

	ngrams := []string{}
	for i := 0; i <= length-n; i++ {
		ngrams = append(ngrams, input[i:i+n])
	}

	return strings.Join(ngrams, ",")
}
