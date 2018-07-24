package mssqldb

import (
	"database/sql"
	"github.com/pkg/errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	_ "github.com/denisenkom/go-mssqldb" //MSSQL driver initialization
)

const (
	dbDriverName = "mssql"
)

// InitDb creates and tests a database connection using the configuration given in dbConfig.
// After it establishes a connection it also ensures that the table exists.
func InitDb(dbConfig Config) (*sql.DB, error) {
	connectionURL := newSQLServerConnectionURL(
		dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Name, dbConfig.Port)
	createTableQuery := sqlServerTableCreationQuery

	log.Debugf("Establishing connection with '%s'. Connection string: '%q'", dbDriverName,
		strings.Replace(connectionURL.String(), connectionURL.User.String() + "@", "***:***@", 1))

	db, err := sql.Open(dbDriverName, connectionURL.String())
	if err != nil {
		return nil, errors.Wrapf(err, "while establishing connection to '%s'", dbDriverName)
	}

	log.Debug("Testing connection")
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "while testing DB connection")
	}

	q := strings.Replace(createTableQuery, "{name}", SanitizeSQLArg(dbConfig.DbOrdersTableName), -1)
	log.Debugf("Ensuring table exists. Running query: '%q'.", q)
	if _, err := db.Exec(q); err != nil {
		return nil, errors.Wrap(err, "while initiating DB table")
	}

	return db, nil
}
