window.addEventListener("DOMContentLoaded", () => {
  loadNotes();
  setupForm();
});

async function loadNotes() {
  const response = await fetch("/notes");
  if (!response.ok) {
    console.error("failed to fetch notes");
    return;
  }

  const notes = await response.json();
  const list = document.getElementById("notes-list");
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

async function createNote(author, text) {
  const response = await fetch("/notes", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ author, text }),
  });
  console.log(author, text);

  if (!response.ok) {
    console.error("Failed to create note");
    return;
  }
  loadNotes();
}

async function setupForm() {
  const form = document.getElementById("note-form");
  const authorInput = document.getElementById("author-input");
  const textInput = document.getElementById("text-input");

  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    const author = authorInput.value.trim();
    const text = textInput.value.trim();

    if (!author || !text) {
      return;
    }
    await createNote(author, text);

    authorInput.value = "";
    textInput.value = "";
  });
}

async function deleteNote(id) {
  const response = await fetch("/notes?id=" + encodeURIComponent(id), {
    method: "DELETE",
  });

  if (!response.ok) {
    console.error("Failed to delete note");
		return
  }

  loadNotes();
}
