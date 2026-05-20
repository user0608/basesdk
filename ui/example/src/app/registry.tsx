import { defineComponentRegistry } from "@basesdk/ui";
import { SecurityGroupsPage } from "./views/SecurityGroupsPage";
import { SecurityPermissionsPage } from "./views/SecurityPermissionsPage";
import { SecurityRolesPage } from "./views/SecurityRolesPage";
import { SecurityUsersPage } from "./views/SecurityUsersPage";

export const componentRegistry = defineComponentRegistry({
  "security-users-page": SecurityUsersPage,
  "security-roles-page": SecurityRolesPage,
  "security-groups-page": SecurityGroupsPage,
  "security-permissions-page": SecurityPermissionsPage,
});
