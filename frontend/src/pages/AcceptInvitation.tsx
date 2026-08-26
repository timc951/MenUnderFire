import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { RoleBadge } from '../rbac/components/RoleBadge';
import { RoleType } from '../rbac/types';

interface InvitationValidation {
  valid: boolean;
  email?: string;
  type?: 'org_admin' | 'group_owner' | 'group_member';
  organizationName?: string | null;
  groupName?: string | null;
  inviterName?: string;
  expiresAt?: string;
  capabilities?: string[];
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const roleLabels: Record<string, string> = {
  org_admin: 'Organization Admin',
  group_owner: 'Group Owner',
  group_member: 'Group Member',
};

const roleTypeMap: Record<string, RoleType> = {
  org_admin: 'org_admin',
  group_owner: 'group_owner',
  group_member: 'member',
  // Handle possible variations from backend
  ORG_ADMIN: 'org_admin',
  GROUP_OWNER: 'group_owner',
  GROUP_MEMBER: 'member',
};

function getRoleType(type: string | undefined): RoleType {
  if (!type) return 'member';
  return roleTypeMap[type] || roleTypeMap[type.toLowerCase()] || 'member';
}

export function AcceptInvitationPage() {
  const { token } = useParams<{ token: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading, user, login, getToken } = useAuth();

  const [invitation, setInvitation] = useState<InvitationValidation | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isAccepting, setIsAccepting] = useState(false);
  const [acceptError, setAcceptError] = useState<string | null>(null);

  // Validate token on mount (public endpoint - no auth required)
  useEffect(() => {
    if (!token) {
      setError('No invitation token provided');
      setIsLoading(false);
      return;
    }

    async function validateToken() {
      try {
        // Public endpoint - no auth required
        const response = await fetch(`${API_BASE_URL}/invitations/validate/${token}`);
        if (!response.ok) {
          throw new Error('Failed to validate invitation');
        }
        const data: InvitationValidation = await response.json();
        setInvitation(data);
        if (!data.valid) {
          setError('This invitation is invalid or expired');
        }
      } catch {
        setError('Unable to validate invitation');
      } finally {
        setIsLoading(false);
      }
    }

    validateToken();
  }, [token]);

  const handleAccept = async () => {
    if (!token || !user) return;

    setIsAccepting(true);
    setAcceptError(null);

    try {
      const accessToken = await getToken();

      const response = await fetch(`${API_BASE_URL}/invitations/accept`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          token,
          externalId: user.id,
          displayName: user.name || user.email || 'User',
        }),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: 'Failed to accept invitation' }));
        throw new Error(errorData.message || 'Failed to accept invitation');
      }

      navigate('/');
    } catch (err) {
      setAcceptError(err instanceof Error ? err.message : 'Failed to accept invitation');
    } finally {
      setIsAccepting(false);
    }
  };

  const handleSignUp = () => {
    // Redirect back to this invitation page after sign-in
    login(window.location.href);
  };

  // Show loading while auth is initializing
  if (authLoading || isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-50 dark:bg-gray-900">
        <div role="status" className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
      </div>
    );
  }

  if (error || !invitation?.valid) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
        <div className="max-w-md w-full p-6 bg-white dark:bg-gray-800 rounded-lg shadow text-center">
          <div className="text-red-500 mb-4">
            <svg className="w-16 h-16 mx-auto" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                clipRule="evenodd"
              />
            </svg>
          </div>
          <h1 className="text-2xl font-bold mb-2 text-gray-900 dark:text-white">Invalid Invitation</h1>
          <p className="text-gray-600 dark:text-gray-400">{error || 'This invitation is invalid or expired.'}</p>
          <button
            onClick={() => navigate('/')}
            className="mt-6 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
          >
            Go Home
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center p-4">
      <div className="max-w-md w-full p-6 bg-white dark:bg-gray-800 rounded-lg shadow">
        <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-white">You've Been Invited!</h1>

        <div className="mb-4">
          <p className="text-gray-600 dark:text-gray-400">
            <span className="font-semibold text-gray-900 dark:text-white">{invitation.inviterName}</span> has invited you to join as:
          </p>
          <div className="mt-2">
            <RoleBadge role={getRoleType(invitation.type)} />
          </div>
        </div>

        {invitation.organizationName && (
          <p className="mb-2 text-gray-700 dark:text-gray-300">
            <span className="font-medium">Organization:</span> {invitation.organizationName}
          </p>
        )}

        {invitation.groupName && (
          <p className="mb-2 text-gray-700 dark:text-gray-300">
            <span className="font-medium">Group:</span> {invitation.groupName}
          </p>
        )}

        {invitation.capabilities && invitation.capabilities.length > 0 && (
          <div className="mt-4 p-3 bg-gray-50 dark:bg-gray-700 rounded">
            <p className="font-medium mb-2 text-gray-900 dark:text-white">
              As a {roleLabels[invitation.type || 'group_member']}, you will be able to:
            </p>
            <ul className="list-disc list-inside text-sm text-gray-600 dark:text-gray-400 space-y-1">
              {invitation.capabilities.map((cap, i) => (
                <li key={i}>{cap}</li>
              ))}
            </ul>
          </div>
        )}

        {acceptError && (
          <div className="mt-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-md text-sm">
            {acceptError}
          </div>
        )}

        {isAuthenticated ? (
          <div className="mt-6">
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
              Signed in as <span className="font-medium text-gray-900 dark:text-white">{user?.email}</span>
            </p>
            <button
              onClick={handleAccept}
              disabled={isAccepting}
              className="w-full bg-primary-600 text-white py-2 rounded-md hover:bg-primary-700 disabled:opacity-50"
            >
              {isAccepting ? 'Accepting...' : 'Accept Invitation'}
            </button>
          </div>
        ) : (
          <div className="mt-6">
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
              Create an account or sign in to accept this invitation.
            </p>
            <button
              onClick={handleSignUp}
              className="w-full bg-primary-600 text-white py-2 rounded-md hover:bg-primary-700"
            >
              Sign Up / Sign In to Accept
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
