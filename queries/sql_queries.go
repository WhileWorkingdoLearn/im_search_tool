package queries

import (
	"database/sql"
	"log"
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

const selectAllFromDB = `SELECT name,prefix,suffix,begin, token, country FROM address;
`
const selectbyPrefix = `SELECT name, prefix,suffix,begin, token, country FROM address WHERE begin = ? AND (prefix = ? OR suffix = ?);`

const (
	Prefix = 3
	Suffix = 3
)

func NewQueryHandler(db *sql.DB) *Query {
	return &Query{db: db}
}

func (q *Query) CreateTable() (sql.Result, error) {
	return q.db.Exec(createDB)
}

func (q *Query) DropTable() {
	q.db.Exec(dropTable)
}

func (q *Query) InstertMany(input []AdressTable) error {
	for _, data := range input {
		token := ProcessString(data.Name)
		if len(token) >= Prefix {
			prefix := token[:Prefix]
			suffix := token[:Prefix][len(token[:Prefix])-Suffix:]
			_, err := q.db.Exec(
				insertDB,     // Query
				data.Name,    // StreetName Human Readable
				prefix,       // normalized Token Chunk Prefix
				suffix,       // normalized Token Chunk Suffix
				prefix[0],    // first letter of prefix
				token,        // normalized TOken of Streetname
				data.Country) // countrycode
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (q *Query) Instert(data AdressTable) error {
	token := ProcessString(data.Name)
	if len(token) >= Prefix {
		prefix := token[:Prefix]
		suffix := token[:Prefix][len(token[:Prefix])-Suffix:]
		_, err := q.db.Exec(
			insertDB,     // Query
			data.Name,    // StreetName Human Readable
			prefix,       // normalized Token Chunk Prefix
			suffix,       // normalized Token Chunk Suffix
			prefix[0],    // first letter of prefix
			token,        // normalized TOken of Streetname
			data.Country) // countrycode
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

func (q *Query) Search(query string) ([]QueryResult, error) {
	table := make([]QueryResult, 0)
	token := ProcessString(query)
	if len(token) >= Prefix {
		prefix := token[:Prefix]
		suffix := token[:Prefix][len(token[:Prefix])-Suffix:]

		rows, err := q.db.Query(selectbyPrefix, prefix[0], prefix, suffix)
		if err != nil {
			log.Fatalf("Fehler beim Ausführen der Query: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var address = QueryResult{}
			placeholder := ""
			if err := rows.Scan(&address.Name, placeholder, placeholder, placeholder, &address.Token, &address.Country); err != nil {
				return table, err
			}
			table = append(table, address)
		}
		if err := rows.Err(); err != nil {
			return table, err
		}
	} else {
		return nil, nil
	}
	return table, nil
}
