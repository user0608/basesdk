import { useId } from "react";
import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useController } from "react-hook-form";
import classNames from "classnames";
import { FormField } from "../shared/FormField";

export type CheckboxProps<TFormValues extends FieldValues> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  label?: string;
  info?: string;
  required?: boolean;
  readOnly?: boolean;
  className?: string;
  onChange?: (value: boolean) => void;
  onRefresh?: () => Promise<void> | void;
};

export const Checkbox = <TFormValues extends FieldValues>({
  form,
  name,
  label,
  info,
  required,
  readOnly,
  className,
  onChange,
  onRefresh,
}: CheckboxProps<TFormValues>) => {
  const controller = useController({
    name,
    control: form.control,
  });

  const autoId = useId();
  const id = `${String(name)}-${autoId}`;
  const error = controller.fieldState.error;
  const checked = Boolean(controller.field.value);

  return (
    <FormField
      id={id}
      label={label}
      info={info}
      required={required}
      onRefresh={onRefresh}
      error={error?.message ? String(error.message) : undefined}
    >
      <label
        htmlFor={id}
        className={classNames(
          "flex items-center gap-2 rounded-xl px-3 py-2 text-sm shadow-sm ring-1 ring-inset transition-colors",
          error ? "ring-ui-danger" : "ring-ui-border/70",
          readOnly
            ? "bg-ui-surface-muted text-ui-text-soft"
            : "bg-ui-surface text-ui-text-muted",
          !readOnly && "hover:bg-ui-surface-hover",
          className,
        )}
      >
        <input
          id={id}
          type="checkbox"
          checked={checked}
          disabled={readOnly}
          onBlur={() => controller.field.onBlur()}
          onChange={(event) => {
            const next = event.target.checked;
            controller.field.onChange(next);
            onChange?.(next);
          }}
          className="h-4 w-4 accent-ui-accent"
        />
        <span className="select-none">{checked ? "Sí" : "No"}</span>
      </label>
    </FormField>
  );
};
