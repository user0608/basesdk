import type { ButtonHTMLAttributes, PropsWithChildren } from "react";
import { FiLoader } from "react-icons/fi";
import { useHasAllPermissions } from "../../platform/permissions";

export type ButtonProps = PropsWithChildren<ButtonHTMLAttributes<HTMLButtonElement>> & {
  variant?: "primary" | "secondary" | "danger";
  loading?: boolean;
  permissions?: readonly string[];
};

const baseClassName =
  "inline-flex items-center justify-center gap-2 rounded-xl px-4 py-2 text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-ui-focus disabled:cursor-not-allowed disabled:opacity-60";

const variantClassName: Record<NonNullable<ButtonProps["variant"]>, string> = {
  primary: "bg-ui-primary text-ui-text-inverse shadow-sm hover:bg-ui-primary-hover",
  secondary:
    "bg-ui-surface text-ui-text shadow-sm ring-1 ring-inset ring-ui-border/70 hover:bg-ui-surface-hover",
  danger: "bg-ui-danger text-ui-text-inverse shadow-sm hover:bg-ui-danger/90",
};

export function Button({ children, className = "", variant = "primary", ...props }: ButtonProps) {
  const { loading = false, disabled, permissions, ...restProps } = props;
  const canRender = useHasAllPermissions(permissions);

  if (!canRender) return null;

  return (
    <button
      className={`${baseClassName} ${variantClassName[variant]} ${className}`.trim()}
      disabled={disabled || loading}
      {...restProps}
    >
      {loading && <FiLoader size={16} className="animate-spin" aria-hidden="true" />}
      {children}
    </button>
  );
}
