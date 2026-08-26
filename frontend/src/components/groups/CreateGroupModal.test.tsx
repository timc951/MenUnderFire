import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { CreateGroupModal } from './CreateGroupModal';

// Mock useOrganizations - return a single org so it auto-selects
vi.mock('../../hooks/useOrganizations', () => ({
  useOrganizations: () => ({
    organizations: [{ id: 'org1', name: 'Test Org', createdAt: '2024-01-01' }],
    isLoading: false,
  }),
}));

// Mock useRbac - org admin with access to the single org
vi.mock('../../rbac', () => ({
  useRbac: () => ({
    isSiteAdmin: () => false,
    isOrgAdmin: () => true,
    permissions: { adminOfOrganizationIds: ['org1'] },
  }),
}));

describe('CreateGroupModal', () => {
  const defaultProps = {
    isOpen: true,
    onClose: vi.fn(),
    onSubmit: vi.fn().mockResolvedValue(undefined),
  };

  beforeEach(() => {
    vi.clearAllMocks();
    defaultProps.onSubmit.mockResolvedValue(undefined);
  });

  it('should render form fields when open', () => {
    render(<CreateGroupModal {...defaultProps} />);
    expect(screen.getByLabelText(/group name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument();
  });

  it('should validate required name field', async () => {
    render(<CreateGroupModal {...defaultProps} />);
    await userEvent.click(screen.getByRole('button', { name: /create group/i }));
    expect(screen.getByText(/group name is required/i)).toBeInTheDocument();
  });

  it('should submit form with data', async () => {
    render(<CreateGroupModal {...defaultProps} />);
    await userEvent.type(screen.getByLabelText(/group name/i), 'New Group');
    await userEvent.type(screen.getByLabelText(/description/i), 'A description');
    await userEvent.click(screen.getByRole('button', { name: /create group/i }));

    await waitFor(() => {
      expect(defaultProps.onSubmit).toHaveBeenCalledWith({
        name: 'New Group',
        description: 'A description',
        organizationId: 'org1',
      });
    });
  });

  it('should call onClose when cancel is clicked', async () => {
    render(<CreateGroupModal {...defaultProps} />);
    await userEvent.click(screen.getByRole('button', { name: /cancel/i }));
    expect(defaultProps.onClose).toHaveBeenCalled();
  });

  it('should not render when isOpen is false', () => {
    render(<CreateGroupModal {...defaultProps} isOpen={false} />);
    expect(screen.queryByLabelText(/group name/i)).not.toBeInTheDocument();
  });
});
