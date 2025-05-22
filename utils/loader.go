package utils

import (
	"bufio"
	"github/WhileCodingDoLearn/searchtool/queries"
	"log"
	"os"
	"strings"
)

func Handler(queryHandler *queries.Query) func(data string) error {
	return func(data string) error {
		err := queryHandler.Instert(queries.AdressTable{Name: data, Country: "DE"})
		return err
	}
}

/*"data/DE-addresses.tsv"*/
func LoadTable(path string, queryHandler func(data string) error) {

	inputFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	limit := 0
	for scanner.Scan() {
		line := scanner.Text()
		column := strings.Split(line, "\t")
		if len(column) >= 2 {
			if len(column[2]) > 3 {
				errInsert := queryHandler(column[2])
				if errInsert != nil {
					continue
				}
				limit++
			}
		}
		if limit >= 1000 {
			break
		}
	}

}
