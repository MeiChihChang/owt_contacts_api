package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MeiChihChang/owt_contacts_api/internal/repository"
	"github.com/MeiChihChang/owt_contacts_api/internal/repository/dbrepo"
)

const portNumber = 8080

type application struct {
	DSN          string
	DB           repository.DatabaseRepo
	auth         Auth
	Domain       string
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

var app application

func main() {
	// set application config

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=owt sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "owt.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "owt.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "localhost", "domain")
	flag.Parse()

	// connect to the database
	log.Println("Connecting to SQL database...")
	conn, err := app.ConnectDB(app.DSN)
	if err != nil {
		// retry
		for i := 0; i < 5; i++ {
			conn, err = app.ConnectDB(app.DSN)
			if err == nil {
				break
			}
			log.Println("retrying after error:", err)
			time.Sleep(60 * time.Second)
		}
		log.Fatal("Cannot connect to database! Dying...")
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer conn.Close()

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 60,
		RefreshExpiry: time.Hour * 24,
		CookieDomain:  app.CookieDomain,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
	}

	// start a web server
	srv := &http.Server{
		Addr:    app.Domain + fmt.Sprintf(":%v", portNumber),
		Handler: routes(&app),
	}

	log.Printf("Starting application at %s", app.Domain+fmt.Sprintf(":%v", portNumber))

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
