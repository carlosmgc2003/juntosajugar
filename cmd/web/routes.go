package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Metodo de la aplicacion donde coloco todas las rutas

	// Middleware por el que pasa cada request/response
	standardMiddleware := alice.New(app.withCORS, app.logRequest, app.session.Enable, secureHeaders)
	securityMiddleware := alice.New(app.restrictedEndpoint)
	// Rutas de regla de negocio
	mux := pat.New()
	mux.Get("/health_check", http.HandlerFunc(app.healthCheck))
	// User Routes
	mux.Options("/login", http.HandlerFunc(app.userLogin))
	mux.Post("/login", http.HandlerFunc(app.userLogin))
	mux.Get("/logout", http.HandlerFunc(app.userLogout))
	mux.Get("/user", securityMiddleware.ThenFunc(app.userList))
	mux.Options("/user", http.HandlerFunc(app.userCreation))
	mux.Post("/user", http.HandlerFunc(app.userCreation))
	mux.Get("/user/boardgames", securityMiddleware.ThenFunc(app.userListBoardgames))
	mux.Get("/user/:id", http.HandlerFunc(app.userRetrieval))
	mux.Get("/user/email/:email", http.HandlerFunc(app.userRetrievalByEmail))
	mux.Options("/user/:id", securityMiddleware.ThenFunc(app.userDeletion))
	mux.Del("/user/:id", securityMiddleware.ThenFunc(app.userDeletion))
	mux.Post("/user/boardgames", securityMiddleware.ThenFunc(app.userAddBoardgame))
	mux.Del("/user/boardgames/:id", securityMiddleware.ThenFunc(app.userDelBoardgames))

	// Boardgame Routes
	mux.Post("/boardgame", securityMiddleware.ThenFunc(app.boardgameCreation))
	mux.Get("/boardgame", securityMiddleware.ThenFunc(app.boardgameList))
	mux.Get("/boardgame/:id", http.HandlerFunc(app.boardgameRetrieval))
	mux.Del("/boardgame/:id", http.HandlerFunc(app.boardgameDeletion))

	// Gameemeting Routes
	mux.Post("/gamemeeting", securityMiddleware.ThenFunc(app.gamemeetingCreation))
	mux.Get("/gamemeeting", securityMiddleware.ThenFunc(app.gamemeetingList))

	return standardMiddleware.Then(mux)
}
