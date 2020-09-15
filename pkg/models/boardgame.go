package models

import (
	"encoding/json"
	"errors"
)

// Errores especificos del modelo boardgame
var (
	InvalidGameName  = errors.New("BoardGame Model: Invalid game name")
	InvalidGameClass = errors.New("BoardGame Model: Invalid game class")
	InvalidGamePic   = errors.New("BoardGame Model: Invalid game picture filename")
)

func (B *Boardgame) FromJson(requestBody []byte) error {
	// Si el array de byte contiene un JSON correcto de instancia de Boardgame
	// ingresa los datos del mismo en la instancia.
	var tempBoardgame Boardgame
	err := json.Unmarshal(requestBody, &tempBoardgame)
	if err != nil {
		return err
	}
	B.ID = tempBoardgame.ID
	if !validGameName(tempBoardgame.Name) {
		return InvalidGameName
	}
	B.Name = tempBoardgame.Name
	if !validGameClass(tempBoardgame.Class) {
		return InvalidGameClass
	}
	B.Class = tempBoardgame.Class
	if !validGamePic(tempBoardgame.DisplayPic) {
		return InvalidGamePic
	}
	B.DisplayPic = tempBoardgame.DisplayPic
	return err
}

func validGameName(gamename string) bool {
	// Validar un string de nombre de juego
	return len(gamename) <= 100 && len(gamename) > 2
}

func validGameClass(gameclass string) bool {
	// Validar si el string de clase esta dentro de las opciones
	gameclasses := []string{"Ingenio", "Estrategia", "Clasico", "Dados", "Palabras", "Cartas"}
	_, valid := find(gameclasses, gameclass)
	return valid
}

func find(slice []string, val string) (int, bool) {
	// Encontar una string en un slice de strings
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func validGamePic(filename string) bool {
	// Validar el nombre de archivo
	return len(filename) >= 10 && len(filename) <= 50
}
