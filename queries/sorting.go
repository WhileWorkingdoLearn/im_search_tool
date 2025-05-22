package queries

import (
	"sort"
)

func StringDistance(s1 string, s2 string) int {
	matrix := make([][]int, len(s1)+1)
	for i, _ := range matrix {
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

func SortByScore(data []QueryResult, searchTerm string) []QueryResult {

	for _, d := range data {
		score := StringDistance(d.Name[:min(len(searchTerm), len(d.Name))], searchTerm)
		d.Score = score
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Score < data[j].Score
	})

	return data

}
