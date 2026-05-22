import { createContext, useMemo, useState } from "react";
import type { PropsWithChildren, ReactNode } from "react";
import { createPortal } from "react-dom";
import { FiAlertCircle, FiCheckCircle, FiInfo, FiLoader } from "react-icons/fi";
import type { ToastContextValue, ToastId, ToastOptions, ToastType } from "./types";

type ToastItem = {
  id: ToastId;
  message: ReactNode;
  type: ToastType;
};

export const ToastContext = createContext<ToastContextValue | null>(null);

const iconByType: Record<ToastType, ReactNode> = {
  loading: <FiLoader className="animate-spin" size={16} />,
  success: <FiCheckCircle size={16} />,
  error: <FiAlertCircle size={16} />,
  info: <FiInfo size={16} />,
};

const typeClassName: Record<ToastType, string> = {
  loading: "text-ui-text-muted",
  success: "text-ui-primary",
  error: "text-ui-danger",
  info: "text-ui-accent",
};

const createToastId = () => `toast-${Date.now()}-${Math.random().toString(36).slice(2)}`;

export const ToastProvider = ({ children }: PropsWithChildren) => {
  const [toasts, setToasts] = useState<ToastItem[]>([]);

  const dismiss = (id: ToastId) => {
    setToasts((current) => current.filter((toast) => toast.id !== id));
  };

  const show = (message: ReactNode, type: ToastType = "info", options?: ToastOptions) => {
    const id = options?.id ?? createToastId();

    setToasts((current) => {
      const next = current.filter((toast) => toast.id !== id);
      return [...next, { id, message, type }];
    });

    if (type !== "loading") {
      window.setTimeout(() => dismiss(id), options?.duration ?? 3500);
    }

    return id;
  };

  const value = useMemo<ToastContextValue>(
    () => ({
      show,
      loading: (message, options) => show(message, "loading", options),
      success: (message, options) => show(message, "success", options),
      error: (message, options) => show(message, "error", options),
      dismiss,
    }),
    [],
  );

  return (
    <ToastContext.Provider value={value}>
      {children}

      {createPortal(
        <div className="pointer-events-none fixed inset-x-0 top-3 z-[60] grid justify-items-center gap-2 px-3 sm:inset-x-auto sm:right-4 sm:justify-items-end">
          {toasts.map((toast) => (
            <div
              key={toast.id}
              className="pointer-events-auto flex min-h-10 w-full max-w-sm items-center gap-2 rounded-xl bg-ui-panel px-3 py-2 text-sm font-medium text-ui-text shadow-xl ring-1 ring-ui-border/50"
            >
              <span className={["grid size-5 shrink-0 place-items-center", typeClassName[toast.type]].join(" ")}>{iconByType[toast.type]}</span>
              <div className="min-w-0 flex-1 leading-5">{toast.message}</div>
            </div>
          ))}
        </div>,
        document.body,
      )}
    </ToastContext.Provider>
  );
};
