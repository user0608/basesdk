package setup

import (
	"basesdk"
	"basesdk/auth"
	"basesdk/auth/jwt"
	"basesdk/configs"
	"basesdk/connection"
	"basesdk/httpapi"
	"basesdk/properties"
	propertieshandlers "basesdk/properties/handlers"
	"basesdk/security/handlers"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/setup/migrations"
	setuppermissions "basesdk/setup/permissions"
	"context"
	"io/fs"
	"os"

	"go.uber.org/fx"
)

type Service struct {
	version                    string
	migrations                 []fs.FS
	permissions                []fs.FS
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
	options = append(options, fx.Invoke(syncPermissions))
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
			setuppermissions.ProvideFSSources(basesdk.PermissionsFS),
			setuppermissions.ProvideFSSources(s.permissions...),
			fx.Annotate(
				migrations.NewMigrationRunner,
				fx.ParamTags(``, migrations.GroupFSSources),
			),
			fx.Annotate(
				setuppermissions.NewPermissionSynchronizer,
				fx.ParamTags(``, setuppermissions.GroupFSSources),
			),
		),
	}
}

func syncPermissions(synchronizer *setuppermissions.PermissionSynchronizer) error {
	return synchronizer.Run(context.Background())
}

func (s *Service) applicationOptions() []fx.Option {
	return []fx.Option{
		propertiesModule(),
		securityModule(),
		httpapi.Module,
	}
}

func propertiesModule() fx.Option {
	return fx.Module("properties",
		fx.Provide(
			properties.NewSystemProperties,
			properties.NewTenantSystemProperties,
			httpapi.AsRoute(propertieshandlers.SystemPropertiesListHandler),
			httpapi.AsRoute(propertieshandlers.SystemPropertyCreateHandler),
			httpapi.AsRoute(propertieshandlers.SystemPropertyFindHandler),
			httpapi.AsRoute(propertieshandlers.SystemPropertyUpdateHandler),
			httpapi.AsRoute(propertieshandlers.SystemPropertyDeleteHandler),
			httpapi.AsRoute(propertieshandlers.SystemTenantPropertiesListHandler),
			httpapi.AsRoute(propertieshandlers.SystemTenantPropertyCreateHandler),
			httpapi.AsRoute(propertieshandlers.SystemTenantPropertyFindHandler),
			httpapi.AsRoute(propertieshandlers.SystemTenantPropertyUpdateHandler),
			httpapi.AsRoute(propertieshandlers.SystemTenantPropertyDeleteHandler),
		),
	)
}

func securityModule() fx.Option {
	return fx.Module("security",
		fx.Provide(securityProviders()...),
	)
}

func securityProviders() []any {
	providers := []any{
		repositories.NewSystemUserRepository,
		repositories.NewAppUserRepository,
		repositories.NewRoleRepository,
		repositories.NewGroupRepository,
		repositories.NewPermissionRepository,
		repositories.NewTenantRepository,
		fx.Annotate(
			usecases.NewAuthorizationUsecase,
			fx.As(new(auth.PermissionValidator)),
		),
		jwt.NewKeyStore,
		jwt.NewTokenService,
		usecases.NewSecurityUsecase,
		usecases.NewSystemUsersUsecase,
		usecases.NewTenantUsersUsecase,
		usecases.NewTenantRolesUsecase,
		usecases.NewTenantGroupsUsecase,
		usecases.NewTenantsUsecase,
		usecases.NewPermissionsUsecase,
		httpapi.AsRoute(handlers.SystemUserHandler),
		httpapi.AsRoute(handlers.TenantUserHandler),
	}

	providers = append(providers, systemRoutes()...)
	providers = append(providers, tenantRoutes()...)

	return providers
}

func systemRoutes() []any {
	return []any{
		httpapi.AsRoute(handlers.SystemUsersListHandler),
		httpapi.AsRoute(handlers.SystemUserCreateHandler),
		httpapi.AsRoute(handlers.SystemUserFindHandler),
		httpapi.AsRoute(handlers.SystemUserUpdateHandler),
		httpapi.AsRoute(handlers.SystemUsersEnableHandler),
		httpapi.AsRoute(handlers.SystemUsersDisableHandler),
		httpapi.AsRoute(handlers.SystemUsersDeleteHandler),
		httpapi.AsRoute(handlers.SystemTenantsListHandler),
		httpapi.AsRoute(handlers.SystemTenantCreateHandler),
		httpapi.AsRoute(handlers.SystemTenantFindHandler),
		httpapi.AsRoute(handlers.SystemTenantUpdateHandler),
		httpapi.AsRoute(handlers.SystemTenantsEnableHandler),
		httpapi.AsRoute(handlers.SystemTenantsDisableHandler),
		httpapi.AsRoute(handlers.SystemTenantUsersListHandler),
		httpapi.AsRoute(handlers.SystemTenantUserCreateHandler),
		httpapi.AsRoute(handlers.SystemTenantUserFindHandler),
		httpapi.AsRoute(handlers.SystemTenantUserUpdateHandler),
		httpapi.AsRoute(handlers.SystemTenantUserPasswordHandler),
		httpapi.AsRoute(handlers.SystemTenantUsersEnableHandler),
		httpapi.AsRoute(handlers.SystemTenantUsersDisableHandler),
		httpapi.AsRoute(handlers.SystemTenantUsersDeleteHandler),
		httpapi.AsRoute(handlers.SystemTenantUserPermissionsHandler),
		httpapi.AsRoute(handlers.SystemTenantUserRolesHandler),
		httpapi.AsRoute(handlers.SystemTenantUserReplaceRolesHandler),
		httpapi.AsRoute(handlers.SystemTenantUserGroupsHandler),
		httpapi.AsRoute(handlers.SystemTenantUserReplaceGroupsHandler),
		httpapi.AsRoute(handlers.SystemTenantRolesListHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleCreateHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleFindHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleUpdateHandler),
		httpapi.AsRoute(handlers.SystemTenantRolesEnableHandler),
		httpapi.AsRoute(handlers.SystemTenantRolesDisableHandler),
		httpapi.AsRoute(handlers.SystemTenantRolesDeleteHandler),
		httpapi.AsRoute(handlers.SystemTenantRolePermissionsHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleReplacePermissionsHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleUsersHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleReplaceUsersHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleGroupsHandler),
		httpapi.AsRoute(handlers.SystemTenantRoleReplaceGroupsHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupsListHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupCreateHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupFindHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupUpdateHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupsEnableHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupsDisableHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupsDeleteHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupUsersHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupReplaceUsersHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupRolesHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupReplaceRolesHandler),
		httpapi.AsRoute(handlers.SystemTenantGroupPermissionsHandler),
		httpapi.AsRoute(handlers.SystemPermissionsListHandler),
		httpapi.AsRoute(handlers.SystemPermissionFindHandler),
		httpapi.AsRoute(handlers.SystemPermissionRolesHandler),
		httpapi.AsRoute(handlers.SystemPermissionGroupsHandler),
		httpapi.AsRoute(handlers.SystemPermissionUsersHandler),
	}
}

func tenantRoutes() []any {
	return []any{
		httpapi.AsRoute(handlers.TenantUsersListHandler),
		httpapi.AsRoute(handlers.TenantUserCreateHandler),
		httpapi.AsRoute(handlers.TenantUserFindHandler),
		httpapi.AsRoute(handlers.TenantUserUpdateHandler),
		httpapi.AsRoute(handlers.TenantUserPasswordHandler),
		httpapi.AsRoute(handlers.TenantUsersEnableHandler),
		httpapi.AsRoute(handlers.TenantUsersDisableHandler),
		httpapi.AsRoute(handlers.TenantUsersDeleteHandler),
		httpapi.AsRoute(handlers.TenantUserPermissionsHandler),
		httpapi.AsRoute(handlers.TenantUserRolesHandler),
		httpapi.AsRoute(handlers.TenantUserReplaceRolesHandler),
		httpapi.AsRoute(handlers.TenantUserGroupsHandler),
		httpapi.AsRoute(handlers.TenantUserReplaceGroupsHandler),
		httpapi.AsRoute(handlers.TenantRolesListHandler),
		httpapi.AsRoute(handlers.TenantRoleCreateHandler),
		httpapi.AsRoute(handlers.TenantRoleFindHandler),
		httpapi.AsRoute(handlers.TenantRoleUpdateHandler),
		httpapi.AsRoute(handlers.TenantRolesEnableHandler),
		httpapi.AsRoute(handlers.TenantRolesDisableHandler),
		httpapi.AsRoute(handlers.TenantRolesDeleteHandler),
		httpapi.AsRoute(handlers.TenantRolePermissionsHandler),
		httpapi.AsRoute(handlers.TenantRoleReplacePermissionsHandler),
		httpapi.AsRoute(handlers.TenantRoleUsersHandler),
		httpapi.AsRoute(handlers.TenantRoleReplaceUsersHandler),
		httpapi.AsRoute(handlers.TenantRoleGroupsHandler),
		httpapi.AsRoute(handlers.TenantRoleReplaceGroupsHandler),
		httpapi.AsRoute(handlers.TenantGroupsListHandler),
		httpapi.AsRoute(handlers.TenantGroupCreateHandler),
		httpapi.AsRoute(handlers.TenantGroupFindHandler),
		httpapi.AsRoute(handlers.TenantGroupUpdateHandler),
		httpapi.AsRoute(handlers.TenantGroupsEnableHandler),
		httpapi.AsRoute(handlers.TenantGroupsDisableHandler),
		httpapi.AsRoute(handlers.TenantGroupsDeleteHandler),
		httpapi.AsRoute(handlers.TenantGroupUsersHandler),
		httpapi.AsRoute(handlers.TenantGroupReplaceUsersHandler),
		httpapi.AsRoute(handlers.TenantGroupRolesHandler),
		httpapi.AsRoute(handlers.TenantGroupReplaceRolesHandler),
		httpapi.AsRoute(handlers.TenantGroupPermissionsHandler),
		httpapi.AsRoute(handlers.TenantPermissionsListHandler),
		httpapi.AsRoute(handlers.TenantPermissionFindHandler),
		httpapi.AsRoute(handlers.TenantPermissionRolesHandler),
		httpapi.AsRoute(handlers.TenantPermissionGroupsHandler),
		httpapi.AsRoute(handlers.TenantPermissionUsersHandler),
		httpapi.AsRoute(handlers.TenantMeHandler),
		httpapi.AsRoute(handlers.TenantMePasswordHandler),
		httpapi.AsRoute(handlers.TenantMePermissionsHandler),
	}
}
