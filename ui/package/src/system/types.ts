export type PropertyDataType = "string" | "int" | "float" | "bool" | "json";

export type SystemUserResponse = {
  username: string;
  disabled: boolean;
  createdBy: string;
  createdAt: string;
  updatedBy: string | null;
  updatedAt: string | null;
};

export type CreateSystemUserInput = {
  username: string;
  password: string;
};

export type UpdateSystemUserInput = {
  disabled: boolean;
};

export type PropertyResponse = {
  key: string;
  value: string;
  dataType: PropertyDataType;
  description: string | null;
};

export type TenantPropertyResponse = PropertyResponse & {
  tenantCodigo: string;
};

export type CreatePropertyInput = {
  key: string;
  value: string;
  dataType: PropertyDataType;
  description?: string | null;
};

export type UpdatePropertyInput = {
  key: string;
  value: string;
  dataType: PropertyDataType;
  description?: string | null;
};
