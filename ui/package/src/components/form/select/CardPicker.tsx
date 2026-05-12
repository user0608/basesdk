import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useController } from "react-hook-form";
import classNames from "classnames";
import { FormField } from "../shared/FormField";
import type { SelectOption } from "../shared/SelectOption";
import { toStringArray } from "../../../utils/toStringArray";

export type CardPickerProps<TFormValues extends FieldValues, TOption extends SelectOption = SelectOption> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  options: TOption[];
  label?: string;
  info?: string;
  required?: boolean;
  multiple?: boolean;
  readOnly?: boolean;
  loading?: boolean;
  onChange?: (value: string | string[]) => void;
  onRefresh?: () => Promise<void> | void;
  className?: string;
};

export const CardPicker = <TFormValues extends FieldValues, TOption extends SelectOption = SelectOption>({
  form,
  name,
  options,
  label,
  info,
  required,
  multiple = false,
  readOnly,
  loading = false,
  onChange,
  onRefresh,
  className,
}: CardPickerProps<TFormValues, TOption>) => {
  const controller = useController({
    name,
    control: form.control,
  });

  const id = String(name);
  const error = controller.fieldState.error;

  const selectedValues = multiple
    ? toStringArray(controller.field.value)
    : controller.field.value
      ? [String(controller.field.value)]
      : [];

  const updateValue = (value: string | string[]) => {
    controller.field.onChange(value);
    onChange?.(value);
  };

  const toggleValue = (value: string) => {
    if (readOnly || loading) return;

    if (!multiple) {
      // Same click deselects.
      updateValue(selectedValues.includes(value) ? "" : value);
      controller.field.onBlur();
      return;
    }

    const nextValue = selectedValues.includes(value)
      ? selectedValues.filter((selectedValue) => selectedValue !== value)
      : [...selectedValues, value];

    updateValue(nextValue);
    controller.field.onBlur();
  };

  return (
    <FormField
      id={id}
      label={label}
      info={info}
      required={required}
      onRefresh={onRefresh}
      error={error?.message ? String(error.message) : undefined}
      loading={loading}
    >
      <div
        role={multiple ? "group" : "radiogroup"}
        aria-labelledby={label ? id : undefined}
        className={classNames("flex flex-wrap gap-2", className)}
      >
        {options.map((option) => {
          const selected = selectedValues.includes(option.value);

          return (
            <button
              key={option.value}
              type="button"
              disabled={readOnly || loading}
              onClick={() => toggleValue(option.value)}
              className={classNames(
                "inline-flex items-center rounded-xl px-3 py-2 text-sm shadow-sm ring-1 ring-inset transition-colors",
                "focus:outline-none focus:ring-2 focus:ring-ui-focus disabled:cursor-not-allowed disabled:opacity-60",
                selected
                  ? "bg-ui-surface-selected text-ui-text ring-ui-accent/40"
                  : "bg-ui-surface text-ui-text-muted ring-ui-border/70 hover:bg-ui-surface-hover",
              )}
              aria-pressed={selected}
            >
              {option.label}
            </button>
          );
        })}
      </div>
    </FormField>
  );
};
