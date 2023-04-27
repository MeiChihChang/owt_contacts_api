package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MeiChihChang/owt_contacts_api/internal/repository"
	"github.com/MeiChihChang/owt_contacts_api/internal/repository/dbrepo"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"

	_ "github.com/MeiChihChang/owt_contacts_api/cmd/api/docs"
)

// @title           OWT Swagger API
// @version         1.0
// @description     This is a sample contacts server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

const portNumber = 8080

type application struct {
	DSN          string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

var app application

func main() {
	swagger := os.Getenv("ENABLE_SWAGGER")
	if swagger == "true" {
		Prepare()
	}

	// read from command line
	flag.StringVar(&app.DSN, "dsn", "host=database port=5432 user=postgres password=postgres dbname=owt sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "owt.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "owt.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
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
		Addr:    fmt.Sprintf(":%v", portNumber),
		Handler: routes(&app),
	}

	log.Printf("Starting application on port %s", fmt.Sprintf(":%v", portNumber))

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func Prepare() {

	// Use gin to create the server
	router := gin.Default()

	// Create other urls, etc but add swagger like this - notice url is removed
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8000")
}
