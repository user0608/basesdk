import { createContext, useCallback, useEffect, useMemo, useState } from "react";
import type { PropsWithChildren } from "react";
import { createPortal } from "react-dom";
import { FiX } from "react-icons/fi";
import type { DialogContextValue, DialogOptions, DialogSize } from "./types";

export const DialogContext = createContext<DialogContextValue | null>(null);

const sizeClassName: Record<DialogSize, string> = {
  sm: "sm:max-w-md",
  md: "sm:max-w-xl",
  lg: "sm:max-w-3xl",
  xl: "sm:max-w-5xl",
  fullscreen: "sm:max-w-[calc(100vw-2rem)] sm:min-h-[calc(100vh-2rem)]",
};

export const DialogProvider = ({ children }: PropsWithChildren) => {
  const [dialog, setDialog] = useState<DialogOptions | null>(null);

  const close = useCallback(() => {
    setDialog((current) => {
      current?.onClose?.();
      return null;
    });
  }, []);

  useEffect(() => {
    if (!dialog) return;

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        close();
      }
    };

    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [dialog]);

  const value = useMemo<DialogContextValue>(
    () => ({
      open: setDialog,
      close,
      isOpen: Boolean(dialog),
    }),
    [close, dialog],
  );

  const content = typeof dialog?.content === "function" ? dialog.content({ close }) : dialog?.content;
  const canCloseOnOverlay = dialog?.closeOnOverlayClick ?? true;
  const panelSize = dialog?.size ?? "md";

  return (
    <DialogContext.Provider value={value}>
      {children}

      {dialog &&
        createPortal(
          <div className="fixed inset-0 z-50 grid items-end bg-ui-text/35 p-0 backdrop-blur-sm sm:items-center sm:p-4">
            <button
              type="button"
              aria-label="Cerrar dialogo"
              className="absolute inset-0 cursor-default"
              onClick={canCloseOnOverlay ? close : undefined}
            />

            <section
              role="dialog"
              aria-modal="true"
              className={[
                "relative z-10 max-h-[92vh] w-full overflow-hidden rounded-t-2xl bg-ui-panel shadow-2xl ring-1 ring-ui-border/40 sm:mx-auto sm:rounded-2xl",
                sizeClassName[panelSize],
              ].join(" ")}
            >
              {(dialog.title || dialog.description) && (
                <header className="flex items-start justify-between gap-3 border-b border-ui-border/30 px-4 py-3">
                  <div className="min-w-0">
                    {dialog.title && <h2 className="truncate text-base font-semibold text-ui-text">{dialog.title}</h2>}
                    {dialog.description && <p className="mt-1 text-sm leading-5 text-ui-text-muted">{dialog.description}</p>}
                  </div>

                  <button
                    type="button"
                    aria-label="Cerrar dialogo"
                    className="grid size-8 shrink-0 place-items-center rounded-lg text-ui-text-muted transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
                    onClick={close}
                  >
                    <FiX size={18} />
                  </button>
                </header>
              )}

              {!dialog.title && !dialog.description && (
                <button
                  type="button"
                  aria-label="Cerrar dialogo"
                  className="absolute right-3 top-3 z-10 grid size-8 place-items-center rounded-lg text-ui-text-muted transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
                  onClick={close}
                >
                  <FiX size={18} />
                </button>
              )}

              <div className="max-h-[calc(92vh-57px)] overflow-auto p-4">{content}</div>
            </section>
          </div>,
          document.body,
        )}
    </DialogContext.Provider>
  );
};
