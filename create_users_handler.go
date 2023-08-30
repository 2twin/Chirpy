package main

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Password string
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	emptyUser := User{}
	err := decoder.Decode(&emptyUser)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode user")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(emptyUser.Password), 2)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
	}

	// проверка есть ли такой email уже

	user, err := cfg.DB.CreateUser(string(hashedPassword), emptyUser.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondeWithJson(w, http.StatusCreated, User{
		ID:       user.ID,
		Email:    user.Email,
	})
}
