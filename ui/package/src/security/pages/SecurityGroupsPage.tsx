import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { FiEdit2, FiKey, FiPlus, FiShield, FiSlash, FiTrash2, FiUnlock, FiUsers } from "react-icons/fi";
import { DataTable } from "../../components/data/DataTable";
import { Permissions } from "../../generated/permissions";
import { useFormModal } from "../../platform/useFormModal";
import { useMutate } from "../../query/useMutate";
import { useSecurityService } from "../useSecurityService";
import type { GroupResponse } from "../types";

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

export const SecurityGroupsPage = () => {
  const security = useSecurityService();
  const queryKey = [...security.queryKeyPrefix, "groups"];
  const formModal = useFormModal();
  const columns: ColumnDef<GroupResponse, unknown>[] = [
    { accessorKey: "code", header: "Codigo", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.code}</span> },
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
              formModal.open("security-group-assignments-form", {
                title: "Usuarios del grupo",
                description: row.original.code,
                size: "lg",
                props: { group: row.original, mode: "users" },
              })
            }
          />
          <AssignmentChip
            label="roles"
            value={row.original.rolesCount}
            onClick={() =>
              formModal.open("security-group-assignments-form", {
                title: "Roles del grupo",
                description: row.original.code,
                size: "lg",
                props: { group: row.original, mode: "roles" },
              })
            }
          />
          <AssignmentChip
            label="permisos"
            value={row.original.permissionsCount}
            onClick={() =>
              formModal.open("security-group-permissions-view", {
                title: "Permisos del grupo",
                description: row.original.code,
                size: "lg",
                props: { group: row.original },
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
        <span className={row.original.disabled ? "text-ui-danger" : "text-ui-primary"}>{row.original.disabled ? "Inactivo" : "Activo"}</span>
      ),
    },
  ];
  const groups = useQuery({ queryKey, queryFn: ({ signal }) => security.groups.list(signal) });
  const enable = useMutate({ mutationFn: (codes: string[]) => security.groups.enable(codes), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres habilitar los grupos seleccionados?", successMessage: "Grupos habilitados." });
  const disable = useMutate({ mutationFn: (codes: string[]) => security.groups.disable(codes), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres deshabilitar los grupos seleccionados?", successMessage: "Grupos deshabilitados." });
  const remove = useMutate({ mutationFn: (codes: string[]) => security.groups.delete(codes), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres eliminar los grupos seleccionados?", confirmLabel: "Eliminar", successMessage: "Grupos eliminados." });

  return (
    <DataTable
      tableId="security.groups"
      title="Grupos"
      description="Grupos del tenant activo."
      data={groups.data ?? []}
      columns={columns}
      initialHiddenColumns={["code"]}
      loading={groups.isLoading}
      emptyMessage={groups.isError ? "No se pudieron cargar los grupos" : "No hay grupos registrados"}
      selectable
      getRowId={(row) => row.code}
      searchPlaceholder="Buscar grupos..."
      searchKeys={["code", "description"]}
      pageSize={20}
      maxItemsRange={100}
      actions={[
        {
          icon: FiPlus,
          label: "Nuevo",
          permissions: [Permissions.securityGroupsCreate],
          onClick: () =>
            formModal.open("security-group-form", {
              title: "Nuevo grupo",
              description: "Crea un grupo para el tenant activo.",
              size: "md",
            }),
        },
      ]}
      options={[
        { icon: FiUnlock, label: "Habilitar", permissions: [Permissions.securityGroupsEnable], disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => enable.mutate(selectedRows.map((row) => row.code)).then(clearSelection) },
        { icon: FiSlash, label: "Deshabilitar", permissions: [Permissions.securityGroupsDisable], disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => disable.mutate(selectedRows.map((row) => row.code)).then(clearSelection) },
        { icon: FiTrash2, label: "Eliminar", permissions: [Permissions.securityGroupsDelete], variant: "danger", disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => remove.mutate(selectedRows.map((row) => row.code)).then(clearSelection) },
      ]}
      rowOptions={[
        {
          icon: FiEdit2,
          label: "Editar",
          permissions: [Permissions.securityGroupsUpdate],
          onClick: (row) =>
            formModal.open("security-group-form", {
              title: "Editar grupo",
              description: row.code,
              size: "md",
              props: { group: row },
            }),
        },
        { icon: FiUnlock, label: "Habilitar", permissions: [Permissions.securityGroupsEnable], disabled: (row) => !row.disabled, onClick: (row) => enable.mutate([row.code]) },
        { icon: FiSlash, label: "Deshabilitar", permissions: [Permissions.securityGroupsDisable], disabled: (row) => row.disabled, onClick: (row) => disable.mutate([row.code]) },
        { icon: FiTrash2, label: "Eliminar", permissions: [Permissions.securityGroupsDelete], variant: "danger", onClick: (row) => remove.mutate([row.code]) },
      ]}
    />
  );
};
