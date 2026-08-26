export interface UserPermissions {
  isSiteAdmin: boolean;
  adminOfOrganizationIds: string[];
  ownedGroupIds: string[];
  memberGroupIds: string[];
}

export interface RbacContextValue {
  permissions: UserPermissions | null;
  loading: boolean;
  error: string | null;

  // Role checks
  isSiteAdmin: () => boolean;
  isOrgAdmin: () => boolean;
  isGroupOwner: () => boolean;

  // Scoped checks
  isAdminOf: (orgId: string) => boolean;
  isOwnerOf: (groupId: string) => boolean;
  isMemberOf: (groupId: string) => boolean;
  hasAccessTo: (groupId: string) => boolean;

  // Capability checks (key differences between roles)
  canCreateOrganizations: () => boolean; // Site Admin only
  canInviteOrgAdmins: () => boolean; // Site Admin only
  canCreateGroups: () => boolean; // Site Admin or Org Admin (NOT Group Owner!)
  canInviteGroupOwners: () => boolean; // Site Admin or Org Admin (NOT Group Owner!)
  canInviteToGroup: (groupId: string) => boolean;
  canSeeReporterIdentity: (groupId: string) => boolean;
}

export type RoleType = 'site_admin' | 'org_admin' | 'group_owner' | 'member';

export type CapabilityCheck =
  | 'canCreateOrganizations'
  | 'canInviteOrgAdmins'
  | 'canCreateGroups'
  | 'canInviteGroupOwners';

export type ScopedCheck =
  | 'isAdminOf'
  | 'isOwnerOf'
  | 'canInviteToGroup'
  | 'canSeeReporterIdentity';
