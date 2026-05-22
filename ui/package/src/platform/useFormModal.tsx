import { createElement } from "react";
import type { ComponentProps, ReactNode } from "react";
import { useModal } from "../components/dialog/useModal";
import { useRegistry } from "./RegistryProvider";
import type { ApplicationRegistry, FormId } from "./types";
import type { DialogSize } from "../components/dialog/types";

type FormComponent<TRegistry extends ApplicationRegistry, TFormId extends FormId<TRegistry>> = NonNullable<
  TRegistry["forms"]
>[TFormId];

export type OpenFormModalOptions<TProps> = {
  title?: ReactNode;
  description?: ReactNode;
  size?: DialogSize;
  props?: Omit<TProps, "close">;
};

export const useFormModal = <TRegistry extends ApplicationRegistry = ApplicationRegistry>() => {
  const registry = useRegistry<TRegistry>();
  const modal = useModal();

  return {
    open: <TFormId extends FormId<TRegistry>>(
      formId: TFormId,
      options: OpenFormModalOptions<ComponentProps<FormComponent<TRegistry, TFormId>>> = {},
    ) => {
      const Form = registry.forms?.[formId];

      if (!Form) {
        throw new Error(`Formulario no registrado: ${String(formId)}`);
      }

      modal.open({
        title: options.title,
        description: options.description,
        size: options.size,
        content: ({ close }) => createElement(Form, { ...(options.props ?? {}), close }),
      });
    },
  };
};
