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
