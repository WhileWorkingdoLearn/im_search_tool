package main

import (
	"database/sql"
	"fmt"
	"github/WhileCodingDoLearn/searchtool/queries"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Hell world")
	db, err := sql.Open("sqlite3", "database/geodb.db")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	fmt.Println("Connected to the Database")
	var sqlVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqlVersion)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sqlVersion)
	queryHandler := queries.NewQueryHandler(db)
	_, err = queryHandler.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	errRegister := queryHandler.RegisterNgram()
	if errRegister != nil {
		log.Fatal(errRegister)
	}

	testData := []queries.AdressTable{
		{Name: "Berlin Hauptbahnhof", Token: "", Country: "DE"},
		{Name: "München Marienplatz", Token: "", Country: "DE"},
		{Name: "Paris Gare du Nord", Token: "", Country: "FR"},
		{Name: "London King's Cross", Token: "", Country: "GB"},
		{Name: "New York Penn Station", Token: "", Country: "US"},
	}
	queryHandler.Instert(testData)

	fmt.Println("Testdaten wurden erfolgreich in die Datenbank eingefügt.")
}

/*

    // Beispielhafter Suchstring – dieser könnte z.B. von einer HTTP-Request stammen.
    suchString := "dein Suchbegriff" // Ersetze dies durch den Input, den du erhältst.

    // Erzeuge N-Gramme aus dem Suchstring (hier nehmen wir 2-Gramme wie beim Update).
    searchNgramsStr := generateNgrams(suchString, 2)
    searchTokens := strings.Split(searchNgramsStr, ", ")

    // Dynamisch eine SQL-Query zusammenbauen.
    // Zuerst bauen wir eine WHERE-Bedingung, die prüft, ob mindestens ein N-Gramm im Token-Feld enthalten ist.
    conditions := []string{}
    args := []interface{}{}
    for _, token := range searchTokens {
        conditions = append(conditions, "token LIKE ?")
        args = append(args, "%"+token+"%")
    }
    whereClause := strings.Join(conditions, " OR ")

    // Zusätzlich berechnen wir einen Score: Für jedes N-Gramm, das gefunden wurde, erhöhen wir den Zähler.
    scoreParts := []string{}
    for range searchTokens {
        scoreParts = append(scoreParts, "CASE WHEN token LIKE ? THEN 1 ELSE 0 END")
    }
    // Die Parameter für die Score-Berechnung sind dieselben wie für die WHERE-Bedingung.
    scoreParams := []interface{}{}
    for _, token := range searchTokens {
        scoreParams = append(scoreParams, "%"+token+"%")
    }

    // Die finale Query kombiniert beide Ansätze und sortiert anschließend nach Score absteigend.
    query := fmt.Sprintf(`
        SELECT id, name, token, country,
            (%s) as score
        FROM address
        WHERE %s
        ORDER BY score DESC
    `, strings.Join(scoreParts, " + "), whereClause)

    // Alle Parameter werden zusammengesetzt: zuerst die Bedingungen und dann die Score-Parameter.
    args = append(args, scoreParams...)

    // Query ausführen
    rows, err := db.Query(query, args...)
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
            log.Fatalf("Fehler beim Scannen der Zeile: %v", err)
        }
        fmt.Printf("ID: %d, Name: %s, Score: %d, Country: %s\n", id, name, score, country)
    }
    if err := rows.Err(); err != nil {
        log.Fatalf("Fehler nach dem Iterieren der Zeilen: %v", err)
    }
}

*/
