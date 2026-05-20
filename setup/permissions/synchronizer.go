package permissions

import (
	"basesdk/connection"
	"context"
	"fmt"
	"io/fs"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PermissionSynchronizer struct {
	storageManager connection.StorageManager
	sources        []fs.FS
}

func NewPermissionSynchronizer(
	storageManager connection.StorageManager,
	sourceGroups []FileSystemSources,
) *PermissionSynchronizer {
	return &PermissionSynchronizer{
		storageManager: storageManager,
		sources:        flattenSources(sourceGroups),
	}
}

func flattenSources(sourceGroups []FileSystemSources) []fs.FS {
	var sources []fs.FS

	for _, group := range sourceGroups {
		for _, src := range group {
			if src == nil {
				continue
			}

			sources = append(sources, src)
		}
	}

	return sources
}

func (ps *PermissionSynchronizer) Run(ctx context.Context) error {
	definitions, err := ps.Read()
	if err != nil {
		return err
	}

	tx := ps.storageManager.Conn(ctx).Session(
		&gorm.Session{
			Logger: logger.Default.LogMode(logger.Error),
		},
	)

	if len(definitions) == 0 {
		return nil
	}

	var values strings.Builder
	args := make([]any, 0, len(definitions)*2)
	for i, definition := range definitions {
		if i > 0 {
			values.WriteString(",")
		}

		values.WriteString("(?, ?)")
		args = append(args, definition.Code, definition.Description)
	}

	query := fmt.Sprintf(`
		insert into permission (
			code,
			description
		)
		values %s
		on conflict (code) do update
		set
			description = excluded.description
	`, values.String())

	rs := tx.Exec(query, args...)
	if rs.Error != nil {
		return fmt.Errorf("sync permissions: %w", rs.Error)
	}

	return nil
}
