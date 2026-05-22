import { useEffect, useMemo, useState } from "react";
import { createPortal } from "react-dom";
import { FiColumns, FiSearch } from "react-icons/fi";
import { flexRender, getCoreRowModel, getFilteredRowModel, useReactTable } from "@tanstack/react-table";
import type { Row, Table, VisibilityState } from "@tanstack/react-table";
import classNames from "classnames";
import { hasAllPermissions, useCurrentPermissions } from "../../platform/permissions";
import { ActionMenu } from "./ActionMenu";
import { PaginationRangeInput } from "./PaginationRangeInput";
import { renderActionIcon, resolveValue } from "./shared";
import type { DataTableContext, DataTableProps } from "./types";

const getColumnLabel = <TData,>(column: { id: string; columnDef: { header?: unknown } }) => {
  return typeof column.columnDef.header === "string" ? column.columnDef.header : column.id;
};

const renderSlot = <TData,>(slot: DataTableProps<TData>["content"], context: DataTableContext<TData>) => {
  if (!slot) return null;
  return typeof slot === "function" ? slot(context) : slot;
};

const columnVisibilityStorageKey = (tableId: string) => `basesdk:data-table:${tableId}:column-visibility`;

const hiddenColumnsToVisibility = (hiddenColumns?: readonly string[]): VisibilityState => {
  const visibility: VisibilityState = {};
  for (const columnId of hiddenColumns ?? []) {
    visibility[columnId] = false;
  }
  return visibility;
};

const readColumnVisibility = (tableId: string | undefined, initialHiddenColumns?: readonly string[]) => {
  if (!tableId || typeof window === "undefined") return hiddenColumnsToVisibility(initialHiddenColumns);

  try {
    const stored = window.localStorage.getItem(columnVisibilityStorageKey(tableId));
    if (!stored) return hiddenColumnsToVisibility(initialHiddenColumns);

    const parsed = JSON.parse(stored);
    return parsed && typeof parsed === "object" ? (parsed as VisibilityState) : hiddenColumnsToVisibility(initialHiddenColumns);
  } catch {
    return hiddenColumnsToVisibility(initialHiddenColumns);
  }
};

const ColumnVisibilityMenu = <TData,>({ table }: { table: Table<TData> }) => {
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0 });
  const buttonRef = useState<HTMLButtonElement | null>(null);
  const menuRef = useState<HTMLDivElement | null>(null);
  const [button, setButton] = buttonRef;
  const [menu, setMenu] = menuRef;
  const columns = table.getAllLeafColumns().filter((column) => column.getCanHide());

  const updatePosition = () => {
    const rect = button?.getBoundingClientRect();
    if (!rect) return;

    const width = 224;
    const estimatedHeight = Math.min(320, columns.length * 36 + 52);
    const left = Math.max(8, Math.min(window.innerWidth - width - 8, rect.right - width));
    const bottomTop = rect.bottom + 4;
    const top = bottomTop + estimatedHeight > window.innerHeight - 8 ? Math.max(8, rect.top - estimatedHeight - 4) : bottomTop;
    setPosition({ top, left });
  };

  useEffect(() => {
    if (!open) return;
    updatePosition();

    const onPointerDown = (event: PointerEvent) => {
      const target = event.target as Node;
      if (!menu?.contains(target) && !button?.contains(target)) setOpen(false);
    };
    const onReposition = () => updatePosition();

    document.addEventListener("pointerdown", onPointerDown);
    window.addEventListener("scroll", onReposition, true);
    window.addEventListener("resize", onReposition);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      window.removeEventListener("scroll", onReposition, true);
      window.removeEventListener("resize", onReposition);
    };
  }, [button, menu, open]);

  if (columns.length === 0) return null;

  return (
    <div className="relative inline-flex">
      <button
        ref={setButton}
        type="button"
        title="Columnas"
        className="inline-flex size-8 items-center justify-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
        onClick={() => {
          updatePosition();
          setOpen((value) => !value);
        }}
      >
        <FiColumns size={16} />
      </button>

      {open &&
        createPortal(
          <div
            ref={setMenu}
            className="fixed z-[70] w-56 overflow-hidden rounded-xl bg-ui-panel py-1 shadow-xl ring-1 ring-ui-border/60"
            style={{ top: position.top, left: position.left }}
          >
            <div className="border-b border-ui-border/30 px-3 py-2 text-xs font-semibold uppercase tracking-[0.08em] text-ui-text-soft">
              Columnas visibles
            </div>
            <div className="max-h-72 overflow-auto py-1">
              {columns.map((column) => (
                <label
                  key={column.id}
                  className="flex cursor-pointer items-center gap-2 px-3 py-2 text-sm text-ui-text-muted transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
                >
                  <input
                    type="checkbox"
                    checked={column.getIsVisible()}
                    onChange={column.getToggleVisibilityHandler()}
                    className="size-4 accent-ui-accent"
                  />
                  <span className="truncate">{getColumnLabel(column)}</span>
                </label>
              ))}
            </div>
          </div>,
          document.body,
        )}
    </div>
  );
};

export const DataTable = <TData,>({
  tableId,
  title,
  description,
  data,
  columns,
  initialHiddenColumns,
  selectable = false,
  getRowId,
  searchable = true,
  searchPlaceholder = "Buscar...",
  searchKeys,
  actions,
  options,
  rowOptions,
  content,
  loading = false,
  emptyMessage = "No hay resultados",
  pagination = true,
  pageSize = 25,
  maxItemsRange = 100,
  onRowClick,
}: DataTableProps<TData>) => {
  const currentPermissions = useCurrentPermissions();
  const [globalFilter, setGlobalFilter] = useState("");
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>(() =>
    readColumnVisibility(tableId, initialHiddenColumns),
  );
  const [selectedRowIds, setSelectedRowIds] = useState<Record<string, boolean>>({});
  const [page, setPage] = useState({ limit: Math.max(1, pageSize), offset: 0 });

  const table = useReactTable({
    data,
    columns,
    getRowId: getRowId ? (row) => getRowId(row) : undefined,
    state: {
      globalFilter,
      columnVisibility,
    },
    onColumnVisibilityChange: setColumnVisibility,
    onGlobalFilterChange: (value) => {
      setGlobalFilter(String(value));
      setPage((current) => ({ ...current, offset: 0 }));
    },
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    globalFilterFn: (row, _columnId, filterValue) => {
      const query = String(filterValue).trim().toLowerCase();
      if (!query) return true;

      if (searchKeys?.length) {
        return searchKeys.some((key) => {
          const value = (row.original as Record<string, unknown>)[key];
          return value != null && String(value).toLowerCase().includes(query);
        });
      }

      return row.getAllCells().some((cell) => {
        const value = cell.getValue();
        return value != null && String(value).toLowerCase().includes(query);
      });
    },
  });

  const filteredRows = table.getFilteredRowModel().rows;
  const visibleRows = pagination ? filteredRows.slice(page.offset, page.offset + page.limit) : filteredRows;
  const selectedRows = useMemo(
    () => table.getCoreRowModel().rows.filter((row) => selectedRowIds[row.id]).map((row) => row.original),
    [selectedRowIds, table],
  );
  const selectedIds = useMemo(() => Object.keys(selectedRowIds).filter((id) => selectedRowIds[id]), [selectedRowIds]);
  const context = useMemo<DataTableContext<TData>>(
    () => ({
      selectedRows,
      selectedRowIds: selectedIds,
      clearSelection: () => setSelectedRowIds({}),
      totalRows: data.length,
      filteredRows: filteredRows.length,
    }),
    [data.length, filteredRows.length, selectedIds, selectedRows],
  );

  useEffect(() => {
    if (page.offset >= filteredRows.length) {
      setPage((current) => ({ ...current, offset: 0 }));
    }
  }, [filteredRows.length, page.offset]);

  useEffect(() => {
    if (!tableId || typeof window === "undefined") return;
    window.localStorage.setItem(columnVisibilityStorageKey(tableId), JSON.stringify(columnVisibility));
  }, [columnVisibility, tableId]);

  const allVisibleSelected = visibleRows.length > 0 && visibleRows.every((row) => selectedRowIds[row.id]);
  const someVisibleSelected = visibleRows.some((row) => selectedRowIds[row.id]);

  const toggleRow = (row: Row<TData>) => {
    setSelectedRowIds((current) => ({
      ...current,
      [row.id]: !current[row.id],
    }));
  };

  const toggleVisibleRows = () => {
    setSelectedRowIds((current) => {
      const next = { ...current };
      for (const row of visibleRows) {
        if (allVisibleSelected) {
          delete next[row.id];
        } else {
          next[row.id] = true;
        }
      }
      return next;
    });
  };

  const visibleActions = (actions ?? [])
    .filter((action) => hasAllPermissions(currentPermissions, action.permissions))
    .filter((action) => !resolveValue(action.hidden, context))
    .map((action) => ({
      ...action,
      disabled: Boolean(resolveValue(action.disabled, context)),
    }));
  const visibleOptions = (options ?? [])
    .filter((option) => hasAllPermissions(currentPermissions, option.permissions))
    .filter((option) => !resolveValue(option.hidden, context))
    .map((option) => ({
      ...option,
      disabled: Boolean(resolveValue(option.disabled, context)),
    }));

  const renderRowOptions = (row: Row<TData>) => {
    const items = (rowOptions ?? [])
      .filter((option) => hasAllPermissions(currentPermissions, option.permissions))
      .filter((option) => !resolveValue(option.hidden, row.original))
      .map((option) => ({
        icon: option.icon,
        label: option.label,
        variant: option.variant,
        disabled: Boolean(resolveValue(option.disabled, row.original)),
        onClick: () => option.onClick(row.original),
      }));

    return <ActionMenu label="" items={items} />;
  };

  const colSpan = table.getVisibleLeafColumns().length + (selectable ? 1 : 0) + (rowOptions?.length ? 1 : 0);

  return (
    <section className="flex h-full min-h-0 flex-col overflow-hidden rounded-2xl bg-ui-panel shadow-sm ring-1 ring-ui-border/35">
      {(title || description) && (
        <header className="border-b border-ui-border/25 px-3 py-3 sm:px-4">
          {title && <h2 className="text-base font-semibold tracking-tight text-ui-text">{title}</h2>}
          {description && <p className="mt-1 text-sm leading-5 text-ui-text-muted">{description}</p>}
        </header>
      )}

      {(searchable || visibleActions.length > 0 || visibleOptions.length > 0 || table.getAllLeafColumns().length > 0) && (
        <div className="grid gap-2 border-b border-ui-border/25 px-3 py-2 sm:px-4 lg:flex lg:items-center lg:justify-between">
          {searchable && (
            <label className="flex min-h-9 min-w-0 flex-1 items-center gap-2 rounded-xl border border-ui-border/60 bg-ui-surface px-3 text-sm text-ui-text">
              <FiSearch size={16} className="shrink-0 text-ui-text-soft" />
              <input
                value={globalFilter}
                disabled={loading}
                placeholder={searchPlaceholder}
                className="h-8 min-w-0 flex-1 bg-transparent outline-none placeholder:text-ui-text-soft disabled:cursor-not-allowed"
                onChange={(event) => table.setGlobalFilter(event.target.value)}
              />
            </label>
          )}

          {(visibleActions.length > 0 || visibleOptions.length > 0 || table.getAllLeafColumns().length > 0) && (
            <div className="flex min-w-0 flex-wrap items-center gap-2 lg:justify-end">
              {visibleActions.map((action, index) => (
                <button
                  key={index}
                  type="button"
                  disabled={action.disabled}
                  className={classNames(
                    "inline-flex h-8 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50",
                    action.variant === "danger"
                      ? "bg-ui-danger/10 text-ui-danger hover:bg-ui-danger/15"
                      : "bg-ui-primary text-ui-text-inverse hover:bg-ui-primary-hover",
                  )}
                  onClick={() => action.onClick(context)}
                >
                  {renderActionIcon(action.icon, "size-4 shrink-0 object-contain")}
                  <span>{action.label}</span>
                </button>
              ))}

              <ActionMenu
                items={visibleOptions.map((option) => ({
                  icon: option.icon,
                  label: option.label,
                  variant: option.variant,
                  disabled: option.disabled,
                  onClick: () => option.onClick(context),
                }))}
              />
              <ColumnVisibilityMenu table={table} />
            </div>
          )}
        </div>
      )}

      {content && <div className="border-b border-ui-border/25 px-3 py-2 sm:px-4">{renderSlot(content, context)}</div>}

      <div className="hidden min-h-0 flex-1 overflow-auto lg:block">
        <table className="w-full border-collapse text-sm">
          <thead className="sticky top-0 z-30 bg-ui-surface-muted text-xs font-semibold uppercase tracking-[0.08em] text-ui-text-soft">
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {selectable && (
                  <th className="sticky left-0 z-40 w-10 border-b border-ui-border/40 bg-ui-surface-muted px-3 py-2 text-left shadow-[8px_0_12px_-14px_rgba(15,23,42,0.45)]">
                    <input
                      type="checkbox"
                      checked={allVisibleSelected}
                      ref={(input) => {
                        if (input) input.indeterminate = !allVisibleSelected && someVisibleSelected;
                      }}
                      onChange={toggleVisibleRows}
                      className="size-4 accent-ui-accent"
                    />
                  </th>
                )}
                {headerGroup.headers.map((header) => (
                  <th key={header.id} className="border-b border-ui-border/40 px-3 py-2 text-left">
                    {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                  </th>
                ))}
                {Boolean(rowOptions?.length) && (
                  <th className="sticky right-0 z-40 w-12 border-b border-ui-border/40 bg-ui-surface-muted px-3 py-2 shadow-[-8px_0_12px_-14px_rgba(15,23,42,0.45)]" />
                )}
              </tr>
            ))}
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={colSpan} className="px-3 py-8 text-center text-sm text-ui-text-muted">
                  Cargando...
                </td>
              </tr>
            ) : visibleRows.length ? (
              visibleRows.map((row) => (
                <tr
                  key={row.id}
                  className={classNames(
                    "border-b border-ui-border/30 last:border-b-0",
                    onRowClick && "cursor-pointer hover:bg-ui-surface-hover",
                    selectedRowIds[row.id] && "bg-ui-surface-selected/70",
                  )}
                  onClick={() => onRowClick?.(row.original)}
                >
                  {selectable && (
                    <td
                      className={classNames(
                        "sticky left-0 z-10 w-10 px-3 py-1.5 align-middle shadow-[8px_0_12px_-14px_rgba(15,23,42,0.45)]",
                        selectedRowIds[row.id] ? "bg-ui-surface-selected" : "bg-ui-panel",
                      )}
                      onClick={(event) => event.stopPropagation()}
                    >
                      <input
                        type="checkbox"
                        checked={Boolean(selectedRowIds[row.id])}
                        onChange={() => toggleRow(row)}
                        className="size-4 accent-ui-accent"
                      />
                    </td>
                  )}
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id} className="px-3 py-1.5 align-middle text-ui-text-muted">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                  {Boolean(rowOptions?.length) && (
                    <td
                      className={classNames(
                        "sticky right-0 z-10 w-12 px-3 py-1.5 text-right align-middle shadow-[-8px_0_12px_-14px_rgba(15,23,42,0.45)]",
                        selectedRowIds[row.id] ? "bg-ui-surface-selected" : "bg-ui-panel",
                      )}
                      onClick={(event) => event.stopPropagation()}
                    >
                      {renderRowOptions(row)}
                    </td>
                  )}
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={colSpan} className="px-3 py-8 text-center text-sm text-ui-text-muted">
                  {emptyMessage}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="min-h-0 flex-1 overflow-auto p-3 lg:hidden">
        <div className="grid gap-2">
        {loading ? (
          <div className="rounded-xl border border-ui-border/40 bg-ui-surface px-3 py-6 text-center text-sm text-ui-text-muted">
            Cargando...
          </div>
        ) : visibleRows.length ? (
          visibleRows.map((row) => {
            const cells = row.getVisibleCells();
            const titleCell = cells[0];

            return (
              <article key={row.id} className="rounded-xl border border-ui-border/45 bg-ui-surface px-3 py-2 shadow-sm">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex min-w-0 items-start gap-2">
                    {selectable && (
                      <input
                        type="checkbox"
                        checked={Boolean(selectedRowIds[row.id])}
                        onChange={() => toggleRow(row)}
                        className="mt-1 size-4 shrink-0 accent-ui-accent"
                      />
                    )}
                    <div className="min-w-0 text-sm font-semibold text-ui-text">
                      {titleCell ? flexRender(titleCell.column.columnDef.cell, titleCell.getContext()) : row.id}
                    </div>
                  </div>
                  {Boolean(rowOptions?.length) && renderRowOptions(row)}
                </div>

                {cells.slice(1).length > 0 && (
                  <dl className="mt-2 grid gap-1 text-sm">
                    {cells.slice(1).map((cell) => (
                      <div key={cell.id} className="grid grid-cols-[96px_minmax(0,1fr)] gap-2">
                        <dt className="truncate text-xs font-medium text-ui-text-soft">{getColumnLabel(cell.column)}</dt>
                        <dd className="min-w-0 truncate text-ui-text-muted">{flexRender(cell.column.columnDef.cell, cell.getContext())}</dd>
                      </div>
                    ))}
                  </dl>
                )}
              </article>
            );
          })
        ) : (
          <div className="rounded-xl border border-ui-border/40 bg-ui-surface px-3 py-6 text-center text-sm text-ui-text-muted">
            {emptyMessage}
          </div>
        )}
        </div>
      </div>

      {pagination && (
        <footer className="flex flex-col gap-2 border-t border-ui-border/25 bg-ui-surface-muted px-3 py-2 sm:flex-row sm:items-center sm:justify-between sm:px-4">
          <span className="text-xs text-ui-text-soft">
            {filteredRows.length} resultado{filteredRows.length === 1 ? "" : "s"}
          </span>
          <PaginationRangeInput
            pageData={{
              totalItems: filteredRows.length,
              limitItems: page.limit,
              offsetItems: page.offset,
            }}
            maxItemsRange={maxItemsRange}
            onChange={(limit, offset) => setPage({ limit: Math.max(1, limit), offset })}
          />
        </footer>
      )}
    </section>
  );
};
