import { useAuth } from "../auth/useAuth";
import { Permissions } from "../generated/permissions";
import { useLocation } from "react-router-dom";

export const hasAllPermissions = <TPermission extends string>(
  grantedPermissions: readonly TPermission[] | null | undefined,
  requiredPermissions?: readonly TPermission[],
) => {
  if (!requiredPermissions || requiredPermissions.length === 0) return true;
  if (!grantedPermissions || grantedPermissions.length === 0) return false;
  if (grantedPermissions.includes(Permissions.securityAdmin as TPermission)) return true;

  return requiredPermissions.every((permission) => grantedPermissions.includes(permission));
};

export const hasPermission = <TPermission extends string>(
  grantedPermissions: readonly TPermission[] | null | undefined,
  requiredPermission?: TPermission,
) => {
  if (!requiredPermission) return true;
  if (!grantedPermissions || grantedPermissions.length === 0) return false;
  if (grantedPermissions.includes(Permissions.securityAdmin as TPermission)) return true;

  return grantedPermissions.includes(requiredPermission);
};

export const hasAnyPermission = <TPermission extends string>(
  grantedPermissions: readonly TPermission[] | null | undefined,
  requiredPermissions?: readonly TPermission[],
) => {
  if (!requiredPermissions || requiredPermissions.length === 0) return true;
  if (!grantedPermissions || grantedPermissions.length === 0) return false;
  if (grantedPermissions.includes(Permissions.securityAdmin as TPermission)) return true;

  return requiredPermissions.some((permission) => grantedPermissions.includes(permission));
};

export const useTenantPermissions = () => {
  const { tenantSession } = useAuth();
  return tenantSession?.permissions ?? null;
};

export const useCurrentPermissions = () => {
  const { systemSession, tenantSession } = useAuth();
  const location = useLocation();
  if (location.pathname.startsWith("/system") && systemSession) return [Permissions.securityAdmin];
  return tenantSession?.permissions ?? null;
};

export const useHasPermission = <TPermission extends string>(requiredPermission?: TPermission) => {
  const currentPermissions = useCurrentPermissions() as readonly TPermission[] | null;
  return hasPermission(currentPermissions, requiredPermission);
};

export const useHasAllPermissions = <TPermission extends string>(requiredPermissions?: readonly TPermission[]) => {
  const currentPermissions = useCurrentPermissions() as readonly TPermission[] | null;
  return hasAllPermissions(currentPermissions, requiredPermissions);
};

export const useHasAnyPermission = <TPermission extends string>(requiredPermissions?: readonly TPermission[]) => {
  const currentPermissions = useCurrentPermissions() as readonly TPermission[] | null;
  return hasAnyPermission(currentPermissions, requiredPermissions);
};
