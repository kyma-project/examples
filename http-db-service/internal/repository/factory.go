package repository

import (
	"github.com/pkg/errors"
)

const (
	// MemoryDatabase value can be used to start the service using an in-memory DB. See Service/DbType.
	MemoryDatabase      = "memory"
	// SQLServerDriverName value can be used to start the service using an external MsSql DB. See Service/DbType.
	SQLServerDriverName = "mssql"
)

// Create is used to create an OrderRepository based on the given dbtype.
// Currently the `MemoryDatabase` and `SQLServerDriverName` are supported.
func Create(dbtype string) (OrderRepository, error) {
	switch dbtype {
		case MemoryDatabase:
			return NewOrderRepositoryMemory(), nil
		case SQLServerDriverName:
			return NewOrderRepositoryDb()
		default:
			return nil, errors.Errorf("Unsupported database type %s", dbtype)
	}
}