package setup

import (
	"basesdk"
	"basesdk/configs"
	"basesdk/connection"
	"basesdk/httpapi"
	"basesdk/properties"
	"basesdk/security/handlers"
	"basesdk/security/jwt"
	"basesdk/security/postgres"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/setup/migrations"
	"context"
	"io"
	"io/fs"
	"os"

	"go.uber.org/fx"
)

type Service struct {
	version                    string
	migrations                 []fs.FS
	configPathProvider         configs.ConfigPathProvider
	applicationConfigsProvider configs.ApplicationConfigsProvider
}

func NewService(opts ...Option) *Service {
	s := &Service{
		configPathProvider:         configs.DefaultConfigPathProvider,
		applicationConfigsProvider: configs.DefaultConfigsProvider,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) Run(opts ...fx.Option) {
	action, ok := migrations.ParseMigrateCommand(os.Args[1:])
	if ok {
		s.runMigration(action)
		return
	}

	var preOptions = []fx.Option{
		fx.Provide(s.configPathProvider, s.applicationConfigsProvider),
		fx.Provide(connection.NewConnection),
		fx.Provide(
			migrations.ProvideFSSources(basesdk.MigrationsFS),
			migrations.ProvideFSSources(s.migrations...),
			fx.Annotate(migrations.NewMigrationRunner,
				fx.ParamTags(``, migrations.GroupFSSources),
			),
		),
		fx.Module("properties",
			fx.Provide(properties.NewSystemProperties),
		),
		fx.Module("security",
			fx.Provide(
				fx.Annotate(
					postgres.NewSystemUserRepository,
					fx.As(new(repositories.SystemUserRepository)),
				),
			),
			fx.Provide(
				jwt.NewKeyStore,
				jwt.NewTokenService,
				usecases.NewSecurityUsecase,
			),
			fx.Provide(
				httpapi.AsRoute(handlers.SystemUserHandler),
			),
		),
		httpapi.Module,
	}

	var postOptions = []fx.Option{
		fx.Invoke(
			httpapi.StartWebServer,
		),
	}

	opts = append(preOptions, opts...)
	opts = append(opts, postOptions...)

	app := fx.New(opts...)
	app.Run()
}

func (s *Service) runMigration(action string) {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(s.configPathProvider, s.applicationConfigsProvider),
		fx.Provide(connection.NewConnection),
		fx.Provide(
			migrations.ProvideFSSources(s.migrations...),
			migrations.ProvideFSSources(basesdk.MigrationsFS),
			fx.Annotate(migrations.NewMigrationRunner,
				fx.ParamTags(``, migrations.GroupFSSources),
			),
		),
		fx.Invoke(func(mr *migrations.MigrationRunner, shutdowner fx.Shutdowner) {
			ctx := context.Background()

			exitCode := 0

			switch action {
			case "up":
				if err := mr.Up(ctx); err != nil {
					exitCode = 1
				}
			case "down":
				if err := mr.Down(ctx); err != nil {
					exitCode = 1
				}
			case "status":
				if err := mr.Status(ctx); err != nil {
					exitCode = 1
				}
			case "script":
				str, err := mr.SQLScript()
				if err != nil {
					exitCode = 1
					break
				}
				if _, err := io.WriteString(os.Stdout, str); err != nil {
					exitCode = 1
				}
			default:
				exitCode = 1
			}

			_ = shutdowner.Shutdown(fx.ExitCode(exitCode))
		}),
	)
	app.Run()
}
