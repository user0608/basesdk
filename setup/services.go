package setup

import (
	"basesdk"
	"basesdk/auth/jwt"
	"basesdk/configs"
	"basesdk/connection"
	"basesdk/httpapi"
	"basesdk/properties"
	"basesdk/security/handlers"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/setup/migrations"
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
	if action, ok := migrations.ParseMigrateCommand(os.Args[1:]); ok {
		s.runMigration(action)
		return
	}

	options := append(s.baseOptions(), s.applicationOptions()...)
	options = append(options, opts...)
	options = append(options, fx.Invoke(httpapi.StartWebServer))

	fx.New(options...).Run()
}

func (s *Service) baseOptions() []fx.Option {
	return []fx.Option{
		fx.Provide(s.configPathProvider, s.applicationConfigsProvider),
		fx.Provide(connection.NewConnection),
		fx.Provide(
			migrations.ProvideFSSources(basesdk.MigrationsFS),
			migrations.ProvideFSSources(s.migrations...),
			fx.Annotate(
				migrations.NewMigrationRunner,
				fx.ParamTags(``, migrations.GroupFSSources),
			),
		),
	}
}

func (s *Service) applicationOptions() []fx.Option {
	return []fx.Option{
		fx.Module("properties",
			fx.Provide(
				properties.NewSystemProperties,
				properties.NewTenantSystemProperties,
			),
		),
		fx.Module("security",
			fx.Provide(
				repositories.NewSystemUserRepository,
				repositories.NewAppUserRepository,
				jwt.NewKeyStore,
				jwt.NewTokenService,
				usecases.NewSecurityUsecase,
				httpapi.AsRoute(handlers.SystemUserHandler),
				httpapi.AsRoute(handlers.TenantUserHandler),
			),
		),
		httpapi.Module,
	}
}
