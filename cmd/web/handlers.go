package main

import (
	"encoding/json"
	"net/http"
)



func (app *application) home(w http.ResponseWriter, r *http.Request) {
	response := struct{
		Key string
		Value string
	}{
		"servidor",
		"ok",
	}

	js, err := json.Marshal(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
