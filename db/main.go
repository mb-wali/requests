package db

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/cyverse-de/dbutil"
	"github.com/pkg/errors"
)

// A convenient shorthand for a statement builder that works with Postgres.
var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// InitDatabase establishes a database connection and verifies that the database can be reached.
func InitDatabase(driverName, databaseURI string) (*sql.DB, error) {
	wrapMsg := "unable to initialize the database"

	// Create a database connector to establish the connection for us.
	connector, err := dbutil.NewDefaultConnector("1m")
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	// Establish the database connection.
	db, err := connector.Connect(driverName, databaseURI)
	if err != nil {
		return nil, errors.Wrap(err, wrapMsg)
	}

	return db, nil
}
