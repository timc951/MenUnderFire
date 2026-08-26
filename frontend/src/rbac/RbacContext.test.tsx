import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ReactNode } from 'react';
import { RbacProvider, useRbac } from './RbacContext';
import { UserPermissions } from './types';

const mockApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('../hooks/useApi', () => ({
  useApi: () => mockApi,
}));

// Mock helper functions
function mockSiteAdmin() {
  const permissions: UserPermissions = {
    isSiteAdmin: true,
    adminOfOrganizationIds: [],
    ownedGroupIds: [],
    memberGroupIds: [],
  };
  mockApi.get.mockResolvedValue(permissions);
}

function mockOrgAdmin(orgIds: string[]) {
  const permissions: UserPermissions = {
    isSiteAdmin: false,
    adminOfOrganizationIds: orgIds,
    ownedGroupIds: [],
    memberGroupIds: [],
  };
  mockApi.get.mockResolvedValue(permissions);
}

function mockGroupOwner(ownedGroupIds: string[], memberGroupIds: string[] = []) {
  const permissions: UserPermissions = {
    isSiteAdmin: false,
    adminOfOrganizationIds: [],
    ownedGroupIds,
    memberGroupIds,
  };
  mockApi.get.mockResolvedValue(permissions);
}

function mockGroupMember(memberGroupIds: string[]) {
  const permissions: UserPermissions = {
    isSiteAdmin: false,
    adminOfOrganizationIds: [],
    ownedGroupIds: [],
    memberGroupIds,
  };
  mockApi.get.mockResolvedValue(permissions);
}

function wrapper({ children }: { children: ReactNode }) {
  return <RbacProvider>{children}</RbacProvider>;
}

describe('RbacProvider', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Site Admin', () => {
    beforeEach(() => mockSiteAdmin());

    it('should be able to create organizations', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canCreateOrganizations()).toBe(true);
      });
    });

    it('should be able to invite org admins', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteOrgAdmins()).toBe(true);
      });
    });

    it('should be able to create groups', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canCreateGroups()).toBe(true);
      });
    });

    it('should be able to invite group owners', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteGroupOwners()).toBe(true);
      });
    });

    it('should identify as site admin', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isSiteAdmin()).toBe(true);
      });
    });

    it('should be able to invite to any group', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteToGroup('any-group-id')).toBe(true);
      });
    });

    it('should be able to see reporter identity in any group', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canSeeReporterIdentity('any-group-id')).toBe(true);
      });
    });
  });

  describe('Org Admin', () => {
    const orgId = 'org-123';
    beforeEach(() => mockOrgAdmin([orgId]));

    it('should NOT be able to create organizations', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canCreateOrganizations()).toBe(false);
      });
    });

    it('should NOT be able to invite org admins', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteOrgAdmins()).toBe(false);
      });
    });

    it('should be able to create groups', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canCreateGroups()).toBe(true);
      });
    });

    it('should be able to invite group owners', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteGroupOwners()).toBe(true);
      });
    });

    it('should identify as org admin', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isOrgAdmin()).toBe(true);
      });
    });

    it('should be admin of the correct organization', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isAdminOf(orgId)).toBe(true);
        expect(result.current.isAdminOf('other-org')).toBe(false);
      });
    });

    it('should NOT identify as site admin', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isSiteAdmin()).toBe(false);
      });
    });
  });

  describe('Group Owner', () => {
    const ownedGroupId = 'group-owned';
    const memberGroupId = 'group-member';
    beforeEach(() => mockGroupOwner([ownedGroupId], [memberGroupId]));

    it('should NOT be able to create groups', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canCreateGroups()).toBe(false);
      });
    });

    it('should NOT be able to invite group owners', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteGroupOwners()).toBe(false);
      });
    });

    it('should be able to invite to owned group only', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteToGroup(ownedGroupId)).toBe(true);
        expect(result.current.canInviteToGroup(memberGroupId)).toBe(false);
      });
    });

    it('should see reporter identity only in owned group', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canSeeReporterIdentity(ownedGroupId)).toBe(true);
        expect(result.current.canSeeReporterIdentity(memberGroupId)).toBe(false);
      });
    });

    it('should identify as group owner', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isGroupOwner()).toBe(true);
      });
    });

    it('should be owner of the correct group', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isOwnerOf(ownedGroupId)).toBe(true);
        expect(result.current.isOwnerOf(memberGroupId)).toBe(false);
      });
    });

    it('should have access to both owned and member groups', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.hasAccessTo(ownedGroupId)).toBe(true);
        expect(result.current.hasAccessTo(memberGroupId)).toBe(true);
        expect(result.current.hasAccessTo('other-group')).toBe(false);
      });
    });

    it('should NOT be able to create organizations', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canCreateOrganizations()).toBe(false);
      });
    });

    it('should NOT be able to invite org admins', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteOrgAdmins()).toBe(false);
      });
    });
  });

  describe('Group Member', () => {
    const memberGroupId = 'group-member';
    beforeEach(() => mockGroupMember([memberGroupId]));

    it('should NOT be able to create groups', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canCreateGroups()).toBe(false);
      });
    });

    it('should NOT be able to invite anyone', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteToGroup(memberGroupId)).toBe(false);
      });
    });

    it('should NOT see reporter identity', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canSeeReporterIdentity(memberGroupId)).toBe(false);
      });
    });

    it('should NOT identify as group owner', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isGroupOwner()).toBe(false);
      });
    });

    it('should be member of the correct group', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.isMemberOf(memberGroupId)).toBe(true);
        expect(result.current.isMemberOf('other-group')).toBe(false);
      });
    });

    it('should have access to member groups', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.hasAccessTo(memberGroupId)).toBe(true);
        expect(result.current.hasAccessTo('other-group')).toBe(false);
      });
    });

    it('should NOT be able to invite group owners', async () => {
      const { result } = renderHook(() => useRbac(), { wrapper });
      await waitFor(() => {
        expect(result.current.canInviteGroupOwners()).toBe(false);
      });
    });
  });

  describe('Loading and Error States', () => {
    it('should show loading state initially', () => {
      mockSiteAdmin();
      const { result } = renderHook(() => useRbac(), { wrapper });
      expect(result.current.loading).toBe(true);
    });

    it('should handle error state', async () => {
      mockApi.get.mockRejectedValue(new Error('Network error'));
      const { result } = renderHook(() => useRbac(), { wrapper });

      await waitFor(() => {
        expect(result.current.error).toBe('Network error');
      });
    });

    it('should return false for all capabilities when permissions not loaded', () => {
      mockApi.get.mockImplementation(() => new Promise(() => {})); // Never resolves
      const { result } = renderHook(() => useRbac(), { wrapper });

      expect(result.current.canCreateOrganizations()).toBe(false);
      expect(result.current.canInviteOrgAdmins()).toBe(false);
      expect(result.current.canCreateGroups()).toBe(false);
      expect(result.current.canInviteGroupOwners()).toBe(false);
    });
  });
});
