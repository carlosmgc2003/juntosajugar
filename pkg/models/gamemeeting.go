package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	InvalidMeetingPlace    = errors.New("Gamemeeting Model: Invalid meeting place")
	InvalidMeetingSchedule = errors.New("Gamemeeting Model: Invalid meeting time")
	InvalidMeetingOwner    = errors.New("Gamemeeting Model: Invalid Owner ID")
	InvalidBoardgame       = errors.New("Gamemeeting Model: Invalid Boardgame ID")
	InvalidMaxPlayers      = errors.New("Gamemeeting Model: Invalid Max Players Quantity")
	InvalidTooManyPlayers  = errors.New("Gamemeeting Model: Too many players JSON")
	InvalidPlayer          = errors.New("Gamemeeting Model: Invalid Player")
)

func (G *Gamemeeting) FromJson(requestBody []byte, db *gorm.DB) error {
	type tempStruct struct {
		Place      string `json:"place"`
		Scheduled  string `json:"scheduled"`
		OwnerId    uint   `json:"owner"`
		GameId     uint   `json:"game"`
		Players    []uint `json:"players"`
		MaxPlayers uint   `json:"max_players"`
	}
	var temp tempStruct
	err := json.Unmarshal(requestBody, &temp)
	if err != nil {
		return err
	}
	if !validMeetingPlace(temp.Place) {
		return InvalidMeetingPlace
	}
	G.Place = temp.Place
	G.Scheduled, err = time.Parse(time.RFC3339, temp.Scheduled)
	if err != nil || !validScheduledTime(G.Scheduled) {
		return InvalidMeetingSchedule
	}
	if !validUser(temp.OwnerId, db) {
		return InvalidMeetingOwner
	}
	// Asignar el owner al boardgame
	var tempOwner User
	err = db.First(&tempOwner, temp.OwnerId).Error
	if err != nil {
		return err
	}
	G.Owner = tempOwner

	if !validGame(temp.GameId, db) {
		return InvalidBoardgame
	}

	var tempBoardgame Boardgame
	err = db.First(&tempBoardgame, temp.GameId).Error
	if err != nil {
		return err
	}
	G.Boardgame = tempBoardgame

	if !validMaxPlayers(temp.MaxPlayers) {
		return InvalidMaxPlayers
	}
	G.MaxPlayers = temp.MaxPlayers

	if len(temp.Players) > int(temp.MaxPlayers) {
		return InvalidTooManyPlayers
	}
	for _, playerId := range temp.Players {
		var tempPlayer User
		if !validUser(playerId, db) {
			return InvalidPlayer
		}
		err = db.First(&tempPlayer, playerId).Error
		if err != nil {
			return err
		}
		G.Players = append(G.Players, tempPlayer)
	}

	return err
}

func validMeetingPlace(meetingplace string) bool {
	var characters = len(meetingplace)
	return characters <= 100 && characters >= 5
}

func validScheduledTime(scheduledtime time.Time) bool {
	return time.Now().Before(scheduledtime)
}

func validUser(ownerid uint, db *gorm.DB) bool {
	var owner User
	result := db.First(&owner, ownerid)
	// Hay que devolver not error is por que la logica de la func es preguntar si es valido el owner
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func validGame(gameid uint, db *gorm.DB) bool {
	var game Boardgame
	result := db.First(&game, gameid)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func validMaxPlayers(maxplayers uint) bool {
	return maxplayers <= 12
}
