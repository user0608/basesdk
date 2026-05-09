package testdb

import (
	"basesdk"
	"basesdk/configs"
	"basesdk/connection"
	"basesdk/setup/migrations"
	"context"
	"io/fs"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	postgresImage    = "postgres:18-alpine"
	testDatabaseName = "testdb"
	testUsername     = "testuser"
	testPassword     = "testpass"
	testLogLevel     = "error"
)

type Config struct {
	host     string
	port     int
	username string
	password string
	database string
	logLevel string
}

var _ configs.DatabaseConfig = (*Config)(nil)

func (c *Config) GetHost() string {
	return c.host
}

func (c *Config) GetPort() int {
	return c.port
}

func (c *Config) GetUsername() string {
	return c.username
}

func (c *Config) GetPassword() string {
	return c.password
}

func (c *Config) GetDBName() string {
	return c.database
}

func (c *Config) GetLogLevel() string {
	return c.logLevel
}

var (
	postgresOnce sync.Once

	sharedStorage   connection.StorageManager
	sharedContainer *tcpostgres.PostgresContainer
	sharedErr       error
)

func NewPostgresStorage(t *testing.T, additionalFileSystems ...fs.FS) connection.StorageManager {
	t.Helper()

	postgresOnce.Do(func() {
		sharedStorage, sharedContainer, sharedErr = startPostgres()
	})

	require.NoError(t, sharedErr)
	require.NotNil(t, sharedStorage)

	migrationFileSystems := defaultMigrationFileSystems(additionalFileSystems...)
	require.NoError(t, runMigrations(sharedStorage, migrationFileSystems...))

	t.Cleanup(func() {
		dropPublicTables(t, sharedStorage)
	})

	return sharedStorage
}

func TerminatePostgres(ctx context.Context) error {
	if sharedContainer == nil {
		return nil
	}

	return testcontainers.TerminateContainer(sharedContainer)
}

func startPostgres() (connection.StorageManager, *tcpostgres.PostgresContainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	container, err := tcpostgres.Run(
		ctx,
		postgresImage,
		tcpostgres.WithDatabase(testDatabaseName),
		tcpostgres.WithUsername(testUsername),
		tcpostgres.WithPassword(testPassword),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, err
	}

	storage, err := newStorageFromContainer(ctx, container)
	if err != nil {
		_ = testcontainers.TerminateContainer(container)
		return nil, nil, err
	}

	return storage, container, nil
}

func newStorageFromContainer(ctx context.Context, container *tcpostgres.PostgresContainer) (connection.StorageManager, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(mappedPort.Port())
	if err != nil {
		return nil, err
	}

	return connection.NewConnection(&Config{
		host:     host,
		port:     port,
		username: testUsername,
		password: testPassword,
		database: testDatabaseName,
		logLevel: testLogLevel,
	})
}

func defaultMigrationFileSystems(additionalFileSystems ...fs.FS) []fs.FS {
	fileSystems := append([]fs.FS{}, basesdk.MigrationsFS)
	fileSystems = append(fileSystems, additionalFileSystems...)

	return fileSystems
}

func runMigrations(storage connection.StorageManager, fileSystems ...fs.FS) error {
	if len(fileSystems) == 0 {
		return nil
	}

	runner := migrations.NewMigrationRunner(
		storage,
		[]migrations.FileSystemSources{
			fileSystems,
		},
	)

	return runner.Up(context.Background())
}

func dropPublicTables(t *testing.T, storage connection.StorageManager) {
	t.Helper()

	err := storage.Conn(context.Background()).Exec(`
		DO $$
		DECLARE
			table_record RECORD;
		BEGIN
			FOR table_record IN (
				SELECT tablename
				FROM pg_tables
				WHERE schemaname = 'public'
			) LOOP
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(table_record.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`).Error
	require.NoError(t, err)
}
