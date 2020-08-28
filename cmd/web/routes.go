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
	mux.Get("/health_check", http.HandlerFunc(app.healthCheck))

	mux.Post("/user", http.HandlerFunc(app.userCreation))
	mux.Get("/user/:id", http.HandlerFunc(app.userRetrieval))
	mux.Get("/user/email/:email", http.HandlerFunc(app.userRetrievalByEmail))
	mux.Del("/user/:id", http.HandlerFunc(app.userDeletion))

	mux.Post("/boardgame", http.HandlerFunc(app.boardgameCreation))
	mux.Get("/boardgame/:id", http.HandlerFunc(app.boardgameRetrieval))
	mux.Del("/boardgame/:id", http.HandlerFunc(app.boardgameDeletion))

	return standardMiddleware.Then(mux)
}
