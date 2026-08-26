import { render, RenderOptions } from '@testing-library/react';
import { ReactElement, ReactNode } from 'react';
import { vi } from 'vitest';

export interface MockAuthOptions {
  isAuthenticated?: boolean;
  isLoading?: boolean;
  error?: Error | null;
  user?: {
    access_token?: string;
    profile: {
      sub?: string;
      email?: string;
      name?: string;
      preferred_username?: string;
      picture?: string;
    };
  } | null;
}

const defaultMockAuth: MockAuthOptions = {
  isAuthenticated: true,
  isLoading: false,
  error: null,
  user: {
    access_token: 'mock-access-token',
    profile: {
      sub: 'keycloak|test123',
      email: 'test@example.com',
      name: 'Test User',
    },
  },
};

export function createMockOidcAuth(options: MockAuthOptions = {}) {
  const merged = { ...defaultMockAuth, ...options };

  return {
    useAuth: () => ({
      isAuthenticated: merged.isAuthenticated,
      isLoading: merged.isLoading,
      error: merged.error,
      user: merged.user,
      signinRedirect: vi.fn(),
      signoutRedirect: vi.fn(),
    }),
    AuthProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
  };
}

export function renderWithAuth(
  ui: ReactElement,
  options?: {
    authOptions?: MockAuthOptions;
    renderOptions?: Omit<RenderOptions, 'wrapper'>;
  }
) {
  const { authOptions, renderOptions } = options || {};
  const mockAuth = createMockOidcAuth(authOptions);

  // Mock the oidc module before rendering
  vi.doMock('react-oidc-context', () => mockAuth);

  return render(ui, renderOptions);
}
