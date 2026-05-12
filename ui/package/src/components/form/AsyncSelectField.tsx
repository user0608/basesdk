import { useCallback } from "react";
import type { QueryKey } from "@tanstack/react-query";
import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useWatch } from "react-hook-form";
import type { ColumnDef } from "@tanstack/react-table";
import { SelectField } from "./SelectField";
import type { SelectFieldVariant } from "./SelectField";
import type { SelectOption } from "./shared/SelectOption";
import { useAsyncQueryOptions } from "./shared/useAsyncQueryOptions";
import { toStringArray } from "../../utils/toStringArray";

export type AsyncSelectFieldProps<TFormValues extends FieldValues, TOption extends SelectOption = SelectOption> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  variant?: SelectFieldVariant;
  queryKey: QueryKey;
  loadOptions: () => Promise<TOption[]>;
  label?: string;
  info?: string;
  required?: boolean;
  multiple?: boolean;
  readOnly?: boolean;
  className?: string;
  placeholder?: string;
  searchPlaceholder?: string;
  emptyMessage?: string;
  small?: boolean;
  pageSize?: number;
  onChange?: (value: string | string[]) => void;
  onRefresh?: () => Promise<void> | void;
  columns?: ColumnDef<TOption, unknown>[];
  searchKeys?: readonly string[];
  loadOnMount?: boolean;
};

export const AsyncSelectField = <TFormValues extends FieldValues, TOption extends SelectOption = SelectOption>({
  form,
  name,
  variant = "native",
  queryKey,
  loadOptions,
  label,
  info,
  required,
  multiple = false,
  readOnly,
  className,
  placeholder,
  searchPlaceholder,
  emptyMessage,
  small,
  pageSize,
  onChange,
  onRefresh,
  columns,
  searchKeys,
  loadOnMount,
}: AsyncSelectFieldProps<TFormValues, TOption>) => {
  const value = useWatch({
    control: form.control,
    name,
  });

  const selectedValues = multiple ? toStringArray(value) : value ? [String(value)] : [];
  const hasInitialValue = selectedValues.length > 0;
  const shouldLoadOnMount = loadOnMount ?? (variant !== "combobox");

  const { options, loaded, error, isPending, refresh } = useAsyncQueryOptions<TOption>({
    queryKey,
    loadOptions,
    loadOnMount: shouldLoadOnMount,
    loadOnInitialValue: !shouldLoadOnMount,
    hasInitialValue,
  });

  const handleOpen = useCallback(() => {
    if (variant !== "combobox" || loaded || isPending) return;

    refresh();
  }, [isPending, loaded, refresh, variant]);

  return (
    <SelectField
      form={form}
      name={name}
      variant={variant}
      options={options}
      label={label}
      info={error ?? info}
      required={required}
      multiple={multiple}
      readOnly={readOnly || isPending}
      className={className}
      placeholder={isPending && variant === "combobox" ? "Cargando..." : placeholder}
      searchPlaceholder={searchPlaceholder}
      emptyMessage={error ?? emptyMessage}
      small={small}
      pageSize={pageSize}
      loading={isPending}
      onChange={onChange}
      onRefresh={refresh}
      onOpen={handleOpen}
      columns={columns}
      searchKeys={searchKeys}
    />
  );
};
