import type { ReactNode } from "react";
import type { HttpApi } from "../api/http/types";

export type TenantSession = {
  scope: "tenant";
  token: string;
  username: string;
  tenant: string;
  timeZone: string;
  expiresAt: string | null;
  permissions: string[];
  mustChangePassword: boolean;
};

export type SystemSession = {
  scope: "system";
  token: string;
  username: string;
  timeZone: string;
  expiresAt: string | null;
  permissions: string[] | null;
};

export type AuthState = {
  tenantSession: TenantSession | null;
  systemSession: SystemSession | null;
  isReady: boolean;
};

export type TenantLoginInput = {
  tenantCodigo: string;
  username: string;
  password: string;
};

export type SystemLoginInput = {
  username: string;
  password: string;
};

export type AuthContextValue = AuthState & {
  httpApi: HttpApi;
  loginTenant: (input: TenantLoginInput) => Promise<TenantSession>;
  loginSystem: (input: SystemLoginInput) => Promise<SystemSession>;
  completeTenantPasswordChange: () => void;
  logoutTenant: () => void;
  logoutSystem: () => void;
};

export type AuthProviderProps = {
  children: ReactNode;
  getBaseUrl: () => string | Promise<string>;
  storageKeyPrefix?: string;
};

export type JwtClaims = {
  type: "tenant" | "system";
  tenant?: string;
  time_zone?: string;
  sub: string;
  exp?: number;
};
