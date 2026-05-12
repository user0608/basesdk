export type SelectOption<TExtra extends object = object> = {
  label: string;
  value: string;
} & TExtra;
