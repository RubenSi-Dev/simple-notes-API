import { useState, type JSX } from "react";
import type { NoteFormProps } from "../App";

export function NoteForm({ onAddNote }: NoteFormProps): JSX.Element {
  const [author, setAuthor] = useState<string>("");
  const [text, setText] = useState<string>("");

  return (
    <>
      <input
        type="text"
				className="entry-fields"
        value={author}
        onChange={(e) => setAuthor(e.target.value)}
        placeholder="author"
      ></input>
      <input
        type="text"
				className="entry-fields"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="text"
      ></input>
      <button
				className="buttons"	
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
