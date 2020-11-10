package main

import (
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	// Le agrega a cada respuesta que emite el servidor parametros de aseguramiento de header.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Este codigo se ejecuta antes de llegar al Application Handler!!!
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
		//El codigo aca se ejecuta despues de pasar por el Application handler
	})
}

func (app *application) withCORS(next http.Handler) http.Handler {
	// para mostrar por la salida de log cada request que se le haga al server
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// La peticion debe venir de un origen determinado, no vale *
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// Stop here for a Preflighted OPTIONS request.
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	// para mostrar por la salida de log cada request que se le haga al server
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})

}

func (app *application) authenticateUser(r *http.Request) bool {
	var email = app.session.GetString(r, "email")
	if len(email) <= 0 {
		return false
	}
	return true
}

func (app *application) restrictedEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.authenticateUser(r) {
			app.clientError(w, "Unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})

}

/*
func (app *application) recoverPanic(next http.Handler) http.Handler {
	// Para cerrar las gorutinas que fallen y devolver informacion para debug al log.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
*/
