package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var (
	ErrDuplicate = errors.New("models: duplicate insertion in DB")
)

// Usuario de juntos a jugar
type User struct {
	gorm.Model        //ya tiene el id
	Name       string `gorm:"unique;size:100" json:"name"`
	Email      string `gorm:"unique;not null;size:100" json:"email"`
	DisplayPic string `gorm:"size:100" json:"display_pic_route"`
}

// Juegos de mesa
type Boardgame struct {
	gorm.Model        //ya tiene el id
	Name       string `gorm:"unique;not null;size:100" json:"name"`
	Class      string `gorm:"not null;size:40" json:"class"`
	DisplayPic string `gorm:"unique;size:100" json:"display_pic_route"`
}

// Reuniones para jugar
type GameMeeting struct {
	gorm.Model           //ya tiene el id
	Place      string    `gorm:"unique;not null;size:100" json:"name"`
	Scheduled  time.Time `json:"scheduled"`
	Owner      User      `gorm:"foreignkey:ID;not null" json:"owner_id"`
	Players    []User    `gorm:"foreignkey:ID;not null" json:"players"`
	Game       Boardgame `gorm:"foreignkey:ID;not null" json:"game"`
	MaxPlayers int       `gorm:"not null" json:"max_players"`
}
