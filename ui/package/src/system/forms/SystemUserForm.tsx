import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import type { SystemUserResponse } from "../types";
import { useSystemService } from "../useSystemService";

const schema = createFormSchema((validators) => ({
  username: validators.requiredString(),
  password: validators.nullableString(),
  disabled: validators.requiredBoolean(),
}));

export type SystemUserFormProps = {
  user?: SystemUserResponse;
  close?: () => void;
};

type SystemUserFormValues = {
  username: string;
  password: string | null;
  disabled: boolean;
};

export const SystemUserForm = ({ user, close }: SystemUserFormProps) => {
  const system = useSystemService();
  const isEdit = Boolean(user);
  const form = useCustomForm(schema, {
    defaultValues: {
      username: user?.username ?? "",
      password: null,
      disabled: user?.disabled ?? false,
    },
  });
  const save = useMutate({
    mutationFn: async (values: SystemUserFormValues) => {
      if (isEdit && user) {
        await system.users.update(user.username, { disabled: values.disabled });
        return;
      }

      await system.users.create({ username: values.username, password: values.password ?? "" });
    },
    invalidateQueryKey: ["system", "users"],
    successMessage: isEdit ? "Usuario system actualizado." : "Usuario system creado.",
  });

  return (
    <form
      className="grid gap-3"
      onSubmit={form.handleSubmit(async (values) => {
        await save.mutate(values as SystemUserFormValues);
        close?.();
      })}
    >
      <InputField form={form} name="username" label="Usuario" readOnly={isEdit} />
      {!isEdit && <InputField form={form} name="password" label="Password" type="password" />}
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
