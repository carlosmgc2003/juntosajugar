package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Metodo para manejar los errores del servidor y mostrar informacion util para el debug (muestra el archivo)
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s,\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) responseJson(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(body)
	if err != nil {
		app.serverError(w, err)
		return
	}
}
