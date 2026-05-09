package permissions

import "io/fs"

type FileSystemSources []fs.FS

type PermissionDefinition struct {
	Code        string
	Description string
}
