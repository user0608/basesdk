import { useEffect, useMemo, useRef, useState } from "react";
import { FiChevronLeft, FiChevronRight } from "react-icons/fi";
import type { PaginationRangeData, PaginationRangeInputProps } from "./types";

type Page = {
  limit: number;
  offset: number;
};

const rangePattern = /^[1-9]\d*\s*-\s*[1-9]\d*$/;

const buildRangeValue = ({ limit, offset }: Page, totalItems: number) => {
  if (totalItems <= 0) return "0-0";

  const start = Math.min(offset + 1, totalItems);
  const end = Math.min(offset + limit, totalItems);

  return `${start}-${end}`;
};

export const PaginationRangeInput = ({ pageData, maxItemsRange, onChange }: PaginationRangeInputProps) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const navigationRangeSizeRef = useRef(
    Math.max(1, Math.min(pageData.limitItems || 1, Math.max(1, maxItemsRange || 1))),
  );

  const [draftValue, setDraftValue] = useState<string | null>(null);
  const safeMaxItemsRange = Math.max(1, maxItemsRange || 1);

  const pageValue = useMemo(() => {
    return buildRangeValue(
      {
        limit: pageData.limitItems,
        offset: pageData.offsetItems,
      },
      pageData.totalItems,
    );
  }, [pageData.limitItems, pageData.offsetItems, pageData.totalItems]);

  const value = draftValue ?? pageValue;

  useEffect(() => {
    const visibleCount = Math.max(0, Math.min(pageData.limitItems, Math.max(pageData.totalItems - pageData.offsetItems, 0)));

    const isPartialBoundaryPage =
      pageData.totalItems > 0 &&
      visibleCount > 0 &&
      visibleCount < navigationRangeSizeRef.current &&
      (pageData.offsetItems === 0 || pageData.offsetItems + visibleCount >= pageData.totalItems);

    if (!isPartialBoundaryPage && pageData.limitItems > 0) {
      navigationRangeSizeRef.current = Math.max(1, Math.min(pageData.limitItems, safeMaxItemsRange));
    }
  }, [pageData.limitItems, pageData.offsetItems, pageData.totalItems, safeMaxItemsRange]);

  useEffect(() => {
    if (!inputRef.current) return;
    inputRef.current.style.width = `${(value.length + 1) * 6.3 + 8}px`;
  }, [value]);

  const getEnsuredValue = (target: string): Page & { value: string } => {
    if (pageData.totalItems <= 0) {
      return {
        value: "0-0",
        offset: 0,
        limit: 0,
      };
    }

    let finalValue = rangePattern.test(target) ? target.replace(/\s+/g, "") : pageValue;
    let [left, right] = finalValue.split("-").map(Number);

    if (Number.isNaN(left) || Number.isNaN(right)) {
      finalValue = pageValue;
      [left, right] = finalValue.split("-").map(Number);
    }

    if (left < 1) left = 1;

    if (left > pageData.totalItems) {
      finalValue = pageValue;
      [left, right] = finalValue.split("-").map(Number);
    } else {
      if (right < left) right = left;

      const requestedSize = right - left + 1;

      if (requestedSize > safeMaxItemsRange) {
        right = left + safeMaxItemsRange - 1;
      }

      if (right > pageData.totalItems) {
        right = pageData.totalItems;
      }

      finalValue = `${left}-${right}`;
    }

    return {
      value: finalValue,
      offset: left - 1,
      limit: right - left + 1,
    };
  };

  const commitValue = () => {
    const cleanValue = value.trim().replace(/ /g, "");
    const { value: ensuredValue, limit, offset } = getEnsuredValue(cleanValue);

    setDraftValue(null);

    if (pageValue === ensuredValue) return;

    onChange?.(limit, offset);
  };

  const previousRange = () => {
    if (pageData.totalItems <= 0) return;

    const pageSize = Math.max(1, Math.min(navigationRangeSizeRef.current, safeMaxItemsRange));
    const currentStart = pageData.offsetItems + 1;
    const newEnd = currentStart - 1;

    if (newEnd < 1) return;

    const newStart = Math.max(newEnd - pageSize + 1, 1);
    onChange?.(newEnd - newStart + 1, newStart - 1);
  };

  const nextRange = () => {
    if (pageData.totalItems <= 0) return;

    const pageSize = Math.max(1, Math.min(navigationRangeSizeRef.current, safeMaxItemsRange));
    const currentEnd = Math.min(pageData.offsetItems + pageData.limitItems, pageData.totalItems);
    const newStart = currentEnd + 1;

    if (newStart > pageData.totalItems) return;

    const newEnd = Math.min(currentEnd + pageSize, pageData.totalItems);
    onChange?.(newEnd - newStart + 1, newStart - 1);
  };

  const canPrevious = pageData.totalItems > 0 && pageData.offsetItems > 0;
  const canNext = pageData.totalItems > 0 && pageData.offsetItems + pageData.limitItems < pageData.totalItems;

  return (
    <div className="flex select-none items-center justify-end gap-1">
      <div className="flex items-center gap-1 rounded-lg border border-ui-border/60 bg-ui-surface px-2 py-1">
        <input
          ref={inputRef}
          className="h-5 min-w-10 bg-transparent p-0 text-right text-xs font-medium text-ui-text outline-none focus:text-ui-primary"
          onChange={({ target }) => setDraftValue(target.value)}
          onBlur={commitValue}
          onKeyDown={({ key }) => {
            if (key !== "Enter") return;
            commitValue();
            inputRef.current?.blur();
          }}
          value={value}
          type="text"
          autoCorrect="off"
          autoComplete="off"
        />
        <span className="text-xs text-ui-text-soft">/</span>
        <span className="text-xs font-medium text-ui-text-muted">{pageData.totalItems}</span>
      </div>

      <button
        type="button"
        disabled={!canPrevious}
        className="grid size-7 place-items-center rounded-lg border border-ui-border/60 text-ui-text-muted transition-colors hover:bg-ui-surface-hover hover:text-ui-text disabled:cursor-not-allowed disabled:opacity-50"
        onClick={previousRange}
      >
        <FiChevronLeft size={16} />
      </button>

      <button
        type="button"
        disabled={!canNext}
        className="grid size-7 place-items-center rounded-lg border border-ui-border/60 text-ui-text-muted transition-colors hover:bg-ui-surface-hover hover:text-ui-text disabled:cursor-not-allowed disabled:opacity-50"
        onClick={nextRange}
      >
        <FiChevronRight size={16} />
      </button>
    </div>
  );
};
