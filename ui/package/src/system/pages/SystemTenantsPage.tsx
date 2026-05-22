import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { Link, useNavigate } from "react-router-dom";
import { FiArrowLeft, FiEdit2, FiPlus, FiSettings, FiShield, FiSlash, FiUnlock } from "react-icons/fi";
import { DataTable } from "../../components/data/DataTable";
import { useFormModal } from "../../platform/useFormModal";
import { useMutate } from "../../query/useMutate";
import { useSecurityService } from "../../security/useSecurityService";
import type { TenantResponse } from "../../security/types";
import { SystemBreadcrumbs } from "../SystemBreadcrumbs";

const queryKey = ["system", "tenants"];

const CountChip = ({ label, value }: { label: string; value: number }) => (
  <span className="inline-flex items-center gap-1 rounded-full bg-ui-surface px-2 py-1 text-xs font-medium text-ui-text-muted ring-1 ring-inset ring-ui-border/55">
    <span className="font-semibold tabular-nums text-ui-text">{value}</span>
    <span>{label}</span>
  </span>
);

export const SystemTenantsPage = () => {
  const security = useSecurityService();
  const formModal = useFormModal();
  const navigate = useNavigate();
  const tenants = useQuery({ queryKey, queryFn: ({ signal }) => security.tenants.list(signal) });
  const enable = useMutate({ mutationFn: (codigos: string[]) => security.tenants.enable(codigos), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres habilitar los tenants seleccionados?", successMessage: "Tenants habilitados." });
  const disable = useMutate({ mutationFn: (codigos: string[]) => security.tenants.disable(codigos), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres deshabilitar los tenants seleccionados?", successMessage: "Tenants deshabilitados." });
  const columns: ColumnDef<TenantResponse, unknown>[] = [
    { accessorKey: "codigo", header: "Codigo", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.codigo}</span> },
    { accessorKey: "name", header: "Nombre" },
    { accessorKey: "timezone", header: "Zona horaria" },
    {
      id: "counts",
      header: "Resumen",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-1.5">
          <CountChip label="usuarios" value={row.original.usersCount} />
          <CountChip label="roles" value={row.original.rolesCount} />
          <CountChip label="grupos" value={row.original.groupsCount} />
        </div>
      ),
    },
    { accessorKey: "maxActiveUsers", header: "Max. usuarios" },
    { accessorKey: "expiresAt", header: "Expira", cell: ({ row }) => row.original.expiresAt ?? <span className="text-ui-text-soft">Sin fecha</span> },
    {
      accessorKey: "disabled",
      header: "Estado",
      cell: ({ row }) => <span className={row.original.disabled ? "text-ui-danger" : "text-ui-primary"}>{row.original.disabled ? "Inactivo" : "Activo"}</span>,
    },
  ];

  return (
    <main className="flex h-screen min-h-0 flex-col overflow-hidden bg-ui-panel-muted p-3 text-ui-text lg:p-4">
      <div className="mb-3 flex items-center justify-between gap-3">
        <div className="flex min-w-0 items-center gap-2">
          <Link
            to="/system"
            aria-label="Volver a system"
            title="Volver a system"
            className="inline-grid size-9 shrink-0 place-items-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
          >
            <FiArrowLeft size={18} />
          </Link>
          <SystemBreadcrumbs items={[{ label: "System", to: "/system" }, { label: "Tenants" }]} />
        </div>
      </div>
      <div className="min-h-0 flex-1">
        <DataTable
          tableId="system.tenants"
          title="Tenants"
          description="Clientes y entornos disponibles en la plataforma."
          data={tenants.data ?? []}
          columns={columns}
          loading={tenants.isLoading}
          emptyMessage={tenants.isError ? "No se pudieron cargar los tenants" : "No hay tenants registrados"}
          selectable
          getRowId={(row) => row.codigo}
          searchPlaceholder="Buscar tenants..."
          searchKeys={["codigo", "name", "timezone"]}
          pageSize={20}
          maxItemsRange={100}
          onRowClick={(row) => navigate(`/system/tenants/${encodeURIComponent(row.codigo)}/security/users`)}
          actions={[
            {
              icon: FiPlus,
              label: "Nuevo",
              onClick: () => formModal.open("system-tenant-form", { title: "Nuevo tenant", description: "Crea un tenant en la plataforma.", size: "md" }),
            },
          ]}
          options={[
            { icon: FiUnlock, label: "Habilitar", disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => enable.mutate(selectedRows.map((row) => row.codigo)).then(clearSelection) },
            { icon: FiSlash, label: "Deshabilitar", disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => disable.mutate(selectedRows.map((row) => row.codigo)).then(clearSelection) },
          ]}
          rowOptions={[
            { icon: FiShield, label: "Seguridad", onClick: (row) => navigate(`/system/tenants/${encodeURIComponent(row.codigo)}/security/users`) },
            { icon: FiSettings, label: "Properties", onClick: (row) => navigate(`/system/tenants/${encodeURIComponent(row.codigo)}/properties`) },
            { icon: FiEdit2, label: "Editar", onClick: (row) => formModal.open("system-tenant-form", { title: "Editar tenant", description: row.codigo, size: "md", props: { tenant: row } }) },
            { icon: FiUnlock, label: "Habilitar", disabled: (row) => !row.disabled, onClick: (row) => enable.mutate([row.codigo]) },
            { icon: FiSlash, label: "Deshabilitar", disabled: (row) => row.disabled, onClick: (row) => disable.mutate([row.codigo]) },
          ]}
        />
      </div>
    </main>
  );
};
