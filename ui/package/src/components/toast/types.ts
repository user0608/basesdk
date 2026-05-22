import type { ReactNode } from "react";

export type ToastType = "loading" | "success" | "error" | "info";

export type ToastId = string;

export type ToastOptions = {
  id?: ToastId;
  duration?: number;
};

export type ToastContextValue = {
  show: (message: ReactNode, type?: ToastType, options?: ToastOptions) => ToastId;
  loading: (message: ReactNode, options?: ToastOptions) => ToastId;
  success: (message: ReactNode, options?: ToastOptions) => ToastId;
  error: (message: ReactNode, options?: ToastOptions) => ToastId;
  dismiss: (id: ToastId) => void;
};
