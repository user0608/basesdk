import { useMemo, useState } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Link, Navigate, Outlet, useParams, useRoutes } from "react-router-dom";
import { AuthProvider } from "../auth/AuthProvider";
import { DialogProvider } from "../components/dialog/DialogProvider";
import { useConfirmDialog } from "../components/dialog/useConfirmDialog";
import { ToastProvider } from "../components/toast/ToastProvider";
import { FiDatabase, FiList, FiLogOut, FiUserCheck } from "react-icons/fi";
import { RegistryProvider } from "../platform/RegistryProvider";
import { RequireGuest } from "../auth/guards/RequireGuest";
import { RequireSystem } from "../auth/guards/RequireSystem";
import { RequireTenant } from "../auth/guards/RequireTenant";
import { SystemLoginPage } from "../auth/pages/SystemLoginPage";
import { TenantPasswordChangePage } from "../auth/pages/TenantPasswordChangePage";
import { TenantLoginPage } from "../auth/pages/TenantLoginPage";
import { useAuth, useHttpApi } from "../auth/useAuth";
import { createTenantModuleRoutes } from "../platform/routes";
import { filterMenuTree } from "../platform/menu";
import { TenantWorkspaceLayout } from "../platform/WorkspaceLayout";
import { createSecurityService } from "../security/SecurityService";
import { securityModules } from "../security/modules";
import { securityRegistry } from "../security/registry";
import { SecurityServiceProvider } from "../security/useSecurityService";
import { SecurityGroupsPage } from "../security/pages/SecurityGroupsPage";
import { SecurityPermissionsPage } from "../security/pages/SecurityPermissionsPage";
import { SecurityRolesPage } from "../security/pages/SecurityRolesPage";
import { SecurityUsersPage } from "../security/pages/SecurityUsersPage";
import { ServiceProvider } from "../services/ServiceProvider";
import { SystemPropertyForm } from "../system/forms/SystemPropertyForm";
import { SystemTenantForm } from "../system/forms/SystemTenantForm";
import { SystemUserForm } from "../system/forms/SystemUserForm";
import { SystemPropertiesPage } from "../system/pages/SystemPropertiesPage";
import { SystemBreadcrumbs } from "../system/SystemBreadcrumbs";
import { SystemTenantsPage } from "../system/pages/SystemTenantsPage";
import { SystemUsersPage } from "../system/pages/SystemUsersPage";
import type { ApplicationRegistry, MenuIcon, MenuTree, PageId } from "../platform/types";

const renderIcon = (icon: MenuIcon | undefined, className: string) => {
  if (!icon) return null;

  if (typeof icon === "string") {
    return <img src={icon} alt="" className={className} />;
  }

  const Icon = icon;
  return <Icon className={className} />;
};

type ApplicationProps<TRegistry extends ApplicationRegistry, TPermission extends string = string> = {
  getBaseUrl?: () => string | Promise<string>;
  modules: MenuTree<PageId<TRegistry>, TPermission>;
  registry: TRegistry;
  queryClient?: QueryClient;
  storageKeyPrefix?: string;
  appTitle?: string;
  appSubtitle?: string;
  loginTitle?: string;
  loginSubtitle?: string;
  systemLoginTitle?: string;
  systemLoginSubtitle?: string;
  defaultTenantCodigo?: string;
  unauthorizedElement?: React.ReactNode;
};

const DefaultUnauthorized = () => {
  return (
    <section className="rounded-2xl bg-ui-panel p-6 shadow-sm ring-1 ring-inset ring-ui-border/35">
      <h2 className="text-xl font-semibold text-ui-text">Sin permiso</h2>
      <p className="mt-2 text-sm leading-6 text-ui-text-muted">
        La sesion actual no tiene permisos para abrir esta operacion.
      </p>
    </section>
  );
};

const RootRedirect = () => {
  const { isReady, tenantSession } = useAuth();

  if (!isReady) return null;
  if (tenantSession) return <Navigate to="/app" replace />;
  return <Navigate to="/login" replace />;
};

const TenantLauncher = <TRegistry extends ApplicationRegistry, TPermission extends string = string>({
  modules,
  title,
}: {
  modules: MenuTree<PageId<TRegistry>, TPermission>;
  title: string;
}) => {
  const { tenantSession, logoutTenant } = useAuth();
  const confirm = useConfirmDialog();
  const visibleModules = filterMenuTree(modules, tenantSession?.permissions ?? []);

  const onLogout = async () => {
    const shouldLogout = await confirm({
      title: "Cerrar sesion",
      content: "Quieres cerrar la sesion actual?",
      confirmLabel: "Cerrar sesion",
      cancelLabel: "Continuar trabajando",
    });

    if (shouldLogout) {
      logoutTenant();
    }
  };

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,rgba(193,122,43,0.16),transparent_32%),linear-gradient(135deg,#f7f2e8_0%,#e9edf1_58%,#dfe9e8_100%)] px-4 py-5 text-ui-text sm:px-6 lg:px-8 lg:py-8">
      <div className="grid gap-6 lg:gap-8">
        <header className="flex items-center justify-between gap-6 border-b border-ui-primary/15 pb-4">
          <h1 className="min-w-0 truncate text-2xl font-semibold tracking-tight text-ui-primary sm:text-3xl">{title}</h1>
          <button
            type="button"
            aria-label="Cerrar sesion"
            title="Cerrar sesion"
            className="grid size-9 shrink-0 place-items-center rounded-xl text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
            onClick={onLogout}
          >
            <FiLogOut size={18} />
          </button>
        </header>

        {visibleModules.length === 0 ? (
          <section className="rounded-3xl bg-ui-panel p-8 shadow-sm ring-1 ring-inset ring-ui-border/30">
            <h2 className="text-2xl font-semibold tracking-tight text-ui-text">Sin modulos disponibles</h2>
            <p className="mt-3 max-w-2xl text-sm leading-6 text-ui-text-muted">
              La sesion actual no tiene permisos para ver ningun modulo registrado.
            </p>
          </section>
        ) : (
          <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
            {visibleModules.map((module) => (
              <Link
                key={module.id}
                to={module.path}
                className="group rounded-2xl bg-ui-panel/95 p-4 shadow-sm shadow-ui-primary/5 ring-1 ring-inset ring-ui-border/45 transition-all hover:-translate-y-0.5 hover:shadow-lg hover:shadow-ui-primary/10"
              >
                <div className="flex items-center gap-3">
                  {module.icon && (
                    <span className="grid size-10 shrink-0 place-items-center rounded-xl bg-ui-primary text-ui-text-inverse shadow-sm shadow-ui-primary/20">
                      {renderIcon(module.icon, "size-5 object-contain")}
                    </span>
                  )}
                  <div className="min-w-0">
                    <h3 className="truncate text-lg font-semibold tracking-tight text-ui-text">{module.label}</h3>
                    <div className="text-sm font-medium text-ui-text-soft transition-colors group-hover:text-ui-accent">Abrir modulo</div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </main>
  );
};

const TenantWorkspaceShell = <TRegistry extends ApplicationRegistry, TPermission extends string = string>({
  modules,
  title,
  subtitle,
}: {
  modules: MenuTree<PageId<TRegistry>, TPermission>;
  title: string;
  subtitle: string;
}) => {
  const { tenantSession } = useAuth();

  return (
    <TenantWorkspaceLayout
      title={title}
      subtitle={subtitle}
      modules={modules}
      permissions={tenantSession?.permissions ?? []}
    >
      <Outlet />
    </TenantWorkspaceLayout>
  );
};

const SystemLauncher = ({ title }: { title: string }) => {
  const { systemSession, logoutSystem } = useAuth();
  const confirm = useConfirmDialog();
  const cards = [
    {
      to: "/system/users",
      icon: FiUserCheck,
      label: "Usuarios system",
      description: "Administra las cuentas con acceso al panel administrativo global.",
      action: "Abrir usuarios",
    },
    {
      to: "/system/tenants",
      icon: FiList,
      label: "Tenants",
      description: "Administra clientes, estados, limites y entra a la seguridad de cada tenant.",
      action: "Abrir tenants",
    },
    {
      to: "/system/properties",
      icon: FiDatabase,
      label: "Properties system",
      description: "Configura parametros globales como JWT, integraciones y valores de plataforma.",
      action: "Abrir properties",
    },
  ];

  const onLogout = async () => {
    const shouldLogout = await confirm({
      title: "Cerrar sesion",
      content: "Quieres cerrar la sesion system actual?",
      confirmLabel: "Cerrar sesion",
      cancelLabel: "Continuar trabajando",
    });

    if (shouldLogout) {
      logoutSystem();
    }
  };

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top_left,rgba(193,122,43,0.16),transparent_32%),linear-gradient(135deg,#f7f2e8_0%,#e9edf1_58%,#dfe9e8_100%)] px-4 py-5 text-ui-text sm:px-6 lg:px-8 lg:py-8">
      <div className="grid gap-6 lg:gap-8">
        <header className="flex items-center justify-between gap-6 border-b border-ui-primary/15 pb-4">
          <h1 className="min-w-0 truncate text-2xl font-semibold tracking-tight text-ui-primary sm:text-3xl">{title}</h1>
          <div className="flex shrink-0 items-center gap-1">
            <button
              type="button"
              aria-label="Cerrar sesion"
              title="Cerrar sesion"
              className="grid size-9 place-items-center rounded-xl text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
              onClick={onLogout}
            >
              <FiLogOut size={18} />
            </button>
          </div>
        </header>

        <section className="overflow-hidden rounded-3xl bg-ui-panel/95 shadow-xl shadow-ui-primary/10 ring-1 ring-inset ring-ui-border/45">
          <div className="relative overflow-hidden border-b border-ui-primary/15 bg-ui-primary px-5 py-6 text-ui-text-inverse sm:px-6">
            <div className="absolute inset-y-0 right-0 w-1/2 bg-[radial-gradient(circle_at_top_right,rgba(193,122,43,0.36),transparent_56%)]" />
            <div className="relative">
              <div className="text-xs font-semibold uppercase tracking-[0.16em] text-ui-accent">System Admin</div>
              <h2 className="mt-2 text-2xl font-semibold tracking-tight text-ui-text-inverse">Centro de administracion</h2>
              <p className="mt-2 max-w-3xl text-sm leading-6 text-ui-text-inverse/75">
                Sesion activa para <span className="font-medium text-ui-text-inverse">{systemSession?.username ?? "usuario system"}</span>. Gestiona configuracion global y opera tenants sin entrar como usuario tenant.
              </p>
            </div>
          </div>

          <div className="grid gap-4 p-4 sm:p-5 md:grid-cols-2">
            {cards.map((card) => {
              const Icon = card.icon;

              return (
                <Link
                  key={card.to}
                  to={card.to}
                  className="group rounded-2xl bg-ui-surface p-4 shadow-sm shadow-ui-primary/5 ring-1 ring-inset ring-ui-border/45 transition-all hover:-translate-y-0.5 hover:bg-ui-surface-hover hover:shadow-lg hover:shadow-ui-primary/10"
                >
                  <div className="flex items-start gap-3">
                    <span className="grid size-11 shrink-0 place-items-center rounded-2xl bg-ui-primary text-ui-text-inverse shadow-sm shadow-ui-primary/20 ring-1 ring-inset ring-ui-text-inverse/10">
                      <Icon className="size-5" />
                    </span>
                    <div className="min-w-0">
                      <h3 className="text-lg font-semibold tracking-tight text-ui-text">{card.label}</h3>
                      <p className="mt-1 text-sm leading-6 text-ui-text-muted">{card.description}</p>
                      <div className="mt-4 text-sm font-semibold text-ui-primary transition-colors group-hover:text-ui-primary-hover">
                        {card.action}
                      </div>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>
        </section>
      </div>
    </main>
  );
};

const createSystemSecurityModules = (tenantCodigo: string) => [
  {
    id: "system-tenant-security",
    label: "Seguridad",
    icon: FiLogOut,
    path: `/system/tenants/${encodeURIComponent(tenantCodigo)}/security`,
    order: 1,
    children: [
      {
        id: "system-tenant-security-main",
        label: "Seguridad",
        path: "",
        order: 1,
        children: [
          { id: "system-tenant-users", label: "Usuarios", path: `/system/tenants/${encodeURIComponent(tenantCodigo)}/security/users`, componentId: "security-users-page", order: 1 },
          { id: "system-tenant-groups", label: "Grupos", path: `/system/tenants/${encodeURIComponent(tenantCodigo)}/security/groups`, componentId: "security-groups-page", order: 2 },
          { id: "system-tenant-roles", label: "Roles", path: `/system/tenants/${encodeURIComponent(tenantCodigo)}/security/roles`, componentId: "security-roles-page", order: 3 },
          { id: "system-tenant-permissions", label: "Permisos", path: `/system/tenants/${encodeURIComponent(tenantCodigo)}/security/permissions`, componentId: "security-permissions-page", order: 4 },
        ],
      },
    ],
  },
] satisfies MenuTree<string, string>;

const SystemTenantSecurityShell = ({ title }: { title: string }) => {
  const api = useHttpApi();
  const { tenantCodigo = "" } = useParams();
  const modules = useMemo(() => createSystemSecurityModules(tenantCodigo), [tenantCodigo]);
  const service = useMemo(
    () =>
      createSecurityService(api, {
        auth: "system",
        resourcePrefix: `/api/v1/system/tenants/${encodeURIComponent(tenantCodigo)}`,
        permissionsPrefix: "/api/v1/system/permissions",
        queryKeyPrefix: ["security", "system", tenantCodigo],
      }),
    [api, tenantCodigo],
  );

  return (
    <SecurityServiceProvider service={service}>
      <TenantWorkspaceLayout
        title={title}
        subtitle={tenantCodigo}
        homePath="/system"
        modules={modules}
        permissions={[]}
        breadcrumbs={
          <SystemBreadcrumbs
            items={[
              { label: "System", to: "/system" },
              { label: "Tenants", to: "/system/tenants" },
              { label: tenantCodigo },
              { label: "Seguridad" },
            ]}
          />
        }
        actions={
          <div className="flex items-center gap-2">
            <Link
              to="/system/tenants"
              aria-label="Volver a tenants"
              title="Volver a tenants"
              className="inline-grid size-8 place-items-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
            >
              <FiList size={17} />
            </Link>
            <Link
              to={`/system/tenants/${encodeURIComponent(tenantCodigo)}/properties`}
              aria-label="Abrir properties"
              title="Abrir properties"
              className="inline-grid size-8 place-items-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
            >
              <FiDatabase size={17} />
            </Link>
          </div>
        }
      >
        <Outlet />
      </TenantWorkspaceLayout>
    </SecurityServiceProvider>
  );
};

const ApplicationRoutes = <TRegistry extends ApplicationRegistry, TPermission extends string = string>({
  modules,
  registry,
  appTitle = "Base ERP",
  appSubtitle = "Selecciona un modulo para entrar a su espacio de trabajo.",
  loginTitle = "Tenant Login",
  loginSubtitle = "Accede al entorno principal del ERP.",
  systemLoginTitle = "System Login",
  systemLoginSubtitle = "Accede al entorno administrativo global.",
  defaultTenantCodigo,
  unauthorizedElement,
}: Omit<ApplicationProps<TRegistry, TPermission>, "getBaseUrl" | "queryClient" | "storageKeyPrefix">) => {
  return useRoutes([
    { path: "/", element: <RootRedirect /> },
    {
      path: "/login",
      element: <RequireGuest scope="tenant" />,
      children: [
        {
          index: true,
          element: (
            <TenantLoginPage
              title={loginTitle}
              subtitle={loginSubtitle}
              defaultTenantCodigo={defaultTenantCodigo}
              redirectTo="/app"
            />
          ),
        },
      ],
    },
    {
      path: "/system/login",
      element: <RequireGuest scope="system" />,
      children: [
        {
          index: true,
          element: <SystemLoginPage title={systemLoginTitle} subtitle={systemLoginSubtitle} />,
        },
      ],
    },
    {
      path: "/system",
      element: <RequireSystem />,
      children: [
        {
          index: true,
          element: <SystemLauncher title={systemLoginTitle} />,
        },
        { path: "users", element: <SystemUsersPage /> },
        { path: "properties", element: <SystemPropertiesPage /> },
        { path: "tenants", element: <SystemTenantsPage /> },
        { path: "tenants/:tenantCodigo/properties", element: <SystemPropertiesPage tenantScoped /> },
        {
          path: "tenants/:tenantCodigo/security",
          element: <SystemTenantSecurityShell title={systemLoginTitle} />,
          children: [
            { index: true, element: <Navigate to="users" replace /> },
            { path: "users", element: <SecurityUsersPage /> },
            { path: "groups", element: <SecurityGroupsPage /> },
            { path: "roles", element: <SecurityRolesPage /> },
            { path: "permissions", element: <SecurityPermissionsPage /> },
          ],
        },
      ],
    },
    {
      element: <RequireTenant />,
      children: [
        { path: "/change-password", element: <TenantPasswordChangePage /> },
        { path: "/app", element: <TenantLauncher modules={modules} title={appTitle} /> },
        {
          element: <TenantWorkspaceShell modules={modules} title={appTitle} subtitle="Tenant workspace" />,
          children: createTenantModuleRoutes({
            modules,
            registry,
            unauthorizedElement: unauthorizedElement ?? <DefaultUnauthorized />,
          }),
        },
      ],
    },
  ]);
};

export const Application = <TRegistry extends ApplicationRegistry, TPermission extends string = string>({
  getBaseUrl = () => "/",
  queryClient,
  storageKeyPrefix,
  ...routesProps
}: ApplicationProps<TRegistry, TPermission>) => {
  const [client] = useState(() => queryClient ?? new QueryClient());
  const finalRegistry = {
    pages: {
      ...securityRegistry.pages,
      ...routesProps.registry.pages,
    },
    forms: {
      ...securityRegistry.forms,
      "system-tenant-form": SystemTenantForm,
      "system-property-form": SystemPropertyForm,
      "system-user-form": SystemUserForm,
      ...(routesProps.registry.forms ?? {}),
    },
  } as ApplicationRegistry;
  const finalModules = [...securityModules, ...routesProps.modules] as MenuTree<string, TPermission>;

  return (
    <QueryClientProvider client={client}>
      <AuthProvider getBaseUrl={getBaseUrl} storageKeyPrefix={storageKeyPrefix}>
        <BrowserRouter>
          <ToastProvider>
            <ServiceProvider>
              <RegistryProvider registry={finalRegistry}>
                <DialogProvider>
                  <ApplicationRoutes {...routesProps} registry={finalRegistry} modules={finalModules} />
                </DialogProvider>
              </RegistryProvider>
            </ServiceProvider>
          </ToastProvider>
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  );
};
