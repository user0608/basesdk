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

type tenantPropertyInput struct {
	Key         string  `json:"key" chk:"nonil"`
	Value       string  `json:"value"`
	DataType    string  `json:"dataType" chk:"nonil"`
	Description *string `json:"description"`
}

func toTenantProperty(tenantCodigo string, input tenantPropertyInput) *properties.TenantSystemProperty {
	return &properties.TenantSystemProperty{
		TenantCodigo: strings.TrimSpace(tenantCodigo),
		Key:          strings.TrimSpace(input.Key),
		Value:        input.Value,
		DataType:     strings.TrimSpace(input.DataType),
		Description:  input.Description,
	}
}

func SystemTenantPropertiesListHandler(props *properties.TenantSystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodGet, Path: "/api/v1/system/tenants/:tenantCodigo/properties", Handler: func(c echo.Context) error {
		response, err := props.GetAll(c.Request().Context(), c.Param("tenantCodigo"))
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Ok(c, response)
	}}
}

func SystemTenantPropertyCreateHandler(props *properties.TenantSystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodPost, Path: "/api/v1/system/tenants/:tenantCodigo/properties", Handler: func(c echo.Context) error {
		var payload tenantPropertyInput
		if err := binds.JSON(c, &payload); err != nil {
			return answer.Err(c, err)
		}
		if err := kcheck.Valid(payload); err != nil {
			return answer.Err(c, err)
		}
		property := toTenantProperty(c.Param("tenantCodigo"), payload)
		if err := properties.ValidatePropertyValue(property.DataType, property.Value); err != nil {
			return answer.Err(c, err)
		}
		if err := props.Create(c.Request().Context(), property); err != nil {
			return answer.Err(c, err)
		}
		return answer.Created(c)
	}}
}

func SystemTenantPropertyFindHandler(props *properties.TenantSystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodGet, Path: "/api/v1/system/tenants/:tenantCodigo/properties/:key", Handler: func(c echo.Context) error {
		response, err := props.Get(c.Request().Context(), c.Param("tenantCodigo"), c.Param("key"))
		if err != nil {
			return answer.Err(c, err)
		}
		return answer.Ok(c, response)
	}}
}

func SystemTenantPropertyUpdateHandler(props *properties.TenantSystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodPut, Path: "/api/v1/system/tenants/:tenantCodigo/properties/:key", Handler: func(c echo.Context) error {
		var payload tenantPropertyInput
		if err := binds.JSON(c, &payload); err != nil {
			return answer.Err(c, err)
		}
		if err := kcheck.Valid(payload); err != nil {
			return answer.Err(c, err)
		}
		property := toTenantProperty(c.Param("tenantCodigo"), payload)
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

func SystemTenantPropertyDeleteHandler(props *properties.TenantSystemProperties) httpapi.Route {
	return &httpapi.SystemHandler{Method: http.MethodDelete, Path: "/api/v1/system/tenants/:tenantCodigo/properties/:key", Handler: func(c echo.Context) error {
		if err := props.Delete(c.Request().Context(), c.Param("tenantCodigo"), c.Param("key")); err != nil {
			return answer.Err(c, err)
		}
		return answer.Success(c)
	}}
}
