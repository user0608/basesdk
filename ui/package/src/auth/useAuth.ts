import { useContext } from "react";
import { AuthContext } from "./AuthProvider";

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe usarse dentro de AuthProvider");
  }

  return context;
};

export const useHttpApi = () => useAuth().httpApi;
export const useTenantSession = () => useAuth().tenantSession;
export const useSystemSession = () => useAuth().systemSession;
