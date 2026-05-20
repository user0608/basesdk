import type { ReactNode } from "react";

export type DialogSize = "sm" | "md" | "lg" | "xl" | "fullscreen";

export type DialogApi = {
  close: () => void;
};

export type DialogContent = ReactNode | ((api: DialogApi) => ReactNode);

export type DialogOptions = {
  title?: ReactNode;
  description?: ReactNode;
  content: DialogContent;
  size?: DialogSize;
  closeOnOverlayClick?: boolean;
  onClose?: () => void;
};

export type DialogContextValue = DialogApi & {
  open: (options: DialogOptions) => void;
  isOpen: boolean;
};
