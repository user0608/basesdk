import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { FiEdit2, FiLock, FiPlus, FiSlash, FiTrash2, FiUnlock } from "react-icons/fi";
import { DataTable } from "../../components/data/DataTable";
import { Permissions } from "../../generated/permissions";
import { useFormModal } from "../../platform/useFormModal";
import { useMutate } from "../../query/useMutate";
import { useSecurityService } from "../useSecurityService";
import type { TenantUserResponse } from "../types";

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

export const SecurityUsersPage = () => {
  const security = useSecurityService();
  const queryKey = [...security.queryKeyPrefix, "users"];
  const formModal = useFormModal();
  const users = useQuery({
    queryKey,
    queryFn: ({ signal }) => security.users.list(signal),
  });
  const columns: ColumnDef<TenantUserResponse, unknown>[] = [
    {
      accessorKey: "username",
      header: "Usuario",
      cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.username}</span>,
    },
    {
      accessorKey: "fullName",
      header: "Nombre",
      cell: ({ row }) => row.original.fullName ?? <span className="text-ui-text-soft">Sin nombre</span>,
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
              formModal.open("security-user-assignments-form", {
                title: "Roles del usuario",
                description: row.original.username,
                size: "lg",
                props: { user: row.original, mode: "roles" },
              })
            }
          />
          <AssignmentChip
            label="grupos"
            value={row.original.groupsCount}
            onClick={() =>
              formModal.open("security-user-assignments-form", {
                title: "Grupos del usuario",
                description: row.original.username,
                size: "lg",
                props: { user: row.original, mode: "groups" },
              })
            }
          />
          <AssignmentChip
            label="permisos"
            value={row.original.permissionsCount}
            onClick={() =>
              formModal.open("security-user-permissions-view", {
                title: "Permisos efectivos",
                description: row.original.username,
                size: "lg",
                props: { user: row.original },
              })
            }
          />
        </div>
      ),
    },
    {
      accessorKey: "mustChangePassword",
      header: "Cambio password",
      cell: ({ row }) => (row.original.mustChangePassword ? "Requerido" : "No"),
    },
    {
      accessorKey: "lastLoginAt",
      header: "Ultimo acceso",
      cell: ({ row }) => row.original.lastLoginAt ?? <span className="text-ui-text-soft">Nunca</span>,
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
  const enable = useMutate({
    mutationFn: (usernames: string[]) => security.users.enable(usernames),
    invalidateQueryKey: queryKey,
    requireConfirm: true,
    confirmMessage: "Quieres habilitar los usuarios seleccionados?",
    successMessage: "Usuarios habilitados.",
  });
  const disable = useMutate({
    mutationFn: (usernames: string[]) => security.users.disable(usernames),
    invalidateQueryKey: queryKey,
    requireConfirm: true,
    confirmMessage: "Quieres deshabilitar los usuarios seleccionados?",
    successMessage: "Usuarios deshabilitados.",
  });
  const remove = useMutate({
    mutationFn: (usernames: string[]) => security.users.delete(usernames),
    invalidateQueryKey: queryKey,
    requireConfirm: true,
    confirmMessage: "Quieres eliminar los usuarios seleccionados?",
    confirmLabel: "Eliminar",
    successMessage: "Usuarios eliminados.",
  });

  return (
    <DataTable
      title="Usuarios"
      description="Usuarios del tenant activo."
      data={users.data ?? []}
      columns={columns}
      loading={users.isLoading}
      emptyMessage={users.isError ? "No se pudieron cargar los usuarios" : "No hay usuarios registrados"}
      selectable
      getRowId={(row) => row.username}
      searchPlaceholder="Buscar usuarios..."
      searchKeys={["username", "fullName"]}
      pageSize={20}
      maxItemsRange={100}
      actions={[
        {
          icon: FiPlus,
          label: "Nuevo",
          permissions: [Permissions.securityUsersCreate],
          onClick: () =>
            formModal.open("security-user-form", {
              title: "Nuevo usuario",
              description: "Crea un usuario para el tenant activo.",
              size: "md",
            }),
        },
      ]}
      options={[
        {
          icon: FiUnlock,
          label: "Habilitar",
          permissions: [Permissions.securityUsersEnable],
          disabled: ({ selectedRows }) => selectedRows.length === 0,
          onClick: ({ selectedRows, clearSelection }) => enable.mutate(selectedRows.map((row) => row.username)).then(clearSelection),
        },
        {
          icon: FiSlash,
          label: "Deshabilitar",
          permissions: [Permissions.securityUsersDisable],
          disabled: ({ selectedRows }) => selectedRows.length === 0,
          onClick: ({ selectedRows, clearSelection }) => disable.mutate(selectedRows.map((row) => row.username)).then(clearSelection),
        },
        {
          icon: FiTrash2,
          label: "Eliminar",
          permissions: [Permissions.securityUsersDelete],
          variant: "danger",
          disabled: ({ selectedRows }) => selectedRows.length === 0,
          onClick: ({ selectedRows, clearSelection }) => remove.mutate(selectedRows.map((row) => row.username)).then(clearSelection),
        },
      ]}
      rowOptions={[
        {
          icon: FiEdit2,
          label: "Editar",
          permissions: [Permissions.securityUsersUpdate],
          onClick: (row) =>
            formModal.open("security-user-form", {
              title: "Editar usuario",
              description: row.username,
              size: "md",
              props: { user: row },
            }),
        },
        {
          icon: FiLock,
          label: "Password",
          permissions: [Permissions.securityUsersPasswordUpdate],
          onClick: (row) =>
            formModal.open("security-user-password-form", {
              title: "Cambiar password",
              description: row.username,
              size: "sm",
              props: { user: row },
            }),
        },
        {
          icon: FiUnlock,
          label: "Habilitar",
          permissions: [Permissions.securityUsersEnable],
          disabled: (row) => !row.disabled,
          onClick: (row) => enable.mutate([row.username]),
        },
        {
          icon: FiSlash,
          label: "Deshabilitar",
          permissions: [Permissions.securityUsersDisable],
          disabled: (row) => row.disabled,
          onClick: (row) => disable.mutate([row.username]),
        },
        {
          icon: FiTrash2,
          label: "Eliminar",
          permissions: [Permissions.securityUsersDelete],
          variant: "danger",
          onClick: (row) => remove.mutate([row.username]),
        },
      ]}
    />
  );
};
