package mssqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionString(t *testing.T) {
    // when
	cURL := newSQLServerConnectionURL("u", "p", "host", "mydb", 5432)
    //then
	assert.Equal(t, cURL.String(), "sqlserver://u:p@host:5432?database=mydb")
}
