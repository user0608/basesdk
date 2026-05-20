import type { ReactNode } from "react";
import { hasAllPermissions, useCurrentPermissions } from "./permissions";

export const RequirePermissions = <TPermission extends string>({
  permissions,
  fallback = null,
  children,
}: {
  permissions?: readonly TPermission[];
  fallback?: ReactNode;
  children: ReactNode;
}) => {
  const currentPermissions = useCurrentPermissions() as readonly TPermission[] | null;

  if (!hasAllPermissions(currentPermissions, permissions)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
};
