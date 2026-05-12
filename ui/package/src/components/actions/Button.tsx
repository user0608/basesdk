import type { ButtonHTMLAttributes, PropsWithChildren } from "react";

export type ButtonProps = PropsWithChildren<ButtonHTMLAttributes<HTMLButtonElement>> & {
  variant?: "primary" | "secondary";
};

const baseClassName =
  "inline-flex items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-[color:var(--ui-focus-ring)] disabled:cursor-not-allowed disabled:opacity-60";

const variantClassName: Record<NonNullable<ButtonProps["variant"]>, string> = {
  primary: "bg-[color:var(--ui-primary)] text-[color:var(--ui-text-inverse)] hover:bg-[color:var(--ui-primary-hover)]",
  secondary:
    "border border-[color:var(--ui-border)] bg-[color:var(--ui-surface)] text-[color:var(--ui-text)] hover:bg-[color:var(--ui-surface-hover)]"
};

export function Button({ children, className = "", variant = "primary", ...props }: ButtonProps) {
  return (
    <button
      className={`${baseClassName} ${variantClassName[variant]} ${className}`.trim()}
      {...props}
    >
      {children}
    </button>
  );
}
