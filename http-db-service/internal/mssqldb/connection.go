package mssqldb

import (
	"fmt"
	"net/url"
)

func newSQLServerConnectionURL(user, pass, host, dbName string, port int) *url.URL {
	query := url.Values{}
	query.Add("database", dbName)

	connURL := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(user, pass),
		Host:     fmt.Sprintf("%s:%d", host, port),
		RawQuery: query.Encode(),
	}

	return connURL
}
