package mssqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeSQLArg(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "NothingToSanitize",
			input:  "Catalog.Table-1",
			output: "Catalog.Table-1",
		},
		{
			name:   "SanitizeForbiddenCharacters",
			input:  "Table ,;!?%(·",
			output: "Table",
		},
		{
			name:   "EmptyInput",
			input:  "",
			output: "",
		},
		{
			name:   "OnlyForbiddenCharacters",
			input:  " ,;!?%(·",
			output: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.output, SanitizeSQLArg(test.input))
		})
	}
}
