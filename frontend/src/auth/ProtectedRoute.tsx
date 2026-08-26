import { RedirectToSignIn } from '@clerk/react';
import { useAuth } from '../hooks/useAuth';
import { ReactNode, useEffect, useRef, useState } from 'react';
import { TestingAgreement } from '../components/TestingAgreement';

interface ProtectedRouteProps {
  children: ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, getToken, logout } = useAuth();
  const [isSyncing, setIsSyncing] = useState(false);
  const [hasSynced, setHasSynced] = useState(false);
  const [syncError, setSyncError] = useState<string | null>(null);
  const [hasAcceptedAgreement, setHasAcceptedAgreement] = useState<boolean | null>(null);
  const prevAuthRef = useRef(isAuthenticated);

  // Reset sync state when auth state changes (e.g. logout then re-login)
  // so the agreement check runs fresh for each session
  useEffect(() => {
    if (isAuthenticated && !prevAuthRef.current) {
      setHasSynced(false);
      setHasAcceptedAgreement(null);
      setSyncError(null);
    }
    prevAuthRef.current = isAuthenticated;
  }, [isAuthenticated]);

  useEffect(() => {
    async function syncUser() {
      if (!isAuthenticated || hasSynced || isSyncing || syncError) return;

      setIsSyncing(true);

      try {
        const token = await getToken();

        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

        const response = await fetch(`${apiUrl}/users/me`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        // If 401, token is invalid/expired - redirect to login
        if (response.status === 401) {
          console.warn('Token expired or invalid, redirecting to login...');
          logout();
          return;
        }

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`Failed to sync user: ${response.status} ${errorText}`);
        }

        const userData = await response.json();
        setHasAcceptedAgreement(!!userData.agreementAcceptedAt);
        setHasSynced(true);
      } catch (error) {
        // Check if error is from token refresh failure
        if (error instanceof Error && error.message.includes('No access token')) {
          console.warn('Token refresh failed, redirecting to login...');
          // Clerk will handle redirect automatically via RedirectToSignIn
          return;
        }
        console.error('User sync failed:', error);
        setSyncError(error instanceof Error ? error.message : 'Unknown error');
      } finally {
        setIsSyncing(false);
      }
    }

    syncUser();
  }, [isAuthenticated, hasSynced, isSyncing, syncError, getToken, logout]);

  if (isLoading) {
    return <div data-testid="auth-loading">Loading...</div>;
  }

  if (!isAuthenticated) {
    return <RedirectToSignIn />;
  }

  if (syncError) {
    return (
      <div data-testid="sync-error" className="p-4 text-red-600">
        <p>Failed to sync your account: {syncError}</p>
        <button
          onClick={() => {
            setHasSynced(false);
            setSyncError(null);
          }}
          className="mt-2 underline"
        >
          Retry
        </button>
      </div>
    );
  }

  if (isSyncing || !hasSynced) {
    return <div data-testid="user-syncing">Setting up your account...</div>;
  }

  if (hasAcceptedAgreement !== true) {
    return (
      <TestingAgreement
        onAccepted={() => setHasAcceptedAgreement(true)}
        onDecline={logout}
        getToken={getToken}
      />
    );
  }

  return <>{children}</>;
}
