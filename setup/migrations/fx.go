package migrations

import (
	"io/fs"

	"go.uber.org/fx"
)

type FileSystemSources []fs.FS

const GroupFSSources = `group:"migration-fs-sources"`

func ProvideFSSources(sources ...fs.FS) any {
	return fx.Annotate(
		func() FileSystemSources {
			return sources
		},
		fx.ResultTags(GroupFSSources),
	)
}
