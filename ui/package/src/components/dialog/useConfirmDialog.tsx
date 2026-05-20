import type { ReactNode } from "react";
import { Button } from "../actions/Button";
import { deferPromise } from "../../utils/deferPromise";
import type { DialogSize } from "./types";
import { useModal } from "./useModal";

export type ConfirmDialogOptions = {
  title?: ReactNode;
  description?: ReactNode;
  content?: ReactNode;
  confirmLabel?: ReactNode;
  cancelLabel?: ReactNode;
  size?: DialogSize;
};

export const useConfirmDialog = () => {
  const dialog = useModal();

  return (options: ConfirmDialogOptions) => {
    const deferred = deferPromise<boolean>();
    let settled = false;

    const settle = (value: boolean, close: () => void) => {
      settled = true;
      deferred.resolve(value);
      close();
    };

    dialog.open({
      title: options.title ?? "Confirmar",
      description: options.description,
      size: options.size ?? "sm",
      onClose: () => {
        if (!settled) {
          deferred.resolve(false);
        }
      },
      content: ({ close }) => (
        <div className="grid gap-4">
          {options.content && <div className="text-sm leading-6 text-ui-text-muted">{options.content}</div>}

          <div className="flex justify-end gap-2">
            <Button type="button" variant="secondary" onClick={() => settle(false, close)}>
              {options.cancelLabel ?? "Cancelar"}
            </Button>
            <Button type="button" onClick={() => settle(true, close)}>
              {options.confirmLabel ?? "Confirmar"}
            </Button>
          </div>
        </div>
      ),
    });

    return deferred.promise;
  };
};
