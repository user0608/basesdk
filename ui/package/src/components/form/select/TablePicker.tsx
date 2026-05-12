"use no memo";

import { useMemo, useState } from "react";
import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useController } from "react-hook-form";
import { FiSearch, FiX } from "react-icons/fi";
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  useReactTable,
  type ColumnDef,
  type PaginationState,
} from "@tanstack/react-table";
import classNames from "classnames";
import { FormField } from "../shared/FormField";
import { toStringArray } from "../../../utils/toStringArray";
import type { SelectOption } from "../shared/SelectOption";

export type TablePickerProps<TFormValues extends FieldValues, TOption extends SelectOption = SelectOption> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  options: TOption[];
  columns?: ColumnDef<TOption, unknown>[];
  searchKeys?: readonly string[];
  label?: string;
  info?: string;
  required?: boolean;
  multiple?: boolean;
  emptyMessage?: string;
  readOnly?: boolean;
  loading?: boolean;
  pageSize?: number;
  onChange?: (value: string | string[]) => void;
  onRefresh?: () => Promise<void> | void;
};

const defaultSearchKeys = ["label"] as const;

export const TablePicker = <TFormValues extends FieldValues, TOption extends SelectOption = SelectOption>({
  form,
  name,
  options,
  columns: extraColumns,
  searchKeys = defaultSearchKeys,
  label,
  info,
  required,
  multiple = false,
  emptyMessage = "No hay resultados",
  readOnly,
  loading = false,
  pageSize = 5,
  onChange,
  onRefresh,
}: TablePickerProps<TFormValues, TOption>) => {
  const controller = useController({
    name,
    control: form.control,
  });

  const id = String(name);
  const error = controller.fieldState.error;

  const [globalFilter, setGlobalFilter] = useState("");
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize,
  });

  const selectedValues = useMemo(() => {
    if (multiple) {
      return toStringArray(controller.field.value);
    }

    return controller.field.value ? [String(controller.field.value)] : [];
  }, [controller.field.value, multiple]);

  const selectedOptions = useMemo(() => {
    return options.filter((option) => selectedValues.includes(option.value));
  }, [options, selectedValues]);

  const updateValue = (value: string | string[]) => {
    controller.field.onChange(value);
    onChange?.(value);
  };

  const selectValue = (value: string) => {
    if (readOnly || loading) return;

    if (!multiple) {
      updateValue(value);
      controller.field.onBlur();
      return;
    }

    const nextValue = selectedValues.includes(value)
      ? selectedValues.filter((selectedValue) => selectedValue !== value)
      : [...selectedValues, value];

    updateValue(nextValue);
    controller.field.onBlur();
  };

  const removeValue = (value: string) => {
    if (readOnly || loading) return;

    if (!multiple) {
      updateValue("");
      controller.field.onBlur();
      return;
    }

    updateValue(selectedValues.filter((selectedValue) => selectedValue !== value));
    controller.field.onBlur();
  };

  const clearValue = () => {
    if (readOnly || loading) return;

    updateValue(multiple ? [] : "");
    controller.field.onBlur();
  };

  const columns = useMemo<ColumnDef<TOption, unknown>[]>(() => {
    const selectedColumn: ColumnDef<TOption, unknown> = {
      id: "selected",
      header: "",
      cell: ({ row }) => {
        const selected = selectedValues.includes(row.original.value);

        return (
          <input
            type={multiple ? "checkbox" : "radio"}
            checked={selected}
            disabled={readOnly}
            readOnly
            tabIndex={-1}
            className="h-4 w-4 accent-slate-900"
          />
        );
      },
      size: 40,
      enableGlobalFilter: false,
    };

    const labelColumn: ColumnDef<TOption, unknown> = {
      accessorKey: "label",
      header: "",
      cell: ({ row }) => <span className="block truncate">{row.original.label}</span>,
    };

    return [selectedColumn, labelColumn, ...(extraColumns ?? [])];
  }, [extraColumns, multiple, readOnly, selectedValues]);

  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data: options,
    columns,
    state: {
      globalFilter,
      pagination,
    },
    onGlobalFilterChange: (value) => {
      setGlobalFilter(String(value));
      setPagination((current) => ({
        ...current,
        pageIndex: 0,
      }));
    },
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    globalFilterFn: (row, _columnId, filterValue) => {
      const query = String(filterValue).trim().toLowerCase();

      if (!query) return true;

      return searchKeys.some((key) => {
        const value = (row.original as Record<string, unknown>)[key];
        if (value == null) return false;
        return String(value).toLowerCase().includes(query);
      });
    },
  });

  const rows = table.getRowModel().rows;
  const filteredRowsCount = table.getFilteredRowModel().rows.length;
  const emptyRowsCount = Math.max(0, pageSize - rows.length);
  const currentPage = table.getState().pagination.pageIndex + 1;
  const totalPages = Math.max(1, table.getPageCount());
  const columnCount = table.getAllLeafColumns().length;

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
        className={classNames(
          "overflow-hidden rounded-xl border shadow-sm",
          error ? "border-[color:var(--ui-border-error)]" : "border-[color:var(--ui-border)]",
          readOnly ? "bg-[color:var(--ui-surface-muted)]" : "bg-[color:var(--ui-panel)]",
        )}
      >
        <div className="flex items-center gap-2 border-b border-[color:var(--ui-border)] px-3 py-2">
          <FiSearch size={16} className="shrink-0 text-[color:var(--ui-text-soft)]" />

          <input
            id={`${id}-search`}
            type="text"
            value={globalFilter}
            disabled={readOnly || loading}
            placeholder={loading ? "Cargando..." : "Buscar..."}
            onBlur={() => controller.field.onBlur()}
            onChange={(event) => table.setGlobalFilter(event.target.value)}
            className="h-8 w-full bg-transparent text-sm text-[color:var(--ui-text)] outline-none placeholder:text-[color:var(--ui-text-soft)] disabled:cursor-not-allowed"
          />
        </div>

        <div className="flex min-h-10 flex-wrap items-start justify-between gap-3 border-b border-[color:var(--ui-border)] bg-[color:var(--ui-surface-muted)] px-3 py-2">
          <div className="flex min-w-0 flex-1 flex-wrap items-start gap-2">
            <span className="shrink-0 pt-0.5 text-xs font-medium text-[color:var(--ui-text-soft)]">
              {selectedOptions.length === 0
                ? "Sin selección"
                : selectedOptions.length === 1
                  ? "1 seleccionado"
                  : `${selectedOptions.length} seleccionados`}
            </span>

            {selectedOptions.length > 0 && (
              <div className="flex min-w-0 flex-1 flex-wrap gap-1">
                {selectedOptions.map((option) => (
                  <span
                    key={option.value}
                  className="inline-flex min-w-0 max-w-36 items-center gap-1 rounded-full border border-[color:var(--ui-border)] bg-[color:var(--ui-surface)] px-2 py-0.5 text-xs text-[color:var(--ui-text-muted)]"
                  >
                    <span className="truncate">{option.label}</span>

                    {!readOnly && (
                      <button
                        type="button"
                        onClick={() => removeValue(option.value)}
                        className="shrink-0 rounded text-[color:var(--ui-text-soft)] hover:text-[color:var(--ui-text)]"
                      >
                        <FiX size={12} />
                      </button>
                    )}
                  </span>
                ))}

              </div>
            )}
          </div>

          {selectedOptions.length > 0 && !readOnly && !loading && (
            <button
              type="button"
              onClick={clearValue}
              className="shrink-0 text-xs font-medium text-[color:var(--ui-text-soft)] hover:text-[color:var(--ui-text)]"
            >
              Limpiar
            </button>
          )}
        </div>

        <table className="w-full table-fixed border-collapse text-sm">
          <tbody>
            {rows.length ? (
              <>
                {rows.map((row) => {
                  const selected = selectedValues.includes(row.original.value);

                  return (
                    <tr
                      key={row.id}
                      onClick={() => selectValue(row.original.value)}
                      className={classNames(
                        "h-10 border-b border-[color:var(--ui-border)]/50 last:border-b-0",
                        readOnly || loading ? "cursor-default" : "cursor-pointer hover:bg-[color:var(--ui-surface-hover)]",
                        selected && "bg-[color:var(--ui-surface-selected)]",
                      )}
                    >
                      {row.getVisibleCells().map((cell) => (
                        <td
                          key={cell.id}
                          className={classNames(
                            "px-3 py-2 align-middle text-[color:var(--ui-text-muted)]",
                            cell.column.id === "selected" && "w-10",
                          )}
                        >
                          {flexRender(cell.column.columnDef.cell, cell.getContext())}
                        </td>
                      ))}
                    </tr>
                  );
                })}

                {Array.from({ length: emptyRowsCount }).map((_, index) => (
                  <tr key={`empty-${index}`} className="h-10 border-b border-slate-100 last:border-b-0">
                    {Array.from({ length: columnCount }).map((__, cellIndex) => (
                      <td
                        key={cellIndex}
                        className={classNames(
                          "px-3 py-2 align-middle",
                          cellIndex === 0 && "w-10",
                        )}
                      >
                        &nbsp;
                      </td>
                    ))}
                  </tr>
                ))}
              </>
            ) : (
              <>
                <tr className="h-10 border-b border-slate-100">
                  <td colSpan={columnCount} className="px-3 py-3 text-sm text-[color:var(--ui-text-soft)]">
                    {loading ? "Cargando..." : emptyMessage}
                  </td>
                </tr>

                {Array.from({ length: Math.max(0, pageSize - 1) }).map((_, index) => (
                  <tr key={`empty-${index}`} className="h-10 border-b border-slate-100 last:border-b-0">
                    {Array.from({ length: columnCount }).map((__, cellIndex) => (
                      <td
                        key={cellIndex}
                        className={classNames(
                          "px-3 py-2 align-middle",
                          cellIndex === 0 && "w-10",
                        )}
                      >
                        &nbsp;
                      </td>
                    ))}
                  </tr>
                ))}
              </>
            )}
          </tbody>
        </table>

        <div className="flex items-center justify-between border-t border-[color:var(--ui-border)] bg-[color:var(--ui-surface-muted)] px-3 py-2">
          <span className="text-xs text-[color:var(--ui-text-soft)]">
            {filteredRowsCount} resultado
            {filteredRowsCount === 1 ? "" : "s"}
          </span>

          <div className="flex items-center gap-2">
            <button
              type="button"
              disabled={readOnly || loading || !table.getCanPreviousPage()}
              onClick={() => table.previousPage()}
             className="rounded border border-[color:var(--ui-border)] bg-[color:var(--ui-surface)] px-2 py-1 text-xs text-[color:var(--ui-text-muted)] disabled:cursor-not-allowed disabled:opacity-50"
            >
              Anterior
            </button>

             <span className="text-xs text-[color:var(--ui-text-soft)]">
               {currentPage} / {totalPages}
             </span>

            <button
              type="button"
              disabled={readOnly || loading || !table.getCanNextPage()}
              onClick={() => table.nextPage()}
             className="rounded border border-[color:var(--ui-border)] bg-[color:var(--ui-surface)] px-2 py-1 text-xs text-[color:var(--ui-text-muted)] disabled:cursor-not-allowed disabled:opacity-50"
            >
              Siguiente
            </button>
          </div>
        </div>
      </div>
      </FormField>
  );
};
