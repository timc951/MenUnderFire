import { createContext, useContext, useEffect, useState, useCallback, ReactNode } from 'react';
import { useApi } from '../hooks/useApi';
import { RbacContextValue, UserPermissions } from './types';

const defaultContextValue: RbacContextValue = {
  permissions: null,
  loading: true,
  error: null,
  isSiteAdmin: () => false,
  isOrgAdmin: () => false,
  isGroupOwner: () => false,
  isAdminOf: () => false,
  isOwnerOf: () => false,
  isMemberOf: () => false,
  hasAccessTo: () => false,
  canCreateOrganizations: () => false,
  canInviteOrgAdmins: () => false,
  canCreateGroups: () => false,
  canInviteGroupOwners: () => false,
  canInviteToGroup: () => false,
  canSeeReporterIdentity: () => false,
};

const RbacContext = createContext<RbacContextValue>(defaultContextValue);

interface RbacProviderProps {
  children: ReactNode;
}

export function RbacProvider({ children }: RbacProviderProps) {
  const api = useApi();
  const [permissions, setPermissions] = useState<UserPermissions | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;

    async function fetchPermissions() {
      try {
        const data = await api.get<UserPermissions>('/users/me/permissions');
        if (mounted) {
          setPermissions(data);
          setLoading(false);
        }
      } catch (err) {
        if (mounted) {
          setError(err instanceof Error ? err.message : 'Failed to load permissions');
          setLoading(false);
        }
      }
    }

    fetchPermissions();

    return () => {
      mounted = false;
    };
  }, [api]);

  // Role checks
  const isSiteAdmin = useCallback(() => {
    return permissions?.isSiteAdmin ?? false;
  }, [permissions]);

  const isOrgAdmin = useCallback(() => {
    return (permissions?.adminOfOrganizationIds?.length ?? 0) > 0;
  }, [permissions]);

  const isGroupOwner = useCallback(() => {
    return (permissions?.ownedGroupIds?.length ?? 0) > 0;
  }, [permissions]);

  // Scoped checks
  const isAdminOf = useCallback(
    (orgId: string) => {
      return permissions?.adminOfOrganizationIds?.includes(orgId) ?? false;
    },
    [permissions]
  );

  const isOwnerOf = useCallback(
    (groupId: string) => {
      return permissions?.ownedGroupIds?.includes(groupId) ?? false;
    },
    [permissions]
  );

  const isMemberOf = useCallback(
    (groupId: string) => {
      return permissions?.memberGroupIds?.includes(groupId) ?? false;
    },
    [permissions]
  );

  const hasAccessTo = useCallback(
    (groupId: string) => {
      if (permissions?.isSiteAdmin) return true;
      return (
        permissions?.ownedGroupIds?.includes(groupId) ||
        permissions?.memberGroupIds?.includes(groupId) ||
        false
      );
    },
    [permissions]
  );

  // Capability checks - Site Admin only
  const canCreateOrganizations = useCallback(() => {
    return permissions?.isSiteAdmin ?? false;
  }, [permissions]);

  const canInviteOrgAdmins = useCallback(() => {
    return permissions?.isSiteAdmin ?? false;
  }, [permissions]);

  // Capability checks - Site Admin OR Org Admin (NOT Group Owner!)
  const canCreateGroups = useCallback(() => {
    return (
      permissions?.isSiteAdmin ||
      (permissions?.adminOfOrganizationIds?.length ?? 0) > 0
    );
  }, [permissions]);

  const canInviteGroupOwners = useCallback(() => {
    return (
      permissions?.isSiteAdmin ||
      (permissions?.adminOfOrganizationIds?.length ?? 0) > 0
    );
  }, [permissions]);

  // Capability checks - Scoped to group
  const canInviteToGroup = useCallback(
    (groupId: string) => {
      // Site Admin can invite to any group
      if (permissions?.isSiteAdmin) return true;
      // Org Admin can invite to any group (checked at service level)
      if ((permissions?.adminOfOrganizationIds?.length ?? 0) > 0) return true;
      // Group Owner can invite to their owned groups only
      return permissions?.ownedGroupIds?.includes(groupId) ?? false;
    },
    [permissions]
  );

  const canSeeReporterIdentity = useCallback(
    (groupId: string) => {
      // Site Admin can see reporter identity in any group
      if (permissions?.isSiteAdmin) return true;
      // Org Admin can see reporter identity (checked at service level)
      if ((permissions?.adminOfOrganizationIds?.length ?? 0) > 0) return true;
      // Group Owner can see reporter identity in their owned groups only
      return permissions?.ownedGroupIds?.includes(groupId) ?? false;
    },
    [permissions]
  );

  const contextValue: RbacContextValue = {
    permissions,
    loading,
    error,
    isSiteAdmin,
    isOrgAdmin,
    isGroupOwner,
    isAdminOf,
    isOwnerOf,
    isMemberOf,
    hasAccessTo,
    canCreateOrganizations,
    canInviteOrgAdmins,
    canCreateGroups,
    canInviteGroupOwners,
    canInviteToGroup,
    canSeeReporterIdentity,
  };

  return <RbacContext.Provider value={contextValue}>{children}</RbacContext.Provider>;
}

export function useRbac(): RbacContextValue {
  const context = useContext(RbacContext);
  if (!context) {
    throw new Error('useRbac must be used within an RbacProvider');
  }
  return context;
}
