package graphql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFieldName(t *testing.T) {
	tests := []struct {
		fieldName      string
		expectedTokens []string
	}{
		{"c", []string{"c"}},
		{"C", []string{"C"}},
		{"camel", []string{"camel"}},
		{"Camel", []string{"Camel"}},
		{"camelCase", []string{"camel", "Case"}},
		{"CamelCase", []string{"Camel", "Case"}},
		{"OneTwoThree", []string{"One", "Two", "Three"}},
		{"DatabaseURL", []string{"Database", "URL"}},
		{"DBEndpoints", []string{"DB", "Endpoints"}},
	}

	for _, tc := range tests {
		tokens := parseFieldName(tc.fieldName)
		fmt.Println(tokens)
		assert.Equal(t, tc.expectedTokens, tokens)
	}
}
