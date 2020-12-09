package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	// Extiende a gorm para que use mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	//ErrDuplicate es el Error de duplicado en la BD
	ErrDuplicate = errors.New("models: duplicate insertion in DB")
)

// User es el struct de Usuario de juntos a jugar
type User struct {
	gorm.Model                 //ya tiene el id
	Name           string      `gorm:"unique;size:100" json:"name"`
	Email          string      `gorm:"unique;not null;size:100" json:"email"`
	DisplayPic     string      `gorm:"size:100" json:"display_pic_route"`
	IDtoken        string      `gorm:"size:1500" json:"tokenId"`
	HashedPassword string      `gorm:"size:60" json:"password"`
	UserBoardgames []Boardgame `gorm:"many2many:user_boardgames;"`
}

// Boardgame es el Juego de mesa
type Boardgame struct {
	gorm.Model        //ya tiene el id
	Name       string `gorm:"unique;not null;size:100" json:"name"`
	Class      string `gorm:"not null;size:40" json:"class"`
	DisplayPic string `gorm:"unique;size:100" json:"display_pic_route"`
}

// Gamemeeting Reuniones para jugar
type Gamemeeting struct {
	gorm.Model            //ya tiene el id
	Place       string    `gorm:"unique;not null;size:100" json:"place"`
	Scheduled   time.Time `gorm:"not null" json:"scheduled"`
	OwnerID     uint
	Owner       User `gorm:"foreignKey:OwnerId"`
	BoardgameID uint
	Boardgame   Boardgame `gorm:"foreignKey:BoardgameId"`
	Players     []User    `gorm:"many2many:user_gamemeeting;"`
	MaxPlayers  uint      `gorm:"not null;" json:"max_players"`
}
