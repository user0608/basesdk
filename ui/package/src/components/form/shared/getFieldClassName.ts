import classNames from "classnames";

export const getFieldClassName = ({
  small,
  error,
  className,
  withRightIcon,
  disabled,
  readOnly,
}: {
  small?: boolean;
  error?: boolean;
  className?: string;
  withRightIcon?: boolean;
  disabled?: boolean;
  readOnly?: boolean;
}) =>
  classNames(
    "w-full rounded-lg border text-sm shadow-sm outline-none transition-colors",
    "text-[color:var(--ui-text)] placeholder:text-[color:var(--ui-text-soft)]",
    "border-[color:var(--ui-border)] bg-[color:var(--ui-surface)]",
    "focus:border-[color:var(--ui-accent)] focus:ring-2 focus:ring-[color:var(--ui-focus-ring)]",
    "disabled:cursor-not-allowed disabled:opacity-60",
    small ? "h-8 px-2" : "h-10 px-3",
    withRightIcon && "pr-10",
    disabled || readOnly ? "bg-[color:var(--ui-surface-muted)] text-[color:var(--ui-text-soft)]" : "",
    error
      ? "border-[color:var(--ui-border-error)] focus:border-[color:var(--ui-border-error)] focus:ring-[color:var(--ui-danger-ring)]"
      : "",
    className,
  );
