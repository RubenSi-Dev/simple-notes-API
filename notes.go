package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)


type Note struct{
	ID 			int 		`json:"id"`
	Author 	string	`json:"author"`
	Text 		string	`json:"text"`
}

type NoteRequest struct {
	Author 	string	`json:"author"`
	Text 		string	`json:"text"`
}

type Notes struct{
	Entries []*Note
	NextID 	int
}

var notes = Notes{
	Entries: []*Note{},
	NextID: 0,
}

func (n *Notes) addNote(nr *NoteRequest) (result *Note) {
	result = &Note{
		ID: n.NextID,
		Author: nr.Author,
		Text: nr.Text,
	}
	n.Entries = append(n.Entries, result)
	n.NextID++
	return 
}

func (n *Notes) removeByID(id int) (*Note, error) {
	for i, note := range n.Entries {
		if note.ID == id {
			result := n.Entries[i]
			n.Entries = append(n.Entries[:i], n.Entries[i+1:]...)
			return result, nil
		}
	}
	return nil, fmt.Errorf("note of ID %v not found", id)
}

func handleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet: // GET
		handleNotesGet(w, r)	

	case http.MethodPost: // POST
		handleNotesPost(w, r)

	case http.MethodDelete: // DELETE
		handleNotesDelete(w, r)

	case http.MethodPut: // PUT
		handleNotesPut(w, r)	

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}


func handleNotesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(notes.Entries)
	if err != nil {
		http.Error(w, "Couldn't encode notes", http.StatusInternalServerError)
		return 
	}
}

func handleNotesPost(w http.ResponseWriter, r *http.Request) {
	var received NoteRequest
	err := json.NewDecoder(r.Body).Decode(&received)
	if err != nil {
		http.Error(w, "Couldn't decode note", http.StatusBadRequest)
		return 
	}

	newNote := notes.addNote(&received)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

func handleNotesDelete(w http.ResponseWriter, r *http.Request) {
	id, errConv := strconv.Atoi(r.URL.Query().Get("id"))
	
	if errConv != nil {
		http.Error(w, "Bad format", http.StatusBadRequest)
		return 
	} 
	deleted, errID := notes.removeByID(id)

	if errID != nil {
		http.Error(w, "No such ID", http.StatusNotFound)
		return 
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(deleted)
}

func handleNotesPut(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
