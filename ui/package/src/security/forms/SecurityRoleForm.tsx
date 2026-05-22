import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { Permissions } from "../../generated/permissions";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import { uniqueCode } from "../../utils/uniqueCode";
import type { RoleResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const queryKey = ["security", "roles"];

const schema = createFormSchema((validators) => ({
  code: validators.requiredString(),
  description: validators.nullableString(),
}));

export type SecurityRoleFormProps = {
  role?: RoleResponse;
  close?: () => void;
};

export const SecurityRoleForm = ({ role, close }: SecurityRoleFormProps) => {
  const security = useSecurityService();
  const isEditing = Boolean(role);
  const form = useCustomForm(schema, {
    defaultValues: {
      code: role?.code ?? uniqueCode(),
      description: role?.description ?? null,
    },
  });
  const save = useMutate({
    mutationFn: (data: { code: string; description: string | null }) => {
      if (role) return security.roles.update(role.code, { description: data.description, disabled: role.disabled });
      return security.roles.create(data);
    },
    invalidateQueryKey: queryKey,
    successMessage: isEditing ? "Rol actualizado." : "Rol creado.",
  });

  return (
    <form
      className="grid gap-3"
      onSubmit={form.handleSubmit(async (values) => {
        await save.mutate(values);
        close?.();
      })}
    >
      {isEditing && <InputField form={form} name="code" label="Codigo" readOnly />}
      <InputField form={form} name="description" label="Descripcion" />
      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button
          type="submit"
          loading={save.isPending}
          permissions={[isEditing ? Permissions.securityRolesUpdate : Permissions.securityRolesCreate]}
        >
          Guardar
        </Button>
      </div>
    </form>
  );
};
