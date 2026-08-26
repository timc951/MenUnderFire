import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';

const mockClerkProvider = vi.fn();

vi.mock('@clerk/react', () => ({
  ClerkProvider: (props: { children: React.ReactNode; [key: string]: unknown }) => {
    mockClerkProvider(props);
    return <>{props.children}</>;
  },
}));

import { AuthProvider } from '../../auth/AuthProvider';

describe('AuthProvider', () => {
  it('should render children inside ClerkProvider', () => {
    render(
      <AuthProvider>
        <div data-testid="child">Hello</div>
      </AuthProvider>
    );

    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('should configure ClerkProvider with correct props', () => {
    render(
      <AuthProvider>
        <div>Test</div>
      </AuthProvider>
    );

    expect(mockClerkProvider).toHaveBeenCalled();
    const props = mockClerkProvider.mock.calls[0][0];
    expect(props.publishableKey).toBeDefined();
    expect(props.signInFallbackRedirectUrl).toBe('/dashboard');
    expect(props.signUpFallbackRedirectUrl).toBe('/dashboard');
  });
});
