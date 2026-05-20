import { Link, Navigate, Outlet, useRoutes } from "react-router-dom";
import { Button } from "../components/actions/Button";
import { RequireGuest } from "../auth/guards/RequireGuest";
import { RequireTenant } from "../auth/guards/RequireTenant";
import { TenantLoginPage } from "../auth/pages/TenantLoginPage";
import { useAuth } from "../auth/useAuth";
import { createTenantModuleRoutes } from "../platform/routes";
import { filterMenuTree } from "../platform/menu";
import { TenantWorkspaceLayout } from "../platform/WorkspaceLayout";
import type { ComponentId, ComponentRegistry, MenuIcon, MenuTree } from "../platform/types";

const renderIcon = (icon: MenuIcon | undefined, className: string) => {
  if (!icon) return null;

  if (typeof icon === "string") {
    return <img src={icon} alt="" className={className} />;
  }

  const Icon = icon;
  return <Icon className={className} />;
};

type TenantApplicationProps<TRegistry extends ComponentRegistry, TPermission extends string = string> = {
  modules: MenuTree<ComponentId<TRegistry>, TPermission>;
  registry: TRegistry;
  appTitle?: string;
  appSubtitle?: string;
  loginTitle?: string;
  loginSubtitle?: string;
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

const TenantLauncher = <TRegistry extends ComponentRegistry, TPermission extends string = string>({
  modules,
  title,
}: {
  modules: MenuTree<ComponentId<TRegistry>, TPermission>;
  title: string;
}) => {
  const { tenantSession } = useAuth();
  const visibleModules = filterMenuTree(modules, tenantSession?.permissions ?? []);

  return (
    <main className="min-h-screen bg-ui-bg px-8 py-8 text-ui-text">
      <div className="grid gap-8">
        <header className="flex items-center justify-between gap-6 border-b border-ui-border/25 pb-4">
          <h1 className="text-3xl font-semibold tracking-tight text-ui-text">{title}</h1>
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
                className="group rounded-2xl bg-ui-panel p-4 shadow-sm ring-1 ring-inset ring-ui-border/30 transition-transform hover:-translate-y-0.5 hover:shadow-md"
              >
                <div className="flex items-center gap-3">
                  {module.icon && (
                    <span className="grid size-10 shrink-0 place-items-center rounded-xl bg-ui-surface-selected text-ui-accent">
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

const TenantWorkspaceShell = <TRegistry extends ComponentRegistry, TPermission extends string = string>({
  modules,
  title,
  subtitle,
}: {
  modules: MenuTree<ComponentId<TRegistry>, TPermission>;
  title: string;
  subtitle: string;
}) => {
  const { tenantSession, logoutTenant } = useAuth();

  return (
    <TenantWorkspaceLayout
      title={title}
      subtitle={subtitle}
      modules={modules}
      permissions={tenantSession?.permissions ?? []}
      actions={
        <Button type="button" variant="secondary" onClick={logoutTenant}>
          Cerrar sesion
        </Button>
      }
    >
      <Outlet />
    </TenantWorkspaceLayout>
  );
};

export const TenantApplication = <TRegistry extends ComponentRegistry, TPermission extends string = string>({
  modules,
  registry,
  appTitle = "Base ERP",
  appSubtitle = "Selecciona un modulo para entrar a su espacio de trabajo.",
  loginTitle = "Tenant Login",
  loginSubtitle = "Accede al entorno principal del ERP.",
  defaultTenantCodigo,
  unauthorizedElement,
}: TenantApplicationProps<TRegistry, TPermission>) => {
  return useRoutes([
    { path: "/", element: <RootRedirect /> },
    {
      element: <RequireGuest scope="tenant" />,
      children: [
        {
          path: "/login",
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
      element: <RequireTenant />,
      children: [
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
