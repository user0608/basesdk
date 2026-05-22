import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "../../components/actions/Button";
import { SelectField } from "../../components/form/SelectField";
import { Permissions } from "../../generated/permissions";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import type { RoleResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const schema = createFormSchema((validators) => ({
  values: validators.nullableStringArray(),
}));

export type SecurityRolePermissionsFormProps = {
  role: RoleResponse;
  close?: () => void;
};

export const SecurityRolePermissionsForm = ({ role, close }: SecurityRolePermissionsFormProps) => {
  const security = useSecurityService();
  const queryKeyPrefix = security.queryKeyPrefix;
  const form = useCustomForm(schema, { defaultValues: { values: [] } });
  const { setValue } = form;
  const permissions = useQuery({ queryKey: [...queryKeyPrefix, "permissions"], queryFn: ({ signal }) => security.permissions.list(signal) });
  const assigned = useQuery({
    queryKey: [...queryKeyPrefix, "roles", role.code, "permissions"],
    queryFn: ({ signal }) => security.roles.permissions.list(role.code, signal),
  });
  const save = useMutate({
    mutationFn: (values: string[]) => security.roles.permissions.replace(role.code, values),
    invalidateQueryKey: [[...queryKeyPrefix, "roles", role.code, "permissions"], [...queryKeyPrefix, "roles"]],
    successMessage: "Permisos del rol actualizados.",
  });

  useEffect(() => {
    if (!assigned.data) return;
    setValue("values", assigned.data.map((permission) => permission.code));
  }, [assigned.data, setValue]);

  const options = (permissions.data ?? []).map((permission) => ({
    value: permission.code,
    label: permission.description ?? permission.code,
  }));

  return (
    <form
      className="grid gap-3"
      onSubmit={form.handleSubmit(async (values) => {
        await save.mutate(values.values ?? []);
        close?.();
      })}
    >
      <SelectField
        form={form}
        name="values"
        variant="table"
        multiple
        label="Permisos"
        info="Selecciona los permisos que concede este rol."
        options={options}
        loading={permissions.isLoading || assigned.isLoading}
        emptyMessage="No hay permisos disponibles"
        pageSize={8}
      />

      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button type="submit" loading={save.isPending} permissions={[Permissions.securityRolesPermissionsReplace]}>
          Guardar
        </Button>
      </div>
    </form>
  );
};
