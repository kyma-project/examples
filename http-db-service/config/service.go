package config

import (
	"encoding/json"
	"fmt"
)

// Service struct is used for configuring how the service will run
// by reading the values from the environment or using the default values.
type Service struct {
	Port   string `envconfig:"serviceport,default=8017" json:"Port"`
	// DbType set to 'mssql' will start the service using an MsSql datbase
	// and it will require extra configuration. See https://github.com/kyma-project/examples/blob/master/http-db-service/internal/mssqldb/config.go
	DbType string `envconfig:"dbtype,default=memory" json:"DBType"` // [memory | mssql]
}

// String returns a printable representation of the config as JSON.
// Use the struct field tag `json:"-"` to hide fields that should not be revealed such as credentials and secrets.
func (s Service) String() string {
	json, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("Error marshalling service configuration JSON: %v", err)
	}
	return fmt.Sprintf("Service Configuration: %s", json)
}
