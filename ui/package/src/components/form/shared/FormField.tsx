import type React from "react";
import { FiInfo, FiRefreshCw } from "react-icons/fi";
import classNames from "classnames";

type FormFieldProps = {
  id: string;
  label?: string;
  info?: string;
  error?: string;
  required?: boolean;
  onRefresh?: () => Promise<void> | void;
  loading?: boolean;
  children: React.ReactNode;
};

export const FormField = ({
  id,
  label,
  info,
  error,
  required = true,
  onRefresh,
  loading = false,
  children,
}: FormFieldProps) => {
  return (
    <div className="grid gap-1.5">
      {label && (
        <div className="flex items-center gap-1.5">
          <label htmlFor={id} className="text-sm font-medium text-ui-text">
            {label}
          </label>

          {required && <span className="text-sm font-medium text-ui-danger">*</span>}

          {info && (
            <span className="group relative inline-flex">
              <button
                type="button"
                aria-label={info}
                className="inline-flex h-4 w-4 items-center justify-center rounded-full text-ui-text-soft hover:text-ui-text focus:outline-none focus:ring-2 focus:ring-ui-focus"
              >
                <FiInfo size={14} />
              </button>

              <span
                role="tooltip"
                className={classNames(
                  "pointer-events-none absolute left-1/2 top-full z-50 mt-1 hidden w-max max-w-64 -translate-x-1/2 rounded-md px-2 py-1.5 text-xs font-normal leading-relaxed shadow-md",
                  "bg-ui-tooltip text-ui-tooltip-text",
                  "group-hover:block group-focus-within:block",
                )}
              >
                {info}
              </span>
            </span>
          )}

          {onRefresh && (
            <button
                type="button"
                onClick={onRefresh}
                aria-label="Actualizar"
                className="inline-flex h-4 w-4 items-center justify-center rounded-full text-ui-text-soft hover:text-ui-text focus:outline-none focus:ring-2 focus:ring-ui-focus"
                disabled={loading}
              >
              <FiRefreshCw size={14} />
            </button>
          )}
        </div>
      )}

      {children}

      {error && <p className="text-xs text-ui-danger">{error}</p>}
    </div>
  );
};
