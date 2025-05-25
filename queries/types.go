package queries

import "database/sql"

type AdressTable struct {
	Name    string
	Country string
}

type Sort struct {
	Address AdressTable
	Score   int
}

type Query struct {
	db        *sql.DB
	ngramSize int
}

type QueryResult struct {
	Name    string
	Token   string
	Score   int
	Country string
}
