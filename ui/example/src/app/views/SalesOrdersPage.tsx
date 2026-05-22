import { DataTable, useFormModal } from "@basesdk/ui";
import type { ColumnDef } from "@tanstack/react-table";
import { FiEdit2, FiPlus } from "react-icons/fi";
import type { appRegistry } from "../registry";

type SalesOrder = {
  id: string;
  number: string;
  customer: string;
  amount: string;
  status: "Abierta" | "Facturada" | "Cancelada";
};

const orders: SalesOrder[] = Array.from({ length: 18 }).map((_, index) => ({
  id: `order-${index + 1}`,
  number: `SO-${String(index + 1).padStart(4, "0")}`,
  customer: index % 2 === 0 ? "Comercial Rivera" : "Distribuidora Norte",
  amount: `$${(350 + index * 42).toFixed(2)}`,
  status: index % 6 === 0 ? "Cancelada" : index % 3 === 0 ? "Facturada" : "Abierta",
}));

const columns: ColumnDef<SalesOrder, unknown>[] = [
  { accessorKey: "number", header: "Orden", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.number}</span> },
  { accessorKey: "customer", header: "Cliente" },
  { accessorKey: "amount", header: "Monto" },
  { accessorKey: "status", header: "Estado" },
];

export const SalesOrdersPage = () => {
  const formModal = useFormModal<typeof appRegistry>();

  return (
    <DataTable
      title="Ventas"
      data={orders}
      columns={columns}
      getRowId={(row) => row.id}
      searchKeys={["number", "customer", "status"]}
      actions={[
        {
          icon: FiPlus,
          label: "Nueva orden",
          onClick: () => formModal.open("sales-order-form", { title: "Nueva orden" }),
        },
      ]}
      rowOptions={[
        {
          icon: FiEdit2,
          label: "Editar",
          onClick: (row) => formModal.open("sales-order-form", { title: "Editar orden", props: { orderId: row.number } }),
        },
      ]}
      pageSize={8}
    />
  );
};
