package query

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/query/adapters"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
)

// Load the query from Create-Tables.sql

// database is the package global db  - this reference is not exported outside the package.
var database adapters.Database

// OpenDatabase opens the database with the given options
func OpenDatabase(opts map[string]string, mu *sync.RWMutex) error {

	// If we already have a db, return it
	if database != nil {
		return fmt.Errorf("query: database already open - %s", database)
	}

	// Assign the db global in query package
	switch opts["adapter"] {
	case "sqlite3":
		database = &adapters.SqliteAdapter{
			Mutex: mu,
		}
	case "mysql":
		database = &adapters.MysqlAdapter{}
	case "postgres":
		database = &adapters.PostgresqlAdapter{}
	default:
		database = nil // fail
	}

	if database == nil {
		return fmt.Errorf("query: database adapter not recognised - %s", opts)
	}

	// Ask the db adapter to open
	err := database.Open(opts)

	if err == nil {
		// Create table if it doesn't exist

		b, err := os.ReadFile("db/Create-Tables.sql")

		if err == nil {
			statement, err := database.SQLDB().Prepare(string(b))

			if err == nil {
				statement.Exec()
				log.Info(log.V{"msg": "Finished creating tables"})
			} else {
				log.Error(log.V{"Error creating database tables": err})
			}
		} else {
			log.Error(log.V{"Unable to read the database file": err})
		}

		// Set max connections for sqlite
		if opts["adapter"] == "sqlite3" {
			SetMaxOpenConns(1)
		}

	} else {
		log.Error(log.V{"Error opening database": err})
	}

	return err
}

// CloseDatabase closes the database opened by OpenDatabase
func CloseDatabase() error {
	var err error
	if database != nil {
		err = database.Close()
		database = nil
	}

	return err
}

// SetMaxOpenConns sets the maximum number of open connections
func SetMaxOpenConns(max int) {
	database.SQLDB().SetMaxOpenConns(max)
}

// QuerySQL executes the given sql Query against our database, with arbitrary args
func QuerySQL(query string, args ...interface{}) (*sql.Rows, error) {
	if database == nil {
		return nil, fmt.Errorf("query: QuerySQL called with nil database")
	}
	results, err := database.Query(query, args...)
	return results, err
}

// ExecSQL executes the given sql against our database with arbitrary args
// NB returns sql.Result - not to be used when rows expected
func ExecSQL(query string, args ...interface{}) (sql.Result, error) {
	if database == nil {
		return nil, fmt.Errorf("query: ExecSQL called with nil database")
	}
	results, err := database.Exec(query, args...)
	return results, err
}

// TimeString returns a string formatted as a time for this db
// if the database is nil, an empty string is returned.
func TimeString(t time.Time) string {
	if database != nil {
		return database.TimeString(t)
	}
	return ""
}
