import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "../../components/actions/Button";
import { SelectField } from "../../components/form/SelectField";
import { Permissions } from "../../generated/permissions";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import type { TenantUserResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const schema = createFormSchema((validators) => ({
  values: validators.nullableStringArray(),
}));

export type SecurityUserAssignmentsFormProps = {
  user: TenantUserResponse;
  mode: "roles" | "groups";
  close?: () => void;
};

export const SecurityUserAssignmentsForm = ({ user, mode, close }: SecurityUserAssignmentsFormProps) => {
  const security = useSecurityService();
  const queryKeyPrefix = security.queryKeyPrefix;
  const form = useCustomForm(schema, { defaultValues: { values: [] } });
  const { setValue } = form;
  const isRolesMode = mode === "roles";
  const allItems = useQuery({
    queryKey: [...queryKeyPrefix, isRolesMode ? "roles" : "groups"],
    queryFn: ({ signal }) => (isRolesMode ? security.roles.list(signal) : security.groups.list(signal)),
  });
  const assignedItems = useQuery({
    queryKey: [...queryKeyPrefix, "users", user.username, mode],
    queryFn: ({ signal }) => (isRolesMode ? security.users.roles.list(user.username, signal) : security.users.groups.list(user.username, signal)),
  });
  const save = useMutate({
    mutationFn: (values: string[]) =>
      isRolesMode ? security.users.roles.replace(user.username, values) : security.users.groups.replace(user.username, values),
    invalidateQueryKey: [
      [...queryKeyPrefix, "users", user.username, mode],
      [...queryKeyPrefix, "users", user.username, "permissions"],
      [...queryKeyPrefix, "users"],
    ],
    successMessage: isRolesMode ? "Roles del usuario actualizados." : "Grupos del usuario actualizados.",
  });

  useEffect(() => {
    if (!assignedItems.data) return;
    setValue("values", assignedItems.data.map((item) => item.code));
  }, [assignedItems.data, setValue]);

  const options = (allItems.data ?? []).map((item) => ({
    value: item.code,
    label: item.description ?? item.code,
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
        label={isRolesMode ? "Roles" : "Grupos"}
        info={isRolesMode ? "Selecciona roles directos para el usuario." : "Selecciona grupos a los que pertenece el usuario."}
        options={options}
        loading={allItems.isLoading || assignedItems.isLoading}
        emptyMessage={isRolesMode ? "No hay roles disponibles" : "No hay grupos disponibles"}
        pageSize={8}
      />

      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button
          type="submit"
          loading={save.isPending}
          permissions={[isRolesMode ? Permissions.securityUsersRolesReplace : Permissions.securityUsersGroupsReplace]}
        >
          Guardar
        </Button>
      </div>
    </form>
  );
};
