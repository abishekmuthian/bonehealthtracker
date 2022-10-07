package app

import (
	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux/middleware/gzip"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux/middleware/secure"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/session"
	useractions "github.com/abishekmuthian/bonehealthtracker/src/users/actions"

	// Resource Actions
	appactions "github.com/abishekmuthian/bonehealthtracker/src/app/actions"
)

// SetupRoutes creates a new router and adds the routes for this app to it.
func SetupRoutes() *mux.Mux {
	router := mux.New()
	mux.SetDefault(router)

	// Add the home page route
	router.Get("/", appactions.HandleHome)

	// Add user actions
	router.Post("/users/upload", useractions.HandleUpload)

	// Add the legal page route
	router.Get("/legal", HandleLegal)

	// Add a route to handle static files
	router.Get("/favicon.ico", fileHandler)
	router.Get("/icons/{path:.*}", fileHandler)
	router.Get("/files/{path:.*}", fileHandler)
	router.Get("/assets/{path:.*}", fileHandler)

	// Set the default file handler
	router.FileHandler = fileHandler
	router.ErrorHandler = errHandler

	// Add middleware
	router.AddMiddleware(log.Middleware)
	router.AddMiddleware(session.Middleware)
	router.AddMiddleware(gzip.Middleware)
	router.AddMiddleware(secure.Middleware)

	return router
}
