import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import type { ColumnDef } from "@tanstack/react-table";
import { NativeSelect } from "./select/NativeSelect";
import { Combobox } from "./select/Combobox";
import { CardPicker } from "./select/CardPicker";
import { TablePicker } from "./select/TablePicker";
import type { SelectOption } from "./shared/SelectOption";

export type SelectFieldVariant = "native" | "combobox" | "card" | "table";

export type SelectFieldProps<TFormValues extends FieldValues, TOption extends SelectOption = SelectOption> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  variant?: SelectFieldVariant;
  options: TOption[];
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
  loading?: boolean;
  onChange?: (value: string | string[]) => void;
  onRefresh?: () => Promise<void> | void;
  onOpen?: () => void;
  columns?: ColumnDef<TOption, unknown>[];
  searchKeys?: readonly string[];
};

export const SelectField = <TFormValues extends FieldValues, TOption extends SelectOption = SelectOption>({
  variant = "native",
  ...props
}: SelectFieldProps<TFormValues, TOption>) => {
  switch (variant) {
    case "combobox":
      return <Combobox {...props} />;
    case "card":
      return <CardPicker {...props} />;
    case "table":
      return <TablePicker {...props} />;
    case "native":
    default:
      return <NativeSelect {...props} />;
  }
};
