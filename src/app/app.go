package app

import (
	"os"
	"sync"
	"time"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/assets"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/config"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
)

// appAssets holds a reference to our assets for use in asset setup
var appAssets *assets.Collection

// Setup sets up our application
func Setup(mu *sync.RWMutex) {
	// Setup log
	err := SetupLog()
	if err != nil {
		println("failed to set up logs %s", err)
		os.Exit(1)
	}

	// Log server startup
	msg := "Starting server"
	if config.Production() {
		msg = msg + " in production"
	}

	log.Info(log.Values{"msg": msg, "port": config.Get("port")})
	defer log.Time(time.Now(), log.Values{"msg": "Finished loading server"})

	// Set up external service interfaces (twitter)
	SetupServices()

	// Set up our assets
	SetupAssets()

	// Setup our view templates
	SetupView()

	// Setup our database
	SetupDatabase(mu)

	// Setup our authentication and authorisation
	SetupAuth()

	// Setup our router and handlers
	SetupRoutes()

}

// SetupLog sets up logging
func SetupLog() error {

	// Set up a stderr logger with time prefix
	logger, err := log.NewStdErr(log.PrefixDateTime)
	if err != nil {
		return err
	}
	log.Add(logger)

	// Set up a file logger pointing at the right location for this config.
	fileLog, err := log.NewFile(config.Get("log"))
	if err != nil {
		return err
	}
	log.Add(fileLog)

	return nil
}
