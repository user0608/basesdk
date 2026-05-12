import type React from "react";
import type { FieldValues, Path, UseFormReturn } from "react-hook-form";
import type { IconType } from "react-icons";
import { Checkbox } from "./input/Checkbox";
import { BooleanSelect } from "./input/BooleanSelect";
import { TextInput } from "./input/TextInput";

type BaseInputFieldProps<TFormValues extends FieldValues> = {
  form: UseFormReturn<TFormValues>;
  name: Path<TFormValues>;
  label?: string;
  info?: string;
  required?: boolean;
  readOnly?: boolean;
  onRefresh?: () => Promise<void> | void;
};

type TextFieldProps<TFormValues extends FieldValues> = BaseInputFieldProps<TFormValues> & {
  variant?: "text";
  type?: React.HTMLInputTypeAttribute;
  className?: string;
  placeholder?: string;
  onChange?: (value: string) => void;
  icon?: IconType;
  small?: boolean;
  onIconClick?: () => Promise<void> | void;
};

type CheckboxFieldProps<TFormValues extends FieldValues> = BaseInputFieldProps<TFormValues> & {
  variant: "checkbox";
  className?: string;
  onChange?: (value: boolean) => void;
};

type BooleanFieldProps<TFormValues extends FieldValues> = BaseInputFieldProps<TFormValues> & {
  variant: "boolean";
  placeholder?: string;
  small?: boolean;
  onChange?: (value: boolean) => void;
  yesLabel?: string;
  noLabel?: string;
};

export type InputFieldProps<TFormValues extends FieldValues> =
  | TextFieldProps<TFormValues>
  | CheckboxFieldProps<TFormValues>
  | BooleanFieldProps<TFormValues>;

export const InputField = <TFormValues extends FieldValues>(props: InputFieldProps<TFormValues>) => {
  switch (props.variant) {
    case "checkbox":
      return <Checkbox {...props} />;
    case "boolean":
      return <BooleanSelect {...props} />;
    case "text":
    case undefined:
      return <TextInput {...props} />;
    default:
      return null;
  }
};
