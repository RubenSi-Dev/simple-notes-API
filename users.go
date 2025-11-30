package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct{
	UID int
	Username string
	PasswordHashed []byte
}

type RegistrationRequest struct{
	Username 	string	`json:"username"`
	Password 	string	`json:"password"`
}

func (rr *RegistrationRequest) addUser() (*User, error) {
	pwHashed, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcrypt.DefaultCost)
	if err != nil {return nil, err}


	res, err := db.Exec(`INSERT INTO users (username, password) VALUES (?, ?);`, rr.Username, pwHashed)
	if err != nil {
		return nil, err
	}

	uid, err := res.LastInsertId()
	if err != nil {return nil, err}

	return &User{
		UID: int(uid),
		Username: rr.Username,
		PasswordHashed: pwHashed,
	}, nil
}

func handleRegistrations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var received RegistrationRequest 
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			http.Error(w, "Couldn't decode note", http.StatusBadRequest)
			return 
		}

		_, err = received.addUser()
		if err != nil {
			http.Error(w, "couldn't store locally", http.StatusInternalServerError)
		}
		 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}
