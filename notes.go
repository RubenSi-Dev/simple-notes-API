package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Note - how notes are saved in memory
type Note struct{
	ID 			int 		`json:"id"`
	Author 	string	`json:"author"`
	Text 		string	`json:"text"`
	Edited  bool    `json:"edited"`
}

// NoteRequest - how new notes are requested
type NoteRequest struct {
	Author 	string	`json:"author"`
	Text 		string	`json:"text"`
}

// NoteUpdateRequest - how updates for notes are requested
type NoteUpdateRequest struct {
	Text		string 	`json:"text"`
}

// toNote - converts a NoteRequest into a Note that can be saved in memory
func (nr *NoteRequest) toNote(id int) (result *Note) {
	return &Note{
		ID: id,
		Author: nr.Author,
		Text: nr.Text,
		Edited: false,
	}
}

// addNote - add a new Note to the notes by passing a NoteRequest
func addNote(nr *NoteRequest) (*Note, error) {
	res, err := db.Exec(
		`INSERT INTO notes (author, text) VALUES (?, ?);`, 
		nr.Author,
		nr.Text,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting %v in db not successful", nr)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("inserting %v in db not successful", nr)
	}

	return nr.toNote(int(id)), nil
}

func getNotes() ([]*Note, error) {
	rows, err := db.Query(`SELECT id, author, text, edited FROM notes ORDER BY id;`)
	defer rows.Close() 

	if err != nil {
		return nil, fmt.Errorf("couldn't read rows %w", err)
	}

	result := []*Note{}

	for rows.Next() {
		note := &Note{}
		var editedInt int;
		if err := rows.Scan(&note.ID, &note.Author, &note.Text, &editedInt); err != nil {
			return nil, fmt.Errorf("couldn't read all rows %w", err)
		}

		if editedInt == 1 {
			note.Edited = true
		}	else {
			note.Edited = false
		}
		result = append(result, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading problem %w", err)
	}

	return result, nil
}

func getNoteByID(id int) (*Note, error) {
	result := &Note{}
	row := db.QueryRow(`SELECT id, author, text, edited FROM notes WHERE id=?`, id)
	var editedInt int;
	err := row.Scan(&result.ID, &result.Author, &result.Text, &editedInt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("note with id %v not found", id)
	}

	if editedInt == 1 {
		result.Edited = true
	} else {
		result.Edited = false
	}

	if err != nil {
		return nil, fmt.Errorf("error reading note %w", err)
	}
	return result, nil
}

func removeNoteByID(id int) error {
	res, err := db.Exec(
		`DELETE FROM notes WHERE id=?`,
		id,
	)

	if err != nil {  // check whether sql query was valid
		return fmt.Errorf("SQL error %w", err)
	}

	affected, err := res.RowsAffected()

	if affected == 0 || err != nil  {
		return fmt.Errorf("invalid id %v", id)
	}

	return nil
}


func updateNoteByID(id int, nur *NoteUpdateRequest) error {
	res, err := db.Exec(`UPDATE notes SET text = ?, edited = 1 WHERE id = ?`, nur.Text, id)
	if err != nil {
		return fmt.Errorf("SQL error %w", err)
	}

	affected, err := res.RowsAffected()

	if affected == 0 || err != nil {
		return fmt.Errorf("invalid id %v", id)
	}

	return nil
}

// handleNotes - main handler for /notes endpoint
func handleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet: // GET
		handleNotesGet(w, r)	

	case http.MethodPost: // POST
		handleNotesPost(w, r)

	case http.MethodDelete: // DELETE
		handleNotesDelete(w, r)

	case http.MethodPatch: // PATCH
		handleNotesPatch(w, r)	

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}


// handleNotesGet - handles GET requests to /notes
func handleNotesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	notes, err := getNotes()
	if err != nil {
		http.Error(w, "Couldn't fetch notes from DB", http.StatusInternalServerError) 
		return
	}

	err = json.NewEncoder(w).Encode(&notes)
	if err != nil {
		http.Error(w, "Couldn't encode notes", http.StatusInternalServerError)
		return 
	}
}

// handleNotesPost - handles POST requests to /notes
func handleNotesPost(w http.ResponseWriter, r *http.Request) {
	var received NoteRequest
	err := json.NewDecoder(r.Body).Decode(&received)
	if err != nil {
		http.Error(w, "Couldn't decode note", http.StatusBadRequest)
		return 
	}

	newNote, err := addNote(&received)
	if err != nil {
		http.Error(w, "couldn't store locally", http.StatusInternalServerError)
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newNote)
	if err != nil {
		http.Error(w, "couldn't encode new note", http.StatusInternalServerError)
		return 
	}

}

// handleNotesDelete - handles DELETE requests to /notes
func handleNotesDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
	}
	id, err := strconv.Atoi(idStr)
	
	if err != nil {
		http.Error(w, "Bad format", http.StatusBadRequest)
		return 
	} 

	err = removeNoteByID(id)

	if err != nil {
		http.Error(w, "No such ID", http.StatusNotFound)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// handleNotesPatch - handles PATCH requests to /notes - changes the text only, not the author
func handleNotesPatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Bad format", http.StatusBadRequest)
		return 
	}
	
	var received NoteUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&received)
	if err != nil {
		http.Error(w, "Couldn't decode note", http.StatusBadRequest)
		return 
	}

	err = updateNoteByID(id, &received)
	if err != nil {
		http.Error(w, "No such ID", http.StatusNotFound)
		return 
	}

	updated, err := getNoteByID(id)
	if err != nil {
		http.Error(w, "coudln't retrieve updated note", http.StatusNotFound)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updated)
	if err != nil {
		http.Error(w, "couldn't encode updated note", http.StatusInternalServerError)
		return 
	}
}
