import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { Link } from "react-router-dom";
import { FiArrowLeft, FiEdit2, FiPlus, FiSlash, FiTrash2, FiUnlock } from "react-icons/fi";
import { DataTable } from "../../components/data/DataTable";
import { useFormModal } from "../../platform/useFormModal";
import { useMutate } from "../../query/useMutate";
import { SystemBreadcrumbs } from "../SystemBreadcrumbs";
import type { SystemUserResponse } from "../types";
import { useSystemService } from "../useSystemService";

const queryKey = ["system", "users"];

export const SystemUsersPage = () => {
  const system = useSystemService();
  const formModal = useFormModal();
  const users = useQuery({ queryKey, queryFn: ({ signal }) => system.users.list(signal) });
  const enable = useMutate({ mutationFn: (usernames: string[]) => system.users.enable(usernames), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres habilitar los usuarios system seleccionados?", successMessage: "Usuarios habilitados." });
  const disable = useMutate({ mutationFn: (usernames: string[]) => system.users.disable(usernames), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres deshabilitar los usuarios system seleccionados?", successMessage: "Usuarios deshabilitados." });
  const remove = useMutate({ mutationFn: (usernames: string[]) => system.users.delete(usernames), invalidateQueryKey: queryKey, requireConfirm: true, confirmMessage: "Quieres eliminar los usuarios system seleccionados?", confirmLabel: "Eliminar", successMessage: "Usuarios eliminados." });
  const columns: ColumnDef<SystemUserResponse, unknown>[] = [
    { accessorKey: "username", header: "Usuario", cell: ({ row }) => <span className="font-medium text-ui-text">{row.original.username}</span> },
    {
      accessorKey: "disabled",
      header: "Estado",
      cell: ({ row }) => <span className={row.original.disabled ? "text-ui-danger" : "text-ui-primary"}>{row.original.disabled ? "Inactivo" : "Activo"}</span>,
    },
    { accessorKey: "createdBy", header: "Creado por" },
    { accessorKey: "createdAt", header: "Creado", cell: ({ row }) => row.original.createdAt },
    { accessorKey: "updatedAt", header: "Actualizado", cell: ({ row }) => row.original.updatedAt ?? <span className="text-ui-text-soft">Sin cambios</span> },
  ];

  return (
    <main className="flex h-screen min-h-0 flex-col overflow-hidden bg-ui-panel-muted p-3 text-ui-text lg:p-4">
      <div className="mb-3 flex min-w-0 items-center gap-2">
        <Link
          to="/system"
          aria-label="Volver a system"
          title="Volver a system"
          className="inline-grid size-9 place-items-center rounded-lg text-ui-text-soft transition-colors hover:bg-ui-surface-hover hover:text-ui-text"
        >
          <FiArrowLeft size={18} />
        </Link>
        <SystemBreadcrumbs items={[{ label: "System", to: "/system" }, { label: "Usuarios system" }]} />
      </div>

      <div className="min-h-0 flex-1">
        <DataTable
          tableId="system.users"
          title="Usuarios system"
          description="Cuentas con acceso al panel administrativo global."
          data={users.data ?? []}
          columns={columns}
          loading={users.isLoading}
          emptyMessage={users.isError ? "No se pudieron cargar los usuarios system" : "No hay usuarios system registrados"}
          selectable
          getRowId={(row) => row.username}
          searchPlaceholder="Buscar usuarios system..."
          searchKeys={["username", "createdBy"]}
          pageSize={20}
          maxItemsRange={100}
          actions={[
            {
              icon: FiPlus,
              label: "Nuevo",
              onClick: () => formModal.open("system-user-form", { title: "Nuevo usuario system", description: "Crea una cuenta de administrador global.", size: "sm" }),
            },
          ]}
          options={[
            { icon: FiUnlock, label: "Habilitar", disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => enable.mutate(selectedRows.map((row) => row.username)).then(clearSelection) },
            { icon: FiSlash, label: "Deshabilitar", disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => disable.mutate(selectedRows.map((row) => row.username)).then(clearSelection) },
            { icon: FiTrash2, label: "Eliminar", variant: "danger", disabled: ({ selectedRows }) => selectedRows.length === 0, onClick: ({ selectedRows, clearSelection }) => remove.mutate(selectedRows.map((row) => row.username)).then(clearSelection) },
          ]}
          rowOptions={[
            { icon: FiEdit2, label: "Editar", onClick: (row) => formModal.open("system-user-form", { title: "Editar usuario system", description: row.username, size: "sm", props: { user: row } }) },
            { icon: FiUnlock, label: "Habilitar", disabled: (row) => !row.disabled, onClick: (row) => enable.mutate([row.username]) },
            { icon: FiSlash, label: "Deshabilitar", disabled: (row) => row.disabled, onClick: (row) => disable.mutate([row.username]) },
            { icon: FiTrash2, label: "Eliminar", variant: "danger", onClick: (row) => remove.mutate([row.username]) },
          ]}
        />
      </div>
    </main>
  );
};
