package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	config2 "github.com/liomazza/bookings/internal/config"
	handlers2 "github.com/liomazza/bookings/internal/handlers"
	models2 "github.com/liomazza/bookings/internal/models"
	render2 "github.com/liomazza/bookings/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"
var app config2.AppConfig
var session *scs.SessionManager

//main is the main application function
func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {

	//what i am going to put in the session
	gob.Register(models2.Reservation{})

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	myTemplateCache, err := render2.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot crete temlate cache")
		return err
	}

	app.TemplateCache = myTemplateCache
	app.UseCache = false

	repo := handlers2.NewRepo(&app)
	handlers2.NewHandlers(repo)
	render2.NewTemplates(&app)

	return nil
}
