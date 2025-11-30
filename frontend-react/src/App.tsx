import { useEffect, useState, type JSX } from "react";
import { NoteForm } from "./components/NoteForm.js";
import { NoteItem } from "./components/NoteItem.js";
import "./App.css";

export interface Note {
  id: number;
  author: string;
  text: string;
  edited: boolean;
}

export interface NoteRequest {
  author: string;
  text: string;
}

export interface NoteFormProps {
  onAddNote: (author: string, text: string) => Promise<void>;
}

export interface NoteItemProps {
  note: Note;
  onEdit: (id: number, text: string) => Promise<void>;
  onDeleteNote: (id: number) => Promise<void>;
}

function App(): JSX.Element {
  const [notes, setNotes] = useState<Note[]>([]);

  useEffect(() => {
    fetch("/notes")
      .then((res) => res.json())
      .then((data) => setNotes(data));
  });

  const handleNewNote = async (author: string, text: string): Promise<void> => {
    await fetch("/notes", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        author: author,
        text: text,
      }),
    });
  };

  const handleDeleteNote = async (id: number): Promise<void> => {
    await fetch("/notes?id=" + encodeURIComponent(id), {
      method: "DELETE",
    });
  };

  const handleEditNote = async (id: number, text: string): Promise<void> => {
    await fetch("/notes?id=" + encodeURIComponent(id), {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        text: text,
      }),
    });
  };

  return (
    <>
      <h1>Notes</h1>
      <div className="notes-layout">
        <NoteForm onAddNote={handleNewNote} />
        <ul className="note-list">
          {notes.map((note) => {
            return (
              <NoteItem
                key={note.id}
                note={note}
                onDeleteNote={handleDeleteNote}
                onEdit={handleEditNote}
              />
            );
          })}
        </ul>
      </div>
    </>
  );
}

export default App;
