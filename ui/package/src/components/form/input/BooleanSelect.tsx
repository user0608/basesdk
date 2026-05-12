import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useController } from "react-hook-form";
import { FiChevronDown } from "react-icons/fi";
import classNames from "classnames";
import { FormField } from "../shared/FormField";
import { getFieldClassName } from "../shared/getFieldClassName";

export type BooleanSelectProps<TFormValues extends FieldValues> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  label?: string;
  info?: string;
  required?: boolean;
  placeholder?: string;
  readOnly?: boolean;
  small?: boolean;
  onChange?: (value: boolean) => void;
  onRefresh?: () => Promise<void> | void;
  yesLabel?: string;
  noLabel?: string;
};

export const BooleanSelect = <TFormValues extends FieldValues>({
  form,
  name,
  label,
  info,
  required,
  placeholder,
  readOnly,
  small,
  onChange,
  onRefresh,
  yesLabel = "Sí",
  noLabel = "No",
}: BooleanSelectProps<TFormValues>) => {
  const controller = useController({
    name,
    control: form.control,
  });

  const id = String(name);
  const error = controller.fieldState.error;
  const value = controller.field.value;

  const selected = typeof value === "boolean" ? (value ? "true" : "false") : "";

  return (
    <FormField
      id={id}
      label={label}
      info={info}
      required={required}
      onRefresh={onRefresh}
      error={error?.message ? String(error.message) : undefined}
    >
      <div className="relative">
        <select
          id={id}
          name={controller.field.name}
          value={selected}
          disabled={readOnly}
          onBlur={() => controller.field.onBlur()}
          onChange={(event) => {
            const raw = event.target.value;
            const next = raw === "true";
            controller.field.onChange(next);
            onChange?.(next);
          }}
          className={classNames(
            getFieldClassName({
              small,
              error: Boolean(error),
              withRightIcon: true,
              disabled: readOnly,
            }),
            "appearance-none",
          )}
        >
          {placeholder && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}
          <option value="true">{yesLabel}</option>
          <option value="false">{noLabel}</option>
        </select>

        <FiChevronDown
          size={small ? 16 : 18}
          className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 text-ui-text-soft"
        />
      </div>
    </FormField>
  );
};
