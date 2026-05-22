package properties

import (
	"encoding/json"
	"strconv"
	"strings"

	"basesdk/errs"
)

const (
	DataTypeString = "string"
	DataTypeInt    = "int"
	DataTypeFloat  = "float"
	DataTypeBool   = "bool"
	DataTypeJSON   = "json"
)

func ValidatePropertyValue(dataType string, value string) error {
	switch strings.TrimSpace(dataType) {
	case DataTypeString:
		return nil
	case DataTypeInt:
		_, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return errs.BadRequestDirect("el valor debe ser un entero valido")
		}
	case DataTypeFloat:
		_, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return errs.BadRequestDirect("el valor debe ser un numero valido")
		}
	case DataTypeBool:
		_, err := parseBool(value)
		if err != nil {
			return errs.BadRequestDirect("el valor debe ser booleano")
		}
	case DataTypeJSON:
		var target any
		if err := json.Unmarshal([]byte(value), &target); err != nil {
			return errs.BadRequestDirect("el valor debe ser JSON valido")
		}
	default:
		return errs.BadRequestDirect("dataType invalido")
	}

	return nil
}
