package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Password string `json:"-"`
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

	user, err := cfg.DB.CreateUser(string(hashedPassword), emptyUser.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondeWithJson(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	emptyUser := User{}

	err := decoder.Decode(&emptyUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode user")
		return
	}

	foundUser, err := cfg.DB.EnsureUser(emptyUser.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(emptyUser.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	respondeWithJson(w, http.StatusOK, User{
		ID:    foundUser.ID,
		Email: foundUser.Email,
	})
}