package properties

import (
	"basesdk/connection"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"basesdk/errs"
)

const tenantNotFoundFormat = "tenant system property '%s' not found for tenant '%s'"

type TenantSystemProperty struct {
	Key          string  `json:"key" gorm:"primaryKey"`
	Value        string  `json:"value"`
	TenantCodigo string  `json:"tenant_codigo" gorm:"primaryKey"`
	DataType     string  `json:"data_type"`
	Description  *string `json:"description"`
}

type TenantSystemProperties struct {
	manager connection.StorageManager
}

func NewTenantSystemProperties(manager connection.StorageManager) *TenantSystemProperties {
	return &TenantSystemProperties{
		manager: manager,
	}
}

func (sr *TenantSystemProperties) GetAll(ctx context.Context, tenantCodigo string) ([]TenantSystemProperty, error) {
	tx := sr.manager.Conn(ctx)

	var properties []TenantSystemProperty
	rs := tx.Table("tenant_system_properties").
		Where("tenant_codigo = ?", tenantCodigo).
		Order("key").
		Find(&properties)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return properties, nil
}

func (sr *TenantSystemProperties) Get(ctx context.Context, tenantCodigo string, key string) (TenantSystemProperty, error) {
	tx := sr.manager.Conn(ctx)

	var property TenantSystemProperty
	rs := tx.Table("tenant_system_properties").
		Where("tenant_codigo = ? and key = ?", tenantCodigo, key).
		Limit(1).
		Find(&property)
	if rs.Error != nil {
		return TenantSystemProperty{}, errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return TenantSystemProperty{}, errs.NotFoundf(tenantNotFoundFormat, key, tenantCodigo)
	}
	return property, nil
}

func (sr *TenantSystemProperties) Exists(ctx context.Context, tenantCodigo string, key string) (bool, error) {
	_, notFound, err := sr.getOrDefault(ctx, tenantCodigo, key)
	if err != nil {
		return false, err
	}

	return !notFound, nil
}

func (sr *TenantSystemProperties) Create(ctx context.Context, property *TenantSystemProperty) error {
	tx := sr.manager.Conn(ctx)

	rs := tx.Table("tenant_system_properties").Create(property)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (sr *TenantSystemProperties) Update(ctx context.Context, property *TenantSystemProperty) error {
	tx := sr.manager.Conn(ctx)

	rs := tx.Table("tenant_system_properties").
		Where("tenant_codigo = ? and key = ?", property.TenantCodigo, property.Key).
		Updates(property)

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	if rs.RowsAffected == 0 {
		return errs.NotFoundf(tenantNotFoundFormat, property.Key, property.TenantCodigo)
	}

	return nil
}

func (sr *TenantSystemProperties) Delete(ctx context.Context, tenantCodigo string, key string) error {
	tx := sr.manager.Conn(ctx)

	rs := tx.Table("tenant_system_properties").
		Where("tenant_codigo = ? and key = ?", tenantCodigo, key).
		Delete(&TenantSystemProperty{})

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (sr *TenantSystemProperties) getOrDefault(ctx context.Context, tenantCodigo string, key string) (TenantSystemProperty, bool, error) {
	property, err := sr.Get(ctx, tenantCodigo, key)
	if err != nil {
		if errs.ContainsMessage(err, "not found") {
			return TenantSystemProperty{}, true, nil
		}
		return TenantSystemProperty{}, false, err
	}

	return property, false, nil
}

func (sr *TenantSystemProperties) GetString(ctx context.Context, tenantCodigo string, key string, defaultValue string) (string, error) {
	property, notFound, err := sr.getOrDefault(ctx, tenantCodigo, key)
	if err != nil {
		return "", err
	}

	if notFound {
		return defaultValue, nil
	}

	return property.Value, nil
}

func (sr *TenantSystemProperties) GetInt(ctx context.Context, tenantCodigo string, key string, defaultValue int) (int, error) {
	property, notFound, err := sr.getOrDefault(ctx, tenantCodigo, key)
	if err != nil {
		return 0, err
	}

	if notFound {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(strings.TrimSpace(property.Value))
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (sr *TenantSystemProperties) GetFloat(ctx context.Context, tenantCodigo string, key string, defaultValue float64) (float64, error) {
	property, notFound, err := sr.getOrDefault(ctx, tenantCodigo, key)
	if err != nil {
		return 0, err
	}

	if notFound {
		return defaultValue, nil
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(property.Value), 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (sr *TenantSystemProperties) GetBool(ctx context.Context, tenantCodigo string, key string, defaultValue bool) (bool, error) {
	property, notFound, err := sr.getOrDefault(ctx, tenantCodigo, key)
	if err != nil {
		return false, err
	}

	if notFound {
		return defaultValue, nil
	}

	return parseBool(property.Value)
}

func (sr *TenantSystemProperties) GetJSON(ctx context.Context, tenantCodigo string, key string, dest any) error {
	property, err := sr.Get(ctx, tenantCodigo, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(property.Value), dest)
}

func (sr *TenantSystemProperties) GetTime(ctx context.Context, tenantCodigo string, key string, defaultValue time.Time) (time.Time, error) {
	property, notFound, err := sr.getOrDefault(ctx, tenantCodigo, key)
	if err != nil {
		return time.Time{}, err
	}

	if notFound {
		return defaultValue, nil
	}

	value, err := time.Parse(time.RFC3339, strings.TrimSpace(property.Value))
	if err != nil {
		return time.Time{}, err
	}

	return value, nil
}

func (sr *TenantSystemProperties) GetDuration(ctx context.Context, tenantCodigo string, key string, defaultValue time.Duration) (time.Duration, error) {
	property, notFound, err := sr.getOrDefault(ctx, tenantCodigo, key)
	if err != nil {
		return 0, err
	}

	if notFound {
		return defaultValue, nil
	}

	value, err := time.ParseDuration(strings.TrimSpace(property.Value))
	if err != nil {
		return 0, err
	}

	return value, nil
}
