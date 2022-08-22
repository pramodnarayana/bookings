package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pramodnarayana/bookings/pkg/config"
	"github.com/pramodnarayana/bookings/pkg/handlers"
	"github.com/pramodnarayana/bookings/pkg/render"
)

const port = ":8080"

var app config.AppConfig

var session *scs.SessionManager

// This is main function
func main() {

	//Change this to true when in Production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour

	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := render.CreateParsedTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	render.NewTemplates(&app)
	handlers.NewHandlers(repo)

	fmt.Printf("Starting application on port %v", port)

	httpServer := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
