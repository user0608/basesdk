import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { DataTable } from "../../components/data/DataTable";
import type { GroupResponse, PermissionResponse, RoleResponse, TenantUserResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

export type SecurityPermissionAssignmentsViewProps = {
  permission: PermissionResponse;
  mode: "roles" | "groups" | "users";
};

const roleColumns: ColumnDef<RoleResponse, unknown>[] = [
  { accessorKey: "code", header: "Rol", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.code}</span> },
  { accessorKey: "description", header: "Descripcion", cell: ({ row }) => row.original.description ?? <span className="text-ui-text-soft">Sin descripcion</span> },
  { accessorKey: "disabled", header: "Estado", cell: ({ row }) => (row.original.disabled ? "Inactivo" : "Activo") },
];

const groupColumns: ColumnDef<GroupResponse, unknown>[] = [
  { accessorKey: "code", header: "Grupo", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.code}</span> },
  { accessorKey: "description", header: "Descripcion", cell: ({ row }) => row.original.description ?? <span className="text-ui-text-soft">Sin descripcion</span> },
  { accessorKey: "disabled", header: "Estado", cell: ({ row }) => (row.original.disabled ? "Inactivo" : "Activo") },
];

const userColumns: ColumnDef<TenantUserResponse, unknown>[] = [
  { accessorKey: "username", header: "Usuario", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.username}</span> },
  { accessorKey: "fullName", header: "Nombre", cell: ({ row }) => row.original.fullName ?? <span className="text-ui-text-soft">Sin nombre</span> },
  { accessorKey: "disabled", header: "Estado", cell: ({ row }) => (row.original.disabled ? "Inactivo" : "Activo") },
];

export const SecurityPermissionAssignmentsView = ({ permission, mode }: SecurityPermissionAssignmentsViewProps) => {
  const security = useSecurityService();
  const items = useQuery({
    queryKey: [...security.queryKeyPrefix, "permissions", permission.code, mode],
    queryFn: ({ signal }) => {
      if (mode === "roles") return security.permissions.roles.list(permission.code, signal);
      if (mode === "groups") return security.permissions.groups.list(permission.code, signal);
      return security.permissions.users.list(permission.code, signal);
    },
  });

  const columns = mode === "roles" ? roleColumns : mode === "groups" ? groupColumns : userColumns;
  const emptyMessage = mode === "roles" ? "No hay roles con este permiso" : mode === "groups" ? "No hay grupos con este permiso" : "No hay usuarios con este permiso";

  return (
    <DataTable
      data={items.data ?? []}
      columns={columns}
      loading={items.isLoading}
      emptyMessage={items.isError ? "No se pudieron cargar las asignaciones" : emptyMessage}
      searchKeys={mode === "users" ? ["username", "fullName"] : ["code", "description"]}
      pageSize={8}
    />
  );
};
