import { FiKey, FiShield, FiUsers } from "react-icons/fi";
import { Permissions } from "../generated/permissions";
import { defineMenuTree } from "../platform/registry";
import type { MenuTree, PermissionCode } from "../platform/types";
import { securityRegistry } from "./registry";

type SecurityPageId = keyof typeof securityRegistry.pages & string;

export const securityModules = defineMenuTree([
  {
    id: "security",
    label: "Seguridad",
    icon: FiShield,
    path: "/app/security",
    order: -100,
    children: [
      {
        id: "security-main",
        label: "Seguridad",
        icon: FiShield,
        path: "",
        order: 1,
        children: [
          {
            id: "security-users-list",
            label: "Usuarios",
            icon: FiUsers,
            path: "users",
            componentId: "security-users-page",
            permissions: [Permissions.securityUsersRead],
            order: 1,
          },
          {
            id: "security-groups-list",
            label: "Grupos",
            icon: FiUsers,
            path: "groups",
            componentId: "security-groups-page",
            permissions: [Permissions.securityGroupsRead],
            order: 2,
          },
          {
            id: "security-roles-list",
            label: "Roles",
            icon: FiShield,
            path: "roles",
            componentId: "security-roles-page",
            permissions: [Permissions.securityRolesRead],
            order: 3,
          },
          {
            id: "security-permissions-list",
            label: "Permisos",
            icon: FiKey,
            path: "permissions",
            componentId: "security-permissions-page",
            permissions: [Permissions.securityPermissionsRead],
            order: 4,
          },
        ],
      },
    ],
  },
] satisfies MenuTree<SecurityPageId, PermissionCode>);
