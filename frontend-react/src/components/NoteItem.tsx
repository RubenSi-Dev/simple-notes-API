import { useState, type JSX, type SetStateAction } from "react";
import { type Note, type NoteItemProps } from "../App";
import { NoteEditForm } from "./NoteEditForm";

export interface NoteEditFormProps {
  note: Note;
  onEdit: (id: number, text: string) => void;
  setIsEditing: React.Dispatch<SetStateAction<boolean>>;
}

export function NoteItem({
  note,
  onDeleteNote,
  onEdit,
}: NoteItemProps): JSX.Element {
  const [isEditing, setIsEditing] = useState<boolean>(false);
  return isEditing ? (
    <>
      <span className={note.edited ? "note-edited" : ""}>
        {`(${note.id}) ${note.author}: ${note.text}  `}
        <button onClick={() => onDeleteNote(note.id)}>Delete</button>
        {"  "}
        <button
          onClick={() => {
            setIsEditing(!isEditing);
            return (
              <NoteEditForm
                note={note}
                onEdit={(id: number, text: string) => {
                  onEdit(id, text);
                }}
                setIsEditing={setIsEditing}
              />
            );
          }}
        >
          Editing
        </button>
        <NoteEditForm note={note} onEdit={onEdit} setIsEditing={setIsEditing} />
      </span>
    </>
  ) : (
    <>
      <span className={note.edited ? "note-edited" : ""}>
        {`(${note.id}) ${note.author}: ${note.text}  `}
        <button onClick={() => onDeleteNote(note.id)}>Delete</button>
        {"  "}
        <button
          onClick={() => {
            setIsEditing(!isEditing);
            return (
              <NoteEditForm
                note={note}
                onEdit={onEdit}
                setIsEditing={setIsEditing}
              />
            );
          }}
        >
          Edit
        </button>
      </span>
    </>
  );
}
