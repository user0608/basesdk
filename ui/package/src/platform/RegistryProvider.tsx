import { createContext, useContext } from "react";
import type { PropsWithChildren } from "react";
import type { ApplicationRegistry } from "./types";

const RegistryContext = createContext<ApplicationRegistry | null>(null);

export const RegistryProvider = ({ registry, children }: PropsWithChildren<{ registry: ApplicationRegistry }>) => {
  return <RegistryContext.Provider value={registry}>{children}</RegistryContext.Provider>;
};

export const useRegistry = <TRegistry extends ApplicationRegistry = ApplicationRegistry>() => {
  const registry = useContext(RegistryContext);

  if (!registry) {
    throw new Error("useRegistry debe usarse dentro de RegistryProvider");
  }

  return registry as TRegistry;
};
