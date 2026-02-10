package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// global database connection
var db *sql.DB
var verbose = false

func main() {
	if len(os.Args) >= 1 && os.Args[1] == "--verbose" {
		verbose = true
	}

	db_conn, err := sql.Open("sqlite3", "./notes.db") // store in notes.db
	if err != nil {
		panic(err)
	}

	defer db_conn.Close() // close connection at the end

	db = db_conn

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// create the table
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY,
			author TEXT NOT NULL,
			text TEXT NOT NULL,
			edited INTEGER NOT NULL DEFAULT 0
		);`,
	)
	if err != nil {
		panic(err)
	}
	if verbose {
		fmt.Printf("database started\n")
	}

	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/notes", handleNotes)
	fs := http.FileServer(http.Dir("./frontend-react/dist/"))
	http.Handle("/", fs)

	if verbose {
		fmt.Printf("serving: http://localhost:8080/\n")
	}

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
