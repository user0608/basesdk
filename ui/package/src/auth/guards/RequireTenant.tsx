import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../useAuth";

export const RequireTenant = ({ redirectTo = "/login" }: { redirectTo?: string }) => {
  const { isReady, tenantSession } = useAuth();
  const location = useLocation();

  if (!isReady) return null;
  if (!tenantSession) {
    return <Navigate to={redirectTo} replace state={{ from: location }} />;
  }

  return <Outlet />;
};
