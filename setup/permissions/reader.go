package permissions

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
)

func (ps *PermissionSynchronizer) Read() ([]PermissionDefinition, error) {
	byCode := map[string]PermissionDefinition{}

	for i, src := range ps.sources {
		err := fs.WalkDir(src, ".", func(filePath string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if entry.IsDir() || path.Ext(filePath) != ".csv" {
				return nil
			}

			file, err := src.Open(filePath)
			if err != nil {
				return err
			}

			definitions, readErr := readCSV(file)
			closeErr := file.Close()

			if readErr != nil {
				return fmt.Errorf("read %q from source %d: %w", filePath, i, readErr)
			}

			if closeErr != nil {
				return fmt.Errorf("close %q from source %d: %w", filePath, i, closeErr)
			}

			for _, definition := range definitions {
				if existing, exists := byCode[definition.Code]; exists {
					return fmt.Errorf(
						"duplicate permission %q: %q and %q",
						definition.Code,
						existing.Description,
						definition.Description,
					)
				}

				byCode[definition.Code] = definition
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk permission source %d: %w", i, err)
		}
	}

	out := make([]PermissionDefinition, 0, len(byCode))
	for _, definition := range byCode {
		out = append(out, definition)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Code < out[j].Code
	})

	return out, nil
}
