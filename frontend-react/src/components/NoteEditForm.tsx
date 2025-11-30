import { useState, type JSX } from "react";
import { type NoteEditFormProps } from "./NoteItem";

export function NoteEditForm({
  note,
  onEdit,
  setIsEditing,
}: NoteEditFormProps): JSX.Element {
  const [text, setText] = useState<string>(note.text);

  return (
    <>
      <input
        type="text"
        placeholder="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      <button className="edit-buttons"
        onClick={async () => {
          if (!text) return;
          onEdit(note.id, text);
          setIsEditing(false);
        }}
      >
        Submit
      </button>
    </>
  );
}
