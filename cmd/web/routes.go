package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Metodo de la aplicacion donde coloco todas las rutas

	// Middleware por el que pasa cada request/response
	standardMiddleware := alice.New(app.logRequest, secureHeaders)

	// Rutas de regla de negocio
	mux := pat.New()
	mux.Get("/health_check", http.HandlerFunc(app.health_check))

	mux.Post("/user", http.HandlerFunc(app.user_creation))

	return standardMiddleware.Then(mux)
}
