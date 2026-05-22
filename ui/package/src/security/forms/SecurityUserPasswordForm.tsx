import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { Permissions } from "../../generated/permissions";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import type { TenantUserResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const schema = createFormSchema((validators) => ({
  password: validators.requiredString(),
}));

export type SecurityUserPasswordFormProps = {
  user: TenantUserResponse;
  close?: () => void;
};

export const SecurityUserPasswordForm = ({ user, close }: SecurityUserPasswordFormProps) => {
  const security = useSecurityService();
  const form = useCustomForm(schema, { defaultValues: { password: "" } });
  const save = useMutate({
    mutationFn: (values: { password: string }) => security.users.changePassword(user.username, values.password),
    successMessage: "Password actualizado.",
  });

  return (
    <form
      className="grid gap-3"
      onSubmit={form.handleSubmit(async (values) => {
        await save.mutate(values);
        close?.();
      })}
    >
      <InputField form={form} name="password" label="Nuevo password" type="password" required />
      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button type="submit" loading={save.isPending} permissions={[Permissions.securityUsersPasswordUpdate]}>
          Guardar
        </Button>
      </div>
    </form>
  );
};
