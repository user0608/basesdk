import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { filterMenuTree, findActiveMenu, findActiveModule } from "./menu";
import type { MenuIcon, WorkspaceLayoutProps } from "./types";

const renderIcon = (icon: MenuIcon | undefined, className: string) => {
  if (!icon) return null;

  if (typeof icon === "string") {
    return <img src={icon} alt="" className={className} />;
  }

  const Icon = icon;
  return <Icon className={className} />;
};

export const WorkspaceLayout = <TComponentId extends string, TPermission extends string = string>({
  modules,
  permissions,
  title,
  subtitle,
  homePath = "/app",
  breadcrumbs,
  actions,
  children,
}: WorkspaceLayoutProps<TComponentId, TPermission>) => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const location = useLocation();
  const visibleModules = filterMenuTree(modules, permissions);
  const activeModule = findActiveModule(visibleModules, location.pathname);
  const activeMenu = findActiveMenu(activeModule, location.pathname);
  const collapseMenuColumn = activeModule?.children.length === 1;

  return (
    <main className="h-screen overflow-hidden bg-ui-bg text-ui-text">
      <div className="flex h-full min-h-0 flex-col">
        <header className="shrink-0 border-b border-ui-border/35 bg-ui-panel shadow-sm">
          <div className="flex min-h-16 items-center justify-between gap-3 px-3 lg:px-6">
            <div className="flex min-w-0 items-center gap-3">
              {activeModule && (
                <button
                  type="button"
                  className="rounded-lg border border-ui-border/40 px-2.5 py-1.5 text-sm font-medium text-ui-text-muted transition-colors hover:bg-ui-surface-hover hover:text-ui-text lg:hidden"
                  onClick={() => setIsMobileMenuOpen((value) => !value)}
                >
                  Menu
                </button>
              )}

              <div className="min-w-0">
                <div className="flex min-w-0 items-center gap-3">
                  {title && (
                    <Link
                      to={homePath}
                      className="truncate text-xl font-semibold tracking-tight text-ui-text transition-colors hover:text-ui-primary"
                      onClick={() => setIsMobileMenuOpen(false)}
                    >
                      {title}
                    </Link>
                  )}
                  {subtitle && <span className="hidden text-sm text-ui-text-soft xl:inline">{subtitle}</span>}
                </div>
                {breadcrumbs && <div className="mt-0.5 hidden min-w-0 sm:block">{breadcrumbs}</div>}
              </div>
            </div>

            {actions && <div className="flex items-center gap-2">{actions}</div>}
          </div>
        </header>

        {isMobileMenuOpen && activeModule && (
          <div className="fixed inset-x-0 top-16 z-40 border-b border-ui-border/30 bg-ui-panel/95 p-2 shadow-lg backdrop-blur lg:hidden">
            <div className="grid gap-2">
              {!collapseMenuColumn && (
                <nav className="flex gap-1 overflow-x-auto">
                  {activeModule.children.map((menu) => {
                    const isActive = activeMenu?.id === menu.id;
                    const fallbackIcon = menu.label.trim().slice(0, 2).toUpperCase();

                    return (
                      <Link
                        key={menu.id}
                        to={menu.path}
                        className={[
                          "group flex min-h-12 min-w-16 flex-col items-center justify-center gap-1 rounded-xl px-1.5 py-1.5 text-center transition-colors",
                          isActive
                            ? "bg-ui-surface-selected text-ui-text shadow-sm ring-1 ring-inset ring-ui-accent/15"
                            : "text-ui-text-muted hover:bg-ui-surface-hover hover:text-ui-text",
                        ].join(" ")}
                        onClick={() => setIsMobileMenuOpen(false)}
                      >
                        {menu.icon ? (
                          <span className="grid size-6 place-items-center">{renderIcon(menu.icon, "size-5 object-contain")}</span>
                        ) : (
                          <span className="grid size-6 place-items-center rounded-lg bg-ui-surface-hover text-[10px] font-bold tracking-tight text-ui-text-soft group-hover:text-ui-text">
                            {fallbackIcon}
                          </span>
                        )}
                        <span className="line-clamp-2 text-[11px] font-medium leading-3">{menu.label}</span>
                      </Link>
                    );
                  })}
                </nav>
              )}

              {activeMenu && (
                <nav className={`${collapseMenuColumn ? "" : "border-t border-ui-border/20 pt-2"} flex gap-1 overflow-x-auto`}>
                  {activeMenu.children.map((operation) => {
                    const isActive = location.pathname === operation.path;

                    return (
                      <Link
                        key={operation.id}
                        to={operation.path}
                        className={[
                          "flex shrink-0 items-center gap-2 rounded-lg px-2.5 py-1.5 text-sm font-medium transition-colors",
                          isActive
                            ? "bg-ui-primary/10 text-ui-primary ring-1 ring-inset ring-ui-primary/15"
                            : "text-ui-text-muted hover:bg-ui-surface-hover hover:text-ui-text",
                        ].join(" ")}
                        onClick={() => setIsMobileMenuOpen(false)}
                      >
                        {renderIcon(operation.icon, "size-4 shrink-0 object-contain")}
                        <span className="truncate">{operation.label}</span>
                      </Link>
                    );
                  })}
                </nav>
              )}
            </div>
          </div>
        )}

        <div className={`grid min-h-0 flex-1 ${collapseMenuColumn ? "lg:grid-cols-[190px_minmax(0,1fr)]" : "lg:grid-cols-[76px_190px_minmax(0,1fr)]"}`}>
          {!collapseMenuColumn && (
            <aside className="hidden overflow-y-auto border-r border-ui-border/30 bg-ui-panel px-2 py-2 lg:block">
              {activeModule ? (
              <nav className="grid gap-1">
                {activeModule.children.map((menu) => {
                  const isActive = activeMenu?.id === menu.id;
                  const fallbackIcon = menu.label.trim().slice(0, 2).toUpperCase();

                  return (
                    <Link
                      key={menu.id}
                      to={menu.path}
                      className={[
                        "group flex min-h-12 min-w-16 flex-col items-center justify-center gap-1 rounded-xl px-1.5 py-1.5 text-center transition-colors lg:min-h-14 lg:min-w-0 lg:py-2",
                        isActive
                          ? "bg-ui-surface-selected text-ui-text shadow-sm ring-1 ring-inset ring-ui-accent/15"
                          : "text-ui-text-muted hover:bg-ui-surface-hover hover:text-ui-text",
                      ].join(" ")}
                    >
                      {menu.icon ? (
                        <span className="grid size-6 place-items-center">{renderIcon(menu.icon, "size-5 object-contain")}</span>
                      ) : (
                        <span className="grid size-6 place-items-center rounded-lg bg-ui-surface-hover text-[10px] font-bold tracking-tight text-ui-text-soft group-hover:text-ui-text">
                          {fallbackIcon}
                        </span>
                      )}
                      <span className="line-clamp-2 text-[11px] font-medium leading-3">{menu.label}</span>
                    </Link>
                  );
                })}
              </nav>
              ) : null}
            </aside>
          )}

          <aside className="hidden overflow-y-auto border-r border-ui-border/30 bg-ui-panel px-2.5 py-2 lg:block">
            {activeMenu ? (
              <nav className="grid gap-1">
                {activeMenu.children.map((operation) => {
                  const isActive = location.pathname === operation.path;

                  return (
                    <Link
                      key={operation.id}
                      to={operation.path}
                      className={[
                        "flex shrink-0 items-center gap-2 rounded-lg px-2.5 py-1.5 text-sm font-medium transition-colors lg:shrink lg:py-2",
                        isActive
                          ? "bg-ui-primary/10 text-ui-primary ring-1 ring-inset ring-ui-primary/15"
                          : "text-ui-text-muted hover:bg-ui-surface-hover hover:text-ui-text",
                      ].join(" ")}
                    >
                      {renderIcon(operation.icon, "size-4 shrink-0 object-contain")}
                      <span className="truncate">{operation.label}</span>
                    </Link>
                  );
                })}
              </nav>
            ) : null}
          </aside>

          <section className="flex min-h-0 flex-col overflow-hidden bg-ui-panel-muted">
            <div className="min-h-0 flex-1 overflow-hidden px-3 py-3 lg:px-4 lg:py-4">{children}</div>
          </section>
        </div>
      </div>
    </main>
  );
};

export const TenantWorkspaceLayout = WorkspaceLayout;
