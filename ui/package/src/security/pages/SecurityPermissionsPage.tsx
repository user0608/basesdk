import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { DataTable } from "../../components/data/DataTable";
import { useFormModal } from "../../platform/useFormModal";
import { useSecurityService } from "../useSecurityService";
import type { PermissionResponse } from "../types";

const AssignmentChip = ({ label, value, onClick }: { label: string; value: number; onClick: () => void }) => (
  <button
    type="button"
    className="inline-flex items-center gap-1 rounded-full bg-ui-surface px-2 py-1 text-xs font-medium text-ui-text-muted ring-1 ring-inset ring-ui-border/55 transition-colors hover:bg-ui-primary/10 hover:text-ui-primary hover:ring-ui-primary/25"
    onClick={(event) => {
      event.stopPropagation();
      onClick();
    }}
  >
    <span className="font-semibold tabular-nums text-ui-text">{value}</span>
    <span>{label}</span>
  </button>
);

export const SecurityPermissionsPage = () => {
  const security = useSecurityService();
  const formModal = useFormModal();
  const permissions = useQuery({
    queryKey: [...security.queryKeyPrefix, "permissions"],
    queryFn: ({ signal }) => security.permissions.list(signal),
  });
  const columns: ColumnDef<PermissionResponse, unknown>[] = [
    {
      accessorKey: "code",
      header: "Codigo",
      cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.code}</span>,
    },
    {
      accessorKey: "description",
      header: "Descripcion",
      cell: ({ row }) => row.original.description ?? <span className="text-ui-text-soft">Sin descripcion</span>,
    },
    {
      id: "assignments",
      header: "Asignaciones",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-1.5">
          <AssignmentChip
            label="roles"
            value={row.original.rolesCount}
            onClick={() =>
              formModal.open("security-permission-assignments-view", {
                title: "Roles con permiso",
                description: row.original.code,
                size: "lg",
                props: { permission: row.original, mode: "roles" },
              })
            }
          />
          <AssignmentChip
            label="grupos"
            value={row.original.groupsCount}
            onClick={() =>
              formModal.open("security-permission-assignments-view", {
                title: "Grupos con permiso",
                description: row.original.code,
                size: "lg",
                props: { permission: row.original, mode: "groups" },
              })
            }
          />
          <AssignmentChip
            label="usuarios"
            value={row.original.usersCount}
            onClick={() =>
              formModal.open("security-permission-assignments-view", {
                title: "Usuarios con permiso",
                description: row.original.code,
                size: "lg",
                props: { permission: row.original, mode: "users" },
              })
            }
          />
        </div>
      ),
    },
  ];

  return (
    <DataTable
      title="Permisos"
      description="Permisos registrados en la plataforma."
      data={permissions.data ?? []}
      columns={columns}
      loading={permissions.isLoading}
      emptyMessage={permissions.isError ? "No se pudieron cargar los permisos" : "No hay permisos registrados"}
      searchPlaceholder="Buscar permisos..."
      searchKeys={["code", "description"]}
      getRowId={(row) => row.code}
      pagination
      pageSize={20}
      maxItemsRange={100}
    />
  );
};
