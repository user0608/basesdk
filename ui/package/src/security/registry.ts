import { defineRegistry } from "../platform/registry";
import { SecurityGroupsPage } from "./pages/SecurityGroupsPage";
import { SecurityPermissionsPage } from "./pages/SecurityPermissionsPage";
import { SecurityRolesPage } from "./pages/SecurityRolesPage";
import { SecurityUsersPage } from "./pages/SecurityUsersPage";
import { SecurityGroupAssignmentsForm } from "./forms/SecurityGroupAssignmentsForm";
import { SecurityGroupForm } from "./forms/SecurityGroupForm";
import { SecurityGroupPermissionsView } from "./forms/SecurityGroupPermissionsView";
import { SecurityPermissionAssignmentsView } from "./forms/SecurityPermissionAssignmentsView";
import { SecurityRolePermissionsForm } from "./forms/SecurityRolePermissionsForm";
import { SecurityRoleAssignmentsForm } from "./forms/SecurityRoleAssignmentsForm";
import { SecurityUserAssignmentsForm } from "./forms/SecurityUserAssignmentsForm";
import { SecurityRoleForm } from "./forms/SecurityRoleForm";
import { SecurityUserPasswordForm } from "./forms/SecurityUserPasswordForm";
import { SecurityUserPermissionsView } from "./forms/SecurityUserPermissionsView";
import { SecurityUserForm } from "./forms/SecurityUserForm";

export const securityRegistry = defineRegistry({
  pages: {
    "security-users-page": SecurityUsersPage,
    "security-roles-page": SecurityRolesPage,
    "security-groups-page": SecurityGroupsPage,
    "security-permissions-page": SecurityPermissionsPage,
  },
  forms: {
    "security-user-form": SecurityUserForm,
    "security-role-form": SecurityRoleForm,
    "security-group-form": SecurityGroupForm,
    "security-group-assignments-form": SecurityGroupAssignmentsForm,
    "security-group-permissions-view": SecurityGroupPermissionsView,
    "security-permission-assignments-view": SecurityPermissionAssignmentsView,
    "security-user-assignments-form": SecurityUserAssignmentsForm,
    "security-user-permissions-view": SecurityUserPermissionsView,
    "security-user-password-form": SecurityUserPasswordForm,
    "security-role-permissions-form": SecurityRolePermissionsForm,
    "security-role-assignments-form": SecurityRoleAssignmentsForm,
  },
});
