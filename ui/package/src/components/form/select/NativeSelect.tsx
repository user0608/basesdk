// Select.tsx
import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useController } from "react-hook-form";
import { FiChevronDown } from "react-icons/fi";
import classNames from "classnames";
import { FormField } from "../shared/FormField";
import { getFieldClassName } from "../shared/getFieldClassName";
import type { SelectOption } from "../shared/SelectOption";
import { toStringArray } from "../../../utils/toStringArray";

export type NativeSelectProps<TFormValues extends FieldValues> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  options: SelectOption[];
  label?: string;
  info?: string;
  required?: boolean;
  multiple?: boolean;
  className?: string;
  placeholder?: string;
  readOnly?: boolean;
  onChange?: (value: string | string[]) => void;
  small?: boolean;
  onRefresh?: () => Promise<void> | void;
};

export const NativeSelect = <TFormValues extends FieldValues>({
  form,
  name,
  options,
  label,
  info,
  required,
  multiple = false,
  className,
  placeholder,
  readOnly,
  onChange,
  small,
  onRefresh,
}: NativeSelectProps<TFormValues>) => {
  const controller = useController({
    name,
    control: form.control,
  });

  const id = String(name);
  const error = controller.fieldState.error;
  const selectedValue = multiple
    ? toStringArray(controller.field.value)
    : String(controller.field.value ?? "");

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
          multiple={multiple}
          value={selectedValue}
          disabled={readOnly}
          onBlur={() => controller.field.onBlur()}
          onChange={(event) => {
            const value = multiple
              ? Array.from(event.target.selectedOptions, (option) => option.value)
              : event.target.value;

            controller.field.onChange(value);
            onChange?.(value);
          }}
          className={classNames(
            getFieldClassName({
              small,
              error: Boolean(error),
              className,
              withRightIcon: true,
              disabled: readOnly,
            }),
            "appearance-none",
          )}
        >
          {placeholder && !multiple && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}

          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>

        <FiChevronDown
          size={small ? 16 : 18}
          className={classNames(
            "pointer-events-none absolute right-2 text-ui-text-soft",
            multiple ? "top-3" : "top-1/2 -translate-y-1/2",
          )}
        />
      </div>
    </FormField>
  );
};
