package queries

import (
	"cmp"
	"slices"
)

// Damerau–Levenshtein distance for calculating the distance between two words.
// It calculates the nececcary substitution, insertion, deletion or transposition of a strings to be equal to the other.
// The returned value is the sum of this step. The lower the score the higher the match.
// First a matrix is generated with the length of the two strings. The x and y axis are filled with default values.
//
//	Length of S1
//
// L[[0,1,2,3,4,5,6]]
// e[[1,0,0,0,0,0,0]]
// n[[2,0,0,0,0,0,0]]
// g[[3,0,0,0,0,0,0]]
// t[[4,0,0,0,0,0,0]]
// h of S2
//
// After this the algrithm iterates over the matrix  and compares the values.
//
// more can be found under https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance
// or https://www.geeksforgeeks.org/damerau-levenshtein-distance/
//
//	input s1 string, s2 string returns int / score
func StringDistance(s1 string, s2 string) int {
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for x := 1; x <= len(s1); x++ {
		for y := 1; y <= len(s2); y++ {
			if s1[x-1] == s2[y-1] {
				matrix[x][y] = matrix[x-1][y-1]
			} else {
				matrix[x][y] = 1 + min(min(matrix[x-1][y], matrix[x][y-1]), matrix[x-1][y-1])
			}
		}
	}

	return matrix[len(s1)][len(s2)]
}
func reverseSort(a, b QueryResult) int {
	return cmp.Compare(a.Score, b.Score)
}

// Sorts the QueryResult by the Damerau–Levenshtein distance
//
//	Input is an Array of type QueryResult struct { Name    string Token   string Score   int Country string}
//
// Returns an Array of type QueryResult - Copy of input array
func SortByScore(data []QueryResult, searchTerm string) []QueryResult {
	copy := make([]QueryResult, len(data))
	for i, d := range data {
		score := StringDistance(d.Name[:min(len(searchTerm), len(d.Name))], searchTerm[:min(len(searchTerm), min(len(d.Name)))])
		d.Score += score
		copy[i] = d
	}

	slices.SortFunc(copy, reverseSort)

	return copy

}
