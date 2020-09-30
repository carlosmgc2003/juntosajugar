package models

import (
	"encoding/json"
	"errors"
	"regexp"
)

// Errores especificos del modelo user
var (
	InvalidUserName = errors.New("User Model: Invalid Username")
	InvalidEmail    = errors.New("User Model: Invalid Email")
	InvalidUserPic  = errors.New("User Model: Invalid Filename")
	InvalidIdToken  = errors.New("User Model: Invalid Id Token")
)

func (U *User) FromJson(requestBody []byte) error {
	var temp_user User
	err := json.Unmarshal(requestBody, &temp_user)
	if err != nil {
		return err
	}
	U.ID = temp_user.ID
	if !validUsername(temp_user.Name) {
		return InvalidUserName
	}
	U.Name = temp_user.Name
	if !validEmail(temp_user.Email) {
		return InvalidEmail
	}
	U.Email = temp_user.Email
	if !validUserPic(temp_user.DisplayPic) {
		return InvalidUserPic
	}
	U.DisplayPic = temp_user.DisplayPic
	if !validIdToken(temp_user.IdToken) {
		return InvalidIdToken
	}
	U.IdToken = temp_user.IdToken
	return err
}

/*
Nombres Occidentales validos ja!
Mathias d'Arras
Martin Luther King, Jr.
Hector Sausage-Hausen
*/

func validUsername(username string) bool {
	// Validar el nombre de usuario de acuerdo a ese regex.
	var re = regexp.MustCompile("^[[:alpha:] ,.'-]+$")
	return re.MatchString(username) && len(username) <= 30 && len(username) >= 6
}

func validEmail(email string) bool {
	// Validar el formato de la direccion de email.
	var re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email) && len(email) <= 100
}

func validUserPic(filename string) bool {
	// Valida el nombre de archivo de la imagen subida
	return len(filename) <= 100
}

func validIdToken(token string) bool {
	return len(token) <= 1500
}
