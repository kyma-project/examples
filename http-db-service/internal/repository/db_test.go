package repository

import (
	"errors"
	"testing"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	dbPkg "github.com/kyma-project/examples/http-db-service/internal/mssqldb"
)

var newOrder = Order{OrderId: "orderId1", Namespace: "N7", Total: 10}

const (
	parsedInsert = "INSERT INTO tableName (order_id, namespace, total) VALUES (?, ?, ?)"
	parsedGet    = "SELECT * FROM tableName"
	parsedDelete = "DELETE FROM tableName"
)

func TestDbCreateSuccess(t *testing.T) {
	databaseMock := mockDbQuerier{}
	repo := orderRepositorySQL{&databaseMock, "tableName"}

	databaseMock.On("Exec", parsedInsert, newOrder.OrderId, newOrder.Namespace, newOrder.Total).Return((sql.Result)(nil), nil)
	//when
	err := repo.InsertOrder(newOrder)
	//then
	assert.Nil(t, err)
}

type primaryKeyViolationError struct{error}
func (e primaryKeyViolationError) sqlErrorNumber() int32 {
	return dbPkg.PrimaryKeyViolation
}

func TestDbCreateDuplicate(t *testing.T) {
	databaseMock := mockDbQuerier{}
	repo := orderRepositorySQL{&databaseMock, "tableName"}

	databaseMock.On("Exec", parsedInsert, newOrder.OrderId, newOrder.Namespace, newOrder.Total).
		Return((sql.Result)(nil), primaryKeyViolationError{})
	//when
	err := repo.InsertOrder(newOrder)
	//then
	assert.EqualValues(t, ErrDuplicateKey, err)
}

type otherSQLError struct{error}
func (e otherSQLError) SQLErrorNumber() int32 {
	return 2
}
func TestDbRepositoryCreateOtherSqlError(t *testing.T) {
	databaseMock := mockDbQuerier{}
	repo := orderRepositorySQL{&databaseMock, "tableName"}

	databaseMock.On("Exec", parsedInsert, newOrder.OrderId, newOrder.Namespace, newOrder.Total).
		Return((sql.Result)(nil), otherSQLError{})
	//when
	err := repo.InsertOrder(newOrder)
	//then
	assert.NotEqual(t, ErrDuplicateKey, err)
}

func TestDbCreateError(t *testing.T) {
	databaseMock := mockDbQuerier{}
	repo := orderRepositorySQL{&databaseMock, "tableName"}

	databaseMock.On("Exec", parsedInsert, newOrder.OrderId, newOrder.Namespace, newOrder.Total).Return((sql.Result)(nil), errors.New("unexpected error"))
	//when
	err := repo.InsertOrder(newOrder)
	//then
	assert.EqualError(t, err, "while inserting order: unexpected error")
}

func TestDbGetError(t *testing.T) {
	databaseMock := mockDbQuerier{}
	repo := orderRepositorySQL{&databaseMock, "tableName"}

	databaseMock.On("Query", parsedGet).Return(&sql.Rows{}, errors.New("unexpected error"))
	//when
	_, err := repo.GetOrders()
	//then
	assert.Error(t, err)
}

func TestDeleteOrders(t *testing.T) {
	databaseMock := mockDbQuerier{}
	repo := orderRepositorySQL{&databaseMock, "tableName"}
	databaseMock.On("Exec", parsedDelete).Return((sql.Result)(nil), nil)

	//when
	err := repo.DeleteOrders()

	//then
	assert.NoError(t, err)
}
