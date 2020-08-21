package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

// struct que engloba a todos los campos comunes de las demas. Los atributos escapados
// son para la interpretacion de librerias, en este caso gorm y json
type commonModelFields struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Usuario de juntos a jugar
type User struct {
	commonModelFields        //ya tiene el id
	Name              string `gorm:"unique;size:100" json:"name"`
	Email             string `gorm:"unique;not null;size:100" json:"email"`
	Display_pic       string `gorm:"unique; size:100" json:"display_pic_route"`
}

// Juegos de mesa
type Boardgame struct {
	commonModelFields        //ya tiene el id
	Name              string `gorm:"unique;not null;size:100" json:"name"`
	Class             string `gorm:"not null;size:40" json:"class"`
	Display_pic       string `gorm:"unique;size:100" json:"display_pic_route"`
}

// Reuniones para jugar
type GameMeeting struct {
	commonModelFields           //ya tiene el id
	Place             string    `gorm:"unique;not null;size:100" json:"name"`
	scheduled         time.Time `json:"scheduled"`
	Owner             User      `gorm:"foreignkey:ID;not null" json:"owner_id"`
	Players           []User    `gorm:"foreignkey:ID;not null" json:"players"`
	Game              Boardgame `gorm:"foreignkey:ID;not null" json:"game"`
	Max_Players       int       `gorm:"not null" json:"max_players"`
}
