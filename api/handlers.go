package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	// Log authentication
	go app.logRequest("Authentication", fmt.Sprintf("%s logged in", user.Email))

	app.writeJSON(w, http.StatusAccepted, jsonResponse{
		Status: true,
		Message: "Login successful!",
		Data: user,
	})
}

func (app *Config) logRequest(name, data string) error {
	var entry struct{
		Name string `json:"name"`
		Data string `json:"Data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println(err)
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	
	log.Println(err)

	return err
}