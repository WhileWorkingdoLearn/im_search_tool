package queries

import "database/sql"

type AdressTable struct {
	Name    string
	Prefix  string
	Suffix  string
	Beginn  string
	Token   string
	Country string
}

type Sort struct {
	Address AdressTable
	Score   int
}

type Query struct {
	db *sql.DB
}

type QueryResult struct {
	Name    string
	Token   string
	Score   int
	Country string
}
