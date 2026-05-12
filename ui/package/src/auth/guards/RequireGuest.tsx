import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../useAuth";

export const RequireGuest = ({
  scope = "tenant",
  redirectTo,
}: {
  scope?: "tenant" | "system";
  redirectTo?: string;
}) => {
  const { isReady, tenantSession, systemSession } = useAuth();

  if (!isReady) return null;

  if (scope === "system" && systemSession) {
    return <Navigate to={redirectTo ?? "/system"} replace />;
  }

  if (scope === "tenant" && tenantSession) {
    return <Navigate to={redirectTo ?? "/app"} replace />;
  }

  return <Outlet />;
};
