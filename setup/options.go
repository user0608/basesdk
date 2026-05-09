package setup

import (
	"basesdk/configs"
	"io/fs"
)

type Option func(*Service)

func WithVersion(version string) Option {
	return func(s *Service) {
		s.version = version
	}
}

func WithMigrations(m ...fs.FS) Option {
	return func(s *Service) {
		for _, fsys := range m {
			if fsys != nil {
				s.migrations = append(s.migrations, fsys)
			}
		}
	}
}

func WithPermissions(p ...fs.FS) Option {
	return func(s *Service) {
		for _, fsys := range p {
			if fsys != nil {
				s.permissions = append(s.permissions, fsys)
			}
		}
	}
}

func WithConfigPathProvider(fn configs.ConfigPathProvider) Option {
	return func(s *Service) {
		s.configPathProvider = fn
	}
}

func WithApplicationConfigsProvider(fn configs.ApplicationConfigsProvider) Option {
	return func(s *Service) {
		s.applicationConfigsProvider = fn
	}
}
