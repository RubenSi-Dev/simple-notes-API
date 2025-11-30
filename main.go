package main

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// global database connection
var db *sql.DB

func main() {
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


	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/notes", handleNotes)
	fs := http.FileServer(http.Dir("./frontend-react/dist/"))
	http.Handle("/", fs)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}

