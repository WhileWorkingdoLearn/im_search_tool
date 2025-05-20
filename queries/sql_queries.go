package queries

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mattn/go-sqlite3"
)

const CreateDB = `CREATE TABLE IF NOT EXISTs address (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	token TEXT NOT  NULL,
	country TEXT NO NULL
);`

const insertDB = `INSERT INTO address (name, token, country) VALUES (?, ngram(?, 4), ?)`

type Query struct {
	db *sql.DB
}

type AdressTable struct {
	Name    string
	Token   string
	Country string
}

func NewQueryHandler(db *sql.DB) *Query {
	return &Query{db: db}
}

func (q *Query) CreateTable() (sql.Result, error) {
	return q.db.Exec(CreateDB)
}

func (q *Query) Instert(data []AdressTable) error {
	for _, data := range data {
		_, err := q.db.Exec(insertDB, data.Name, data.Token, data.Country)
		if err != nil {
			return err
		}
	}
	return nil
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
			return GenerateNgrams(text, int(n))
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
