/*
Package main provides the application entry point. It spins up a server
and fires off all route handlers.
*/
package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/rlr524/bookings/pkg/config"
	"github.com/rlr524/bookings/pkg/handlers"
	"github.com/rlr524/bookings/pkg/template"
	"log"
	"net/http"
	"time"
)

// In Go, the colon is required when passing the port number in a variable
const portNumber = ":3030"

var app config.AppConfig
var sessionManager *scs.SessionManager

// main is the application entry point
func main() {
	// Set this to true when deploying to prod
	app.InProduction = false

	// Use scs package to initialize a new session manager and configure the session lifetime
	// of 24 hours as well as persist the cookie after the browser window is closed, allow lax
	// treatment for same site cookie enforcement, and not require secure cookies in dev
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = app.InProduction

	app.Session = sessionManager

	tc, err := template.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	// Create a repository, assign it to repo and pass it back to the NewHandlers function
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// Calls the NewTemplates function from the template package and references the AppConfig
	// Need to reference app because we are setting app as a pointer to config.AppConfig above
	// in that it is a var named "app" that is of type config.AppConfig
	template.NewTemplates(&app)

	// Create a server using our routes and pat router in the routes package
	srv := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	err = srv.ListenAndServe()
	log.Fatal(err)
}
