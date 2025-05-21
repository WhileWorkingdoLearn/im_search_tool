package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github/WhileCodingDoLearn/searchtool/queries"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	searchTerm := flag.String("s", "", "Input value for search")
	flag.Parse()

	db, err := sql.Open("sqlite3", "database/geodb.db")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	var sqlVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqlVersion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the Database: SQLite3 -v: ", sqlVersion)

	queryHandler := queries.NewQueryHandler(db)
	queryHandler.DropTable()

	_, err = queryHandler.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	/*
		errRegister := queryHandler.RegisterNgram()
		if errRegister != nil {
			log.Fatal(errRegister)
		}
	*/

	testData := []queries.AdressTable{
		{Name: "Berlin", Token: "", Country: "DE"},
		{Name: "Bertin", Token: "", Country: "DE"},
		{Name: "Berllin", Token: "", Country: "DE"},
		{Name: "Hauptbahnhof", Token: "", Country: "DE"},
		{Name: "München", Token: "", Country: "DE"},
		{Name: "Marienplatz", Token: "", Country: "DE"},
		{Name: "Gare du Nord", Token: "", Country: "FR"},
		{Name: "Paris", Token: "", Country: "FR"},
		{Name: "London", Token: "", Country: "GB"},
		{Name: "King's Cross", Token: "", Country: "GB"},
		{Name: "Penn Station", Token: "", Country: "US"},
		{Name: "New York", Token: "", Country: "US"},
	}
	errInsert := queryHandler.Instert(testData, 3)
	if errInsert != nil {
		log.Fatal(errInsert)
	}

	if len(*searchTerm) < 3 {
		data, err := queryHandler.Search(*searchTerm)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
	} else {
		data, err := queryHandler.SelectAll()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
	}

}
