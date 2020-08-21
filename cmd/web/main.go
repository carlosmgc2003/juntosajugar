package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	"juntosajugar/pkg/models"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	db *gorm.DB
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)


	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := gorm.Open("mysql", "root:admin@/jaj?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.User{},&models.Boardgame{},&models.GameMeeting{})

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
	}


	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on https://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}