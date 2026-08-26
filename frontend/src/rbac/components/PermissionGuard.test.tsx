import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PermissionGuard } from './PermissionGuard';
import { RbacContextValue, UserPermissions } from '../types';

// Mock the useRbac hook
const mockRbac: RbacContextValue = {
  permissions: null,
  loading: false,
  error: null,
  isSiteAdmin: vi.fn(() => false),
  isOrgAdmin: vi.fn(() => false),
  isGroupOwner: vi.fn(() => false),
  isAdminOf: vi.fn(() => false),
  isOwnerOf: vi.fn(() => false),
  isMemberOf: vi.fn(() => false),
  hasAccessTo: vi.fn(() => false),
  canCreateOrganizations: vi.fn(() => false),
  canInviteOrgAdmins: vi.fn(() => false),
  canCreateGroups: vi.fn(() => false),
  canInviteGroupOwners: vi.fn(() => false),
  canInviteToGroup: vi.fn(() => false),
  canSeeReporterIdentity: vi.fn(() => false),
};

vi.mock('../RbacContext', () => ({
  useRbac: () => mockRbac,
}));

function resetMockRbac() {
  mockRbac.loading = false;
  mockRbac.error = null;
  mockRbac.permissions = {
    isSiteAdmin: false,
    adminOfOrganizationIds: [],
    ownedGroupIds: [],
    memberGroupIds: [],
  };
  (mockRbac.isSiteAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isOrgAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isGroupOwner as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isAdminOf as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isOwnerOf as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isMemberOf as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.hasAccessTo as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canCreateOrganizations as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteOrgAdmins as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canCreateGroups as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteGroupOwners as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteToGroup as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canSeeReporterIdentity as ReturnType<typeof vi.fn>).mockReturnValue(false);
}

describe('PermissionGuard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    resetMockRbac();
  });

  describe('capability checks', () => {
    it('should render children when canCreateOrganizations is true', () => {
      (mockRbac.canCreateOrganizations as ReturnType<typeof vi.fn>).mockReturnValue(true);

      render(
        <PermissionGuard check="canCreateOrganizations">
          <div>Create Org Button</div>
        </PermissionGuard>
      );

      expect(screen.getByText('Create Org Button')).toBeInTheDocument();
    });

    it('should not render children when canCreateOrganizations is false', () => {
      (mockRbac.canCreateOrganizations as ReturnType<typeof vi.fn>).mockReturnValue(false);

      render(
        <PermissionGuard check="canCreateOrganizations">
          <div>Create Org Button</div>
        </PermissionGuard>
      );

      expect(screen.queryByText('Create Org Button')).not.toBeInTheDocument();
    });

    it('should render fallback when canCreateGroups is false', () => {
      (mockRbac.canCreateGroups as ReturnType<typeof vi.fn>).mockReturnValue(false);

      render(
        <PermissionGuard check="canCreateGroups" fallback={<div>No permission</div>}>
          <div>Create Group Button</div>
        </PermissionGuard>
      );

      expect(screen.queryByText('Create Group Button')).not.toBeInTheDocument();
      expect(screen.getByText('No permission')).toBeInTheDocument();
    });

    it('should render children when canInviteGroupOwners is true', () => {
      (mockRbac.canInviteGroupOwners as ReturnType<typeof vi.fn>).mockReturnValue(true);

      render(
        <PermissionGuard check="canInviteGroupOwners">
          <div>Invite Owner Button</div>
        </PermissionGuard>
      );

      expect(screen.getByText('Invite Owner Button')).toBeInTheDocument();
    });
  });

  describe('scoped checks', () => {
    it('should render children when canInviteToGroup returns true for the scopeId', () => {
      (mockRbac.canInviteToGroup as ReturnType<typeof vi.fn>).mockImplementation(
        (groupId: string) => groupId === 'group-123'
      );

      render(
        <PermissionGuard check="canInviteToGroup" scopeId="group-123">
          <div>Invite Member Button</div>
        </PermissionGuard>
      );

      expect(screen.getByText('Invite Member Button')).toBeInTheDocument();
    });

    it('should not render children when canInviteToGroup returns false for the scopeId', () => {
      (mockRbac.canInviteToGroup as ReturnType<typeof vi.fn>).mockImplementation(
        (groupId: string) => groupId === 'group-123'
      );

      render(
        <PermissionGuard check="canInviteToGroup" scopeId="other-group">
          <div>Invite Member Button</div>
        </PermissionGuard>
      );

      expect(screen.queryByText('Invite Member Button')).not.toBeInTheDocument();
    });

    it('should render children when isOwnerOf returns true for the scopeId', () => {
      (mockRbac.isOwnerOf as ReturnType<typeof vi.fn>).mockImplementation(
        (groupId: string) => groupId === 'my-group'
      );

      render(
        <PermissionGuard check="isOwnerOf" scopeId="my-group">
          <div>Owner Section</div>
        </PermissionGuard>
      );

      expect(screen.getByText('Owner Section')).toBeInTheDocument();
    });

    it('should render children when canSeeReporterIdentity returns true', () => {
      (mockRbac.canSeeReporterIdentity as ReturnType<typeof vi.fn>).mockReturnValue(true);

      render(
        <PermissionGuard check="canSeeReporterIdentity" scopeId="group-123">
          <div>Reporter Name: John</div>
        </PermissionGuard>
      );

      expect(screen.getByText('Reporter Name: John')).toBeInTheDocument();
    });

    it('should not render children when canSeeReporterIdentity returns false', () => {
      (mockRbac.canSeeReporterIdentity as ReturnType<typeof vi.fn>).mockReturnValue(false);

      render(
        <PermissionGuard check="canSeeReporterIdentity" scopeId="group-123">
          <div>Reporter Name: John</div>
        </PermissionGuard>
      );

      expect(screen.queryByText('Reporter Name: John')).not.toBeInTheDocument();
    });
  });

  describe('loading state', () => {
    it('should render nothing while loading', () => {
      mockRbac.loading = true;

      render(
        <PermissionGuard check="canCreateGroups">
          <div>Content</div>
        </PermissionGuard>
      );

      expect(screen.queryByText('Content')).not.toBeInTheDocument();
    });

    it('should not render fallback while loading', () => {
      mockRbac.loading = true;

      render(
        <PermissionGuard check="canCreateGroups" fallback={<div>Fallback</div>}>
          <div>Content</div>
        </PermissionGuard>
      );

      expect(screen.queryByText('Fallback')).not.toBeInTheDocument();
    });
  });

  describe('isAdminOf check', () => {
    it('should render children when isAdminOf returns true for the scopeId', () => {
      (mockRbac.isAdminOf as ReturnType<typeof vi.fn>).mockImplementation(
        (orgId: string) => orgId === 'org-123'
      );

      render(
        <PermissionGuard check="isAdminOf" scopeId="org-123">
          <div>Admin Section</div>
        </PermissionGuard>
      );

      expect(screen.getByText('Admin Section')).toBeInTheDocument();
    });
  });
});
