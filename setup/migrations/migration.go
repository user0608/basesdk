package migrations

import (
	"basesdk/connection"
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"path"
	"regexp"
	"sort"
	"strings"
	"testing/fstest"

	"github.com/pressly/goose/v3"
)

const baseDir = "migrations"

type MigrationRunner struct {
	manager connection.StorageManager
	sources []fs.FS
}

func NewMigrationRunner(manager connection.StorageManager, sources ...fs.FS) *MigrationRunner {
	return &MigrationRunner{
		manager: manager,
		sources: sources,
	}
}

func (mr *MigrationRunner) MigrationsFS() (fs.FS, error) {
	out := fstest.MapFS{}

	for i, src := range mr.sources {
		rootDirs, err := fs.ReadDir(src, ".")
		if err != nil {
			return nil, fmt.Errorf("read source %d root: %w", i, err)
		}

		for _, rootDir := range rootDirs {
			if !rootDir.IsDir() {
				continue
			}

			entries, err := fs.ReadDir(src, rootDir.Name())
			if err != nil {
				return nil, fmt.Errorf("read source %d dir %q: %w", i, rootDir.Name(), err)
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				name := entry.Name()
				if path.Ext(name) != ".sql" {
					continue
				}

				sourcePath := path.Join(rootDir.Name(), name)
				targetPath := path.Join(baseDir, name)

				if _, exists := out[targetPath]; exists {
					return nil, fmt.Errorf("duplicate migration file %q", targetPath)
				}

				data, err := fs.ReadFile(src, sourcePath)
				if err != nil {
					return nil, fmt.Errorf("read source %d file %q: %w", i, sourcePath, err)
				}

				out[targetPath] = &fstest.MapFile{
					Data: data,
					Mode: 0644,
				}
			}
		}
	}

	return out, nil
}

func (mr *MigrationRunner) SQLScript() (string, error) {
	migrationsFS, err := mr.MigrationsFS()
	if err != nil {
		return "", err
	}

	entries, err := fs.ReadDir(migrationsFS, baseDir)
	if err != nil {
		return "", fmt.Errorf("read migrations dir: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	rgx := regexp.MustCompile(`(?s)-- \+goose Up(.*?)-- \+goose Down`)

	var output strings.Builder

	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != ".sql" {
			continue
		}

		filePath := path.Join(baseDir, entry.Name())

		content, err := fs.ReadFile(migrationsFS, filePath)
		if err != nil {
			return "", fmt.Errorf("read migration file %q: %w", filePath, err)
		}

		match := rgx.FindSubmatch(content)
		if match == nil {
			continue
		}

		fmt.Fprintf(&output, "-- %s\n", entry.Name())

		scanner := bufio.NewScanner(bytes.NewReader(match[1]))
		for scanner.Scan() {
			line := bytes.TrimSpace(scanner.Bytes())

			if len(line) == 0 || bytes.HasPrefix(line, []byte("--")) {
				continue
			}

			output.Write(line)
			output.WriteByte('\n')
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("scan migration file %q: %w", filePath, err)
		}
	}

	return output.String(), nil
}

func (mr *MigrationRunner) setupGoose(ctx context.Context) (*sql.DB, error) {
	migrationsFS, err := mr.MigrationsFS()
	if err != nil {
		return nil, err
	}
	tx := mr.manager.Conn(ctx)
	db, err := tx.DB()
	if err != nil {
		return nil, err
	}
	goose.SetBaseFS(migrationsFS)
	return db, nil
}

func (mr *MigrationRunner) Up(ctx context.Context) error {
	db, err := mr.setupGoose(ctx)
	if err != nil {
		return err
	}
	if err := goose.UpContext(ctx, db, baseDir); err != nil {
		slog.Error(fmt.Sprintf("Database %s failed", "up"), "error", err)
		return err
	}
	return nil
}

func (mr *MigrationRunner) Down(ctx context.Context) error {
	db, err := mr.setupGoose(ctx)
	if err != nil {
		return err
	}
	if err := goose.DownContext(ctx, db, baseDir); err != nil {
		slog.Error(fmt.Sprintf("Database %s failed", "up"), "error", err)
		return err
	}
	return nil
}

func (mr *MigrationRunner) Status(ctx context.Context) error {
	db, err := mr.setupGoose(ctx)
	if err != nil {
		return err
	}
	if err := goose.StatusContext(ctx, db, baseDir); err != nil {
		slog.Error(fmt.Sprintf("Database %s failed", "up"), "error", err)
		return err
	}
	return nil
}
