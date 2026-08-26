import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { http, HttpResponse, delay } from 'msw';
import { server } from '../../__mocks__/server';

// Mock useAuth hook (replaces old react-oidc-context mock)
const mockGetToken = vi.fn();
const mockLogout = vi.fn();
const mockLogin = vi.fn();
const mockUseAuth = vi.fn();

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => mockUseAuth(),
}));

// Mock Clerk's RedirectToSignIn component
vi.mock('@clerk/react', () => ({
  RedirectToSignIn: () => 'Redirecting to sign in',
}));

import { ProtectedRoute } from '../../auth/ProtectedRoute';

// Use the same env var resolution as the component
const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockGetToken.mockResolvedValue('test-token');
  });

  it('should render children when authenticated and synced', async () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      getToken: mockGetToken,
      logout: mockLogout,
    });

    // Add MSW handler for the user sync endpoint
    server.use(
      http.get(`${API_BASE}/users/me`, () => {
        // agreementAcceptedAt must be present, or ProtectedRoute renders the
        // testing-agreement gate instead of the children.
        return HttpResponse.json({
          id: '123',
          email: 'test@test.com',
          agreementAcceptedAt: '2026-01-01T00:00:00Z',
        });
      })
    );

    render(
      <ProtectedRoute>
        <div data-testid="protected-content">Secret Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByTestId('protected-content')).toBeInTheDocument();
    });
  });

  it('should show the testing agreement when it has not been accepted', async () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      getToken: mockGetToken,
      logout: mockLogout,
    });

    server.use(
      http.get(`${API_BASE}/users/me`, () => {
        return HttpResponse.json({ id: '123', email: 'test@test.com', agreementAcceptedAt: null });
      })
    );

    render(
      <ProtectedRoute>
        <div data-testid="protected-content">Secret Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
    });
  });

  it('should redirect to login when not authenticated', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
      getToken: mockGetToken,
      logout: mockLogout,
    });

    render(
      <ProtectedRoute>
        <div>Secret Content</div>
      </ProtectedRoute>
    );

    // RedirectToSignIn should be rendered (mocked to return text)
    expect(screen.getByText('Redirecting to sign in')).toBeInTheDocument();
  });

  it('should show loading state while checking auth', () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: true,
      getToken: mockGetToken,
      logout: mockLogout,
    });

    render(
      <ProtectedRoute>
        <div>Secret Content</div>
      </ProtectedRoute>
    );

    expect(screen.getByTestId('auth-loading')).toBeInTheDocument();
  });

  it('should show syncing state while setting up account', async () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      getToken: mockGetToken,
      logout: mockLogout,
    });

    // Override MSW handler to delay the response
    server.use(
      http.get(`${API_BASE}/users/me`, async () => {
        await delay(100);
        return HttpResponse.json({ id: '123', email: 'test@test.com' });
      })
    );

    render(
      <ProtectedRoute>
        <div data-testid="protected-content">Secret Content</div>
      </ProtectedRoute>
    );

    expect(screen.getByTestId('user-syncing')).toBeInTheDocument();
  });

  it('should show error when sync fails', async () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      getToken: mockGetToken,
      logout: mockLogout,
    });

    // Override MSW handler to return an error
    server.use(
      http.get(`${API_BASE}/users/me`, () => {
        return new HttpResponse('Server error', { status: 500 });
      })
    );

    render(
      <ProtectedRoute>
        <div data-testid="protected-content">Secret Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByTestId('sync-error')).toBeInTheDocument();
    });
  });

  it('should logout when 401 is returned', async () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      getToken: mockGetToken,
      logout: mockLogout,
    });

    // Override MSW handler to return 401 (token expired)
    server.use(
      http.get(`${API_BASE}/users/me`, () => {
        return new HttpResponse('Unauthorized', { status: 401 });
      })
    );

    render(
      <ProtectedRoute>
        <div data-testid="protected-content">Secret Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(mockLogout).toHaveBeenCalled();
    });
  });
});
