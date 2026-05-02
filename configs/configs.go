package configs

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type ConfigPath string

type ApplicationConfigs interface {
	ListenAddress() string
	Database() DatabaseConfig
}

type DatabaseConfig interface {
	GetHost() string
	GetPort() int
	GetUsername() string
	GetPassword() string
	GetDBName() string
	GetLogLevel() string
}

type ConfigPathProvider func() (ConfigPath, error)
type ApplicationConfigsProvider func(configPath ConfigPath) (ApplicationConfigs, DatabaseConfig, error)

type rawConfig struct {
	Address  string            `mapstructure:"address"`
	Database rawDatabaseConfig `mapstructure:"database"`
}

type rawDatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"db_name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	LogLevel string `mapstructure:"log_level"`
}

type config struct {
	address  string
	database databaseConfig
}

type databaseConfig struct {
	host     string
	port     int
	dbName   string
	username string
	password string
	logLevel string
}

var _ ApplicationConfigs = (*config)(nil)
var _ DatabaseConfig = (*databaseConfig)(nil)

func (c *config) ListenAddress() string {
	return c.address
}

func (c *config) Database() DatabaseConfig {
	return &c.database
}

func (d *databaseConfig) GetHost() string {
	return d.host
}

func (d *databaseConfig) GetPort() int {
	return d.port
}

func (d *databaseConfig) GetUsername() string {
	return d.username
}

func (d *databaseConfig) GetPassword() string {
	return d.password
}

func (d *databaseConfig) GetDBName() string {
	return d.dbName
}

func (d *databaseConfig) GetLogLevel() string {
	return d.logLevel
}

func (c *config) validate() error {
	if c.address == "" {
		return errors.New("address is required")
	}

	if c.database.host == "" {
		return errors.New("database host is required")
	}

	if c.database.port <= 0 {
		return fmt.Errorf("database port must be positive: %d", c.database.port)
	}

	if c.database.dbName == "" {
		return errors.New("database db_name is required")
	}

	if c.database.username == "" {
		return errors.New("database username is required")
	}

	if c.database.logLevel == "" {
		return errors.New("database log_level is required")
	}

	return nil
}

var _ ApplicationConfigsProvider = DefaultConfigsProvider

func DefaultConfigsProvider(configPath ConfigPath) (ApplicationConfigs, DatabaseConfig, error) {
	filePath := strings.TrimSpace(string(configPath))
	if filePath == "" {
		slog.Error("config path is empty")
		return nil, nil, errors.New("config path is empty")
	}

	var raw rawConfig

	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		slog.Error("read config", "file", filepath.Base(filePath), "error", err)
		return nil, nil, err
	}

	if err := v.Unmarshal(&raw); err != nil {
		slog.Error("unmarshal config", "file", filepath.Base(filePath), "error", err)
		return nil, nil, err
	}

	cfg := newConfig(raw)

	if err := cfg.validate(); err != nil {
		slog.Error("validate config", "file", filepath.Base(filePath), "error", err)
		return nil, nil, err
	}

	return cfg, cfg.Database(), nil
}

func newConfig(raw rawConfig) *config {
	return &config{
		address: strings.TrimSpace(raw.Address),
		database: databaseConfig{
			host:     strings.TrimSpace(raw.Database.Host),
			port:     raw.Database.Port,
			dbName:   strings.TrimSpace(raw.Database.DBName),
			username: strings.TrimSpace(raw.Database.Username),
			password: raw.Database.Password,
			logLevel: strings.TrimSpace(raw.Database.LogLevel),
		},
	}
}

var _ ConfigPathProvider = DefaultConfigPathProvider

func DefaultConfigPathProvider() (ConfigPath, error) {
	args := os.Args[1:]

	for i := range args {
		arg := strings.TrimSpace(args[i])

		if arg == "--config" || arg == "-c" {
			if i+1 >= len(args) {
				return "", errors.New("missing value for config flag")
			}

			value := strings.TrimSpace(args[i+1])
			if value == "" {
				return "", errors.New("empty config path")
			}

			return ConfigPath(value), nil
		}

		if after, ok := strings.CutPrefix(arg, "--config="); ok {
			value := strings.TrimSpace(after)
			if value == "" {
				return "", errors.New("empty config path")
			}

			return ConfigPath(value), nil
		}

		if after, ok := strings.CutPrefix(arg, "-c="); ok {
			value := strings.TrimSpace(after)
			if value == "" {
				return "", errors.New("empty config path")
			}

			return ConfigPath(value), nil
		}
	}

	return "", errors.New("config flag not provided")
}
