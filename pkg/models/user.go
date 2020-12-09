package models

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

// Errores especificos del modelo user
var (
	ErrInvalidUserName    = errors.New("User Model: Invalid Username")
	ErrInvalidEmail       = errors.New("User Model: Invalid Email")
	ErrInvalidUserPic     = errors.New("User Model: Invalid Filename")
	ErrInvalidIDToken     = errors.New("User Model: Invalid Id Token")
	ErrInvalidPlainPass   = errors.New("User Model: Invalid Pass")
	ErrInvalidHashPass    = errors.New("User Model: Error while hashing password")
	ErrInvalidCredentials = errors.New("User Model: Invalid Credentials")
	ErrInvalidJSON        = errors.New("User Model: Invalid JSON")
)

// Login es un struct auxiliar para recibir los datos de login de usuario
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// FromJSON toma una slice de bytes JSON de login y guarda en el struct los datos.
func (L *Login) FromJSON(requestBody []byte) error {
	var tempLogin Login
	err := json.Unmarshal(requestBody, &tempLogin)
	if err != nil {
		return err
	}
	if !validEmail(tempLogin.Email) {
		return ErrInvalidEmail
	}
	if !validPassword(tempLogin.Password) {
		return ErrInvalidPlainPass
	}
	err = copier.Copy(L, tempLogin)
	return err
}

// Authenticate Valida el password con el hash guardado en la BD.
func (U *User) Authenticate(loginData Login) error {
	err := bcrypt.CompareHashAndPassword([]byte(U.HashedPassword), []byte(loginData.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidCredentials
	}
	return err
}

// FromJSON toma un slice de byte con un JSON de usuario y lo carga en la instancia.
func (U *User) FromJSON(requestBody []byte) error {
	var tempUser User
	err := json.Unmarshal(requestBody, &tempUser)
	if err != nil {
		return ErrInvalidJSON
	}
	if !validUsername(tempUser.Name) {
		return ErrInvalidUserName
	}
	U.Name = tempUser.Name
	if !validEmail(tempUser.Email) {
		return ErrInvalidEmail
	}
	U.Email = tempUser.Email
	if !validUserPic(tempUser.DisplayPic) {
		return ErrInvalidUserPic
	}
	U.DisplayPic = tempUser.DisplayPic
	if !validIDToken(tempUser.IDtoken) {
		return ErrInvalidIDToken
	}
	U.IDtoken = tempUser.IDtoken
	if !validPassword(tempUser.HashedPassword) {
		return ErrInvalidPlainPass
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(tempUser.HashedPassword), 12)
	if err != nil {
		return ErrInvalidHashPass
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
	return re.MatchString(username) && len(username) <= 100 && len(username) >= 6
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

func validIDToken(token string) bool {
	return len(token) <= 1500
}

func validPassword(pass string) bool {
	return len(pass) >= 6
}
