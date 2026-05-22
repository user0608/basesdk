export type PermissionResponse = {
  code: string;
  description: string | null;
  rolesCount: number;
  groupsCount: number;
  usersCount: number;
};

export type TenantResponse = {
  codigo: string;
  name: string;
  timezone: string;
  maxActiveUsers: number;
  disabled: boolean;
  expiresAt: string | null;
  usersCount: number;
  rolesCount: number;
  groupsCount: number;
  createdBy: string;
  createdAt: string;
  updatedBy: string | null;
  updatedAt: string | null;
};

export type CreateTenantInput = {
  codigo: string;
  name: string;
  timezone: string;
  maxActiveUsers: number;
  expiresAt?: string | null;
};

export type UpdateTenantInput = {
  name: string;
  timezone: string;
  maxActiveUsers: number;
  disabled: boolean;
  expiresAt?: string | null;
};

export type TenantUserResponse = {
  tenantCodigo: string;
  username: string;
  fullName: string | null;
  mustChangePassword: boolean;
  lastLoginAt: string | null;
  disabled: boolean;
  rolesCount: number;
  groupsCount: number;
  permissionsCount: number;
};

export type CreateTenantUserInput = {
  username: string;
  fullName?: string | null;
  password: string;
  mustChangePassword: boolean;
};

export type UpdateTenantUserInput = {
  fullName?: string | null;
  mustChangePassword: boolean;
  disabled: boolean;
};

export type RoleResponse = {
  tenantCodigo: string;
  code: string;
  description: string | null;
  disabled: boolean;
  usersCount: number;
  groupsCount: number;
  permissionsCount: number;
};

export type CreateRoleInput = {
  code: string;
  description?: string | null;
};

export type UpdateRoleInput = {
  description?: string | null;
  disabled: boolean;
};

export type GroupResponse = {
  tenantCodigo: string;
  code: string;
  description: string | null;
  disabled: boolean;
  usersCount: number;
  rolesCount: number;
  permissionsCount: number;
};

export type CreateGroupInput = {
  code: string;
  description?: string | null;
};

export type UpdateGroupInput = {
  description?: string | null;
  disabled: boolean;
};
