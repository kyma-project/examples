package mssqldb

import "regexp"

const (

  sqlServerTableCreationQuery = `IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = '{name}')
BEGIN
    CREATE TABLE {name} (
      order_id VARCHAR(64),
      namespace VARCHAR(64),
      total DECIMAL(8,2),
      PRIMARY KEY (order_id, namespace)
    )
END`

  // PrimaryKeyViolation is the SQL code used by MsSql to indicate an attempt
  // to insert an entry which violates the primary key constraint.
  PrimaryKeyViolation = 2627
)

var safeSQLRegex = regexp.MustCompile(`[^a-zA-Z0-9\.\-_]`)

// SanitizeSQLArg returns the input string sanitized for safe use in an SQL query as argument.
func SanitizeSQLArg(s string) string {
	return safeSQLRegex.ReplaceAllString(s, "")
}
