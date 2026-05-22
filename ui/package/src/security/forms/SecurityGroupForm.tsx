import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { Permissions } from "../../generated/permissions";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import { uniqueCode } from "../../utils/uniqueCode";
import type { GroupResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const queryKey = ["security", "groups"];

const schema = createFormSchema((validators) => ({
  code: validators.requiredString(),
  description: validators.nullableString(),
}));

export type SecurityGroupFormProps = {
  group?: GroupResponse;
  close?: () => void;
};

export const SecurityGroupForm = ({ group, close }: SecurityGroupFormProps) => {
  const security = useSecurityService();
  const isEditing = Boolean(group);
  const form = useCustomForm(schema, {
    defaultValues: {
      code: group?.code ?? uniqueCode(),
      description: group?.description ?? null,
    },
  });
  const save = useMutate({
    mutationFn: (data: { code: string; description: string | null }) => {
      if (group) return security.groups.update(group.code, { description: data.description, disabled: group.disabled });
      return security.groups.create(data);
    },
    invalidateQueryKey: queryKey,
    successMessage: isEditing ? "Grupo actualizado." : "Grupo creado.",
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
          permissions={[isEditing ? Permissions.securityGroupsUpdate : Permissions.securityGroupsCreate]}
        >
          Guardar
        </Button>
      </div>
    </form>
  );
};
