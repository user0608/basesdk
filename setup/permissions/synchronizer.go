package permissions

import (
	"basesdk/connection"
	"context"
	"fmt"
	"io/fs"
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

	tx := ps.storageManager.Conn(ctx)

	const query = `
		insert into permission (
			code,
			description
		)
		values (
			?,
			?
		)
		on conflict (code) do update
		set
			description = excluded.description
	`

	for _, definition := range definitions {
		rs := tx.Exec(query, definition.Code, definition.Description)
		if rs.Error != nil {
			return fmt.Errorf("sync permission %q: %w", definition.Code, rs.Error)
		}
	}

	return nil
}
