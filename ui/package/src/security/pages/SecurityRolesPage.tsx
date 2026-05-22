import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { FiEdit2, FiPlus, FiSlash, FiTrash2, FiUnlock } from "react-icons/fi";
import { DataTable } from "../../components/data/DataTable";
import { Permissions } from "../../generated/permissions";
import { useFormModal } from "../../platform/useFormModal";
import { useMutate } from "../../query/useMutate";
import { useSecurityService } from "../useSecurityService";
import type { RoleResponse } from "../types";

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

export const SecurityRolesPage = () => {
  const security = useSecurityService();
  const queryKey = [...security.queryKeyPrefix, "roles"];
  const formModal = useFormModal();
  const roles = useQuery({ queryKey, queryFn: ({ signal }) => security.roles.list(signal) });
  const columns: ColumnDef<RoleResponse, unknown>[] = [
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
            label="usuarios"
            value={row.original.usersCount}
            onClick={() =>
              formModal.open("security-role-assignments-form", {
                title: "Usuarios del rol",
                description: row.original.code,
                size: "lg",
                props: { role: row.original, mode: "users" },
              })
            }
          />
          <AssignmentChip
            label="grupos"
            value={row.original.groupsCount}
            onClick={() =>
              formModal.open("security-role-assignments-form", {
                title: "Grupos del rol",
                description: row.original.code,
                size: "lg",
                props: { role: row.original, mode: "groups" },
              })
            }
          />
          <AssignmentChip
            label="permisos"
            value={row.original.permissionsCount}
            onClick={() =>
              formModal.open("security-role-permissions-form", {
                title: "Permisos del rol",
                description: row.original.code,
                size: "lg",
                props: { role: row.original },
              })
            }
          />
        </div>
      ),
    },
    {
      accessorKey: "disabled",
      header: "Estado",
      cell: ({ row }) => (
        <span className={row.original.disabled ? "text-ui-danger" : "text-ui-primary"}>
          {row.original.disabled ? "Inactivo" : "Activo"}
        </span>
      ),
    },
  ];
  const enable = useMutate({ mutationFn: (codes: string[]) => security.roles.enable(codes), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres habilitar los roles seleccionados?", successMessage: "Roles habilitados." });
  const disable = useMutate({ mutationFn: (codes: string[]) => security.roles.disable(codes), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres deshabilitar los roles seleccionados?", successMessage: "Roles deshabilitados." });
  const remove = useMutate({ mutationFn: (codes: string[]) => security.roles.delete(codes), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres eliminar los roles seleccionados?", confirmLabel: "Eliminar", successMessage: "Roles eliminados." });

  return (
    <DataTable
      tableId="security.roles"
      title="Roles"
      description="Roles del tenant activo."
      data={roles.data ?? []}
      columns={columns}
      initialHiddenColumns={["code"]}
      loading={roles.isLoading}
      emptyMessage={roles.isError ? "No se pudieron cargar los roles" : "No hay roles registrados"}
      selectable
      getRowId={(row) => row.code}
      searchPlaceholder="Buscar roles..."
      searchKeys={["code", "description"]}
      pageSize={20}
      maxItemsRange={100}
      actions={[
        {
          icon: FiPlus,
          label: "Nuevo",
          permissions: [Permissions.securityRolesCreate],
          onClick: () =>
            formModal.open("security-role-form", {
              title: "Nuevo rol",
              description: "Crea un rol para el tenant activo.",
              size: "md",
            }),
        },
      ]}
      options={[
        { icon: FiUnlock, label: "Habilitar", permissions: [Permissions.securityRolesEnable], disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => enable.mutate(selectedRows.map((row) => row.code)).then(clearSelection) },
        { icon: FiSlash, label: "Deshabilitar", permissions: [Permissions.securityRolesDisable], disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => disable.mutate(selectedRows.map((row) => row.code)).then(clearSelection) },
        { icon: FiTrash2, label: "Eliminar", permissions: [Permissions.securityRolesDelete], variant: "danger", disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => remove.mutate(selectedRows.map((row) => row.code)).then(clearSelection) },
      ]}
      rowOptions={[
        {
          icon: FiEdit2,
          label: "Editar",
          permissions: [Permissions.securityRolesUpdate],
          onClick: (row) =>
            formModal.open("security-role-form", { title: "Editar rol", description: row.code, size: "md", props: { role: row } }),
        },
        { icon: FiUnlock, label: "Habilitar", permissions: [Permissions.securityRolesEnable], disabled: (row) => !row.disabled, onClick: (row) => enable.mutate([row.code]) },
        { icon: FiSlash, label: "Deshabilitar", permissions: [Permissions.securityRolesDisable], disabled: (row) => row.disabled, onClick: (row) => disable.mutate([row.code]) },
        { icon: FiTrash2, label: "Eliminar", permissions: [Permissions.securityRolesDelete], variant: "danger", onClick: (row) => remove.mutate([row.code]) },
      ]}
    />
  );
};
