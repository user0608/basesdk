package permissions

import (
	"io/fs"

	"go.uber.org/fx"
)

const GroupFSSources = `group:"permission-fs-sources"`

func ProvideFSSources(sources ...fs.FS) any {
	return fx.Annotate(
		func() FileSystemSources {
			return sources
		},
		fx.ResultTags(GroupFSSources),
	)
}
