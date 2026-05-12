import { useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";
import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import { useController } from "react-hook-form";
import { FiCheck, FiChevronDown, FiSearch, FiX } from "react-icons/fi";
import classNames from "classnames";
import { FormField } from "../shared/FormField";
import { getFieldClassName } from "../shared/getFieldClassName";
import { toStringArray } from "../../../utils/toStringArray";
import type { SelectOption } from "../shared/SelectOption";

export type ComboboxProps<TFormValues extends FieldValues> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  options: SelectOption[];
  label?: string;
  info?: string;
  required?: boolean;
  multiple?: boolean;
  className?: string;
  placeholder?: string;
  searchPlaceholder?: string;
  emptyMessage?: string;
  readOnly?: boolean;
  small?: boolean;
  maxVisibleSelected?: number;
  onChange?: (value: string | string[]) => void;
  onRefresh?: () => Promise<void> | void;
  loading?: boolean;
  onOpen?: () => void;
};

type DropdownPosition = {
  top?: number;
  bottom?: number;
  left: number;
  width: number;
  maxHeight: number;
  direction: "top" | "bottom";
};

const SEARCH_HEIGHT = 49;
const OPTION_ROW_HEIGHT = 40;
const OPTION_OVERSCAN = 6;

export const Combobox = <TFormValues extends FieldValues>({
  form,
  name,
  options,
  label,
  info,
  required,
  multiple = false,
  className,
  placeholder = "Seleccionar",
  searchPlaceholder = "Buscar...",
  emptyMessage = "No hay resultados",
  readOnly,
  small,
  maxVisibleSelected,
  onChange,
  onRefresh,
  onOpen,
  loading = false,
}: ComboboxProps<TFormValues>) => {
  const controller = useController({
    name,
    control: form.control,
  });

  const id = String(name);
  const error = controller.fieldState.error;
  const containerRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);

  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [dropdownPosition, setDropdownPosition] = useState<DropdownPosition>();
  const [optionsScrollTop, setOptionsScrollTop] = useState(0);
  const optionsListRef = useRef<HTMLDivElement>(null);

  const selectedValues = useMemo(() => {
    if (multiple) return toStringArray(controller.field.value);

    return controller.field.value ? [String(controller.field.value)] : [];
  }, [controller.field.value, multiple]);

  const optionByValue = useMemo(() => {
    const map = new Map<string, SelectOption>();
    for (const option of options) map.set(option.value, option);
    return map;
  }, [options]);

  const selectedOptions = useMemo(() => {
    const result: SelectOption[] = [];
    for (const value of selectedValues) {
      const option = optionByValue.get(value);
      if (option) result.push(option);
    }
    return result;
  }, [optionByValue, selectedValues]);

  const filteredOptions = useMemo(() => {
    const value = search.trim().toLowerCase();

    if (!value) return options;

    return options.filter((option) => option.label.toLowerCase().includes(value));
  }, [options, search]);

  useEffect(() => {
    if (!open) return;

    // Reset scroll position when opening or changing search.
    optionsListRef.current?.scrollTo({ top: 0 });
    queueMicrotask(() => {
      setOptionsScrollTop(optionsListRef.current?.scrollTop ?? 0);
    });
  }, [open, search]);

  const updateValue = (value: string | string[]) => {
    controller.field.onChange(value);
    onChange?.(value);
  };

  const closeDropdown = () => {
    setOpen(false);
    setSearch("");
    setDropdownPosition(undefined);
  };

  const selectValue = (value: string) => {
    if (readOnly) return;

    if (!multiple) {
      updateValue(value);
      closeDropdown();
      return;
    }

    const nextValue = selectedValues.includes(value)
      ? selectedValues.filter((selectedValue) => selectedValue !== value)
      : [...selectedValues, value];

    updateValue(nextValue);
  };

  const removeValue = (value: string) => {
    if (readOnly) return;

    if (!multiple) {
      updateValue("");
      return;
    }

    updateValue(selectedValues.filter((selectedValue) => selectedValue !== value));
  };

  const clearValue = () => {
    if (readOnly) return;

    updateValue(multiple ? [] : "");
  };

  const updateDropdownPosition = () => {
    const button = buttonRef.current;

    if (!button) return;

    const rect = button.getBoundingClientRect();
    const gap = 4;
    const viewportPadding = 12;
    const minHeight = 120;
    const preferredHeight = 288;

    const spaceBelow = window.innerHeight - rect.bottom - viewportPadding - gap;
    const spaceAbove = rect.top - viewportPadding - gap;
    const direction = spaceAbove > spaceBelow && spaceBelow < minHeight ? "top" : "bottom";
    const availableSpace = direction === "top" ? spaceAbove : spaceBelow;
    const maxHeight = Math.max(minHeight, Math.min(preferredHeight, availableSpace));

    setDropdownPosition({
      left: rect.left,
      width: rect.width,
      maxHeight,
      direction,
      top: direction === "bottom" ? rect.bottom + gap : undefined,
      bottom: direction === "top" ? window.innerHeight - rect.top + gap : undefined,
    });
  };

  useLayoutEffect(() => {
    if (!open) return;

    updateDropdownPosition();
  }, [open]);

  useEffect(() => {
    if (!open) return;

    const handlePositionUpdate = () => updateDropdownPosition();

    window.addEventListener("resize", handlePositionUpdate);
    window.addEventListener("scroll", handlePositionUpdate, true);

    return () => {
      window.removeEventListener("resize", handlePositionUpdate);
      window.removeEventListener("scroll", handlePositionUpdate, true);
    };
  }, [open]);

  useEffect(() => {
    const handlePointerDown = (event: PointerEvent) => {
      const target = event.target as Node;

      if (!containerRef.current?.contains(target)) {
        closeDropdown();
      }
    };

    document.addEventListener("pointerdown", handlePointerDown);

    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
    };
  }, []);

  const searchInput = (
    <div
      className={classNames(
        "flex shrink-0 items-center gap-2 px-3 py-2",
        dropdownPosition?.direction === "top"
          ? "border-t border-[color:var(--ui-border)]"
          : "border-b border-[color:var(--ui-border)]",
      )}
    >
      <FiSearch size={16} className="text-[color:var(--ui-text-soft)]" />

      <input
        value={search}
        placeholder={searchPlaceholder}
        onChange={(event) => setSearch(event.target.value)}
        className="h-8 w-full bg-transparent text-sm text-[color:var(--ui-text)] outline-none placeholder:text-[color:var(--ui-text-soft)]"
      />
    </div>
  );

  const optionsList = (
    <div
      ref={optionsListRef}
      className="overflow-auto py-1"
      style={{
        maxHeight: Math.max(72, (dropdownPosition?.maxHeight ?? 288) - SEARCH_HEIGHT),
      }}
      onScroll={(event) => {
        setOptionsScrollTop((event.currentTarget as HTMLDivElement).scrollTop);
      }}
    >
      {filteredOptions.length ? (() => {
        const listHeight = Math.max(72, (dropdownPosition?.maxHeight ?? 288) - SEARCH_HEIGHT);
        const visibleCount = Math.max(1, Math.ceil(listHeight / OPTION_ROW_HEIGHT));
        const start = Math.max(0, Math.floor(optionsScrollTop / OPTION_ROW_HEIGHT) - OPTION_OVERSCAN);
        const end = Math.min(filteredOptions.length, start + visibleCount + OPTION_OVERSCAN * 2);
        const topSpacer = start * OPTION_ROW_HEIGHT;
        const bottomSpacer = Math.max(0, (filteredOptions.length - end) * OPTION_ROW_HEIGHT);
        const slice = filteredOptions.slice(start, end);

        return (
          <>
            {topSpacer > 0 && <div style={{ height: topSpacer }} />}

            {slice.map((option) => {
              const selected = selectedValues.includes(option.value);

              return (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => selectValue(option.value)}
                  className={classNames(
                     "flex h-10 w-full items-center gap-2 px-3 text-left text-sm hover:bg-[color:var(--ui-surface-hover)]",
                     selected ? "text-[color:var(--ui-text)]" : "text-[color:var(--ui-text-muted)]",
                   )}
                 >
                   <span
                     className={classNames(
                       "flex h-4 w-4 items-center justify-center border",
                       multiple ? "rounded" : "rounded-full",
                       selected
                         ? "border-[color:var(--ui-accent)] bg-[color:var(--ui-accent)] text-[color:var(--ui-bg)]"
                         : "border-[color:var(--ui-border)] bg-[color:var(--ui-surface)]",
                     )}
                   >
                    {selected && <FiCheck size={12} />}
                  </span>

                  <span className="truncate">{option.label}</span>
                </button>
              );
            })}

            {bottomSpacer > 0 && <div style={{ height: bottomSpacer }} />}
          </>
        );
      })() : (
         <div className="px-3 py-2 text-sm text-[color:var(--ui-text-soft)]">{emptyMessage}</div>
       )}
    </div>
  );

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
      <div ref={containerRef} className="relative">
        <button
          ref={buttonRef}
          id={id}
          type="button"
          disabled={readOnly}
          onBlur={() => controller.field.onBlur()}
          onClick={() => {
            if (readOnly) return;

            setOpen((current) => {
              const nextOpen = !current;

              if (nextOpen) {
                onOpen?.();
              } else {
                setSearch("");
                setDropdownPosition(undefined);
              }

              return nextOpen;
            });
          }}
          className={classNames(
            getFieldClassName({
              small,
              error: Boolean(error),
              className,
              withRightIcon: true,
              disabled: readOnly,
            }),
            "flex items-center gap-1.5 text-left",
            multiple && selectedOptions.length > 1 && "h-auto min-h-10 py-2",
          )}
        >
          <span className="flex min-w-0 flex-1 flex-wrap items-center gap-1">
            {selectedOptions.length ? (
              multiple ? (
                <>
                  {(typeof maxVisibleSelected === "number"
                    ? selectedOptions.slice(0, Math.max(0, maxVisibleSelected))
                    : selectedOptions).map((option) => (
                    <span
                      key={option.value}
                       className="inline-flex max-w-full items-center gap-1 rounded-md border border-[color:var(--ui-border)] bg-[color:var(--ui-surface-muted)] px-1.5 py-0.5 text-xs text-[color:var(--ui-text-muted)]"
                     >
                      <span className="truncate">{option.label}</span>

                      {!readOnly && (
                        <button
                          type="button"
                          onClick={(event) => {
                            event.stopPropagation();
                            removeValue(option.value);
                          }}
                           className="rounded text-[color:var(--ui-text-soft)] hover:text-[color:var(--ui-text)]"
                        >
                          <FiX size={12} />
                        </button>
                      )}
                    </span>
                  ))}

                  {typeof maxVisibleSelected === "number" && selectedOptions.length > maxVisibleSelected && (
                     <span className="inline-flex items-center rounded-md border border-[color:var(--ui-border)] bg-[color:var(--ui-surface-muted)] px-1.5 py-0.5 text-xs text-[color:var(--ui-text-soft)]">
                      +{selectedOptions.length - maxVisibleSelected}
                    </span>
                  )}
                </>
              ) : (
                <span className="truncate">{selectedOptions[0]?.label}</span>
              )
            ) : (
               <span className="text-[color:var(--ui-text-soft)]">{placeholder}</span>
             )}
           </span>

            <span className="flex items-center gap-1">
             {selectedOptions.length > 0 && !readOnly && (
              <button
                type="button"
                onClick={(event) => {
                  event.stopPropagation();
                  clearValue();
                }}
                 className="rounded p-0.5 text-[color:var(--ui-text-soft)] hover:text-[color:var(--ui-text)]"
               >
                <FiX size={14} />
              </button>
             )}

            <FiChevronDown
              size={small ? 16 : 18}
               className={classNames("text-[color:var(--ui-text-soft)] transition-transform", open && "rotate-180")}
             />
           </span>
         </button>

         {open && (
           <div
             className="fixed z-50 flex flex-col overflow-hidden rounded-xl border border-[color:var(--ui-border)] bg-[color:var(--ui-panel)] shadow-lg"
             style={{
              top: dropdownPosition?.top,
              bottom: dropdownPosition?.bottom,
              left: dropdownPosition?.left ?? 0,
              width: dropdownPosition?.width ?? 0,
              maxHeight: dropdownPosition?.maxHeight ?? 0,
              visibility: dropdownPosition ? "visible" : "hidden",
            }}
          >
            {dropdownPosition?.direction === "top" ? (
              <>
                {optionsList}
                {searchInput}
              </>
            ) : (
              <>
                {searchInput}
                {optionsList}
              </>
            )}
          </div>
        )}
      </div>
    </FormField>
  );
};
