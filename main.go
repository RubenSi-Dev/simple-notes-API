package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/notes", handleNotes)
	fs := http.FileServer(http.Dir("./frontend-react/dist/"))
	http.Handle("/", fs)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}

