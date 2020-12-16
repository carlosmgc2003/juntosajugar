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

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err != nil {
		app.serverError(w, err)
		return
	}
	var newLogin models.Login
	var tempUser models.User
	err = newLogin.FromJSON(body)
	if err != nil {
		app.clientError(w, err.Error(), 400)
		return
	}
	err = app.db.Where("email = ?", newLogin.Email).First(&tempUser).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = tempUser.Authenticate(newLogin)
	if err != nil {
		app.clientError(w, err.Error(), 401)
	}
	app.session.Put(r, "email", tempUser.Email)
	app.session.Put(r, "user_id", tempUser.ID)
	w.Write([]byte(`{"status":"OK"}`))
	return
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "email")
	app.session.Remove(r, "user_id")
	response := struct {
		Key   string
		Value string
	}{
		"logout",
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

func (app *application) userList(w http.ResponseWriter, r *http.Request) {
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
	app.responseJSON(w, body)
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

	err = newUser.FromJSON(body)
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

	body, err = json.Marshal(&newUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
	return
}

func (app *application) userDeletion(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id en la URI, elimina al usuario y devuelve un 200 vacío.
	userID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var delUser models.User
	err = app.db.First(&delUser, userID).Error
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
	app.responseJSON(w, body)
	return
}

func (app *application) userRetrieval(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id de usuario en la URI, devuelve los datos del mismo en Json
	userID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var reqUser models.User
	err = app.db.First(&reqUser, userID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	bgCount := app.db.Model(&reqUser).Association("UserBoardgames").Count()
	bgSlice := make([]models.Boardgame, bgCount)
	err = app.db.Model(&reqUser).Association("UserBoardgames").Find(&bgSlice).Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	reqUser.UserBoardgames = bgSlice

	body, err := json.Marshal(&reqUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
	return
}

func (app *application) userAddBoardgame(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := struct {
		BoardgameID int `json:"boardgame_id"`
	}{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		app.clientError(w, err.Error(), http.StatusBadRequest)
	}

	var bgUser models.User
	err = app.db.First(&bgUser, userID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	var bg models.Boardgame
	err = app.db.First(&bg, data.BoardgameID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.Model(&bgUser).Association("UserBoardgames").Append(&bg).Error
	if err != nil {
		app.serverError(w, err)
		return
	}

	body, err = json.Marshal(&bgUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
	return
}

func (app *application) userListBoardgames(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	var bgUser models.User
	err = app.db.First(&bgUser, userID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	bgCount := app.db.Model(&bgUser).Association("UserBoardgames").Count()
	bgSlice := make([]models.Boardgame, bgCount)
	err = app.db.Model(&bgUser).Association("UserBoardgames").Find(&bgSlice).Error
	if err != nil {
		app.serverError(w, err)
		return
	}

	body, err := json.Marshal(&bgSlice)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
	return
}

func (app *application) userDelBoardgames(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get(":user_id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	bgID, err := strconv.Atoi(r.URL.Query().Get(":bg_id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	var bgUser models.User
	err = app.db.First(&bgUser, userID).Error
	if err != nil && err.Error() == "User not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	var bgDel models.Boardgame
	err = app.db.First(&bgDel, bgID).Error
	if err != nil && err.Error() == "Boardgame not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.Model(&bgUser).Association("UserBoardgames").Delete(bgDel).Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	body, err := json.Marshal(&bgDel)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
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

	bgCount := app.db.Model(&reqUser).Association("UserBoardgames").Count()
	bgSlice := make([]models.Boardgame, bgCount)
	err = app.db.Model(&reqUser).Association("UserBoardgames").Find(&bgSlice).Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	reqUser.UserBoardgames = bgSlice

	body, err := json.Marshal(&reqUser)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
	return
}

func (app *application) boardgameCreation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var newBoardgame models.Boardgame
	err = newBoardgame.FromJson(body)
	if err != nil {
		app.clientError(w, err.Error(), 400)
		return
	}

	err = app.db.Create(&newBoardgame).Error
	if duplicateError(err) {
		app.clientError(w, err.Error(), 409)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	body, err = json.Marshal(&newBoardgame)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
	return
}

func (app *application) boardgameRetrieval(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id del boardgame en la URI, devuelve los datos del mismo en Json
	boardgameID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var reqBoardgame models.Boardgame
	err = app.db.First(&reqBoardgame, boardgameID).Error
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
	app.responseJSON(w, body)
	return
}

func (app *application) boardgameDeletion(w http.ResponseWriter, r *http.Request) {
	// Manejador que dada una peticion con el id en la URI, elimina al usuario y devuelve un 200 vacío.
	boardgameID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(boardgameID)

	var delBoardgame models.Boardgame
	err = app.db.First(&delBoardgame, boardgameID).Error
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
	app.responseJSON(w, body)
	return
}

func (app *application) boardgameList(w http.ResponseWriter, r *http.Request) {
	var boardgames []models.Boardgame
	result := app.db.Find(&boardgames)

	if result.Error != nil {
		app.serverError(w, result.Error)
		return
	}
	body, err := json.Marshal(boardgames)
	if err != nil {
		app.serverError(w, err)
	}
	app.responseJSON(w, body)
}

func (app *application) boardgameByName(w http.ResponseWriter, r *http.Request) {
	encodedParameter := r.URL.Query().Get(":name")
	bgName, err := url.QueryUnescape(encodedParameter)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var reqBoardgame models.Boardgame
	err = app.db.Where("name = ?", bgName).First(&reqBoardgame).Error
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
	app.responseJSON(w, body)
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
	if err != nil {
		app.clientError(w, err.Error(), 400)
		return
	}
	err = app.db.Create(&newGameMeeting).Error
	if duplicateError(err) {
		app.clientError(w, err.Error(), 409)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	app.responseJSON(w, body)
	return
}

func (app *application) gamemeetingList(w http.ResponseWriter, r *http.Request) {
	var gamemeetings []models.Gamemeeting
	result := app.db.Find(&gamemeetings)
	if result.Error != nil {
		app.serverError(w, result.Error)
		return
	}
	if result.RowsAffected >= 0 {
		for i := range gamemeetings {
			gamemeetings[i].PopulateGamemeeting(app.db)
		}
	}
	body, err := json.Marshal(gamemeetings)
	if err != nil {
		app.serverError(w, err)
	}
	app.responseJSON(w, body)
}

func (app *application) joinUserToGamemeeting(w http.ResponseWriter, r *http.Request) {
	meetID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	userID, err := strconv.Atoi(r.URL.Query().Get(":user_id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	var meetMod models.Gamemeeting
	var userMod models.User
	err = app.db.First(&meetMod, meetID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.First(&userMod, userID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	err = meetMod.AddUser(app.db, userMod)
	if err != nil && err.Error() == "Gamemeeting Model: Player is the same as owner" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (app *application) disjoinUserToGamemeeting(w http.ResponseWriter, r *http.Request) {
	meetID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	userID, err := strconv.Atoi(r.URL.Query().Get(":user_id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	var meetMod models.Gamemeeting
	var userMod models.User
	err = app.db.First(&meetMod, meetID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.First(&userMod, userID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	err = meetMod.RemoveUser(app.db, userMod)
	if err != nil && err.Error() == "Gamemeeting Model: Player is the same as owner" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (app *application) userParticipatesGamemeetings(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var reqUser models.User
	err = app.db.First(&reqUser, userID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	var gamemeetings []models.Gamemeeting
	var userParticipates []models.Gamemeeting

	result := app.db.Find(&gamemeetings)
	if result.Error != nil {
		app.serverError(w, result.Error)
		return
	}
	if result.RowsAffected >= 0 {
		for i := range gamemeetings {
			gamemeetings[i].PopulateGamemeeting(app.db)
			for _, player := range gamemeetings[i].Players {
				if player.ID == reqUser.ID {
					userParticipates = append(userParticipates, gamemeetings[i])
					continue
				}
			}
		}
	}
	body, err := json.Marshal(userParticipates)
	if err != nil {
		app.serverError(w, err)
	}
	app.responseJSON(w, body)
}

func (app *application) deleteUserGamemeeting(w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	headerData := r.Header.Get("Authorization")
	userID, err := strconv.Atoi(headerData)
	if err != nil {
		app.clientError(w, "Invalid Header Data", 401)
		return
	}
	var gamemeeting models.Gamemeeting
	err = app.db.First(&gamemeeting, gameID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	err = gamemeeting.PopulateGamemeeting(app.db)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var ownerUser models.User
	err = app.db.First(&ownerUser, userID).Error
	if err != nil && err.Error() == "record not found" {
		app.clientError(w, err.Error(), 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	if ownerUser.ID != gamemeeting.OwnerID {
		app.clientError(w, "Unathorized to eliminate gamemeeting", 401)
		return
	}
	err = app.db.Model(&gamemeeting).Association("Players").Clear().Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.db.Delete(&gamemeeting).Error
	if err != nil {
		app.serverError(w, err)
		return
	}
	body, err := json.Marshal(gamemeeting)
	if err != nil {
		app.serverError(w, err)
	}
	app.responseJSON(w, body)
}

func (app *application) preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
