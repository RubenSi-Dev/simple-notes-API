window.addEventListener("DOMContentLoaded", () => {
  loadNotes();
  setupForm();
});

type NoteRequest = {
	author: string;
	text: string;
}

type Note = {
  id: number;
  author: string;
  text: string;
};

async function loadNotes(): Promise<void> {
  const response = await fetch("/notes");
  if (!response.ok) {
    throw new Error("failed to fetch notes");
  }

  const notes: Note[] = await response.json();

  const list = document.getElementById("notes-list");
  if (!(list instanceof HTMLUListElement)) {
    // type guards
    throw new Error("notes-list element not found");
  }

  list.innerHTML = "";

  notes.forEach((note) => {
    const li = document.createElement("li");

    const textSpan = document.createElement("span");
    textSpan.textContent = `(${note.id}) ${note.author}:  ${note.text}			`;
    const deleteButton = document.createElement("button");
    deleteButton.textContent = "Delete";
    deleteButton.addEventListener("click", async () => {
      await deleteNote(note.id);
    });
    li.appendChild(textSpan);
    li.appendChild(deleteButton);
    list.appendChild(li);
  });
}

async function createNote(noteReq: NoteRequest): Promise<void> {
  const response = await fetch("/notes", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(noteReq),
  });

  if (!response.ok) {
    throw new Error("failed to create note");
  }
  loadNotes();
}

function setupForm(): void {
  const form = document.getElementById("note-form");
  const authorInput = document.getElementById("author-input");
  const textInput = document.getElementById("text-input");

  if (!(form instanceof HTMLFormElement))
    throw new Error("form not found or wrong type");

  if (!(authorInput instanceof HTMLInputElement))
    throw new Error("form not found or wrong type");

  if (!(textInput instanceof HTMLInputElement))
    throw new Error("form not found or wrong type");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const author = authorInput.value.trim();
    const text = textInput.value.trim();

    await createNote({author: author, text: text});

    authorInput.value = "";
    textInput.value = "";
  });
}

async function deleteNote(id: number): Promise<void> {
  const response = await fetch("/notes?id=" + encodeURIComponent(id), {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error("Failed to delete note");
  }

  await loadNotes();
}
