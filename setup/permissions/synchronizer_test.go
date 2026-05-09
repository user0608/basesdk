package permissions

import (
	"basesdk/connection"
	"basesdk/testdb"
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestPermissionSynchronizerReadSortsDefinitions(t *testing.T) {
	syncer := NewPermissionSynchronizer(nil, []FileSystemSources{{fstest.MapFS{
		"permissions/users.csv": &fstest.MapFile{Data: []byte(`users.update,Update users
users.read,Read users
`)},
		"permissions/roles.csv": &fstest.MapFile{Data: []byte(`roles.read,Read roles
`)},
	}}})

	definitions, err := syncer.Read()

	require.NoError(t, err)
	require.Equal(t, []PermissionDefinition{
		{Code: "roles.read", Description: "Read roles"},
		{Code: "users.read", Description: "Read users"},
		{Code: "users.update", Description: "Update users"},
	}, definitions)
}

func TestPermissionSynchronizerReadRejectsDuplicateCodes(t *testing.T) {
	syncer := NewPermissionSynchronizer(nil, []FileSystemSources{{fstest.MapFS{
		"one.csv": &fstest.MapFile{Data: []byte(`users.read,Read users
`)},
		"two.csv": &fstest.MapFile{Data: []byte(`users.read,Read users again
`)},
	}}})

	_, err := syncer.Read()

	require.ErrorContains(t, err, `duplicate permission "users.read"`)
}

func TestPermissionSynchronizerRunUpsertsPermissions(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	ctx := context.Background()

	syncer := NewPermissionSynchronizer(storage, []FileSystemSources{{fstest.MapFS{
		"permissions.csv": &fstest.MapFile{Data: []byte(`users.read,Read users
`)},
	}}})

	require.NoError(t, syncer.Run(ctx))
	requirePermissionDescription(t, storage, "users.read", "Read users")

	syncer = NewPermissionSynchronizer(storage, []FileSystemSources{{fstest.MapFS{
		"permissions.csv": &fstest.MapFile{Data: []byte(`users.read,Read users updated
`)},
	}}})

	require.NoError(t, syncer.Run(ctx))
	requirePermissionDescription(t, storage, "users.read", "Read users updated")
}

func requirePermissionDescription(t *testing.T, storage connection.StorageManager, code string, expected string) {
	t.Helper()

	var got string
	err := storage.Conn(context.Background()).Raw(`
		select description
		from permission
		where code = ?
	`, code).Scan(&got).Error
	require.NoError(t, err)
	require.Equal(t, expected, got)
}
