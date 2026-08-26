import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Header } from './Header';

const renderWithRouter = (component: React.ReactElement) => {
  return render(<BrowserRouter>{component}</BrowserRouter>);
};

const mockLogout = vi.fn();
const mockUseAuth = vi.fn();
const mockUseRbac = vi.fn();

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock('../../rbac', () => ({
  useRbac: () => mockUseRbac(),
  RoleBadge: ({ role }: { role: string }) => <span data-testid="role-badge">{role}</span>,
}));

describe('Header', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue({
      user: { id: '1', name: 'Test User', email: 'test@test.com' },
      logout: mockLogout,
    });
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => false,
      isOrgAdmin: () => false,
      isGroupOwner: () => true,
      loading: false,
    });
  });

  it('should display app name', () => {
    renderWithRouter(<Header />);
    expect(screen.getByText('Men Under Fire')).toBeInTheDocument();
  });

  it('should display user name when authenticated', () => {
    renderWithRouter(<Header />);
    expect(screen.getByText('Test User')).toBeInTheDocument();
  });

  it('should call logout when logout button is clicked', async () => {
    renderWithRouter(<Header />);
    await userEvent.click(screen.getByText('Logout'));
    expect(mockLogout).toHaveBeenCalledTimes(1);
  });

  it('should not show user info when not authenticated', () => {
    mockUseAuth.mockReturnValue({ user: null, logout: mockLogout });
    renderWithRouter(<Header />);
    expect(screen.queryByText('Test User')).not.toBeInTheDocument();
  });

  it('should display role badge when permissions are loaded', () => {
    renderWithRouter(<Header />);
    expect(screen.getByTestId('role-badge')).toBeInTheDocument();
    expect(screen.getByText('group_owner')).toBeInTheDocument();
  });

  it('should not display role badge while loading', () => {
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => false,
      isOrgAdmin: () => false,
      isGroupOwner: () => false,
      loading: true,
    });
    renderWithRouter(<Header />);
    expect(screen.queryByTestId('role-badge')).not.toBeInTheDocument();
  });

  it('should display site_admin role for site admins', () => {
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => true,
      isOrgAdmin: () => false,
      isGroupOwner: () => false,
      loading: false,
    });
    renderWithRouter(<Header />);
    expect(screen.getByText('site_admin')).toBeInTheDocument();
  });

  it('should display org_admin role for org admins', () => {
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => false,
      isOrgAdmin: () => true,
      isGroupOwner: () => false,
      loading: false,
    });
    renderWithRouter(<Header />);
    expect(screen.getByText('org_admin')).toBeInTheDocument();
  });

  it('should display member role when user has no special roles', () => {
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => false,
      isOrgAdmin: () => false,
      isGroupOwner: () => false,
      loading: false,
    });
    renderWithRouter(<Header />);
    expect(screen.getByText('member')).toBeInTheDocument();
  });
});
