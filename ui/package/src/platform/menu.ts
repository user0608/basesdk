import { hasAllPermissions } from "./permissions";
import type {
  MenuNode,
  MenuTree,
  OperationNode,
  VisibleMenuTree,
  VisibleModuleNode,
  VisibleMenuNode,
  VisibleOperationNode,
} from "./types";

const sortByOrder = <T extends { order: number }>(values: readonly T[]) => [...values].sort((a, b) => a.order - b.order);

const mergePermissions = <TPermission extends string>(...parts: Array<readonly TPermission[] | undefined>) => {
  const merged = new Set<TPermission>();

  for (const part of parts) {
    if (!part) continue;
    for (const permission of part) {
      merged.add(permission);
    }
  }

  return [...merged] as TPermission[];
};

const toVisibleOperation = <TComponentId extends string, TPermission extends string>(
  operation: OperationNode<TComponentId, TPermission>,
  inheritedPermissions: readonly TPermission[],
  grantedPermissions: readonly TPermission[] | null | undefined,
): VisibleOperationNode<TComponentId, TPermission> | null => {
  const effectivePermissions = mergePermissions(inheritedPermissions, operation.permissions);
  if (!hasAllPermissions(grantedPermissions, effectivePermissions)) return null;

  return {
    ...operation,
    effectivePermissions,
  };
};

const toVisibleMenu = <TComponentId extends string, TPermission extends string>(
  menu: MenuNode<TComponentId, TPermission>,
  inheritedPermissions: readonly TPermission[],
  grantedPermissions: readonly TPermission[] | null | undefined,
): VisibleMenuNode<TComponentId, TPermission> | null => {
  const effectivePermissions = mergePermissions(inheritedPermissions, menu.permissions);
  if (!hasAllPermissions(grantedPermissions, effectivePermissions)) return null;

  const operations = sortByOrder(menu.children)
    .map((operation) => toVisibleOperation(operation, effectivePermissions, grantedPermissions))
    .filter(Boolean) as VisibleOperationNode<TComponentId, TPermission>[];

  if (operations.length === 0) return null;

  return {
    ...menu,
    effectivePermissions,
    children: operations,
  };
};

export const filterMenuTree = <TComponentId extends string, TPermission extends string>(
  modules: MenuTree<TComponentId, TPermission>,
  grantedPermissions: readonly TPermission[] | null | undefined,
): VisibleMenuTree<TComponentId, TPermission> => {
  return sortByOrder(modules)
    .map((module) => {
      const effectivePermissions = mergePermissions(module.permissions);
      if (!hasAllPermissions(grantedPermissions, effectivePermissions)) return null;

      const menus = sortByOrder(module.children)
        .map((menu) => toVisibleMenu(menu, effectivePermissions, grantedPermissions))
        .filter(Boolean) as VisibleMenuNode<TComponentId, TPermission>[];

      if (menus.length === 0) return null;

      const visibleModule: VisibleModuleNode<TComponentId, TPermission> = {
        ...module,
        effectivePermissions,
        children: menus,
      };

      return visibleModule;
    })
    .filter(Boolean) as VisibleMenuTree<TComponentId, TPermission>;
};

const isPathActive = (currentPath: string, candidatePath: string) => {
  return currentPath === candidatePath || currentPath.startsWith(`${candidatePath}/`);
};

export const findActiveModule = <TComponentId extends string, TPermission extends string>(
  modules: VisibleMenuTree<TComponentId, TPermission>,
  currentPath: string,
) => {
  return modules.find((module) => isPathActive(currentPath, module.path)) ?? modules[0] ?? null;
};

export const findActiveMenu = <TComponentId extends string, TPermission extends string>(
  module: VisibleModuleNode<TComponentId, TPermission> | null,
  currentPath: string,
) => {
  if (!module) return null;
  return module.children.find((menu) => isPathActive(currentPath, menu.path)) ?? module.children[0] ?? null;
};
