import { ClerkProvider } from '@clerk/react';
import { ReactNode } from 'react';

interface AuthProviderProps {
  children: ReactNode;
}

const PUBLISHABLE_KEY = import.meta.env.VITE_AUTH_PUBLISHABLE_KEY;

if (!PUBLISHABLE_KEY) {
  throw new Error('Missing VITE_AUTH_PUBLISHABLE_KEY environment variable');
}

export function AuthProvider({ children }: AuthProviderProps) {
  return (
    <ClerkProvider
      publishableKey={PUBLISHABLE_KEY}
      signInFallbackRedirectUrl="/dashboard"
      signUpFallbackRedirectUrl="/dashboard"
    >
      {children}
    </ClerkProvider>
  );
}
