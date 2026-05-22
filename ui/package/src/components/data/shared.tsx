import type { ReactNode } from "react";
import type { DataTableActionIcon, MaybeResolver } from "./types";

export const resolveValue = <TInput, TValue>(value: MaybeResolver<TInput, TValue> | undefined, input: TInput) => {
  if (typeof value === "function") {
    return (value as (input: TInput) => TValue)(input);
  }

  return value;
};

export const renderActionIcon = (icon: DataTableActionIcon | undefined, className: string): ReactNode => {
  if (!icon) return null;

  if (typeof icon === "string") {
    return <img src={icon} alt="" className={className} />;
  }

  const Icon = icon;
  return <Icon className={className} />;
};
