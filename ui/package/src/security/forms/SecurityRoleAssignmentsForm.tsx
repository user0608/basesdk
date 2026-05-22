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

export type SecurityRoleAssignmentsFormProps = {
  role: RoleResponse;
  mode: "users" | "groups";
  close?: () => void;
};

export const SecurityRoleAssignmentsForm = ({ role, mode, close }: SecurityRoleAssignmentsFormProps) => {
  const security = useSecurityService();
  const queryKeyPrefix = security.queryKeyPrefix;
  const form = useCustomForm(schema, { defaultValues: { values: [] } });
  const { setValue } = form;
  const isUsersMode = mode === "users";
  const allItems = useQuery({
    queryKey: [...queryKeyPrefix, isUsersMode ? "users" : "groups"],
    queryFn: ({ signal }) => (isUsersMode ? security.users.list(signal) : security.groups.list(signal)),
  });
  const assignedItems = useQuery({
    queryKey: [...queryKeyPrefix, "roles", role.code, mode],
    queryFn: ({ signal }) => (isUsersMode ? security.roles.users.list(role.code, signal) : security.roles.groups.list(role.code, signal)),
  });
  const save = useMutate({
    mutationFn: (values: string[]) =>
      isUsersMode ? security.roles.users.replace(role.code, values) : security.roles.groups.replace(role.code, values),
    invalidateQueryKey: [
      [...queryKeyPrefix, "roles", role.code, mode],
      [...queryKeyPrefix, "roles"],
    ],
    successMessage: isUsersMode ? "Usuarios del rol actualizados." : "Grupos del rol actualizados.",
  });

  useEffect(() => {
    if (!assignedItems.data) return;
    setValue("values", isUsersMode ? assignedItems.data.map((item) => item.username) : assignedItems.data.map((item) => item.code));
  }, [assignedItems.data, isUsersMode, setValue]);

  const options = (allItems.data ?? []).map((item) => {
    if ("username" in item) {
      return { value: item.username, label: item.fullName ? `${item.username} - ${item.fullName}` : item.username };
    }

    return { value: item.code, label: item.description ?? item.code };
  });

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
        label={isUsersMode ? "Usuarios" : "Grupos"}
        info={isUsersMode ? "Selecciona usuarios con este rol directo." : "Selecciona grupos que reciben este rol."}
        options={options}
        loading={allItems.isLoading || assignedItems.isLoading}
        emptyMessage={isUsersMode ? "No hay usuarios disponibles" : "No hay grupos disponibles"}
        pageSize={8}
      />

      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button
          type="submit"
          loading={save.isPending}
          permissions={[isUsersMode ? Permissions.securityRolesUsersReplace : Permissions.securityRolesGroupsReplace]}
        >
          Guardar
        </Button>
      </div>
    </form>
  );
};
