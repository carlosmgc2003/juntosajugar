package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	ErrDuplicate = errors.New("models: duplicate insertion in DB")
)

// Usuario de juntos a jugar
type User struct {
	gorm.Model            //ya tiene el id
	Name           string `gorm:"unique;size:100" json:"name"`
	Email          string `gorm:"unique;not null;size:100" json:"email"`
	DisplayPic     string `gorm:"size:100" json:"display_pic_route"`
	IdToken        string `gorm:"size:1500" json:"tokenId"`
	HashedPassword string `gorm:"size:60" json:"password"`
}

// Juegos de mesa
type Boardgame struct {
	gorm.Model        //ya tiene el id
	Name       string `gorm:"unique;not null;size:100" json:"name"`
	Class      string `gorm:"not null;size:40" json:"class"`
	DisplayPic string `gorm:"unique;size:100" json:"display_pic_route"`
}

// Reuniones para jugar
type Gamemeeting struct {
	gorm.Model           //ya tiene el id
	Place      string    `gorm:"unique;not null;size:100" json:"place"`
	Scheduled  time.Time `json:"scheduled"`
	OwnerId    uint      `json:"owner"`
	Owner      User      `gorm:"foreignKey:OwnerId"`
	GameId     uint      `json:"game"`
	Boardgame  Boardgame `gorm:"foreignKey:OwnerId"`
	MaxPlayers uint      `gorm:"not null;" json:"max_players"`
}
