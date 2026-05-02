package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigPathProvider(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    ConfigPath
		wantErr bool
	}{
		{
			name: "long flag with value",
			args: []string{"app", "--config", "config.yaml"},
			want: "config.yaml",
		},
		{
			name: "short flag with value",
			args: []string{"app", "-c", "config.yaml"},
			want: "config.yaml",
		},
		{
			name: "long flag equals",
			args: []string{"app", "--config=config.yaml"},
			want: "config.yaml",
		},
		{
			name: "short flag equals",
			args: []string{"app", "-c=config.yaml"},
			want: "config.yaml",
		},
		{
			name:    "missing config flag",
			args:    []string{"app"},
			wantErr: true,
		},
		{
			name:    "missing value after long flag",
			args:    []string{"app", "--config"},
			wantErr: true,
		},
		{
			name:    "missing value after short flag",
			args:    []string{"app", "-c"},
			wantErr: true,
		},
		{
			name:    "empty long flag equals",
			args:    []string{"app", "--config="},
			wantErr: true,
		},
		{
			name:    "empty short flag equals",
			args:    []string{"app", "-c="},
			wantErr: true,
		},
		{
			name: "trims spaces",
			args: []string{"app", "--config", "  config.yaml  "},
			want: "config.yaml",
		},
	}

	originalArgs := os.Args
	t.Cleanup(func() {
		os.Args = originalArgs
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args

			got, err := DefaultConfigPathProvider()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestDefaultConfigsProvider(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `
address: ":8080"
database:
  host: "localhost"
  port: 5432
  db_name: "my_app"
  username: "postgres"
  password: "secret"
  log_level: "debug"
`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	appCfg, dbCfg, err := DefaultConfigsProvider(ConfigPath(configPath))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if appCfg.ListenAddress() != ":8080" {
		t.Fatalf("expected address %q, got %q", ":8080", appCfg.ListenAddress())
	}

	if dbCfg.GetHost() != "localhost" {
		t.Fatalf("expected host %q, got %q", "localhost", dbCfg.GetHost())
	}

	if dbCfg.GetPort() != 5432 {
		t.Fatalf("expected port %d, got %d", 5432, dbCfg.GetPort())
	}

	if dbCfg.GetDBName() != "my_app" {
		t.Fatalf("expected db name %q, got %q", "my_app", dbCfg.GetDBName())
	}

	if dbCfg.GetUsername() != "postgres" {
		t.Fatalf("expected username %q, got %q", "postgres", dbCfg.GetUsername())
	}

	if dbCfg.GetPassword() != "secret" {
		t.Fatalf("expected password %q, got %q", "secret", dbCfg.GetPassword())
	}

	if dbCfg.GetLogLevel() != "debug" {
		t.Fatalf("expected log level %q, got %q", "debug", dbCfg.GetLogLevel())
	}
}

func TestDefaultConfigsProviderReturnsErrorWhenPathIsEmpty(t *testing.T) {
	appCfg, dbCfg, err := DefaultConfigsProvider(" ")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if appCfg != nil {
		t.Fatalf("expected nil app config, got %#v", appCfg)
	}

	if dbCfg != nil {
		t.Fatalf("expected nil database config, got %#v", dbCfg)
	}
}

func TestDefaultConfigsProviderReturnsErrorWhenFileDoesNotExist(t *testing.T) {
	appCfg, dbCfg, err := DefaultConfigsProvider("missing.yaml")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if appCfg != nil {
		t.Fatalf("expected nil app config, got %#v", appCfg)
	}

	if dbCfg != nil {
		t.Fatalf("expected nil database config, got %#v", dbCfg)
	}
}

func TestDefaultConfigsProviderReturnsErrorWhenConfigIsInvalid(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `
address: ""
database:
  host: ""
  port: 0
  db_name: ""
  username: ""
  password: ""
  log_level: ""
`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	appCfg, dbCfg, err := DefaultConfigsProvider(ConfigPath(configPath))

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if appCfg != nil {
		t.Fatalf("expected nil app config, got %#v", appCfg)
	}

	if dbCfg != nil {
		t.Fatalf("expected nil database config, got %#v", dbCfg)
	}
}

func TestNewConfigTrimsValues(t *testing.T) {
	raw := rawConfig{
		Address: "  :8080  ",
		Database: rawDatabaseConfig{
			Host:     "  localhost  ",
			Port:     5432,
			DBName:   "  my_app  ",
			Username: "  postgres  ",
			Password: "  secret  ",
			LogLevel: "  debug  ",
		},
	}

	cfg := newConfig(raw)

	if cfg.ListenAddress() != ":8080" {
		t.Fatalf("expected address %q, got %q", ":8080", cfg.ListenAddress())
	}

	db := cfg.Database()

	if db.GetHost() != "localhost" {
		t.Fatalf("expected host %q, got %q", "localhost", db.GetHost())
	}

	if db.GetDBName() != "my_app" {
		t.Fatalf("expected db name %q, got %q", "my_app", db.GetDBName())
	}

	if db.GetUsername() != "postgres" {
		t.Fatalf("expected username %q, got %q", "postgres", db.GetUsername())
	}

	if db.GetPassword() != "  secret  " {
		t.Fatalf("expected password to keep spaces, got %q", db.GetPassword())
	}

	if db.GetLogLevel() != "debug" {
		t.Fatalf("expected log level %q, got %q", "debug", db.GetLogLevel())
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &config{
				address: ":8080",
				database: databaseConfig{
					host:     "localhost",
					port:     5432,
					dbName:   "my_app",
					username: "postgres",
					logLevel: "debug",
				},
			},
		},
		{
			name: "missing address",
			cfg: &config{
				database: databaseConfig{
					host:     "localhost",
					port:     5432,
					dbName:   "my_app",
					username: "postgres",
					logLevel: "debug",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database host",
			cfg: &config{
				address: ":8080",
				database: databaseConfig{
					port:     5432,
					dbName:   "my_app",
					username: "postgres",
					logLevel: "debug",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid database port",
			cfg: &config{
				address: ":8080",
				database: databaseConfig{
					host:     "localhost",
					port:     0,
					dbName:   "my_app",
					username: "postgres",
					logLevel: "debug",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database name",
			cfg: &config{
				address: ":8080",
				database: databaseConfig{
					host:     "localhost",
					port:     5432,
					username: "postgres",
					logLevel: "debug",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database username",
			cfg: &config{
				address: ":8080",
				database: databaseConfig{
					host:     "localhost",
					port:     5432,
					dbName:   "my_app",
					logLevel: "debug",
				},
			},
			wantErr: true,
		},
		{
			name: "missing database log level",
			cfg: &config{
				address: ":8080",
				database: databaseConfig{
					host:     "localhost",
					port:     5432,
					dbName:   "my_app",
					username: "postgres",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
