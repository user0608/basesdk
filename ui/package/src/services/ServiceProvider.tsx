import { createContext, useContext, useMemo } from "react";
import type { PropsWithChildren } from "react";
import { useHttpApi } from "../auth/useAuth";
import { createSecurityService } from "../security/SecurityService";
import type { SecurityService } from "../security/SecurityService";
import { createSystemService } from "../system/SystemService";
import type { SystemService } from "../system/SystemService";

export type ApplicationServices = {
  security: SecurityService;
  system: SystemService;
};

const ServiceContext = createContext<ApplicationServices | null>(null);

export const ServiceProvider = ({ children }: PropsWithChildren) => {
  const api = useHttpApi();
  const services = useMemo<ApplicationServices>(
    () => ({
      security: createSecurityService(api),
      system: createSystemService(api),
    }),
    [api],
  );

  return <ServiceContext.Provider value={services}>{children}</ServiceContext.Provider>;
};

export const useServices = () => {
  const services = useContext(ServiceContext);

  if (!services) {
    throw new Error("useServices debe usarse dentro de ServiceProvider");
  }

  return services;
};
