import type { HttpApi } from "../api/http/types";
import type {
  CreatePropertyInput,
  CreateSystemUserInput,
  PropertyResponse,
  SystemUserResponse,
  TenantPropertyResponse,
  UpdatePropertyInput,
  UpdateSystemUserInput,
} from "./types";

export const createSystemService = (api: HttpApi) => ({
  users: {
    list: (signal?: AbortSignal) =>
      api.get<SystemUserResponse[]>({ path: "/api/v1/system/users", auth: "system", signal }),
    create: (data: CreateSystemUserInput) =>
      api.post<void>({ path: "/api/v1/system/users", auth: "system", data }),
    update: (username: string, data: UpdateSystemUserInput) =>
      api.put<void>({ path: `/api/v1/system/users/${encodeURIComponent(username)}`, auth: "system", data }),
    enable: (usernames: string[]) =>
      api.patch<void>({ path: "/api/v1/system/users/enable", auth: "system", data: { usernames } }),
    disable: (usernames: string[]) =>
      api.patch<void>({ path: "/api/v1/system/users/disable", auth: "system", data: { usernames } }),
    delete: (usernames: string[]) =>
      api.delete<void>({ path: "/api/v1/system/users", auth: "system", data: { usernames } }),
  },
  properties: {
    list: (signal?: AbortSignal) =>
      api.get<PropertyResponse[]>({ path: "/api/v1/system/properties", auth: "system", signal }),
    create: (data: CreatePropertyInput) =>
      api.post<void>({ path: "/api/v1/system/properties", auth: "system", data }),
    update: (key: string, data: UpdatePropertyInput) =>
      api.put<void>({ path: `/api/v1/system/properties/${encodeURIComponent(key)}`, auth: "system", data }),
    delete: (key: string) =>
      api.delete<void>({ path: `/api/v1/system/properties/${encodeURIComponent(key)}`, auth: "system" }),
  },
  tenantProperties: (tenantCodigo: string) => ({
    list: (signal?: AbortSignal) =>
      api.get<TenantPropertyResponse[]>({
        path: `/api/v1/system/tenants/${encodeURIComponent(tenantCodigo)}/properties`,
        auth: "system",
        signal,
      }),
    create: (data: CreatePropertyInput) =>
      api.post<void>({
        path: `/api/v1/system/tenants/${encodeURIComponent(tenantCodigo)}/properties`,
        auth: "system",
        data,
      }),
    update: (key: string, data: UpdatePropertyInput) =>
      api.put<void>({
        path: `/api/v1/system/tenants/${encodeURIComponent(tenantCodigo)}/properties/${encodeURIComponent(key)}`,
        auth: "system",
        data,
      }),
    delete: (key: string) =>
      api.delete<void>({
        path: `/api/v1/system/tenants/${encodeURIComponent(tenantCodigo)}/properties/${encodeURIComponent(key)}`,
        auth: "system",
      }),
  }),
});

export type SystemService = ReturnType<typeof createSystemService>;
