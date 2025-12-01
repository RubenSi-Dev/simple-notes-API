package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sqlite3 "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
	UID 						int			`json:"id"`
	Username 				string	`json:"username"`
	PasswordHashed 	[]byte	`json:"-"` // hide hash in JSON output
}

type AuthenticationRequest struct{
	Username 	string	`json:"username"`
	Password 	string	`json:"password"`
}

func (ar *AuthenticationRequest) addUser() (*User, error) {
	pwHashed, err := bcrypt.GenerateFromPassword([]byte(ar.Password), bcrypt.DefaultCost)
	if err != nil {return nil, err}

	res, err := db.Exec(`INSERT INTO users (username, password) VALUES (?, ?);`, ar.Username, pwHashed)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("username already taken")
		}

		return nil, fmt.Errorf("database error during registration: %w", err)
	}
	uid, err := res.LastInsertId()
	if err != nil {return nil, err}

	return &User{
		UID: int(uid),
		Username: ar.Username,
		PasswordHashed: pwHashed,
	}, nil
}

func getUserByUsername(username string) (*User, error) {
	user := &User{}
	row := db.QueryRow(`SELECT uid, username, password FROM users WHERE username=?`, username)

	if err := row.Scan(&user.UID, &user.Username, &user.PasswordHashed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user does not exist")
		}
		return nil, fmt.Errorf("database error during user retrieval: %w", err)
	}
	return user, nil
}


func handleRegistrations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var received AuthenticationRequest 
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			http.Error(w, "Couldn't decode request", http.StatusBadRequest)
			return 
		}

		newUser, err := received.addUser()

		if err != nil {
			if errors.Is(err, fmt.Errorf("username already taken")) {
				http.Error(w, "Username already exists", http.StatusConflict)
				return 
			}

			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return 
		}
		 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newUser)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleLogins(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var received AuthenticationRequest
		err := json.NewDecoder(r.Body).Decode(&received)

		if err != nil {
			http.Error(w, "Couldn't deocode request", http.StatusBadRequest)
			return 
		}

		user, err := getUserByUsername(received.Username)
		if err != nil {
			http.Error(w, "wrong credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword(user.PasswordHashed, []byte(received.Password))
		if err != nil {
			http.Error(w, "wrong credentials", http.StatusUnauthorized)
			return
		}

		// JWT be valid for 24h
		expirationTime := time.Now().Add(24 * time.Hour) 
		claims := &Claims{
			UID: user.UID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "couldn't generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "Succesfully logged in",
			"token": tokenString,
			"uid": user.UID,
		})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
