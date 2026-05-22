import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { DataTable } from "../../components/data/DataTable";
import type { PermissionResponse, TenantUserResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const columns: ColumnDef<PermissionResponse, unknown>[] = [
  { accessorKey: "code", header: "Permiso", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.code}</span> },
  { accessorKey: "description", header: "Descripcion", cell: ({ row }) => row.original.description ?? <span className="text-ui-text-soft">Sin descripcion</span> },
];

export type SecurityUserPermissionsViewProps = {
  user: TenantUserResponse;
};

export const SecurityUserPermissionsView = ({ user }: SecurityUserPermissionsViewProps) => {
  const security = useSecurityService();
  const permissions = useQuery({
    queryKey: [...security.queryKeyPrefix, "users", user.username, "permissions"],
    queryFn: ({ signal }) => security.users.permissions.list(user.username, signal),
  });

  return (
    <DataTable
      data={permissions.data?.permissions ?? []}
      columns={columns}
      loading={permissions.isLoading}
      emptyMessage={permissions.isError ? "No se pudieron cargar los permisos" : "El usuario no tiene permisos efectivos"}
      searchKeys={["code", "description"]}
      pageSize={8}
    />
  );
};
