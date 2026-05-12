import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "../useAuth";

export const RequireSystem = ({ redirectTo = "/system/login" }: { redirectTo?: string }) => {
  const { isReady, systemSession } = useAuth();
  const location = useLocation();

  if (!isReady) return null;
  if (!systemSession) {
    return <Navigate to={redirectTo} replace state={{ from: location }} />;
  }

  return <Outlet />;
};
