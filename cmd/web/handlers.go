package main

import (
	"encoding/json"
	"io/ioutil"
	"juntosajugar/pkg/models"
	"net/http"
	"net/url"
	"strconv"
)

func (app *application) healthCheck(w http.ResponseWriter, _ *http.Request) {
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
func (app *application) userList(w http.ResponseWriter, _ *http.Request) {
	enableCors(&w) //insecure
	var users []models.User
	result := app.db.Find(&users)
	if result.Error != nil {
		app.serverError(w, result.Error)
		return
	}
	body, err := json.Marshal(users)
	if err != nil {
		app.serverError(w, err)
	}
	app.responseJson(w, body)
}

func (app *application) userCreation(w http.ResponseWriter, r *http.Request) {
	// Handler que toma el body del request, trata de Unmarshalearlo en una struct de
	// tipo user, y si no hay duplicados responde con el mismo user en el cuerpo.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var newUser models.User
	err = newUser.FromJson(body)
	if err != nil {
		app.clientError(w, err.Error(), 400)
		return
	}

	err = app.db.Create(&newUser).Error
	if duplicateError(err) {
		app.clientError(w, err.Error(), 409)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) userDeletion(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id en la URI, elimina al usuario y devuelve un 200 vacío.
	userId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var delUser models.User
	err = app.db.First(&delUser, userId).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.Unscoped().Delete(&delUser).Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) userRetrieval(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id de usuario en la URI, devuelve los datos del mismo en Json
	userId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var reqUser models.User
	err = app.db.First(&reqUser, userId).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	body, err := json.Marshal(&reqUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) userRetrievalByEmail(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id de usuario en la URI, devuelve los datos del mismo en Json
	encodedParameter := r.URL.Query().Get(":email")
	userEmail, err := url.QueryUnescape(encodedParameter)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var reqUser models.User
	err = app.db.Where("email = ?", userEmail).First(&reqUser).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	body, err := json.Marshal(&reqUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) boardgameCreation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var newBordgame models.Boardgame
	err = newBordgame.FromJson(body)
	if err != nil {
		app.clientError(w, err.Error(), 400)
		return
	}

	err = app.db.Create(&newBordgame).Error
	if duplicateError(err) {
		app.clientError(w, err.Error(), 409)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) boardgameRetrieval(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id de usuario en la URI, devuelve los datos del mismo en Json
	boardgameId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var reqBoardgame models.Boardgame
	err = app.db.First(&reqBoardgame, boardgameId).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	body, err := json.Marshal(&reqBoardgame)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) boardgameDeletion(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id en la URI, elimina al usuario y devuelve un 200 vacío.
	boardgameId, err := strconv.Atoi(r.URL.Query().Get(":id"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(boardgameId)

	var delBoardgame models.Boardgame
	err = app.db.First(&delBoardgame, boardgameId).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.Unscoped().Delete(&delBoardgame).Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func (app *application) gamemeetingCreation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var newGameMeeting models.Gamemeeting
	err = newGameMeeting.FromJson(body, app.db)
	app.infoLog.Println(newGameMeeting)
	if err != nil {
		app.clientError(w, err.Error(), 400)
		return
	}
	app.infoLog.Println(newGameMeeting)
	err = app.db.Create(&newGameMeeting).Error
	if duplicateError(err) {
		app.clientError(w, err.Error(), 409)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJson(w, body)
	return
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
