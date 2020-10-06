package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

// Errores especificos del modelo user
var (
	InvalidUserName    = errors.New("User Model: Invalid Username")
	InvalidEmail       = errors.New("User Model: Invalid Email")
	InvalidUserPic     = errors.New("User Model: Invalid Filename")
	InvalidIdToken     = errors.New("User Model: Invalid Id Token")
	InvalidPlainPass   = errors.New("User Model: Invalid Pass")
	InvalidHashPass    = errors.New("User Model: Error while hashing password")
	InvalidCredentials = errors.New("User Model: Invalid Credentials")
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (L *Login) FromJson(requestBody []byte) error {
	var tempLogin Login
	err := json.Unmarshal(requestBody, &tempLogin)
	if err != nil {
		return err
	}
	if !validEmail(tempLogin.Email) {
		return InvalidEmail
	}
	if !validPassword(tempLogin.Password) {
		return InvalidPlainPass
	}
	err = copier.Copy(L, tempLogin)
	return err
}

func (U *User) Authenticate(loginData Login) error {
	err := bcrypt.CompareHashAndPassword([]byte(U.HashedPassword), []byte(loginData.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return InvalidCredentials
	}
	return err
}

func (U *User) FromJson(requestBody []byte) error {
	var tempUser User
	err := json.Unmarshal(requestBody, &tempUser)
	if err != nil {
		return err
	}
	if !validUsername(tempUser.Name) {
		return InvalidUserName
	}
	U.Name = tempUser.Name
	if !validEmail(tempUser.Email) {
		return InvalidEmail
	}
	U.Email = tempUser.Email
	if !validUserPic(tempUser.DisplayPic) {
		return InvalidUserPic
	}
	U.DisplayPic = tempUser.DisplayPic
	if !validIdToken(tempUser.IdToken) {
		return InvalidIdToken
	}
	U.IdToken = tempUser.IdToken
	if !validPassword(tempUser.HashedPassword) {
		return InvalidPlainPass
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(tempUser.HashedPassword), 12)
	if err != nil {
		return InvalidHashPass
	}
	U.HashedPassword = string(hash)
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

func validPassword(pass string) bool {
	return len(pass) >= 6
}
