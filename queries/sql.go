package queries

import (
	"database/sql"
	"fmt"
)

const CreateAddressTable = `CREATE TABLE IF NOT EXISTS address (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	country TEXT NO NULL
);`

const CreateNgramTable = `CREATE TABLE IF NOT EXISTS ngrams_dict (
            id INTEGER PRIMARY KEY,
            ngram TEXT NOT NULL UNIQUE
        );`

const CreateJoinTable = `CREATE TABLE IF NOT EXISTS string_ngrams (
            string_id INTEGER NOT NULL,
            ngram_id INTEGER NOT NULL,
            PRIMARY KEY (string_id, ngram_id),
            FOREIGN KEY (string_id) REFERENCES strings(id),
            FOREIGN KEY (ngram_id) REFERENCES ngrams_dict(id)
        );`

const CreateIndex = `CREATE INDEX IF NOT EXISTS idx_ngrams_dict_ngram ON ngrams_dict (ngram);`

const DropAddressTable = `DROP TABLE IF EXISTS address`

const DropJoinTable = `DROP TABLE IF EXISTS ngrmas_dict`

const DropNgramTable = `DROP TABLE IF EXISTS ngrmas_dict`

const InsertAdressIntoTable = `INSERT INTO address (name, country) VALUES (?, ?)`

const InsertNgrmaIntoTable = `INSERT INTO ngrams_dict (ngram) VALUES (?)`

func (q *Query) CreateTables() error {

	tables := []string{CreateAddressTable, CreateNgramTable, CreateJoinTable, CreateIndex}
	for _, table := range tables {
		_, err := q.db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *Query) DropTables() error {
	tables := []string{DropAddressTable, DropJoinTable, DropNgramTable}
	for _, table := range tables {
		_, err := q.db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *Query) Insert(name, country string) error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert into address table
	res, err := tx.Exec(InsertAdressIntoTable, name, country)
	if err != nil {
		return fmt.Errorf("insert address: %w", err)
	}

	addressID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	// 2. Generate trigrams (or other n-grams)
	token := ProcessString(name)
	nGrams := GenerateNGrams(token, q.ngramSize)

	for _, gram := range nGrams {
		var ngramID int64

		// Check if ngram already exists
		err = tx.QueryRow(SelectNgram, gram).Scan(&ngramID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Insert new ngram
				res, err := tx.Exec(InsertNgrmaIntoTable, gram)
				if err != nil {
					return fmt.Errorf("insert ngram '%s': %w", gram, err)
				}
				ngramID, err = res.LastInsertId()
				if err != nil {
					return fmt.Errorf("last insert id for ngram '%s': %w", gram, err)
				}
			} else {
				return fmt.Errorf("select ngram '%s': %w", gram, err)
			}
		}

		// Insert into join table
		_, err = tx.Exec(`INSERT OR IGNORE INTO string_ngrams (string_id, ngram_id) VALUES (?, ?)`, addressID, ngramID)
		if err != nil {
			return fmt.Errorf("insert into join table: %w", err)
		}
	}

	return tx.Commit()
}
