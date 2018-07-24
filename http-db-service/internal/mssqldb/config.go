package mssqldb

import (
	"encoding/json"
	"fmt"
)

// Config is a struct used for configuring the connection and the usage of the database.
// It contains the information needed to form a connection string and the name of the table used for storing the data.
type Config struct {
	Name              string `envconfig:"database,default=orderservice" json:"Name"`
	Host              string `envconfig:"host,default=127.0.0.1" json:"Host"`
	Port              int    `envconfig:"port,default=1433" json:"Port"`
	User              string `envconfig:"username,default=test" json:"User"`
	Pass              string `envconfig:"password,default=test" json:"-"` // hidden from logging
	DbOrdersTableName string `envconfig:"tablename,default=orders" json:"OrdersTable"`
}

// String returns a printable representation of the config as JSON.
// Use the struct field tag `json:"-"` to hide fields that should not be revealed such as credentials and secrets.
func (config Config) String() string {
	json, err := json.Marshal(config)
	if err != nil {
		return fmt.Sprintf("Error marshalling DB configuration JSON: %v", err)
	}
	return fmt.Sprintf("DB Configuration: %s", json)
}
