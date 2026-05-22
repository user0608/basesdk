import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { DataTable } from "../../components/data/DataTable";
import type { GroupResponse, PermissionResponse } from "../types";
import { useSecurityService } from "../useSecurityService";

const columns: ColumnDef<PermissionResponse, unknown>[] = [
  { accessorKey: "code", header: "Permiso", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.code}</span> },
  {
    accessorKey: "description",
    header: "Descripcion",
    cell: ({ row }) => row.original.description ?? <span className="text-ui-text-soft">Sin descripcion</span>,
  },
];

export type SecurityGroupPermissionsViewProps = {
  group: GroupResponse;
};

export const SecurityGroupPermissionsView = ({ group }: SecurityGroupPermissionsViewProps) => {
  const security = useSecurityService();
  const permissions = useQuery({
    queryKey: [...security.queryKeyPrefix, "groups", group.code, "permissions"],
    queryFn: ({ signal }) => security.groups.permissions.list(group.code, signal),
  });

  return (
    <DataTable
      data={permissions.data ?? []}
      columns={columns}
      loading={permissions.isLoading}
      emptyMessage={permissions.isError ? "No se pudieron cargar los permisos" : "El grupo no tiene permisos efectivos"}
      searchKeys={["code", "description"]}
      pageSize={8}
    />
  );
};
