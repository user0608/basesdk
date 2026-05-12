import type { SystemSession, TenantSession } from "./types";

const runtimeStorage = () => {
  if (typeof window === "undefined") return null;
  return window.localStorage;
};

const readValue = <T>(key: string): T | null => {
  const storage = runtimeStorage();
  if (!storage) return null;

  const raw = storage.getItem(key);
  if (!raw) return null;

  try {
    return JSON.parse(raw) as T;
  } catch {
    storage.removeItem(key);
    return null;
  }
};

const writeValue = (key: string, value: unknown) => {
  const storage = runtimeStorage();
  if (!storage) return;
  storage.setItem(key, JSON.stringify(value));
};

const removeValue = (key: string) => {
  const storage = runtimeStorage();
  if (!storage) return;
  storage.removeItem(key);
};

export const createAuthStorage = (prefix = "basesdk.ui.auth") => {
  const tenantKey = `${prefix}.tenant`;
  const systemKey = `${prefix}.system`;

  return {
    readTenantSession: () => readValue<TenantSession>(tenantKey),
    readSystemSession: () => readValue<SystemSession>(systemKey),
    writeTenantSession: (session: TenantSession) => writeValue(tenantKey, session),
    writeSystemSession: (session: SystemSession) => writeValue(systemKey, session),
    clearTenantSession: () => removeValue(tenantKey),
    clearSystemSession: () => removeValue(systemKey),
  };
};
