package queries

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mattn/go-sqlite3"
)

const dropTable = `DROP TABLE address`

const createDB = `CREATE TABLE IF NOT EXISTs address (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	prefix TEXT NOT NULL,
	suffix TEXT NOT NULL,
	begin INTEGER NOT NULL,
	token TEXT NOT  NULL,
	country TEXT NO NULL
);`

const insertDB = `INSERT INTO address (name, prefix,suffix,begin, token, country) VALUES (?,?,?,?,?,?)`

const selectAllFromDB = `SELECT name,prefix,suffix,begin, token, country FROM address;`
const selectbyPrefix = `SELECT name, prefix,suffix,begin, token, country FROM address WHERE begin = ? AND (prefix = ? OR suffix = ?);`

type Query struct {
	db *sql.DB
}

const (
	Prefix = 4
	Suffix = 3
)

type AdressTable struct {
	Name    string
	Prefix  string
	Suffix  string
	Beginn  string
	Token   string
	Country string
}

func NewQueryHandler(db *sql.DB) *Query {
	return &Query{db: db}
}

func (q *Query) CreateTable() (sql.Result, error) {
	return q.db.Exec(createDB)
}

func (q *Query) DropTable() {
	q.db.Exec(dropTable)
}

func (q *Query) Instert(data []AdressTable, ngramCount int) error {
	for _, data := range data {
		prefix := InputTranspiler.transpileString(InputNormalizer.normalizeString(data.Name[:Prefix]))
		suffix := InputTranspiler.transpileString(InputNormalizer.normalizeString(data.Name[len(data.Name)-Suffix:]))
		_, err := q.db.Exec(insertDB, data.Name, prefix, suffix, prefix[0], GenerateNGrams(InputNormalizer.normalizeString(data.Name[Prefix:]), ngramCount), data.Country)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *Query) SelectAll() ([]AdressTable, error) {
	table := make([]AdressTable, 0)
	rows, err := q.db.Query(selectAllFromDB)
	if err != nil {
		log.Fatalf("Fehler beim Ausführen der Query: %v", err)
	}
	defer rows.Close()

	// Ergebnisse ausgeben
	for rows.Next() {
		var address = AdressTable{}
		if err := rows.Scan(&address.Name, &address.Prefix, &address.Suffix, &address.Beginn, &address.Token, &address.Country); err != nil {
			return table, err
		}
		table = append(table, address)
	}
	if err := rows.Err(); err != nil {
		return table, err
	}
	return table, nil

}

func (q *Query) Search(query string) ([]AdressTable, error) {
	prefix := InputTranspiler.transpileString(InputNormalizer.normalizeString(query[:Prefix]))
	suffix := InputTranspiler.transpileString(InputNormalizer.normalizeString(query[len(query)-Suffix:]))
	table := make([]AdressTable, 0)
	rows, err := q.db.Query(selectbyPrefix, prefix[0], prefix, suffix)
	if err != nil {
		log.Fatalf("Fehler beim Ausführen der Query: %v", err)
	}
	defer rows.Close()

	// Ergebnisse ausgeben
	for rows.Next() {
		var address = AdressTable{}
		if err := rows.Scan(&address.Name, &address.Prefix, &address.Suffix, &address.Beginn, &address.Token, &address.Country); err != nil {
			return table, err
		}
		table = append(table, address)
	}
	if err := rows.Err(); err != nil {
		return table, err
	}
	return table, nil

}

func (q *Query) RegisterNgram() error {
	ctx := context.Background()
	conn, err := q.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Raw(func(dc interface{}) error {
		sqliteConn, ok := dc.(*sqlite3.SQLiteConn)
		if !ok {
			return fmt.Errorf("unerwarteter Typ: %T", dc)
		}

		return sqliteConn.RegisterFunc("ngram", func(text string, n int64) string {
			return GenerateNGrams(text, int(n))
		}, true)
	})
	if err != nil {
		fmt.Println("Fehler beim Registrieren der Funktion:", err)
		return err
	}

	return nil
}

func BuildQuery(searcktokens string) (string, []interface{}) {
	searchTokens := strings.Split(searcktokens, ", ")
	conditions := []string{}
	args := []interface{}{}
	for _, token := range searchTokens {
		conditions = append(conditions, "token LIKE ?")
		args = append(args, "%"+token+"%")
	}
	whereClause := strings.Join(conditions, " OR ")

	scoreParts := []string{}
	for range searchTokens {
		scoreParts = append(scoreParts, "CASE WHEN token LIKE ? THEN 1 ELSE 0 END")
	}

	scoreParams := []interface{}{}
	for _, token := range searchTokens {
		scoreParams = append(scoreParams, "%"+token+"%")
	}

	query := fmt.Sprintf(`
        SELECT id, name, token, country,
            (%s) as score
        FROM address
        WHERE %s
        ORDER BY score DESC
    `, strings.Join(scoreParts, " + "), whereClause)

	args = append(args, scoreParams...)

	return query, args
}

func (q *Query) FindAdresses(input string, args []interface{}) error {

	// Query ausführen
	rows, err := q.db.Query(input, args...)
	if err != nil {
		log.Fatalf("Fehler beim Ausführen der Query: %v", err)
	}
	defer rows.Close()

	// Ergebnisse ausgeben
	for rows.Next() {
		var id int
		var name, token, country string
		var score int
		if err := rows.Scan(&id, &name, &token, &country, &score); err != nil {
			return err
		}
		fmt.Printf("ID: %d, Name: %s, Score: %d, Country: %s\n", id, name, score, country)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
