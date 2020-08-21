package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Rutas de regla de negocio
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))


	return standardMiddleware.Then(mux)
}