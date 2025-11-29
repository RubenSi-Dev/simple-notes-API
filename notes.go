package main

import (
	"encoding/json"
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

// Notes - datastructure that stores all notes
type Notes struct{
	Entries []*Note
	NextID 	int
}

var notes = Notes{
	Entries: []*Note{},
	NextID: 0,
}

// addNote - add a new Note to the notes by passing a NoteRequest
func (n *Notes) addNote(nr *NoteRequest) (result *Note) {
	result = nr.toNote(n.NextID)
	n.Entries = append(n.Entries, result)
	n.NextID++
	return 
}


// removeByID - removes note with specified ID if it is inside notes
// returns the Note that was deleted
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

// updateById - given a NoteUpdateRequest and the corresponding ID, update the notes text
// Returns the updated note as a Note if id was found
func (n *Notes) updateById(id int, nur *NoteUpdateRequest) (*Note, error) {
	for i, note := range n.Entries {
		if note.ID == id {
			result := &Note{
				ID: id,
				Author: note.Author,
				Text: nur.Text,
				Edited: true,
			}
			n.Entries = append(append(n.Entries[:i], result), n.Entries[i+1:]...)
			return result, nil
		}
	}
	return nil, fmt.Errorf("note of ID %v not found", id)
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

	case http.MethodPut: // PUT
		handleNotesPut(w, r)	

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}


// handleNotesGet - handles GET requests to /notes
func handleNotesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	errEnc := json.NewEncoder(w).Encode(notes.Entries)
	if errEnc != nil {
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

	newNote := notes.addNote(&received)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	errEnc := json.NewEncoder(w).Encode(newNote)
	if errEnc != nil {
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
	id, errConv := strconv.Atoi(idStr)
	
	if errConv != nil {
		http.Error(w, "Bad format", http.StatusBadRequest)
		return 
	} 
	deleted, errDel := notes.removeByID(id)

	if errDel != nil {
		http.Error(w, "No such ID", http.StatusNotFound)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	errEnc := json.NewEncoder(w).Encode(deleted)
	if errEnc != nil {
		http.Error(w, "Couldn't encode deleted note", http.StatusInternalServerError)	
		return 
	}
}

// handleNotesPut - handles PUT requests to /notes
func handleNotesPut(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}
	id, errConv := strconv.Atoi(idStr)

	if errConv != nil {
		http.Error(w, "Bad format", http.StatusBadRequest)
		return 
	}
	
	var received NoteUpdateRequest
	errDec := json.NewDecoder(r.Body).Decode(&received)
	if errDec != nil {
		http.Error(w, "Couldn't decode note", http.StatusBadRequest)
		return 
	}

	updated, errUpd := notes.updateById(id, &received)

	if errUpd != nil {
		http.Error(w, "No such ID", http.StatusNotFound)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	errEnc := json.NewEncoder(w).Encode(updated)
	if errEnc != nil {
		http.Error(w, "Couldn't encode updated note", http.StatusInternalServerError)
		return 
	}
}
