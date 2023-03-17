package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/driver"
	"github.com/oskarpedosk/baltijas-kauss/internal/handlers"
	"github.com/oskarpedosk/baltijas-kauss/internal/helpers"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
	"github.com/oskarpedosk/baltijas-kauss/internal/render"

	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting application on http://localhost%s\n", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}


func run() (*driver.DB, error) {
	// What am I going to put in the session
	gob.Register(models.User{})
	gob.Register(models.Badge{})
	gob.Register(models.Player{})
	gob.Register(models.Team{})
	gob.Register(models.Standings{})
	gob.Register(models.Result{})
	gob.Register(models.Season{})
	gob.Register(models.Positions{})
	gob.Register(models.Attributes{})
	gob.Register(models.PlayerWithTeamInfo{})
	gob.Register(models.Filter{})

	// Read flags
	inProduction := flag.Bool("production", false, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "enable", "Database SSL settings (disable, preger, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	// Change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Set up the session
	session = scs.New()
	session.Lifetime = 168 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	fmt.Println(connectionString)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	log.Println("Connected to database!")

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = templateCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
