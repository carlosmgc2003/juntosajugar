package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Metodo de la aplicacion donde coloco todas las rutas

	// Middleware por el que pasa cada request/response
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Rutas de regla de negocio
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))

	return standardMiddleware.Then(mux)
}
