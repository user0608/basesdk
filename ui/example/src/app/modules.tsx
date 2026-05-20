import { Permissions, defineMenuTree } from "@basesdk/ui";
import type { ComponentId, MenuTree, PermissionCode } from "@basesdk/ui";
import { FiKey, FiSettings, FiShield, FiUsers } from "react-icons/fi";
import { componentRegistry } from "./registry";

type ExampleComponentId = ComponentId<typeof componentRegistry>;

export const appModules = defineMenuTree([
  {
    id: "security",
    label: "Seguridad",
    icon: FiShield,
    path: "/app/security",
    order: 1,
    children: [
      {
        id: "security-users",
        label: "Usuarios y roles",
        icon: FiUsers,
        path: "users",
        order: 1,
        children: [
          {
            id: "security-users-list",
            label: "Usuarios",
            icon: FiUsers,
            path: "list",
            componentId: "security-users-page",
            permissions: [Permissions.securityUsersRead],
            order: 1,
          },
          {
            id: "security-roles-list",
            label: "Roles",
            icon: FiShield,
            path: "roles",
            componentId: "security-roles-page",
            permissions: [Permissions.securityRolesRead],
            order: 2,
          },
        ],
      },
    ],
  },
  {
    id: "governance",
    label: "Gobernanza",
    icon: FiSettings,
    path: "/app/governance",
    order: 2,
    children: [
      {
        id: "governance-access",
        label: "Accesos",
        icon: FiKey,
        path: "access",
        order: 1,
        children: [
          {
            id: "security-groups-list",
            label: "Grupos",
            icon: FiUsers,
            path: "groups",
            componentId: "security-groups-page",
            permissions: [Permissions.securityGroupsRead],
            order: 1,
          },
          {
            id: "security-permissions-list",
            label: "Permisos",
            icon: FiKey,
            path: "permissions",
            componentId: "security-permissions-page",
            permissions: [Permissions.securityPermissionsRead],
            order: 2,
          },
        ],
      },
    ],
  },
] satisfies MenuTree<ExampleComponentId, PermissionCode>);
