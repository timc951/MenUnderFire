import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { CreateInvitationModal } from './CreateInvitationModal';
import { RbacContextValue } from '../types';

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

const mockApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('../../hooks/useApi', () => ({
  useApi: () => mockApi,
}));

function mockSiteAdmin() {
  mockRbac.loading = false;
  mockRbac.permissions = {
    isSiteAdmin: true,
    adminOfOrganizationIds: [],
    ownedGroupIds: [],
    memberGroupIds: [],
  };
  (mockRbac.isSiteAdmin as ReturnType<typeof vi.fn>).mockReturnValue(true);
  (mockRbac.isOrgAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isGroupOwner as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteOrgAdmins as ReturnType<typeof vi.fn>).mockReturnValue(true);
  (mockRbac.canInviteGroupOwners as ReturnType<typeof vi.fn>).mockReturnValue(true);
  (mockRbac.canInviteToGroup as ReturnType<typeof vi.fn>).mockReturnValue(true);
}

function mockOrgAdmin(orgIds: string[] = ['org-123']) {
  mockRbac.loading = false;
  mockRbac.permissions = {
    isSiteAdmin: false,
    adminOfOrganizationIds: orgIds,
    ownedGroupIds: [],
    memberGroupIds: [],
  };
  (mockRbac.isSiteAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isOrgAdmin as ReturnType<typeof vi.fn>).mockReturnValue(true);
  (mockRbac.isGroupOwner as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteOrgAdmins as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteGroupOwners as ReturnType<typeof vi.fn>).mockReturnValue(true);
  (mockRbac.canInviteToGroup as ReturnType<typeof vi.fn>).mockReturnValue(true);
}

function mockGroupOwner(ownedGroupIds: string[] = ['group-owned']) {
  mockRbac.loading = false;
  mockRbac.permissions = {
    isSiteAdmin: false,
    adminOfOrganizationIds: [],
    ownedGroupIds,
    memberGroupIds: [],
  };
  (mockRbac.isSiteAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isOrgAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isGroupOwner as ReturnType<typeof vi.fn>).mockReturnValue(true);
  (mockRbac.canInviteOrgAdmins as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteGroupOwners as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteToGroup as ReturnType<typeof vi.fn>).mockImplementation(
    (groupId: string) => ownedGroupIds.includes(groupId)
  );
}

function mockGroupMember(memberGroupIds: string[] = ['group-member']) {
  mockRbac.loading = false;
  mockRbac.permissions = {
    isSiteAdmin: false,
    adminOfOrganizationIds: [],
    ownedGroupIds: [],
    memberGroupIds,
  };
  (mockRbac.isSiteAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isOrgAdmin as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.isGroupOwner as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteOrgAdmins as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteGroupOwners as ReturnType<typeof vi.fn>).mockReturnValue(false);
  (mockRbac.canInviteToGroup as ReturnType<typeof vi.fn>).mockReturnValue(false);
}

const mockOrganizations = [
  { id: 'org-1', name: 'Test Organization' },
  { id: 'org-2', name: 'Another Organization' },
];

const mockGroups = [
  { id: 'group-1', name: 'Fitness Group', organizationId: 'org-1' },
  { id: 'group-2', name: 'Study Group', organizationId: 'org-1' },
  { id: 'group-owned', name: 'My Owned Group', organizationId: 'org-1' },
];

describe('CreateInvitationModal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockApi.get.mockImplementation((path: string) => {
      if (path === '/organizations') return Promise.resolve(mockOrganizations);
      if (path === '/groups') return Promise.resolve(mockGroups);
      return Promise.resolve([]);
    });
    mockApi.post.mockResolvedValue({
      id: 'inv-1',
      token: 'test-token-abc123',
      email: 'test@example.com',
      type: 'org_admin',
      status: 'pending',
      expiresAt: '2025-02-01T00:00:00Z',
      createdAt: '2025-01-01T00:00:00Z',
    });
  });

  describe('Site Admin View', () => {
    beforeEach(() => mockSiteAdmin());

    it('should show org admin invitation option', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('radio', { name: /organization admin/i })).toBeInTheDocument();
      });
    });

    it('should show all invitation types', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('radio', { name: /organization admin/i })).toBeInTheDocument();
        expect(screen.getByRole('radio', { name: /group owner/i })).toBeInTheDocument();
        expect(screen.getByRole('radio', { name: /group member/i })).toBeInTheDocument();
      });
    });

    it('should show organization selector when org admin is selected', async () => {
      const user = userEvent.setup();
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('radio', { name: /organization admin/i })).toBeInTheDocument();
      });

      const orgAdminRadio = screen.getByRole('radio', { name: /organization admin/i });
      await user.click(orgAdminRadio);

      await waitFor(() => {
        expect(screen.getByRole('combobox', { name: /organization/i })).toBeInTheDocument();
      });
    });
  });

  describe('Org Admin View', () => {
    beforeEach(() => mockOrgAdmin());

    it('should NOT show org admin invitation option', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.queryByRole('radio', { name: /organization admin/i })).not.toBeInTheDocument();
      });
    });

    it('should show group owner invitation option', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('radio', { name: /group owner/i })).toBeInTheDocument();
      });
    });

    it('should show group member invitation option', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('radio', { name: /group member/i })).toBeInTheDocument();
      });
    });
  });

  describe('Group Owner View', () => {
    beforeEach(() => mockGroupOwner(['group-owned']));

    it('should NOT show org admin invitation option', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.queryByRole('radio', { name: /organization admin/i })).not.toBeInTheDocument();
      });
    });

    it('should NOT show group owner invitation option', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.queryByRole('radio', { name: /group owner/i })).not.toBeInTheDocument();
      });
    });

    it('should show ONLY group member invitation option', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('radio', { name: /group member/i })).toBeInTheDocument();
      });
    });
  });

  describe('Group Member View', () => {
    beforeEach(() => mockGroupMember());

    it('should show no permission message', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByText(/don't have permission to send invitations/i)).toBeInTheDocument();
      });
    });
  });

  describe('Modal behavior', () => {
    beforeEach(() => mockSiteAdmin());

    it('should not render when isOpen is false', () => {
      render(<CreateInvitationModal isOpen={false} onClose={vi.fn()} />);
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('should render when isOpen is true', async () => {
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);
      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });

    it('should call onClose when close button is clicked', async () => {
      const onClose = vi.fn();
      const user = userEvent.setup();
      render(<CreateInvitationModal isOpen onClose={onClose} />);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      const closeButton = screen.getByRole('button', { name: /close/i });
      await user.click(closeButton);

      expect(onClose).toHaveBeenCalled();
    });

    it('should call onClose when cancel button is clicked', async () => {
      const onClose = vi.fn();
      const user = userEvent.setup();
      render(<CreateInvitationModal isOpen onClose={onClose} />);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      const cancelButton = screen.getByRole('button', { name: /cancel/i });
      await user.click(cancelButton);

      expect(onClose).toHaveBeenCalled();
    });
  });

  describe('Form submission', () => {
    beforeEach(() => mockSiteAdmin());

    it('should require email field', async () => {
      const user = userEvent.setup();
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Select org admin type
      const orgAdminRadio = screen.getByRole('radio', { name: /organization admin/i });
      await user.click(orgAdminRadio);

      // Try to submit without email
      const submitButton = screen.getByRole('button', { name: /create invitation/i });
      await user.click(submitButton);

      // The form should not submit (no API call made)
      expect(mockApi.post).not.toHaveBeenCalled();
    });
  });

  describe('Success state with invite link', () => {
    beforeEach(() => mockSiteAdmin());

    it('should show invite link after successful creation', async () => {
      const user = userEvent.setup();
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Select group member type and fill form
      const memberRadio = screen.getByRole('radio', { name: /group member/i });
      await user.click(memberRadio);

      await waitFor(() => {
        expect(screen.getByLabelText(/select group/i)).toBeInTheDocument();
      });

      // Select a group
      fireEvent.change(screen.getByLabelText(/select group/i), { target: { value: 'group-1' } });

      // Fill email
      const emailInput = screen.getByPlaceholderText(/user@example.com/i);
      await user.type(emailInput, 'newuser@example.com');

      // Submit
      const submitButton = screen.getByRole('button', { name: /create invitation/i });
      await user.click(submitButton);

      // Should show success state with invite link
      await waitFor(() => {
        expect(screen.getByText(/invitation sent successfully/i)).toBeInTheDocument();
        expect(screen.getByText(/invite link/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /copy/i })).toBeInTheDocument();
      });
    });

    it('should show Create Another button after success', async () => {
      const user = userEvent.setup();
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Select group member type
      const memberRadio = screen.getByRole('radio', { name: /group member/i });
      await user.click(memberRadio);

      await waitFor(() => {
        expect(screen.getByLabelText(/select group/i)).toBeInTheDocument();
      });

      fireEvent.change(screen.getByLabelText(/select group/i), { target: { value: 'group-1' } });

      const emailInput = screen.getByPlaceholderText(/user@example.com/i);
      await user.type(emailInput, 'newuser@example.com');

      const submitButton = screen.getByRole('button', { name: /create invitation/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /create another/i })).toBeInTheDocument();
      });
    });

    it('should reset form when Create Another is clicked', async () => {
      const user = userEvent.setup();
      render(<CreateInvitationModal isOpen onClose={vi.fn()} />);

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });

      // Create an invitation
      const memberRadio = screen.getByRole('radio', { name: /group member/i });
      await user.click(memberRadio);

      await waitFor(() => {
        expect(screen.getByLabelText(/select group/i)).toBeInTheDocument();
      });

      fireEvent.change(screen.getByLabelText(/select group/i), { target: { value: 'group-1' } });

      const emailInput = screen.getByPlaceholderText(/user@example.com/i);
      await user.type(emailInput, 'newuser@example.com');

      const submitButton = screen.getByRole('button', { name: /create invitation/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /create another/i })).toBeInTheDocument();
      });

      // Click Create Another
      await user.click(screen.getByRole('button', { name: /create another/i }));

      // Should be back to form view
      await waitFor(() => {
        expect(screen.getByRole('radio', { name: /group member/i })).toBeInTheDocument();
        expect(screen.getByPlaceholderText(/user@example.com/i)).toHaveValue('');
      });
    });
  });
});
