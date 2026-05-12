import type { JwtClaims, SystemSession, TenantSession } from "./types";

const decodeBase64Url = (value: string) => {
  const normalized = value.replace(/-/g, "+").replace(/_/g, "/");
  const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, "=");

  if (typeof globalThis.atob === "function") {
    return globalThis.atob(padded);
  }

  throw new Error("El entorno no soporta atob");
};

export const decodeJwtClaims = (token: string): JwtClaims => {
  const parts = token.split(".");
  if (parts.length < 2) {
    throw new Error("Token invalido");
  }

  return JSON.parse(decodeBase64Url(parts[1]!)) as JwtClaims;
};

const claimsExpiration = (claims: JwtClaims) => {
  if (!claims.exp) return null;
  return new Date(claims.exp * 1000).toISOString();
};

export const toTenantSession = (token: string, claims: JwtClaims, permissions: string[]): TenantSession => {
  if (claims.type !== "tenant" || !claims.tenant) {
    throw new Error("Token tenant invalido");
  }

  return {
    scope: "tenant",
    token,
    username: claims.sub,
    tenant: claims.tenant,
    timeZone: claims.time_zone ?? "UTC",
    expiresAt: claimsExpiration(claims),
    permissions,
  };
};

export const toSystemSession = (token: string, claims: JwtClaims): SystemSession => {
  if (claims.type !== "system") {
    throw new Error("Token system invalido");
  }

  return {
    scope: "system",
    token,
    username: claims.sub,
    timeZone: claims.time_zone ?? "UTC",
    expiresAt: claimsExpiration(claims),
    permissions: null,
  };
};
