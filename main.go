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

	queryHandler := queries.NewQueryHandler(db)
	//queryHandler.DropTable()
	//utils.LoadTable("data/example.tsv", utils.Handler(queryHandler))

	_, err = queryHandler.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	if len(*searchTerm) >= 3 {
		defer utils.TimeTrack(time.Now(), "db access")
		start := time.Now()
		data, err := queryHandler.Search(*searchTerm)
		if err != nil {
			fmt.Println(err)
		}
		temr := queries.ProcessString(*searchTerm)
		fmt.Println(temr)
		sorted := queries.SortByScore(data[:], temr)
		for _, s := range sorted {
			fmt.Println(s.Name)
		}
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
