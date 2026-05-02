package migrations

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrationsFSCombinesFirstLevelSQLFromUnknownRootDirs(t *testing.T) {
	src1 := fstest.MapFS{
		"seguridad/001_init.sql":          {Data: []byte("create table init;")},
		"seguridad/002_users.sql":         {Data: []byte("create table users;")},
		"seguridad/subdir/003_nested.sql": {Data: []byte("nested")},
		"seguridad/readme.md":             {Data: []byte("ignored")},
	}

	src2 := fstest.MapFS{
		"ventas/001_clientes.sql":      {Data: []byte("create table clientes;")},
		"ventas/002_proveedores.sql":   {Data: []byte("create table proveedores;")},
		"ventas/subdir/003_nested.sql": {Data: []byte("nested")},
		"ventas/notes.txt":             {Data: []byte("ignored")},
	}

	runner := NewMigrationRunner(nil, src1, src2)

	combined, err := runner.MigrationsFS()
	if err != nil {
		t.Fatalf("MigrationsFS() error = %v", err)
	}

	entries, err := fs.ReadDir(combined, "migrations")
	if err != nil {
		t.Fatalf("ReadDir(migrations) error = %v", err)
	}

	got := make([]string, 0, len(entries))
	for _, entry := range entries {
		got = append(got, entry.Name())
	}

	want := []string{
		"001_clientes.sql",
		"001_init.sql",
		"002_proveedores.sql",
		"002_users.sql",
	}

	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestMigrationsFSPreservesFileContent(t *testing.T) {
	src := fstest.MapFS{
		"seguridad/001_init.sql": {Data: []byte("create table users;")},
	}

	runner := NewMigrationRunner(nil, src)

	combined, err := runner.MigrationsFS()
	if err != nil {
		t.Fatalf("MigrationsFS() error = %v", err)
	}

	data, err := fs.ReadFile(combined, "migrations/001_init.sql")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(data) != "create table users;" {
		t.Fatalf("got %q", string(data))
	}
}

func TestMigrationsFSIgnoresNestedSQL(t *testing.T) {
	src := fstest.MapFS{
		"seguridad/001_init.sql":          {Data: []byte("root")},
		"seguridad/subdir/002_nested.sql": {Data: []byte("nested")},
	}

	runner := NewMigrationRunner(nil, src)

	combined, err := runner.MigrationsFS()
	if err != nil {
		t.Fatalf("MigrationsFS() error = %v", err)
	}

	_, err = fs.ReadFile(combined, "migrations/002_nested.sql")
	if err == nil {
		t.Fatal("expected nested SQL file to be ignored")
	}
}

func TestMigrationsFSIgnoresNonSQLFiles(t *testing.T) {
	src := fstest.MapFS{
		"seguridad/001_init.sql": {Data: []byte("sql")},
		"seguridad/readme.md":    {Data: []byte("markdown")},
		"seguridad/config.json":  {Data: []byte("{}")},
	}

	runner := NewMigrationRunner(nil, src)

	combined, err := runner.MigrationsFS()
	if err != nil {
		t.Fatalf("MigrationsFS() error = %v", err)
	}

	entries, err := fs.ReadDir(combined, "migrations")
	if err != nil {
		t.Fatalf("ReadDir(migrations) error = %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("got %d files, want 1", len(entries))
	}

	if entries[0].Name() != "001_init.sql" {
		t.Fatalf("got %q, want %q", entries[0].Name(), "001_init.sql")
	}
}

func TestMigrationsFSReturnsErrorOnDuplicateFile(t *testing.T) {
	src1 := fstest.MapFS{
		"seguridad/001_init.sql": {Data: []byte("one")},
	}

	src2 := fstest.MapFS{
		"ventas/001_init.sql": {Data: []byte("two")},
	}

	runner := NewMigrationRunner(nil, src1, src2)

	_, err := runner.MigrationsFS()
	if err == nil {
		t.Fatal("expected duplicate file error, got nil")
	}

	if !strings.Contains(err.Error(), "duplicate migration file") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// region Sql Script

func TestMigrationRunnerSQLScript(t *testing.T) {
	src := fstest.MapFS{
		"seguridad/001_init.sql": {
			Data: []byte(`
-- +goose Up
CREATE TABLE users (
	id TEXT PRIMARY KEY
);
-- comment ignored
CREATE INDEX idx_users_id ON users(id);
-- +goose Down
DROP TABLE users;
`),
		},
		"ventas/002_orders.sql": {
			Data: []byte(`
-- +goose Up
CREATE TABLE orders (
	id TEXT PRIMARY KEY
);
-- +goose Down
DROP TABLE orders;
`),
		},
	}

	runner := NewMigrationRunner(nil, src)

	got, err := runner.SQLScript()
	if err != nil {
		t.Fatalf("SQLScript() error = %v", err)
	}

	want := `-- 001_init.sql
CREATE TABLE users (
id TEXT PRIMARY KEY
);
CREATE INDEX idx_users_id ON users(id);
-- 002_orders.sql
CREATE TABLE orders (
id TEXT PRIMARY KEY
);
`

	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestMigrationRunnerSQLScriptIgnoresFilesWithoutGooseBlock(t *testing.T) {
	src := fstest.MapFS{
		"seguridad/001_init.sql": {
			Data: []byte(`CREATE TABLE ignored (id TEXT);`),
		},
		"seguridad/002_users.sql": {
			Data: []byte(`
-- +goose Up
CREATE TABLE users (id TEXT);
-- +goose Down
DROP TABLE users;
`),
		},
	}

	runner := NewMigrationRunner(nil, src)

	got, err := runner.SQLScript()
	if err != nil {
		t.Fatalf("SQLScript() error = %v", err)
	}

	if strings.Contains(got, "ignored") {
		t.Fatalf("expected file without goose block to be ignored, got:\n%s", got)
	}

	if !strings.Contains(got, "CREATE TABLE users") {
		t.Fatalf("expected users migration, got:\n%s", got)
	}
}

func TestMigrationRunnerSQLScriptReturnsDuplicateError(t *testing.T) {
	src1 := fstest.MapFS{
		"seguridad/001_init.sql": {
			Data: []byte(`
-- +goose Up
SELECT 1;
-- +goose Down
SELECT 0;
`),
		},
	}

	src2 := fstest.MapFS{
		"ventas/001_init.sql": {
			Data: []byte(`
-- +goose Up
SELECT 2;
-- +goose Down
SELECT 0;
`),
		},
	}

	runner := NewMigrationRunner(nil, src1, src2)

	_, err := runner.SQLScript()
	if err == nil {
		t.Fatal("expected duplicate error, got nil")
	}

	if !strings.Contains(err.Error(), "duplicate migration file") {
		t.Fatalf("unexpected error: %v", err)
	}
}
