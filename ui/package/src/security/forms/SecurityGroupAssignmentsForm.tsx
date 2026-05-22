import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "../../components/actions/Button";
import { SelectField } from "../../components/form/SelectField";
import { Permissions } from "../../generated/permissions";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import type { GroupResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const schema = createFormSchema((validators) => ({
  values: validators.nullableStringArray(),
}));

export type SecurityGroupAssignmentsFormProps = {
  group: GroupResponse;
  mode: "users" | "roles";
  close?: () => void;
};

export const SecurityGroupAssignmentsForm = ({ group, mode, close }: SecurityGroupAssignmentsFormProps) => {
  const security = useSecurityService();
  const queryKeyPrefix = security.queryKeyPrefix;
  const form = useCustomForm(schema, {
    defaultValues: {
      values: [],
    },
  });
  const { setValue } = form;
  const isUsersMode = mode === "users";
  const allItems = useQuery({
    queryKey: [...queryKeyPrefix, isUsersMode ? "users" : "roles"],
    queryFn: ({ signal }) => (isUsersMode ? security.users.list(signal) : security.roles.list(signal)),
  });
  const assignedItems = useQuery({
    queryKey: [...queryKeyPrefix, "groups", group.code, mode],
    queryFn: ({ signal }) => (isUsersMode ? security.groups.users.list(group.code, signal) : security.groups.roles.list(group.code, signal)),
  });
  const save = useMutate({
    mutationFn: (values: string[]) =>
      isUsersMode ? security.groups.users.replace(group.code, values) : security.groups.roles.replace(group.code, values),
    invalidateQueryKey: [
      [...queryKeyPrefix, "groups", group.code, mode],
      [...queryKeyPrefix, "groups"],
    ],
    successMessage: isUsersMode ? "Usuarios del grupo actualizados." : "Roles del grupo actualizados.",
  });

  useEffect(() => {
    if (!assignedItems.data) return;

    setValue(
      "values",
      isUsersMode ? assignedItems.data.map((item) => item.username) : assignedItems.data.map((item) => item.code),
    );
  }, [assignedItems.data, isUsersMode, setValue]);

  const options = (allItems.data ?? []).map((item) => {
    if ("username" in item) {
      return {
        value: item.username,
        label: item.fullName ? `${item.username} - ${item.fullName}` : item.username,
      };
    }

    return {
      value: item.code,
      label: item.description ?? item.code,
    };
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
        label={isUsersMode ? "Usuarios" : "Roles"}
        info={isUsersMode ? "Selecciona los usuarios que pertenecen al grupo." : "Selecciona los roles asignados al grupo."}
        options={options}
        loading={allItems.isLoading || assignedItems.isLoading}
        emptyMessage={isUsersMode ? "No hay usuarios disponibles" : "No hay roles disponibles"}
        pageSize={8}
      />

      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button
          type="submit"
          loading={save.isPending}
          permissions={[isUsersMode ? Permissions.securityGroupsUsersReplace : Permissions.securityGroupsRolesReplace]}
        >
          Guardar
        </Button>
      </div>
    </form>
  );
};
