package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github/WhileCodingDoLearn/searchtool/queries"
	"github/WhileCodingDoLearn/searchtool/utils"
	"log"
	"time"

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

	queryHandler := queries.NewQueryHandler(db, 5)
	/*
		utils.LoadTable("data/DE-addresses.tsv", func(data string) error {
			err := queryHandler.InsertAddressWithNgrams(5, data, "DE")
			if err != nil {
				return err
			}
			return nil
		})
	*/

	if len(*searchTerm) >= 2 {
		defer utils.TimeTrack(time.Now(), "db access")
		start := time.Now()

		result, err := queryHandler.Search(*searchTerm)
		if err != nil {
			log.Fatal(err)
		}
		sorted := queries.SortByScore(result, *searchTerm)

		fmt.Println(sorted)
		timeSince := time.Since(start)
		fmt.Println(timeSince)

	} else {
		data, err := queryHandler.SelectAll()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
	}

}
