package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Println("Get auth request body")

	// Validate user credentials
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return;
	}

	log.Println("Get user")


	// validate password
	valid := user.PasswordMatches(requestPayload.Password)
	if !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	log.Println("Verify user")

	app.writeJSON(w, http.StatusAccepted, jsonResponse{
		Status: true,
		Message: "Login successful!",
		Data: user,
	})
}