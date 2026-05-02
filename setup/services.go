package setup

import (
	"basesdk/configs"
	"basesdk/connection"
	"basesdk/setup/migrations"
	"io/fs"

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
	var preOptions = []fx.Option{
		fx.Provide(s.configPathProvider, s.applicationConfigsProvider),
		fx.Provide(connection.NewConnection),
		fx.Provide(migrations.ProvideFSSources(s.migrations...),
			fx.Annotate(migrations.NewMigrationRunner,
				fx.ParamTags(``, migrations.GroupFSSources),
			),
		),
	}
	var postOptions = []fx.Option{}

	opts = append(preOptions, opts...)
	opts = append(opts, postOptions...)

	app := fx.New(opts...)
	app.Run()
}
