import { DataTable, useFormModal, useToast } from "@basesdk/ui";
import type { ColumnDef } from "@tanstack/react-table";
import { FiEdit2, FiPackage, FiPlus, FiRefreshCw } from "react-icons/fi";
import type { appRegistry } from "../registry";

type InventoryItem = {
  id: string;
  sku: string;
  name: string;
  warehouse: string;
  stock: number;
  reserved: number;
  cost: string;
  status: "Disponible" | "Bajo" | "Agotado";
};

const items: InventoryItem[] = Array.from({ length: 30 }).map((_, index) => ({
  id: `item-${index + 1}`,
  sku: `SKU-${String(index + 1).padStart(4, "0")}`,
  name: `Producto ${index + 1}`,
  warehouse: index % 2 === 0 ? "Almacen central" : "Almacen norte",
  stock: 120 - index * 3,
  reserved: index % 7,
  cost: `$${(15 + index * 1.7).toFixed(2)}`,
  status: index % 9 === 0 ? "Agotado" : index % 4 === 0 ? "Bajo" : "Disponible",
}));

const columns: ColumnDef<InventoryItem, unknown>[] = [
  { accessorKey: "sku", header: "SKU", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.sku}</span> },
  { accessorKey: "name", header: "Producto" },
  { accessorKey: "warehouse", header: "Almacen" },
  { accessorKey: "stock", header: "Stock" },
  { accessorKey: "reserved", header: "Reservado" },
  { accessorKey: "cost", header: "Costo" },
  { accessorKey: "status", header: "Estado" },
];

export const InventoryItemsPage = () => {
  const formModal = useFormModal<typeof appRegistry>();
  const toast = useToast();

  return (
    <DataTable
      title="Inventario"
      description="Ejemplo de modulo externo agregado sobre la base."
      data={items}
      columns={columns}
      selectable
      getRowId={(row) => row.id}
      searchKeys={["sku", "name", "warehouse", "status"]}
      actions={[
        {
          icon: FiPlus,
          label: "Nuevo",
          onClick: () => formModal.open("inventory-item-form", { title: "Nuevo producto" }),
        },
        {
          icon: FiRefreshCw,
          label: "Refrescar",
          onClick: () => toast.success("Inventario actualizado"),
        },
      ]}
      options={[
        {
          icon: FiPackage,
          label: "Ajuste masivo",
          disabled: ({ selectedRows }) => selectedRows.length === 0,
          onClick: ({ selectedRows }) => toast.success(`${selectedRows.length} producto(s) seleccionados`),
        },
      ]}
      rowOptions={[
        {
          icon: FiEdit2,
          label: "Editar",
          onClick: (row) => formModal.open("inventory-item-form", { title: "Editar producto", props: { itemId: row.sku } }),
        },
      ]}
      pageSize={8}
      maxItemsRange={20}
    />
  );
};
