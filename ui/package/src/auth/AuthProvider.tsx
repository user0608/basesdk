import { createContext, useEffect, useMemo, useState } from "react";
import { createHttpApi } from "../api/http/createHttpApi";
import { createAuthStorage } from "./storage";
import { decodeJwtClaims, toSystemSession, toTenantSession } from "./jwt";
import type {
  AuthContextValue,
  AuthProviderProps,
  AuthState,
  SystemLoginInput,
  SystemSession,
  TenantLoginInput,
  TenantSession,
} from "./types";

type LoginResponse = {
  token: string;
};

type TenantPermissionsResponse = {
  tenantCodigo: string;
  username: string;
  permissions: Array<{
    code: string;
  }>;
};

type TenantMeResponse = {
  mustChangePassword: boolean;
};

export const AuthContext = createContext<AuthContextValue | null>(null);

export const AuthProvider = ({ children, getBaseUrl, storageKeyPrefix }: AuthProviderProps) => {
  const storage = useMemo(() => createAuthStorage(storageKeyPrefix), [storageKeyPrefix]);
  const [tenantSession, setTenantSession] = useState<TenantSession | null>(null);
  const [systemSession, setSystemSession] = useState<SystemSession | null>(null);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    setTenantSession(storage.readTenantSession());
    setSystemSession(storage.readSystemSession());
    setIsReady(true);
  }, [storage]);

  const logoutTenant = () => {
    storage.clearTenantSession();
    setTenantSession(null);
  };

  const logoutSystem = () => {
    storage.clearSystemSession();
    setSystemSession(null);
  };

  const completeTenantPasswordChange = () => {
    setTenantSession((current) => {
      if (!current) return current;

      const next = { ...current, mustChangePassword: false };
      storage.writeTenantSession(next);
      return next;
    });
  };

  const httpApi = useMemo(
    () =>
      createHttpApi({
        getBaseUrl,
        getTenantToken: () => tenantSession?.token ?? null,
        getSystemToken: () => systemSession?.token ?? null,
        onUnauthorized: (auth) => {
          if (auth === "tenant") {
            logoutTenant();
            return;
          }

          logoutSystem();
        },
      }),
    [getBaseUrl, systemSession?.token, tenantSession?.token],
  );

  const loginTenant = async ({ tenantCodigo, username, password }: TenantLoginInput) => {
    const { token } = await httpApi.post<LoginResponse>({
      path: `/api/v1/tenants/${encodeURIComponent(tenantCodigo)}/auth/login`,
      auth: "none",
      data: { username, password },
    });

    const permissionsResponse = await httpApi.get<TenantPermissionsResponse>({
      path: "/api/v1/me/permissions",
      auth: "none",
      tokenOverride: token,
    });

    const meResponse = await httpApi.get<TenantMeResponse>({
      path: "/api/v1/me",
      auth: "none",
      tokenOverride: token,
    });

    const session = toTenantSession(
      token,
      decodeJwtClaims(token),
      permissionsResponse.permissions.map((permission) => permission.code),
      meResponse.mustChangePassword,
    );

    storage.writeTenantSession(session);
    setTenantSession(session);
    return session;
  };

  const loginSystem = async ({ username, password }: SystemLoginInput) => {
    const { token } = await httpApi.post<LoginResponse>({
      path: "/api/v1/system/auth/login",
      auth: "none",
      data: { username, password },
    });

    const session = toSystemSession(token, decodeJwtClaims(token));
    storage.writeSystemSession(session);
    setSystemSession(session);
    return session;
  };

  const value: AuthContextValue = {
    ...( { tenantSession, systemSession, isReady } satisfies AuthState),
    httpApi,
    loginTenant,
    loginSystem,
    completeTenantPasswordChange,
    logoutTenant,
    logoutSystem,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
