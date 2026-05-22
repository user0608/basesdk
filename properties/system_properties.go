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

const notFoundFormat = "system property '%s' not found"

type Property struct {
	Key         string  `json:"key" gorm:"primaryKey"`
	Value       string  `json:"value"`
	DataType    string  `json:"dataType"`
	Description *string `json:"description"`
}

type SystemProperties struct {
	manager connection.StorageManager
}

func NewSystemProperties(manager connection.StorageManager) *SystemProperties {
	return &SystemProperties{
		manager: manager,
	}
}

func (sr *SystemProperties) GetAll(ctx context.Context) ([]Property, error) {
	tx := sr.manager.Conn(ctx)

	var properties []Property
	rs := tx.Table("system_properties").Order("key").Find(&properties)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return properties, nil
}

func (sr *SystemProperties) Get(ctx context.Context, key string) (Property, error) {
	tx := sr.manager.Conn(ctx)

	var property Property
	rs := tx.Table("system_properties").Where("key = ?", key).Limit(1).Find(&property)
	if rs.Error != nil {
		return Property{}, errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return Property{}, errs.NotFoundf(notFoundFormat, key)
	}
	return property, nil
}

func (sr *SystemProperties) Exists(ctx context.Context, key string) (bool, error) {
	_, notFound, err := sr.getOrDefault(ctx, key)
	if err != nil {
		return false, err
	}

	return !notFound, nil
}

func (sr *SystemProperties) Create(ctx context.Context, property *Property) error {
	tx := sr.manager.Conn(ctx)

	rs := tx.Table("system_properties").Create(property)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (sr *SystemProperties) Update(ctx context.Context, property *Property) error {
	tx := sr.manager.Conn(ctx)

	rs := tx.Table("system_properties").
		Where("key = ?", property.Key).
		Updates(property)

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	if rs.RowsAffected == 0 {
		return errs.NotFoundf(notFoundFormat, property.Key)
	}

	return nil
}

func (sr *SystemProperties) Delete(ctx context.Context, key string) error {
	tx := sr.manager.Conn(ctx)

	rs := tx.Table("system_properties").
		Where("key = ?", key).
		Delete(&Property{})

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}

func (sr *SystemProperties) getOrDefault(ctx context.Context, key string) (Property, bool, error) {
	property, err := sr.Get(ctx, key)
	if err != nil {
		if errs.ContainsMessage(err, "not found") {
			return Property{}, true, nil
		}
		return Property{}, false, err
	}

	return property, false, nil
}

func (sr *SystemProperties) GetString(ctx context.Context, key string, defaultValue string) (string, error) {
	property, notFound, err := sr.getOrDefault(ctx, key)
	if err != nil {
		return "", err
	}

	if notFound {
		return defaultValue, nil
	}

	return property.Value, nil
}

func (sr *SystemProperties) GetInt(ctx context.Context, key string, defaultValue int) (int, error) {
	property, notFound, err := sr.getOrDefault(ctx, key)
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

func (sr *SystemProperties) GetFloat(ctx context.Context, key string, defaultValue float64) (float64, error) {
	property, notFound, err := sr.getOrDefault(ctx, key)
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

func (sr *SystemProperties) GetBool(ctx context.Context, key string, defaultValue bool) (bool, error) {
	property, notFound, err := sr.getOrDefault(ctx, key)
	if err != nil {
		return false, err
	}

	if notFound {
		return defaultValue, nil
	}

	return parseBool(property.Value)
}

func (sr *SystemProperties) GetJSON(ctx context.Context, key string, dest any) error {
	property, err := sr.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(property.Value), dest)
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "t", "true", "y", "yes", "on", "enabled":
		return true, nil
	case "0", "f", "false", "n", "no", "off", "disabled":
		return false, nil
	default:
		return false, strconv.ErrSyntax
	}
}

func (sr *SystemProperties) GetTime(ctx context.Context, key string, defaultValue time.Time) (time.Time, error) {
	property, notFound, err := sr.getOrDefault(ctx, key)
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

func (sr *SystemProperties) GetDuration(ctx context.Context, key string, defaultValue time.Duration) (time.Duration, error) {
	property, notFound, err := sr.getOrDefault(ctx, key)
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
