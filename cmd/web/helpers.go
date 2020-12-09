package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-sql-driver/mysql"
)

// Metodo para manejar los errores del servidor y mostrar informacion util para el debug (muestra el archivo)
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s,\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, message string, status int) {
	response := struct {
		Error  string
		Status int
	}{
		message,
		status,
	}
	// convierto la string en un json
	js, err := json.Marshal(response)
	if err != nil {
		app.serverError(w, err)
	}
	http.Error(w, string(js), status)
}

func (app *application) responseJSON(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(body)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func duplicateError(err error) bool {
	if err != nil {
		return err.(*mysql.MySQLError).Number == 1062
	}
	return false
}
