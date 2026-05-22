import type { HttpApi } from "../api/http/types";
import type {
  CreateGroupInput,
  CreateRoleInput,
  CreateTenantInput,
  CreateTenantUserInput,
  GroupResponse,
  PermissionResponse,
  RoleResponse,
  TenantResponse,
  TenantUserResponse,
  UpdateGroupInput,
  UpdateRoleInput,
  UpdateTenantInput,
  UpdateTenantUserInput,
} from "./types";

type SecurityServiceOptions = {
  auth?: "tenant" | "system";
  resourcePrefix?: string;
  permissionsPrefix?: string;
  queryKeyPrefix?: readonly string[];
};

export const createSecurityService = (api: HttpApi, options: SecurityServiceOptions = {}) => {
  const auth = options.auth ?? "tenant";
  const resourcePrefix = options.resourcePrefix ?? "/api/v1";
  const permissionsPrefix = options.permissionsPrefix ?? `${resourcePrefix}/permissions`;
  const queryKeyPrefix = options.queryKeyPrefix ?? ["security", auth];

  return {
    queryKeyPrefix,
    tenants: {
    list: (signal?: AbortSignal) =>
      api.get<TenantResponse[]>({
        path: "/api/v1/system/tenants",
        auth: "system",
        signal,
      }),
    find: (codigo: string, signal?: AbortSignal) =>
      api.get<TenantResponse>({
        path: `/api/v1/system/tenants/${encodeURIComponent(codigo)}`,
        auth: "system",
        signal,
      }),
    create: (data: CreateTenantInput) =>
      api.post<void>({
        path: "/api/v1/system/tenants",
        auth: "system",
        data,
      }),
    update: (codigo: string, data: UpdateTenantInput) =>
      api.put<void>({
        path: `/api/v1/system/tenants/${encodeURIComponent(codigo)}`,
        auth: "system",
        data,
      }),
    enable: (codigos: string[]) =>
      api.patch<void>({
        path: "/api/v1/system/tenants/enable",
        auth: "system",
        data: { codigos },
      }),
    disable: (codigos: string[]) =>
      api.patch<void>({
        path: "/api/v1/system/tenants/disable",
        auth: "system",
        data: { codigos },
      }),
    },
    permissions: {
    list: (signal?: AbortSignal) =>
      api.get<PermissionResponse[]>({
        path: permissionsPrefix,
        auth,
        signal,
      }),
    find: (code: string, signal?: AbortSignal) =>
      api.get<PermissionResponse>({
        path: `${permissionsPrefix}/${encodeURIComponent(code)}`,
        auth,
        signal,
      }),
    roles: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<RoleResponse[]>({
          path: `${permissionsPrefix}/${encodeURIComponent(code)}/roles`,
          auth,
          signal,
        }),
    },
    groups: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<GroupResponse[]>({
          path: `${permissionsPrefix}/${encodeURIComponent(code)}/groups`,
          auth,
          signal,
        }),
    },
    users: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<TenantUserResponse[]>({
          path: `${permissionsPrefix}/${encodeURIComponent(code)}/users`,
          auth,
          signal,
        }),
    },
  },
    users: {
    list: (signal?: AbortSignal) =>
      api.get<TenantUserResponse[]>({
        path: `${resourcePrefix}/users`,
        auth,
        signal,
      }),
    create: (data: CreateTenantUserInput) =>
      api.post<void>({
        path: `${resourcePrefix}/users`,
        auth,
        data,
      }),
    update: (username: string, data: UpdateTenantUserInput) =>
      api.put<void>({
        path: `${resourcePrefix}/users/${encodeURIComponent(username)}`,
        auth,
        data,
      }),
    changePassword: (username: string, password: string) =>
      api.patch<void>({
        path: `${resourcePrefix}/users/${encodeURIComponent(username)}/password`,
        auth,
        data: { password },
      }),
    enable: (usernames: string[]) =>
      api.patch<void>({
        path: `${resourcePrefix}/users/enable`,
        auth,
        data: { usernames },
      }),
    disable: (usernames: string[]) =>
      api.patch<void>({
        path: `${resourcePrefix}/users/disable`,
        auth,
        data: { usernames },
      }),
    delete: (usernames: string[]) =>
      api.delete<void>({
        path: `${resourcePrefix}/users`,
        auth,
        data: { usernames },
      }),
    permissions: {
      list: (username: string, signal?: AbortSignal) =>
        api.get<{ tenantCodigo: string; username: string; permissions: PermissionResponse[] }>({
          path: `${resourcePrefix}/users/${encodeURIComponent(username)}/permissions`,
          auth,
          signal,
        }),
    },
    roles: {
      list: (username: string, signal?: AbortSignal) =>
        api.get<RoleResponse[]>({
          path: `${resourcePrefix}/users/${encodeURIComponent(username)}/roles`,
          auth,
          signal,
        }),
      replace: (username: string, codes: string[]) =>
        api.put<void>({
          path: `${resourcePrefix}/users/${encodeURIComponent(username)}/roles`,
          auth,
          data: { codes },
        }),
    },
    groups: {
      list: (username: string, signal?: AbortSignal) =>
        api.get<GroupResponse[]>({
          path: `${resourcePrefix}/users/${encodeURIComponent(username)}/groups`,
          auth,
          signal,
        }),
      replace: (username: string, codes: string[]) =>
        api.put<void>({
          path: `${resourcePrefix}/users/${encodeURIComponent(username)}/groups`,
          auth,
          data: { codes },
        }),
    },
  },
  roles: {
    list: (signal?: AbortSignal) =>
      api.get<RoleResponse[]>({
        path: `${resourcePrefix}/roles`,
        auth,
        signal,
      }),
    create: (data: CreateRoleInput) =>
      api.post<void>({
        path: `${resourcePrefix}/roles`,
        auth,
        data,
      }),
    update: (code: string, data: UpdateRoleInput) =>
      api.put<void>({
        path: `${resourcePrefix}/roles/${encodeURIComponent(code)}`,
        auth,
        data,
      }),
    enable: (codes: string[]) =>
      api.patch<void>({
        path: `${resourcePrefix}/roles/enable`,
        auth,
        data: { codes },
      }),
    disable: (codes: string[]) =>
      api.patch<void>({
        path: `${resourcePrefix}/roles/disable`,
        auth,
        data: { codes },
      }),
    delete: (codes: string[]) =>
      api.delete<void>({
        path: `${resourcePrefix}/roles`,
        auth,
        data: { codes },
      }),
    permissions: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<PermissionResponse[]>({
          path: `${resourcePrefix}/roles/${encodeURIComponent(code)}/permissions`,
          auth,
          signal,
        }),
      replace: (code: string, codes: string[]) =>
        api.put<void>({
          path: `${resourcePrefix}/roles/${encodeURIComponent(code)}/permissions`,
          auth,
          data: { codes },
        }),
    },
    users: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<TenantUserResponse[]>({
          path: `${resourcePrefix}/roles/${encodeURIComponent(code)}/users`,
          auth,
          signal,
        }),
      replace: (code: string, usernames: string[]) =>
        api.put<void>({
          path: `${resourcePrefix}/roles/${encodeURIComponent(code)}/users`,
          auth,
          data: { usernames },
        }),
    },
    groups: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<GroupResponse[]>({
          path: `${resourcePrefix}/roles/${encodeURIComponent(code)}/groups`,
          auth,
          signal,
        }),
      replace: (code: string, codes: string[]) =>
        api.put<void>({
          path: `${resourcePrefix}/roles/${encodeURIComponent(code)}/groups`,
          auth,
          data: { codes },
        }),
    },
  },
  groups: {
    list: (signal?: AbortSignal) =>
      api.get<GroupResponse[]>({
        path: `${resourcePrefix}/groups`,
        auth,
        signal,
      }),
    create: (data: CreateGroupInput) =>
      api.post<void>({
        path: `${resourcePrefix}/groups`,
        auth,
        data,
      }),
    update: (code: string, data: UpdateGroupInput) =>
      api.put<void>({
        path: `${resourcePrefix}/groups/${encodeURIComponent(code)}`,
        auth,
        data,
      }),
    enable: (codes: string[]) =>
      api.patch<void>({
        path: `${resourcePrefix}/groups/enable`,
        auth,
        data: { codes },
      }),
    disable: (codes: string[]) =>
      api.patch<void>({
        path: `${resourcePrefix}/groups/disable`,
        auth,
        data: { codes },
      }),
    delete: (codes: string[]) =>
      api.delete<void>({
        path: `${resourcePrefix}/groups`,
        auth,
        data: { codes },
      }),
    users: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<TenantUserResponse[]>({
          path: `${resourcePrefix}/groups/${encodeURIComponent(code)}/users`,
          auth,
          signal,
        }),
      replace: (code: string, usernames: string[]) =>
        api.put<void>({
          path: `${resourcePrefix}/groups/${encodeURIComponent(code)}/users`,
          auth,
          data: { usernames },
        }),
    },
    roles: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<RoleResponse[]>({
          path: `${resourcePrefix}/groups/${encodeURIComponent(code)}/roles`,
          auth,
          signal,
        }),
      replace: (code: string, codes: string[]) =>
        api.put<void>({
          path: `${resourcePrefix}/groups/${encodeURIComponent(code)}/roles`,
          auth,
          data: { codes },
        }),
    },
    permissions: {
      list: (code: string, signal?: AbortSignal) =>
        api.get<PermissionResponse[]>({
          path: `${resourcePrefix}/groups/${encodeURIComponent(code)}/permissions`,
          auth,
          signal,
        }),
    },
  },
  };
};

export type SecurityService = ReturnType<typeof createSecurityService>;
