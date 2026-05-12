import { Link, Outlet, Route, Routes } from "react-router-dom";
import {
  Button,
  RequireGuest,
  RequireSystem,
  RequireTenant,
  SystemLoginPage,
  TenantLoginPage,
  useAuth,
} from "@basesdk/ui";

const PageShell = () => {
  return (
    <main className="min-h-screen bg-ui-bg px-6 py-16 text-ui-text">
      <div className="mx-auto flex max-w-6xl flex-col gap-8 rounded-3xl bg-ui-panel p-8 shadow-2xl shadow-black/10 ring-1 ring-inset ring-ui-border/60">
        <header className="flex flex-wrap items-start justify-between gap-4 border-b border-ui-border/60 pb-6">
          <div className="space-y-3">
            <span className="inline-flex rounded-full bg-ui-surface-selected px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-ui-accent ring-1 ring-inset ring-ui-accent/15">
              UI package playground
            </span>
            <h1 className="text-4xl font-semibold tracking-tight">Auth + Routing demo</h1>
            <p className="max-w-3xl text-sm leading-6 text-ui-text-muted">
              La libreria UI gestiona `tenant login`, `system login`, sesiones separadas, guards y
              `HttpApi` compartido.
            </p>
          </div>

          <nav className="flex flex-wrap gap-2">
            <Link to="/">
              <Button variant="secondary" type="button">
                Inicio
              </Button>
            </Link>
            <Link to="/login">
              <Button variant="secondary" type="button">
                Tenant login
              </Button>
            </Link>
            <Link to="/system/login">
              <Button variant="secondary" type="button">
                System login
              </Button>
            </Link>
          </nav>
        </header>

        <Outlet />
      </div>
    </main>
  );
};

const HomePage = () => {
  const { tenantSession, systemSession, isReady } = useAuth();

  return (
    <section className="grid gap-8 lg:grid-cols-[minmax(0,1.2fr)_minmax(320px,0.8fr)]">
      <div className="grid gap-6 rounded-2xl bg-ui-panel-muted p-6 shadow-sm ring-1 ring-inset ring-ui-border/50">
        <div className="space-y-2">
          <h2 className="text-2xl font-semibold">Estado actual</h2>
          <p className="text-sm text-ui-text-muted">
            Este demo deja coexistir sesiones `tenant` y `system` en paralelo.
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <SessionCard
            title="Tenant session"
            active={Boolean(tenantSession)}
            body={
              tenantSession
                ? {
                    username: tenantSession.username,
                    tenant: tenantSession.tenant,
                    permissions: tenantSession.permissions.length,
                  }
                : null
            }
            to="/app"
          />
          <SessionCard
            title="System session"
            active={Boolean(systemSession)}
            body={systemSession ? { username: systemSession.username, permissions: systemSession.permissions } : null}
            to="/system"
          />
        </div>

        <div className="text-sm text-ui-text-soft">
          Auth provider listo: {isReady ? "si" : "no"}
        </div>
      </div>

      <section className="grid gap-4 rounded-2xl bg-ui-panel-muted p-6 shadow-sm ring-1 ring-inset ring-ui-border/50">
        <h2 className="text-xl font-semibold">Rutas del demo</h2>
        <ul className="grid gap-3 text-sm text-ui-text-muted">
          <li>`/login` : login tenant</li>
          <li>`/system/login` : login system</li>
          <li>`/app` : area privada tenant</li>
          <li>`/system` : area privada system</li>
        </ul>
      </section>
    </section>
  );
};

const SessionCard = ({
  title,
  active,
  body,
  to,
}: {
  title: string;
  active: boolean;
  body: Record<string, unknown> | null;
  to: string;
}) => {
  return (
    <div className="grid gap-3 rounded-2xl bg-ui-panel p-5 shadow-sm ring-1 ring-inset ring-ui-border/50">
      <div className="flex items-center justify-between gap-3">
        <h3 className="text-base font-semibold">{title}</h3>
        <span className="rounded-full bg-ui-surface-muted px-2.5 py-1 text-xs text-ui-text-soft ring-1 ring-inset ring-ui-border/50">
          {active ? "activa" : "sin sesion"}
        </span>
      </div>

      <pre className="overflow-auto rounded-xl bg-ui-bg p-4 text-xs leading-6 text-ui-text-muted ring-1 ring-inset ring-ui-border/40">
        {JSON.stringify(body, null, 2)}
      </pre>

      <Link to={to}>
        <Button type="button">Ir al area protegida</Button>
      </Link>
    </div>
  );
};

const TenantAppPage = () => {
  const { tenantSession, logoutTenant } = useAuth();

  return (
    <section className="grid gap-6 rounded-2xl bg-ui-panel-muted p-6 shadow-sm ring-1 ring-inset ring-ui-border/50">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="space-y-2">
          <h2 className="text-2xl font-semibold">Area privada tenant</h2>
          <p className="text-sm text-ui-text-muted">
            Esta ruta exige `tenantSession` y ya trae permisos al iniciar sesion.
          </p>
        </div>

        <Button type="button" variant="secondary" onClick={logoutTenant}>
          Cerrar sesion tenant
        </Button>
      </div>

      <pre className="overflow-auto rounded-xl bg-ui-panel p-4 text-xs leading-6 text-ui-text-muted ring-1 ring-inset ring-ui-border/40">
        {JSON.stringify(tenantSession, null, 2)}
      </pre>
    </section>
  );
};

const SystemAppPage = () => {
  const { systemSession, logoutSystem } = useAuth();

  return (
    <section className="grid gap-6 rounded-2xl bg-ui-panel-muted p-6 shadow-sm ring-1 ring-inset ring-ui-border/50">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="space-y-2">
          <h2 className="text-2xl font-semibold">Area privada system</h2>
          <p className="text-sm text-ui-text-muted">
            Esta ruta exige `systemSession`. Los permisos quedan `null` hasta que exista el endpoint backend.
          </p>
        </div>

        <Button type="button" variant="secondary" onClick={logoutSystem}>
          Cerrar sesion system
        </Button>
      </div>

      <pre className="overflow-auto rounded-xl bg-ui-panel p-4 text-xs leading-6 text-ui-text-muted ring-1 ring-inset ring-ui-border/40">
        {JSON.stringify(systemSession, null, 2)}
      </pre>
    </section>
  );
};

export default function App() {
  return (
    <Routes>
      <Route element={<PageShell />}>
        <Route index element={<HomePage />} />

        <Route element={<RequireGuest scope="tenant" />}>
          <Route path="/login" element={<TenantLoginPage defaultTenantCodigo="tenant_default" />} />
        </Route>

        <Route element={<RequireGuest scope="system" />}>
          <Route path="/system/login" element={<SystemLoginPage />} />
        </Route>

        <Route element={<RequireTenant />}>
          <Route path="/app" element={<TenantAppPage />} />
        </Route>

        <Route element={<RequireSystem />}>
          <Route path="/system" element={<SystemAppPage />} />
        </Route>
      </Route>
    </Routes>
  );
}
