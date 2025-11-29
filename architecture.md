# High level architecture

```notes
[ Client (curl/Postman/browser) ]
                |
                v
        HTTP over localhost
                |
                v
[ Go HTTP Server (net/http) ]
   - Routing (which URL → which handler)
   - Handlers (business logic)
   - In-memory storage for notes
```

no database, no frontend, just:

- Go process listening on loclahost:8000
- A few https routes:
    - GET `/notes`
    - POST `/notes`
    - DELETE `/notes/{id}`

## Internals for the go server

1. HTTP Layer (transport)
    - receives http requests, sends http responses
    - uses net/http in the standard library
    - parses:
        - URL path and query (`/notes`, `/notes?id=1`)
        - HTTP method (`GET`, `POST`, `DELETE`)
        - Body (JSON)

2. Handler Layer (Application Logic)
    - Each route has handler logic, eg:
        - `handleGetNotes(w http.ResponseWriter, r *http.Request)`
        - `handleCreateNotes(w http.ResponseWriter, r *http.Request)`
    - Handlers:
        - validate input
        - call storage functions
        - deicde HTTP status codes and response body

3. Storage Layer (In-Memory for now):
    - A simple structure in go, eg:
        - `type Note struct { ID int; Text string }`
        - `var notes = map[int]Note{}` or `[]Note` plus a counter for IDs
    - Functions such as:
        - `func GetNotes() []Note | map[int]Note{}`
        - `func AddNote(text string) Note`
        - `func DeleteNote(id int) error`

--> Later this storage layer can be swapped for a database withouth changing the HTTP/handler much
