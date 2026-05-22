import { Button } from "@basesdk/ui";

export type InventoryItemFormProps = {
  itemId?: string;
  close?: () => void;
};

export const InventoryItemForm = ({ itemId, close }: InventoryItemFormProps) => {
  return (
    <form className="grid gap-3">
      <label className="grid gap-1.5 text-sm font-medium text-ui-text">
        Codigo
        <input
          defaultValue={itemId}
          className="rounded-xl border border-ui-border bg-ui-surface px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ui-focus"
          placeholder="SKU-001"
        />
      </label>
      <label className="grid gap-1.5 text-sm font-medium text-ui-text">
        Nombre
        <input className="rounded-xl border border-ui-border bg-ui-surface px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ui-focus" />
      </label>
      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button type="button" onClick={close}>
          Guardar
        </Button>
      </div>
    </form>
  );
};
