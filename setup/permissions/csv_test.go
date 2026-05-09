package permissions

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadCSV(t *testing.T) {
	definitions, err := readCSV(strings.NewReader(`users.read, Read users
roles.update,Update roles
`))

	require.NoError(t, err)
	require.Equal(t, []PermissionDefinition{
		{Code: "users.read", Description: "Read users"},
		{Code: "roles.update", Description: "Update roles"},
	}, definitions)
}

func TestReadCSVAllowsEmptyDescription(t *testing.T) {
	definitions, err := readCSV(strings.NewReader(`users.read,
`))

	require.NoError(t, err)
	require.Equal(t, []PermissionDefinition{
		{Code: "users.read", Description: ""},
	}, definitions)
}

func TestReadCSVRejectsEmptyCode(t *testing.T) {
	_, err := readCSV(strings.NewReader(` ,Read users
`))

	require.ErrorContains(t, err, "empty permission code")
}
