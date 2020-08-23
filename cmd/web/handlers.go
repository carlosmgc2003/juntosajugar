package main

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"io/ioutil"
	"juntosajugar/pkg/models"
	"net/http"
)

func (app *application) health_check(w http.ResponseWriter, r *http.Request) {
	// Handler de ejemplo que devuele un Json indicando que el servidor esta ok

	// Creo un struct anonima con los valores que quiero mandar
	response := struct {
		Key   string
		Value string
	}{
		"servidor",
		"ok",
	}

	// convierto la string en un json
	js, err := json.Marshal(response)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) user_creation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var t models.User
	err = json.Unmarshal(body, &t)
	if err != nil {
		app.serverError(w, err)
	}

	err = app.db.Create(&t).Error
	if err != nil && err.(*mysql.MySQLError).Number == 1062 {
		app.errorLog.Printf("%s - %s", r.RemoteAddr, err.(*mysql.MySQLError).Message)
		app.clientError(w, 409)
		return
	}

	//Realizo el "pong" de respuesta
	js, err := json.Marshal(t)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		app.serverError(w, err)
		return
	}

}
