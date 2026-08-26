import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { http, HttpResponse } from 'msw';
import { server } from '../__mocks__/server';

// Mock useAuth (replaces old react-oidc-context mock)
const mockLogin = vi.fn();
const mockLogout = vi.fn();
const mockGetToken = vi.fn().mockResolvedValue('mock-token');

const mockAuthState = {
  user: { id: 'user-123', email: 'test@example.com', name: 'Test User' } as { id: string; email: string; name: string } | null,
  isAuthenticated: true,
  isLoading: false,
  login: mockLogin,
  logout: mockLogout,
  getToken: mockGetToken,
};

vi.mock('../hooks/useAuth', () => ({
  useAuth: () => mockAuthState,
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return { ...actual, useNavigate: () => mockNavigate };
});

import { AcceptInvitationPage } from './AcceptInvitation';

// Use the same env var resolution as the component
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

function renderWithRouter(token: string = 'valid-token') {
  return render(
    <MemoryRouter initialEntries={[`/invite/${token}`]}>
      <Routes>
        <Route path="/invite/:token" element={<AcceptInvitationPage />} />
      </Routes>
    </MemoryRouter>
  );
}

describe('AcceptInvitationPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockGetToken.mockResolvedValue('mock-token');
    // Reset auth state
    mockAuthState.isAuthenticated = true;
    mockAuthState.isLoading = false;
    mockAuthState.user = { id: 'user-123', email: 'test@example.com', name: 'Test User' };
  });

  describe('Loading state', () => {
    it('should show loading spinner while validating token', () => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return new Promise(() => {}); // Never resolves
        })
      );
      renderWithRouter();
      expect(screen.getByRole('status')).toBeInTheDocument();
    });
  });

  describe('Invalid invitation', () => {
    it('should show error for invalid token', async () => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return HttpResponse.json({ valid: false });
        })
      );
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText(/invalid or expired/i)).toBeInTheDocument();
      });
    });

    it('should show error when API fails', async () => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return new HttpResponse(null, { status: 500 });
        })
      );
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText(/unable to validate/i)).toBeInTheDocument();
      });
    });
  });

  describe('Valid Org Admin invitation', () => {
    beforeEach(() => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return HttpResponse.json({
            valid: true,
            email: 'test@example.com',
            type: 'org_admin',
            organizationName: 'Test Organization',
            groupName: null,
            inviterName: 'Admin User',
            expiresAt: '2025-01-01T00:00:00Z',
            capabilities: [
              'Create groups within the organization',
              'Invite group owners to manage groups',
              'Invite members to any group',
              'View all reports in the organization',
            ],
          });
        })
      );
    });

    it('should display invitation details', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText(/you've been invited/i)).toBeInTheDocument();
      });

      expect(screen.getByText('Admin User')).toBeInTheDocument();
      expect(screen.getByText('Test Organization')).toBeInTheDocument();
      // RoleBadge renders "Org Admin" for org_admin role
      expect(screen.getByText('Org Admin')).toBeInTheDocument();
    });

    it('should display capabilities list', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText(/create groups within the organization/i)).toBeInTheDocument();
      });

      expect(screen.getByText(/invite group owners/i)).toBeInTheDocument();
    });
  });

  describe('Valid Group Owner invitation', () => {
    beforeEach(() => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return HttpResponse.json({
            valid: true,
            email: 'test@example.com',
            type: 'group_owner',
            organizationName: null,
            groupName: 'Fitness Group',
            inviterName: 'Org Admin',
            expiresAt: '2025-01-01T00:00:00Z',
            capabilities: [
              'Manage the group settings',
              'Invite members to the group',
              'View who submitted anonymous reports',
              'Note: You will NOT be able to create new groups',
            ],
          });
        })
      );
    });

    it('should display group owner badge', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText('Owner')).toBeInTheDocument();
      });
    });

    it('should display group name', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText('Fitness Group')).toBeInTheDocument();
      });
    });

    it('should display capability limitation', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText(/will not be able to create new groups/i)).toBeInTheDocument();
      });
    });
  });

  describe('Valid Group Member invitation', () => {
    beforeEach(() => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return HttpResponse.json({
            valid: true,
            email: 'test@example.com',
            type: 'group_member',
            organizationName: null,
            groupName: 'Study Group',
            inviterName: 'Group Owner',
            expiresAt: '2025-01-01T00:00:00Z',
            capabilities: [
              'Submit accountability reports',
              'View group reports',
              'Participate in group activities',
            ],
          });
        })
      );
    });

    it('should display member badge', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByText('Member')).toBeInTheDocument();
      });
    });
  });

  describe('Accept invitation', () => {
    beforeEach(() => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return HttpResponse.json({
            valid: true,
            email: 'test@example.com',
            type: 'org_admin',
            organizationName: 'Test Organization',
            groupName: null,
            inviterName: 'Admin User',
            expiresAt: '2025-01-01T00:00:00Z',
            capabilities: [],
          });
        }),
        http.post(`${API_BASE}/invitations/accept`, () => {
          return HttpResponse.json({
            invitation: { id: 'inv-1' },
            user: { id: 'user-1', email: 'test@example.com' },
          });
        })
      );
    });

    it('should show accept button', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /accept invitation/i })).toBeInTheDocument();
      });
    });

    it('should call API when accept button is clicked', async () => {
      const user = userEvent.setup();
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /accept invitation/i })).toBeInTheDocument();
      });

      await user.click(screen.getByRole('button', { name: /accept invitation/i }));

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/');
      });
    });

    it('should navigate on success', async () => {
      const user = userEvent.setup();
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /accept invitation/i })).toBeInTheDocument();
      });

      await user.click(screen.getByRole('button', { name: /accept invitation/i }));

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/');
      });
    });

    it('should show error on API failure', async () => {
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return HttpResponse.json({
            valid: true,
            email: 'test@example.com',
            type: 'org_admin',
            organizationName: 'Test Organization',
            groupName: null,
            inviterName: 'Admin User',
            expiresAt: '2025-01-01T00:00:00Z',
            capabilities: [],
          });
        }),
        http.post(`${API_BASE}/invitations/accept`, () => {
          return HttpResponse.json(
            { message: 'Failed to accept invitation' },
            { status: 500 }
          );
        })
      );

      const user = userEvent.setup();
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /accept invitation/i })).toBeInTheDocument();
      });

      await user.click(screen.getByRole('button', { name: /accept invitation/i }));

      await waitFor(() => {
        expect(screen.getByText(/failed to accept/i)).toBeInTheDocument();
      });
    });
  });

  describe('Unauthenticated user', () => {
    beforeEach(() => {
      mockAuthState.isAuthenticated = false;
      mockAuthState.user = null;
      server.use(
        http.get(`${API_BASE}/invitations/validate/:token`, () => {
          return HttpResponse.json({
            valid: true,
            email: 'test@example.com',
            type: 'org_admin',
            organizationName: 'Test Organization',
            groupName: null,
            inviterName: 'Admin User',
            expiresAt: '2025-01-01T00:00:00Z',
            capabilities: [],
          });
        })
      );
    });

    it('should show login button for unauthenticated users', async () => {
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /sign up.*accept/i })).toBeInTheDocument();
      });
    });

    it('should call login when sign up button is clicked', async () => {
      const user = userEvent.setup();
      renderWithRouter();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /sign up.*accept/i })).toBeInTheDocument();
      });

      await user.click(screen.getByRole('button', { name: /sign up.*accept/i }));

      expect(mockLogin).toHaveBeenCalled();
    });
  });
});
