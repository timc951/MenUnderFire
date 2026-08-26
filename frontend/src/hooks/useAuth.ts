import { useUser, useAuth as useClerkAuth, useClerk } from '@clerk/react';
import { useCallback } from 'react';

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  picture?: string;
}

export function useAuth() {
  const { user: clerkUser, isLoaded: userLoaded, isSignedIn } = useUser();
  const { getToken, isLoaded: authLoaded } = useClerkAuth();
  const clerk = useClerk();

  const getTokenFn = useCallback(async (): Promise<string> => {
    // Clerk handles token refresh automatically
    // Use custom template to include email/name claims for backend
    const token = await getToken({ template: 'menunderfire' });
    if (!token) {
      throw new Error('No access token available');
    }
    return token;
  }, [getToken]);

  const login = useCallback((redirectUrl?: string) => {
    clerk.redirectToSignIn({ redirectUrl: redirectUrl || '/dashboard' });
  }, [clerk]);

  const logout = useCallback(() => {
    clerk.signOut();
  }, [clerk]);

  const authUser: AuthUser | null = clerkUser
    ? {
        id: clerkUser.id,
        email: clerkUser.primaryEmailAddress?.emailAddress || '',
        name: clerkUser.fullName || clerkUser.firstName || '',
        picture: clerkUser.imageUrl,
      }
    : null;

  return {
    user: authUser,
    isAuthenticated: !!isSignedIn,
    isLoading: !userLoaded || !authLoaded,
    login,
    logout,
    getToken: getTokenFn,
  };
}
