import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import { useSecurityService } from "../../security/useSecurityService";
import type { TenantResponse } from "../../security/types";

const schema = createFormSchema((validators) => ({
  codigo: validators.requiredString(),
  name: validators.requiredString(),
  timezone: validators.requiredString(),
  maxActiveUsers: validators.requiredPositiveInteger(),
  disabled: validators.requiredBoolean(),
  expiresAt: validators.nullableString(),
}));

export type SystemTenantFormProps = {
  tenant?: TenantResponse;
  close?: () => void;
};

type TenantFormValues = {
  codigo: string;
  name: string;
  timezone: string;
  maxActiveUsers: number;
  disabled: boolean;
  expiresAt: string | null;
};

export const SystemTenantForm = ({ tenant, close }: SystemTenantFormProps) => {
  const security = useSecurityService();
  const form = useCustomForm(schema, {
    defaultValues: {
      codigo: tenant?.codigo ?? "",
      name: tenant?.name ?? "",
      timezone: tenant?.timezone ?? "America/Lima",
      maxActiveUsers: tenant?.maxActiveUsers ?? 1,
      disabled: tenant?.disabled ?? false,
      expiresAt: tenant?.expiresAt ? tenant.expiresAt.slice(0, 16) : null,
    },
  });
  const isEdit = Boolean(tenant);
  const save = useMutate({
    mutationFn: (values: TenantFormValues) => {
      const expiresAt = values.expiresAt ? new Date(values.expiresAt).toISOString() : null;
      if (isEdit && tenant) {
        return security.tenants.update(tenant.codigo, {
          name: values.name,
          timezone: values.timezone,
          maxActiveUsers: values.maxActiveUsers,
          disabled: values.disabled,
          expiresAt,
        });
      }

      return security.tenants.create({
        codigo: values.codigo,
        name: values.name,
        timezone: values.timezone,
        maxActiveUsers: values.maxActiveUsers,
        expiresAt,
      });
    },
    invalidateQueryKey: ["system", "tenants"],
    successMessage: isEdit ? "Tenant actualizado." : "Tenant creado.",
  });

  return (
    <form
      className="grid gap-3"
      onSubmit={form.handleSubmit(async (values) => {
        await save.mutate(values);
        close?.();
      })}
    >
      <InputField form={form} name="codigo" label="Codigo" readOnly={isEdit} />
      <InputField form={form} name="name" label="Nombre" />
      <InputField form={form} name="timezone" label="Zona horaria" placeholder="America/Lima" />
      <InputField form={form} name="maxActiveUsers" label="Usuarios activos maximos" type="number" />
      <InputField form={form} name="expiresAt" label="Expira en" type="datetime-local" />
      {isEdit && <InputField form={form} name="disabled" label="Deshabilitado" variant="checkbox" />}

      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button type="submit" loading={save.isPending}>
          Guardar
        </Button>
      </div>
    </form>
  );
};
