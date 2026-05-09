package permissions

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func readCSV(file io.Reader) ([]PermissionDefinition, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	out := make([]PermissionDefinition, 0, len(rows))

	for i, row := range rows {
		rowNumber := i + 1
		if len(row) != 2 {
			return nil, fmt.Errorf("invalid row %d", rowNumber)
		}

		code := strings.TrimSpace(row[0])
		description := strings.TrimSpace(row[1])

		if code == "" {
			return nil, fmt.Errorf("empty permission code at row %d", rowNumber)
		}

		out = append(out, PermissionDefinition{
			Code:        code,
			Description: description,
		})
	}

	return out, nil
}
