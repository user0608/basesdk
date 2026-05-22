import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { Link, useParams } from "react-router-dom";
import { FiArrowLeft, FiEdit2, FiPlus, FiShield, FiTrash2 } from "react-icons/fi";
import { DataTable } from "../../components/data/DataTable";
import { useFormModal } from "../../platform/useFormModal";
import { useMutate } from "../../query/useMutate";
import type { PropertyResponse, TenantPropertyResponse } from "../types";
import { SystemBreadcrumbs } from "../SystemBreadcrumbs";
import { useSystemService } from "../useSystemService";

type PropertyRow = PropertyResponse | TenantPropertyResponse;

export const SystemPropertiesPage = ({ tenantScoped = false }: { tenantScoped?: boolean }) => {
  const { tenantCodigo = "" } = useParams();
  const system = useSystemService();
  const formModal = useFormModal();
  const service = tenantScoped ? system.tenantProperties(tenantCodigo) : system.properties;
  const queryKey = tenantScoped ? ["system", "tenants", tenantCodigo, "properties"] : ["system", "properties"];
  const properties = useQuery({ queryKey, queryFn: ({ signal }) => service.list(signal) });
  const remove = useMutate({
    mutationFn: (keys: string[]) => Promise.all(keys.map((key) => service.delete(key))),
    invalidateQueryKey: queryKey,
    requireConfirm: true,
    confirmMessage: "Quieres eliminar las propiedades seleccionadas?",
    confirmLabel: "Eliminar",
    successMessage: "Propiedades eliminadas.",
  });
  const columns: ColumnDef<PropertyRow, unknown>[] = [
    { accessorKey: "key", header: "Key", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.key}</span> },
    { accessorKey: "dataType", header: "Tipo" },
    { accessorKey: "value", header: "Valor", cell: ({ row }) => <span className="block max-w-xl truncate">{row.original.value}</span> },
    {
      accessorKey: "description",
      header: "Descripcion",
      cell: ({ row }) => row.original.description ?? <span className="text-ui-text-soft">Sin descripcion</span>,
    },
  ];

  return (
    <main className="flex h-screen min-h-0 flex-col overflow-hidden bg-ui-panel-muted p-3 text-ui-text lg:p-4">
      <div className="mb-3 flex flex-wrap items-center gap-2">
        <div className="flex min-w-0 items-center gap-2">
          <Link
            to={tenantScoped ? "/system/tenants" : "/system"}
            aria-label={tenantScoped ? "Volver a tenants" : "Volver a system"}
            title={tenantScoped ? "Volver a tenants" : "Volver a system"}
            className="inline-grid size-9 shrink-0 place-items-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
          >
            <FiArrowLeft size={18} />
          </Link>
          <SystemBreadcrumbs
            items={
              tenantScoped
                ? [
                    { label: "System", to: "/system" },
                    { label: "Tenants", to: "/system/tenants" },
                    { label: tenantCodigo, to: `/system/tenants/${encodeURIComponent(tenantCodigo)}/security/users` },
                    { label: "Properties" },
                  ]
                : [{ label: "System", to: "/system" }, { label: "Properties" }]
            }
          />
        </div>
        {tenantScoped && (
          <Link
            to={`/system/tenants/${encodeURIComponent(tenantCodigo)}/security/users`}
            aria-label="Abrir seguridad"
            title="Abrir seguridad"
            className="inline-grid size-9 place-items-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
          >
            <FiShield size={18} />
          </Link>
        )}
      </div>
      <div className="min-h-0 flex-1">
        <DataTable
          tableId={tenantScoped ? `system.tenant.${tenantCodigo}.properties` : "system.properties"}
          title={tenantScoped ? "Properties del tenant" : "Properties system"}
          description={tenantScoped ? `Configuracion system del tenant ${tenantCodigo}.` : "Configuracion global de la plataforma."}
          data={properties.data ?? []}
          columns={columns}
          loading={properties.isLoading}
          emptyMessage={properties.isError ? "No se pudieron cargar las propiedades" : "No hay propiedades registradas"}
          selectable
          getRowId={(row) => row.key}
          searchPlaceholder="Buscar properties..."
          searchKeys={["key", "value", "dataType", "description"]}
          pageSize={20}
          maxItemsRange={100}
          actions={[
            {
              icon: FiPlus,
              label: "Nueva",
              onClick: () => formModal.open("system-property-form", { title: "Nueva propiedad", size: "md", props: { tenantCodigo: tenantScoped ? tenantCodigo : undefined } }),
            },
          ]}
          options={[
            {
              icon: FiTrash2,
              label: "Eliminar",
              variant: "danger",
              disabled: ({ selectedRows }) => selectedRows.length === 0,
              onClick: ({ selectedRows, clearSelection }) => remove.mutate(selectedRows.map((row) => row.key)).then(clearSelection),
            },
          ]}
          rowOptions={[
            {
              icon: FiEdit2,
              label: "Editar",
              onClick: (row) => formModal.open("system-property-form", { title: "Editar propiedad", description: row.key, size: "md", props: { property: row, tenantCodigo: tenantScoped ? tenantCodigo : undefined } }),
            },
            { icon: FiTrash2, label: "Eliminar", variant: "danger", onClick: (row) => remove.mutate([row.key]) },
          ]}
        />
      </div>
    </main>
  );
};
