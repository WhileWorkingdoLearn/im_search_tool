package queries

import (
	"fmt"
	"strings"
)

type Mutations map[string]string

var InputConverter = Mutations{
	"Ä": "Ae",
	"ä": "ae",
	"Ö": "Oe",
	"ö": "oe",
	"Ü": "Ue",
	"ü": "ue",
	"ß": "ss",
	"Æ": "Ae",
	"æ": "ae",
	"Ø": "Oe",
	"ø": "oe",
	"Å": "Aa",
	"å": "aa",
	"Ç": "C",
	"ç": "c",
	"Ñ": "N",
	"ñ": "n",
}

// Funktion zur Ersetzung der Umlaute im gegebenen Text
func (mut Mutations) convertMutations(text string) string {
	for key, value := range mut {
		text = strings.ReplaceAll(text, key, value)
	}
	return text
}

type Normalizer map[string]string

var InputNormalizer = Normalizer{
	"-": "",
	"'": "",
	"ä": "ae",
	"ü": "ue",
	"ß": "ss",
}

func (norm Normalizer) normalizeString(input string) string {
	input = strings.TrimSpace(input)
	for k, v := range norm {
		input = strings.ReplaceAll(input, k, v)
	}
	return input
}

type Transpiler map[string]string

var InputTranspiler = Transpiler{
	" ": "",
	"a": "$",
	"e": "=",
	"i": "=",
	"o": "!",
	"u": "!",
	"p": "§",
	"b": "§",
	"m": "+",
	"n": "+",
	"k": "-",
	"c": "-",
}

func (t Transpiler) transpileString(input string) string {
	for k, v := range t {
		input = strings.ReplaceAll(input, k, v)
	}
	return input
}

func ProcessString(input string) string {
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)
	begin := input[0]
	tokenized := InputTranspiler.transpileString(InputNormalizer.normalizeString(InputConverter.convertMutations(input)))
	return fmt.Sprintf("%v%v", string(begin), tokenized[1:])
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
