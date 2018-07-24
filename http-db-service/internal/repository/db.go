package repository

import (
	"database/sql"
	"fmt"
	"io"
	log "github.com/Sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"github.com/pkg/errors"

	dbPkg "github.com/kyma-project/examples/http-db-service/internal/mssqldb"
)

const (
	insertQuery   = "INSERT INTO %s (order_id, namespace, total) VALUES (?, ?, ?)"
	getQuery      = "SELECT * FROM %s"
	getNSQuery    = "SELECT * FROM %s WHERE namespace = ?"
	deleteQuery   = "DELETE FROM %s"
	deleteNSQuery = "DELETE FROM %s WHERE namespace = ?"

)

type orderRepositorySQL struct {
	database        dbQuerier
	ordersTableName string
}

//go:generate mockery -name dbQuerier -inpkg
type dbQuerier interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	io.Closer
}

// NewOrderRepositoryDb is used to instantiate and return the DB implementation of the OrderRepository.
// The connection to the database is created by initiating the configuration defined
// in https://github.com/kyma-project/examples/blob/master/http-db-service/internal/mssqldb/config.go
func NewOrderRepositoryDb() (OrderRepository, error) {
	var (
		dbCfg dbPkg.Config
		database dbQuerier
		err error
	)
	if err = envconfig.Init(&dbCfg); err != nil {
		return nil, errors.Wrap(err, "Error loading db configuration %v.")
	}
	if database, err = dbPkg.InitDb(dbCfg); err != nil {
		return nil, errors.Wrap(err, "Error loading db configuration %v.")
	}
	return &orderRepositorySQL{database, dbCfg.DbOrdersTableName}, nil
}

type sqlError interface {
	sqlErrorNumber() int32
}

func (repository *orderRepositorySQL) InsertOrder(order Order) error {
	q := fmt.Sprintf(insertQuery, dbPkg.SanitizeSQLArg(repository.ordersTableName))
	log.Debugf("Running insert order query: '%q'.", q)
	_, err := repository.database.Exec(q, order.OrderId, order.Namespace, order.Total)

	if errorWithNumber, ok := err.(sqlError); ok {
		if errorWithNumber.sqlErrorNumber() == dbPkg.PrimaryKeyViolation {
			return ErrDuplicateKey
		}
	}

	return errors.Wrap(err, "while inserting order")
}

func (repository *orderRepositorySQL) GetOrders() ([]Order, error) {
	q := fmt.Sprintf(getQuery, dbPkg.SanitizeSQLArg(repository.ordersTableName))
	log.Debugf("Quering orders: '%q'.", q)
	rows, err := repository.database.Query(q)

	if err != nil {
		return nil, errors.Wrap(err, "while reading orders from DB")
	}

	defer rows.Close()
	return readFromResult(rows)
}

func (repository *orderRepositorySQL) GetNamespaceOrders(ns string) ([]Order, error) {
	q := fmt.Sprintf(getNSQuery, dbPkg.SanitizeSQLArg(repository.ordersTableName))
	log.Debugf("Quering orders for namespace: '%q'.", q)
	rows, err := repository.database.Query(q, ns)

	if err != nil {
		return nil, errors.Wrapf(err, "while reading orders for namespace: '%q' from DB", ns)
	}

	defer rows.Close()
	return readFromResult(rows)
}

func (repository *orderRepositorySQL) DeleteOrders() error {
	q := fmt.Sprintf(deleteQuery, dbPkg.SanitizeSQLArg(repository.ordersTableName))
	log.Debugf("Deleting orders: '%q'.", q)
	_, err := repository.database.Exec(q)

	if err != nil {
		return errors.Wrap(err, "while deleting orders")
	}
	return nil
}

func (repository *orderRepositorySQL) DeleteNamespaceOrders(ns string) error {
	q := fmt.Sprintf(deleteNSQuery, dbPkg.SanitizeSQLArg(repository.ordersTableName))
	log.Debugf("Deleting orders: '%q'.", q)
	_, err := repository.database.Exec(q, ns)

	if err != nil {
		return errors.Wrap(err, "while deleting orders")
	}
	return nil
}

func readFromResult(rows *sql.Rows) ([]Order, error) {
	orderList := make([]Order, 0)
	for rows.Next() {
		order := Order{}
		if err := rows.Scan(&order.OrderId, &order.Namespace, &order.Total); err != nil {
			return []Order{}, err
		}
		orderList = append(orderList, order)
	}
	return orderList, nil
}

func (repository *orderRepositorySQL) cleanUp() error {
	log.Debug("Removing DB table")

	if _, err := repository.database.Exec("DROP TABLE " + dbPkg.SanitizeSQLArg(repository.ordersTableName)); err != nil {
		return errors.Wrap(err, "while removing the DB table.")
	}
	if err := repository.database.Close(); err != nil {
		return errors.Wrap(err, "while closing connection to the DB.")
	}
	return nil
}
