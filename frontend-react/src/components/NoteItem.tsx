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
      <div>
        <span className={note.edited ? "note-edited" : "note-non-edited"}>
          <ul>
            <div className="note-item">{`(${note.id}) ${note.author}: ${note.text}`}</div>
            <div className="note-item-buttons">
              <button className="edit-buttons" onClick={() => onDeleteNote(note.id)}>Delete</button>
              <button className="edit-buttons"
                onClick={() => {
                  setIsEditing(!isEditing);
                }}
              >
                Cancel
              </button>
              <NoteEditForm
                note={note}
                onEdit={onEdit}
                setIsEditing={setIsEditing}
              />
            </div>
          </ul>
        </span>
      </div>
    </>
  ) : (
    <>
      <div>
        <span className={note.edited ? "note-edited" : ""}>
          <ul>
            <div className="note-item">{`(${note.id}) ${note.author}: ${note.text}  `}</div>
            <div className="note-item-buttons">
              <button className="edit-buttons" onClick={() => onDeleteNote(note.id)}>Delete</button>
              {"  "}
              <button className="edit-buttons"
                onClick={() => {
                  setIsEditing(!isEditing);
                }}
              >
                Edit
              </button>
            </div>
          </ul>
        </span>
      </div>
    </>
  );
}
