import type { ReactNode } from "react";
import type { ColumnDef } from "@tanstack/react-table";

export type DataTableActionIcon = React.ComponentType<{ className?: string }> | string;

export type DataTableActionVariant = "default" | "danger";

export type MaybeResolver<TInput, TValue> = TValue | ((input: TInput) => TValue);

export type DataTableContext<TData> = {
  selectedRows: TData[];
  selectedRowIds: string[];
  clearSelection: () => void;
  totalRows: number;
  filteredRows: number;
};

export type DataTableAction<TData> = {
  icon?: DataTableActionIcon;
  label: ReactNode;
  permissions?: readonly string[];
  onClick: (context: DataTableContext<TData>) => void | Promise<void>;
  disabled?: MaybeResolver<DataTableContext<TData>, boolean>;
  hidden?: MaybeResolver<DataTableContext<TData>, boolean>;
  variant?: DataTableActionVariant;
};

export type DataTableRowOption<TData> = {
  icon?: DataTableActionIcon;
  label: ReactNode;
  permissions?: readonly string[];
  onClick: (row: TData) => void | Promise<void>;
  disabled?: MaybeResolver<TData, boolean>;
  hidden?: MaybeResolver<TData, boolean>;
  variant?: DataTableActionVariant;
};

export type DataTableProps<TData> = {
  tableId?: string;
  title?: ReactNode;
  description?: ReactNode;
  data: TData[];
  columns: ColumnDef<TData, unknown>[];
  initialHiddenColumns?: readonly string[];
  selectable?: boolean;
  getRowId?: (row: TData) => string;
  searchable?: boolean;
  searchPlaceholder?: string;
  searchKeys?: readonly string[];
  actions?: DataTableAction<TData>[];
  options?: DataTableAction<TData>[];
  rowOptions?: DataTableRowOption<TData>[];
  content?: ReactNode | ((context: DataTableContext<TData>) => ReactNode);
  loading?: boolean;
  emptyMessage?: string;
  pagination?: boolean;
  pageSize?: number;
  maxItemsRange?: number;
  onRowClick?: (row: TData) => void;
};

export type PaginationRangeData = {
  totalItems: number;
  limitItems: number;
  offsetItems: number;
};

export type PaginationRangeInputProps = {
  pageData: PaginationRangeData;
  maxItemsRange: number;
  onChange?: (limit: number, offset: number) => void;
};
