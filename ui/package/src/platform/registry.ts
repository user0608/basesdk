import type { ComponentRegistry, MenuTree } from "./types";

export const defineComponentRegistry = <TRegistry extends ComponentRegistry>(registry: TRegistry) => registry;

const trimSlashes = (value: string) => value.replace(/^\/+|\/+$/g, "");

const resolvePath = (parentPath: string, path: string) => {
  if (path === "") return parentPath;
  if (path.startsWith("/")) return `/${trimSlashes(path)}`;

  const parent = trimSlashes(parentPath);
  const child = trimSlashes(path);

  if (!parent) return `/${child}`;
  return `/${parent}/${child}`;
};

export const defineMenuTree = <TMenuTree extends MenuTree<string, string>>(menuTree: TMenuTree) => {
  return menuTree.map((module) => {
    const modulePath = resolvePath("", module.path);

    return {
      ...module,
      path: modulePath,
      children: module.children.map((menu) => {
        const menuPath = resolvePath(modulePath, menu.path);

        return {
          ...menu,
          path: menuPath,
          children: menu.children.map((operation) => ({
            ...operation,
            path: resolvePath(menuPath, operation.path),
          })),
        };
      }),
    };
  }) as unknown as TMenuTree;
};
