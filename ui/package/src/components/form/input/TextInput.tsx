import type React from "react";
import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useController } from "react-hook-form";
import type { IconType } from "react-icons";
import { FiSearch } from "react-icons/fi";
import classNames from "classnames";
import { FormField } from "../shared/FormField";
import { getFieldClassName } from "../shared/getFieldClassName";

export type TextInputProps<TFormValues extends FieldValues> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  label?: string;
  info?: string;
  required?: boolean;
  type?: React.HTMLInputTypeAttribute;
  className?: string;
  placeholder?: string;
  readOnly?: boolean;
  onChange?: (value: string) => void;
  icon?: IconType;
  small?: boolean;
  onIconClick?: () => Promise<void> | void;
  onRefresh?: () => Promise<void> | void;
};

export const TextInput = <TFormValues extends FieldValues>({
  form,
  name,
  label,
  info,
  required,
  type = "text",
  className,
  placeholder,
  readOnly,
  onChange,
  icon,
  small,
  onIconClick,
  onRefresh,
}: TextInputProps<TFormValues>) => {
  const controller = useController({
    name,
    control: form.control,
  });

  const id = String(name);
  const error = controller.fieldState.error;
  const Icon = icon ?? (onIconClick ? FiSearch : undefined);

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
        <input
          id={id}
          name={controller.field.name}
          type={type}
          value={String(controller.field.value ?? "")}
          placeholder={placeholder}
          readOnly={readOnly}
          onBlur={() => controller.field.onBlur()}
          onChange={(event) => {
            const value = event.target.value;

            controller.field.onChange(value);
            onChange?.(value);
          }}
          className={getFieldClassName({
            small,
            error: Boolean(error),
            className,
            withRightIcon: Boolean(Icon),
            readOnly,
          })}
        />

        {Icon && (
          <button
            type="button"
            onClick={onIconClick}
            disabled={!onIconClick}
            className={classNames(
              "absolute right-2 top-1/2 -translate-y-1/2 rounded p-1 text-ui-text-soft",
              onIconClick && "hover:text-ui-text",
              !onIconClick && "cursor-default",
            )}
          >
            <Icon size={small ? 16 : 18} />
          </button>
        )}
      </div>
    </FormField>
  );
};
