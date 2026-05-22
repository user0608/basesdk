import type { ComponentType, ReactNode } from "react";
import type { PermissionCode } from "../generated/permissions";

export type MenuIcon = ComponentType<{ className?: string }> | string;

export type ComponentRegistry = Record<string, ComponentType<any>>;

export type ApplicationRegistry = {
  pages: ComponentRegistry;
  forms?: ComponentRegistry;
};

export type ComponentId<TRegistry extends ComponentRegistry> = keyof TRegistry & string;
export type PageId<TRegistry extends ApplicationRegistry> = keyof TRegistry["pages"] & string;
export type FormId<TRegistry extends ApplicationRegistry> = keyof NonNullable<TRegistry["forms"]> & string;

export type OperationNode<TComponentId extends string, TPermission extends string = PermissionCode> = {
  id: string;
  label: string;
  path: string;
  componentId: TComponentId;
  permissions?: readonly TPermission[];
  order: number;
  icon?: MenuIcon;
};

export type MenuNode<TComponentId extends string, TPermission extends string = PermissionCode> = {
  id: string;
  label: string;
  path: string;
  permissions?: readonly TPermission[];
  order: number;
  icon?: MenuIcon;
  children: readonly OperationNode<TComponentId, TPermission>[];
};

export type ModuleNode<TComponentId extends string, TPermission extends string = PermissionCode> = {
  id: string;
  label: string;
  description?: string;
  path: string;
  permissions?: readonly TPermission[];
  order: number;
  icon?: MenuIcon;
  children: readonly MenuNode<TComponentId, TPermission>[];
};

export type MenuTree<TComponentId extends string, TPermission extends string = PermissionCode> = readonly ModuleNode<
  TComponentId,
  TPermission
>[];

export type VisibleOperationNode<TComponentId extends string, TPermission extends string = PermissionCode> =
  OperationNode<TComponentId, TPermission> & {
    effectivePermissions: readonly TPermission[];
  };

export type VisibleMenuNode<TComponentId extends string, TPermission extends string = PermissionCode> = Omit<
  MenuNode<TComponentId, TPermission>,
  "children"
> & {
  effectivePermissions: readonly TPermission[];
  children: VisibleOperationNode<TComponentId, TPermission>[];
};

export type VisibleModuleNode<TComponentId extends string, TPermission extends string = PermissionCode> = Omit<
  ModuleNode<TComponentId, TPermission>,
  "children"
> & {
  effectivePermissions: readonly TPermission[];
  children: VisibleMenuNode<TComponentId, TPermission>[];
};

export type VisibleMenuTree<TComponentId extends string, TPermission extends string = PermissionCode> =
  VisibleModuleNode<TComponentId, TPermission>[];

export type WorkspaceLayoutProps<TComponentId extends string, TPermission extends string = PermissionCode> = {
  modules: MenuTree<TComponentId, TPermission>;
  permissions?: readonly TPermission[] | null;
  title?: string;
  subtitle?: string;
  homePath?: string;
  breadcrumbs?: ReactNode;
  actions?: ReactNode;
  children: ReactNode;
};
