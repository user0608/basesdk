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
    "w-full rounded-xl border border-ui-border bg-ui-surface text-sm text-ui-text shadow-sm outline-none transition-colors",
    "placeholder:text-ui-text-soft focus:border-ui-accent focus:ring-2 focus:ring-ui-focus",
    "disabled:cursor-not-allowed disabled:opacity-60",
    small ? "h-8 px-2" : "h-10 px-3",
    withRightIcon && "pr-10",
    disabled || readOnly ? "bg-ui-surface-muted text-ui-text-soft" : "",
    error
      ? "border-ui-danger text-ui-text focus:border-ui-danger focus:ring-ui-danger/20"
      : "",
    className,
  );
