import { useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { FiMoreHorizontal } from "react-icons/fi";
import { renderActionIcon } from "./shared";
import type { DataTableActionIcon, DataTableActionVariant } from "./types";

export type ActionMenuItem = {
  icon?: DataTableActionIcon;
  label: React.ReactNode;
  onClick: () => void | Promise<void>;
  disabled?: boolean;
  variant?: DataTableActionVariant;
};

export const ActionMenu = ({ label = "Opciones", items }: { label?: string; items: ActionMenuItem[] }) => {
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0 });
  const buttonRef = useRef<HTMLButtonElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const visibleItems = items.filter(Boolean);

  const updatePosition = () => {
    const rect = buttonRef.current?.getBoundingClientRect();
    if (!rect) return;

    const width = 176;
    const estimatedHeight = visibleItems.length * 36 + 8;
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
      if (!menuRef.current?.contains(target) && !buttonRef.current?.contains(target)) {
        setOpen(false);
      }
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
  }, [open]);

  if (visibleItems.length === 0) return null;

  return (
    <div className="relative inline-flex">
      <button
        ref={buttonRef}
        type="button"
        title={label || "Opciones"}
        className="inline-flex size-8 items-center justify-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
        onClick={() => {
          updatePosition();
          setOpen((value) => !value);
        }}
      >
        <FiMoreHorizontal size={16} />
      </button>

      {open &&
        createPortal(
          <div
            ref={menuRef}
            className="fixed z-[70] min-w-44 overflow-hidden rounded-xl bg-ui-panel py-1 shadow-xl ring-1 ring-ui-border/60"
            style={{ top: position.top, left: position.left }}
          >
          {visibleItems.map((item, index) => (
            <button
              key={index}
              type="button"
              disabled={item.disabled}
              className={[
                "flex w-full items-center gap-2 px-3 py-2 text-left text-sm transition-colors disabled:cursor-not-allowed disabled:opacity-50",
                item.variant === "danger"
                  ? "text-ui-danger hover:bg-ui-danger/10"
                  : "text-ui-text-muted hover:bg-ui-surface-hover hover:text-ui-text",
              ].join(" ")}
              onClick={async () => {
                if (item.disabled) return;
                await item.onClick();
                setOpen(false);
              }}
            >
              {renderActionIcon(item.icon, "size-4 shrink-0 object-contain")}
              <span className="truncate">{item.label}</span>
            </button>
          ))}
          </div>,
          document.body,
        )}
    </div>
  );
};
