package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// generateNGrams generiert n-Gramme aus dem übergebenen String s
// unter Verwendung von Padding. Beispiel: s = "alice", n = 3
// ergibt: ["$$a", "$al", "ali", "lic", "ice", "ce$", "e$$"]
func generateNGrams(s string, n int) []string {
	pad := strings.Repeat("$", n-1)
	padded := pad + s + pad
	var ngrams []string
	for i := 0; i <= len(padded)-n; i++ {
		ngrams = append(ngrams, padded[i:i+n])
	}
	return ngrams
}

// insertStringAndNgrams fügt einen neuen String in die strings-Tabelle ein,
// berechnet die eindeutigen n-Gramme (mittels generateNGrams) und befüllt
// die normalisierten Tabellen ngrams_dict (Wörterbuch) sowie die Junction-Tabelle
// string_ngrams, die den Zusammenhang zwischen String und n-gram abbildet.
func insertStringAndNgrams(db *sql.DB, input string, n int) error {
	// Normalisiere den String, beispielsweise in Kleinbuchstaben.
	normalized := strings.ToLower(input)

	// Starte eine Transaktion, um alle Inserts atomar durchzuführen.
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 1. Füge den normalisierten String in die Tabelle `strings` ein.
	res, err := tx.Exec("INSERT INTO strings (value) VALUES (?)", normalized)
	if err != nil {
		tx.Rollback()
		return err
	}

	stringID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Generiere die n-Gramme für den String.
	grams := generateNGrams(normalized, n)
	// Nutze eine Map, um Dopplungen innerhalb eines Strings zu vermeiden.
	seenGrams := make(map[string]bool)

	// Bereite die Statements vor:
	// a) Insert-Statement für das Wörterbuch (ngrams_dict). Mit "INSERT OR IGNORE"
	//    wird vermieden, dass bereits vorhandene n-Gramme erneut eingefügt werden.
	stmtInsertDict, err := tx.Prepare("INSERT OR IGNORE INTO ngrams_dict (ngram) VALUES (?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmtInsertDict.Close()

	// b) Selektions-Statement, um die ID eines n-Gramms aus ngrams_dict zu erhalten.
	stmtSelectDict, err := tx.Prepare("SELECT id FROM ngrams_dict WHERE ngram = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmtSelectDict.Close()

	// c) Einfügen in die Junction-Tabelle, die den String mit den n-Grammen verbindet.
	stmtInsertMapping, err := tx.Prepare("INSERT INTO string_ngrams (string_id, ngram_id) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmtInsertMapping.Close()

	// Für jedes (einmalige) n-Gramm: Einfügen in ngrams_dict und in string_ngrams.
	for _, gram := range grams {
		if seenGrams[gram] {
			continue
		}
		seenGrams[gram] = true

		// Einfügen in das Wörterbuch. Ist das n-Gramm bereits vorhanden,
		// sorgt INSERT OR IGNORE dafür, dass kein Fehler generiert wird.
		if _, err := stmtInsertDict.Exec(gram); err != nil {
			tx.Rollback()
			return err
		}

		// Abrufen der ID des n-Gramms.
		var ngramID int64
		if err := stmtSelectDict.QueryRow(gram).Scan(&ngramID); err != nil {
			tx.Rollback()
			return err
		}

		// Einfügen in die Junction-Tabelle.
		if _, err := stmtInsertMapping.Exec(stringID, ngramID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Bei Erfolg wird die Transaktion committed.
	return tx.Commit()
}

func main() {
	// Verbindung zur SQLite-Datenbank herstellen. Passe den Pfad ggf. an.
	db, err := sql.Open("sqlite3", "../database/example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Erstelle die nötigen Tabellen, falls sie noch nicht existieren.
	createStatements := []string{
		`DROP TABLE IF EXISTS strings`,
		`DROP TABLE IF EXISTS ngrams_dict`,
		`DROP TABLE IF EXISTS string_ngrams`,
		`CREATE TABLE IF NOT EXISTS strings (
            id INTEGER PRIMARY KEY,
            value TEXT NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS ngrams_dict (
            id INTEGER PRIMARY KEY,
            ngram TEXT NOT NULL UNIQUE
        );`,
		`CREATE TABLE IF NOT EXISTS string_ngrams (
            string_id INTEGER NOT NULL,
            ngram_id INTEGER NOT NULL,
            PRIMARY KEY (string_id, ngram_id),
            FOREIGN KEY (string_id) REFERENCES strings(id),
            FOREIGN KEY (ngram_id) REFERENCES ngrams_dict(id)
        );`,
		`CREATE INDEX IF NOT EXISTS idx_ngrams_dict_ngram ON ngrams_dict (ngram);`,
	}
	for _, stmt := range createStatements {
		if _, err := db.Exec(stmt); err != nil {
			log.Fatalf("Fehler beim Erstellen der Tabelle: %v", err)
		}
	}

	// Beispiel: Füge einige Teststrings ein.
	testStrings := []string{"Alice", "Alias", "Andrew", "Beth", "Bob", "Charlie", "Natalie"}
	n := 3

	for _, s := range testStrings {
		if err := insertStringAndNgrams(db, s, n); err != nil {
			log.Fatalf("Fehler beim Einfügen von %s: %v", s, err)
		} else {
			fmt.Printf("Eingefügt: %s\n", s)
		}
	}

	fmt.Println("Alle Testdaten wurden erfolgreich eingefügt.")

	scan(db)
}

func scan(db *sql.DB) {

	searchNgrams := []string{"ali"} // Beispiel für "alice"

	// Wir erstellen den Platzhalter-String für die IN-Klausel:
	// ("?", "?", "?", ...) abhängig von der Anzahl an n-Grammen.
	placeholders := make([]string, len(searchNgrams))
	for i := range searchNgrams {
		placeholders[i] = "?"
	}
	inClause := "(" + strings.Join(placeholders, ", ") + ")"

	// Definiere den Schwellenwert. Dieser muss natürlich aus deiner Applikationslogik oder
	// aus der Rechnung der Ungleichung berechnet werden. Hier als Beispiel:
	threshold := 1

	// Aufbau der Query. Diese Query ermittelt, wie oft jeder String ein n-Gramm mit dem Query teilt.
	query := fmt.Sprintf(`SELECT s.id, s.value, COUNT(*) AS common_ngrams 
	FROM string_ngrams AS sn 
	JOIN ngrams_dict AS nd ON sn.ngram_id = nd.id 
	JOIN strings AS s ON sn.string_id = s.id 
	WHERE nd.ngram IN %s GROUP BY s.id 
	HAVING common_ngrams >= ? 
	ORDER BY common_ngrams DESC;`, inClause)

	// Kombiniere die Argumente: zuerst die n-Gramme, dann als letztes den Schwellenwert.
	args := make([]interface{}, 0, len(searchNgrams)+1)
	for _, gram := range searchNgrams {
		args = append(args, gram)
	}
	args = append(args, threshold)

	// Führe die Query aus.
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iteriere über die Ergebnismenge.
	for rows.Next() {
		var id int
		var value string
		var commonNgrams int

		if err := rows.Scan(&id, &value, &commonNgrams); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Wert: %s, Gemeinsame n-Gramme: %d\n", id, value, commonNgrams)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
