import { createContext, createElement, useContext } from "react";
import type { PropsWithChildren } from "react";
import { useServices } from "../services/ServiceProvider";
import type { SecurityService } from "./SecurityService";

const SecurityServiceContext = createContext<SecurityService | null>(null);

export const SecurityServiceProvider = ({ service, children }: PropsWithChildren<{ service: SecurityService }>) => {
  return createElement(SecurityServiceContext.Provider, { value: service }, children);
};

export const useSecurityService = () => useContext(SecurityServiceContext) ?? useServices().security;
