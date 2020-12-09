package main

import (
	"flag"
	"juntosajugar/pkg/models"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	db       *gorm.DB
}

// RETRIES es la cantidad de veces que la API intenta conectarse con la BD
const RETRIES int = 10

func main() {
	// Obtener y parsear el numero de puerto de la linea de comandos
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	// Loggers de informacion y errores para obtener visibilidad de servidor
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Inicializacion de la base de datos con el ORM
	// Intentamos varias veces para solucionar la demora que requiere el contenedor de mySql para arrancar.
	var db *gorm.DB
	var err error
	for i := 0; i < RETRIES; i++ {
		db, err = gorm.Open("mysql", "api_web:api_web_pass@tcp(db:3306)/juntosajugar?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			infoLog.Printf("Intento de conexion %d de %d", i+1, RETRIES)
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}
	if err != nil {
		errorLog.Fatal(err)
	}

	db.DropTableIfExists(&models.Gamemeeting{}, &models.Boardgame{}, &models.User{})
	db.AutoMigrate(&models.Gamemeeting{}, &models.Boardgame{}, &models.User{})
	defer db.Close()

	// Migracion de los modelos de datos
	db.AutoMigrate(&models.User{}, &models.Boardgame{}, &models.Gamemeeting{})
	db.Model(&models.Gamemeeting{}).AddForeignKey("owner_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.Gamemeeting{}).AddForeignKey("boardgame_id", "boardgames(id)", "CASCADE", "CASCADE")
	db.Table("user_gamemeeting").AddForeignKey("gamemeeting_id", "gamemeetings(id)", "CASCADE", "CASCADE")
	db.Table("user_gamemeeting").AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	// Objeto de tipo aplicacion con sus dos atributos.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.HttpOnly = false

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		db:       db,
	}

	// Creacion del usuario Admin de JaJ
	hashedAdminPass, _ := bcrypt.GenerateFromPassword([]byte("123456"), 12)

	var adminUser models.User = models.User{
		Name:           "Administrador",
		Email:          "admin@juntosajugar.com",
		HashedPassword: string(hashedAdminPass),
	}

	app.db.Create(&adminUser)

	// Parametros para el objeto http.Server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Inicializacion del servidor.
	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
