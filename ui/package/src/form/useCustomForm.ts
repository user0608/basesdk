import { yupResolver } from "@hookform/resolvers/yup";
import type { DefaultValues } from "react-hook-form";
import { useForm } from "react-hook-form";
import * as yup from "yup";

const isDevelopment = () => {
  const runtime = globalThis as typeof globalThis & {
    process?: {
      env?: {
        NODE_ENV?: string;
      };
    };
  };

  return !runtime.process?.env?.NODE_ENV || runtime.process.env.NODE_ENV === "development";
};

export const useCustomForm = <T extends yup.ObjectSchema<any>>(
  schema: T,
  props?: {
    defaultValues?:
      | DefaultValues<yup.Asserts<T>>
      | ((payload?: unknown) => Promise<yup.Asserts<T>>);
  },
) => {
  const form = useForm<yup.Asserts<T>>({
    resolver: yupResolver(schema),
    defaultValues: props?.defaultValues,
  });

  type HandleType = typeof form.handleSubmit;

  const newHandleSubmit: HandleType = (onSuccess) => {
    return (event) => {
      return form.handleSubmit(onSuccess, (error) => {
        if (isDevelopment()) {
          console.log(error);
        }
      })(event);
    };
  };

  return { ...form, handleSubmit: newHandleSubmit };
};
