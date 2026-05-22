package handlers

import (
	"basesdk/binds"
	"basesdk/httpapi"
	"basesdk/properties"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/user0608/goones/answer"
	"github.com/user0608/goones/kcheck"
)

type propertyInput struct {
	Key         string  `json:"key" chk:"nonil"`
	Value       string  `json:"value"`
	DataType    string  `json:"dataType" chk:"nonil"`
	Description *string `json:"description"`
}

func toProperty(input propertyInput) *properties.Property {
	return &properties.Property{
		Key:         strings.TrimSpace(input.Key),
		Value:       input.Value,
		DataType:    strings.TrimSpace(input.DataType),
		Description: input.Description,
	}
}

func SystemPropertiesListHandler(props *properties.SystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodGet, Path: "/api/v1/system/properties", Handler: func(c echo.Context) error {
		response, err := props.GetAll(c.Request().Context())
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Ok(c, response)
	}}
}

func SystemPropertyCreateHandler(props *properties.SystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodPost, Path: "/api/v1/system/properties", Handler: func(c echo.Context) error {
		var payload propertyInput
		if err := binds.JSON(c, &payload); err != nil {
			return answer.Err(c, err)
		}
		if err := kcheck.Valid(payload); err != nil {
			return answer.Err(c, err)
		}
		property := toProperty(payload)
		if err := properties.ValidatePropertyValue(property.DataType, property.Value); err != nil {
			return answer.Err(c, err)
		}
		if err := props.Create(c.Request().Context(), property); err != nil {
			return answer.Err(c, err)
		}
		return answer.Created(c)
	}}
}

func SystemPropertyFindHandler(props *properties.SystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodGet, Path: "/api/v1/system/properties/:key", Handler: func(c echo.Context) error {
		response, err := props.Get(c.Request().Context(), c.Param("key"))
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Ok(c, response)
	}}
}

func SystemPropertyUpdateHandler(props *properties.SystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodPut, Path: "/api/v1/system/properties/:key", Handler: func(c echo.Context) error {
		var payload propertyInput
		if err := binds.JSON(c, &payload); err != nil {
			return answer.Err(c, err)
		}
		if err := kcheck.Valid(payload); err != nil {
			return answer.Err(c, err)
		}
		property := toProperty(payload)
		property.Key = strings.TrimSpace(c.Param("key"))
		if err := properties.ValidatePropertyValue(property.DataType, property.Value); err != nil {
			return answer.Err(c, err)
		}
		if err := props.Update(c.Request().Context(), property); err != nil {
			return answer.Err(c, err)
		}
		return answer.Success(c)
	}}
}

func SystemPropertyDeleteHandler(props *properties.SystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodDelete, Path: "/api/v1/system/properties/:key", Handler: func(c echo.Context) error {
		if err := props.Delete(c.Request().Context(), c.Param("key")); err != nil {
			return answer.Err(c, err)
		}
		return answer.Success(c)
	}}
}
