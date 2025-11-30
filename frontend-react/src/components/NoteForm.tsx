import { useState, type JSX } from "react";
import type { NoteFormProps } from "../App";

export function NoteForm({ onAddNote }: NoteFormProps): JSX.Element {
  const [author, setAuthor] = useState<string>("");
  const [text, setText] = useState<string>("");

  return (
    <>
      <input
        type="text"
        value={author}
        onChange={(e) => setAuthor(e.target.value)}
        placeholder="author"
      ></input>
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="text"
      ></input>
      <button
        onClick={() => {
          onAddNote(author, text);
          setAuthor("");
          setText("");
        }}
      >
        Submit
      </button>
    </>
  );
}
