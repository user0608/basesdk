import { createElement } from "react";
import { Navigate } from "react-router-dom";
import type { RouteObject } from "react-router-dom";
import { RequirePermissions } from "./RequirePermissions";
import type { ComponentId, ComponentRegistry, MenuTree } from "./types";

const combinePermissions = <TPermission extends string>(...parts: Array<readonly TPermission[] | undefined>) => {
  const merged = new Set<TPermission>();

  for (const part of parts) {
    if (!part) continue;
    for (const permission of part) {
      merged.add(permission);
    }
  }

  return [...merged] as TPermission[];
};

export const createModuleRoutes = <TRegistry extends ComponentRegistry, TPermission extends string = string>({
  modules,
  registry,
  unauthorizedElement,
}: {
  modules: MenuTree<ComponentId<TRegistry>, TPermission>;
  registry: TRegistry;
  unauthorizedElement?: React.ReactNode;
}): RouteObject[] => {
  return modules.flatMap((module) => {
    const firstMenu = module.children[0];
    const firstOperation = firstMenu?.children[0];

    const moduleRoutes: RouteObject[] = firstOperation
      ? [{ path: module.path, element: <Navigate to={firstOperation.path} replace /> }]
      : [];

    const menuRoutes = module.children.flatMap((menu) => {
      const firstMenuOperation = menu.children[0];
      const redirectRoute = firstMenuOperation
        ? [{ path: menu.path, element: <Navigate to={firstMenuOperation.path} replace /> } satisfies RouteObject]
        : [];

      const operationRoutes = menu.children.map((operation) => {
        const component = registry[operation.componentId];
        const permissions = combinePermissions(module.permissions, menu.permissions, operation.permissions);

        return {
          path: operation.path,
          element: (
            <RequirePermissions permissions={permissions} fallback={unauthorizedElement}>
              {createElement(component)}
            </RequirePermissions>
          ),
        } satisfies RouteObject;
      });

      return [...redirectRoute, ...operationRoutes];
    });

    return [...moduleRoutes, ...menuRoutes];
  });
};

export const createTenantModuleRoutes = createModuleRoutes;
