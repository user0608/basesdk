import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { Permissions } from "../../generated/permissions";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import type { TenantUserResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const queryKey = ["security", "users"];

const createSchema = createFormSchema((validators) => ({
  username: validators.requiredString(),
  fullName: validators.nullableString(),
  password: validators.requiredString(),
  mustChangePassword: validators.requiredBoolean(),
}));

const updateSchema = createFormSchema((validators) => ({
  username: validators.requiredString(),
  fullName: validators.nullableString(),
  password: validators.nullableString(),
  mustChangePassword: validators.requiredBoolean(),
}));

export type SecurityUserFormProps = {
  user?: TenantUserResponse;
  close?: () => void;
};

export const SecurityUserForm = ({ user, close }: SecurityUserFormProps) => {
  const security = useSecurityService();
  const isEditing = Boolean(user);
  const form = useCustomForm(isEditing ? updateSchema : createSchema, {
    defaultValues: {
      username: user?.username ?? "",
      fullName: user?.fullName ?? null,
      password: "",
      mustChangePassword: user?.mustChangePassword ?? true,
    },
  });
  const save = useMutate({
    mutationFn: (data: {
      username: string;
      fullName: string | null;
      password: string;
      mustChangePassword: boolean;
    }) => {
      if (user) {
        return security.users.update(user.username, {
          fullName: data.fullName,
          mustChangePassword: data.mustChangePassword,
          disabled: user.disabled,
        });
      }

      return security.users.create(data);
    },
    invalidateQueryKey: queryKey,
    successMessage: isEditing ? "Usuario actualizado." : "Usuario creado.",
  });

  return (
    <form
      className="grid gap-3"
      onSubmit={form.handleSubmit(async (values) => {
        await save.mutate(values);
        close?.();
      })}
    >
      <InputField form={form} name="username" label="Usuario" required readOnly={isEditing} />
      <InputField form={form} name="fullName" label="Nombre" />
      {!isEditing && (
        <InputField form={form} name="password" label="Password" type="password" required />
      )}
      <InputField form={form} name="mustChangePassword" label="Requiere cambio de password" variant="checkbox" />
      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button
          type="submit"
          loading={save.isPending}
          permissions={[isEditing ? Permissions.securityUsersUpdate : Permissions.securityUsersCreate]}
        >
          Guardar
        </Button>
      </div>
    </form>
  );
};
