package queries

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

const (
	Prefix = 3
	Suffix = 3
)

const SelectNgram = `SELECT id FROM ngrams_dict WHERE ngram = ?`

func NewQueryHandler(db *sql.DB, ngramSize int) *Query {
	return &Query{db: db, ngramSize: ngramSize}
}

func (q *Query) SelectAll() ([]AdressTable, error) {
	table := make([]AdressTable, 0)
	rows, err := q.db.Query("select * from address")
	if err != nil {
		log.Fatalf("Fehler beim Ausführen der Query: %v", err)
	}
	defer rows.Close()

	// Ergebnisse ausgeben
	for rows.Next() {
		var address = AdressTable{}
		id := 0
		if err := rows.Scan(&id, &address.Name, &address.Country); err != nil {
			return table, err
		}
		table = append(table, address)
	}
	if err := rows.Err(); err != nil {
		return table, err
	}
	return table, nil

}

func anySlice(strs []string) []any {
	out := make([]any, len(strs))
	for i, v := range strs {
		out[i] = v
	}
	return out
}

func (q *Query) Search(input string) ([]QueryResult, error) {
	token := ProcessString(input)
	ngram := GenerateNGrams(token, 5)
	placeholders := strings.Repeat("?,", len(ngram))
	placeholders = placeholders[:len(placeholders)-1] // remove trailing comma

	query := fmt.Sprintf(`
	SELECT a.name, a.country, COUNT(*) AS score
	FROM ngrams_dict nd
	JOIN string_ngrams sg ON nd.id = sg.ngram_id
	JOIN address a ON a.id = sg.string_id
	WHERE nd.ngram IN (%s)
	GROUP BY a.id, a.name, a.country
	ORDER BY score DESC
	LIMIT 10;
`, placeholders)

	rows, err := q.db.Query(query, anySlice(ngram)...)
	if err != nil {
		return nil, err
	}

	Adresses := make([]QueryResult, 0)
	for rows.Next() {
		found := QueryResult{}
		rows.Scan(&found.Name, &found.Country, &found.Score)
		Adresses = append(Adresses, found)
	}

	return Adresses, nil

}
