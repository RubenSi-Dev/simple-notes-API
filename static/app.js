window.addEventListener("DOMContentLoaded", () => {
  loadNotes();
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

  notes.forEach(note => {
    const li = document.createElement("li");
    li.textContent = `${note.id}: [${note.author}		${note.text}]`;
    list.appendChild(li);
  });
}
