"use strict";
(() => {
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __commonJS = (cb, mod) => function __require() {
    return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
  };

  // frontend/main.ts
  var require_main = __commonJS({
    "frontend/main.ts"() {
      window.addEventListener("DOMContentLoaded", () => {
        loadNotes();
        setupForm();
      });
      async function loadNotes() {
        const response = await fetch("/notes");
        if (!response.ok) {
          throw new Error("failed to fetch notes");
        }
        const notes = await response.json();
        const list = document.getElementById("notes-list");
        if (!(list instanceof HTMLUListElement)) {
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
      async function createNote(author, text) {
        const response = await fetch("/notes", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({ author, text })
        });
        if (!response.ok) {
          throw new Error("failed to create note");
        }
        loadNotes();
      }
      function setupForm() {
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
          await createNote(author, text);
          authorInput.value = "";
          textInput.value = "";
        });
      }
      async function deleteNote(id) {
        const response = await fetch("/notes?id=" + encodeURIComponent(id), {
          method: "DELETE"
        });
        if (!response.ok) {
          throw new Error("Failed to delete note");
        }
        await loadNotes();
      }
    }
  });
  require_main();
})();
//# sourceMappingURL=app.js.map
