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
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *gorm.DB
}

func main() {
	// Obtener y parsear el numero de puerto de la linea de comandos
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Loggers de informacion y errores para obtener visibilidad de servidor
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Inicializacion de la base de datos con el ORM
	db, err := gorm.Open("mysql", "root:admin@/jaj?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Migracion de los modelos de datos
	db.AutoMigrate(&models.User{}, &models.Boardgame{}, &models.GameMeeting{})

	// Objeto de tipo aplicacion con sus dos atributos.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Parametros para el objeto http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Inicializacion del servidor.
	infoLog.Printf("Starting server on https://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
