import { useState, useEffect, FormEvent } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';

interface JoinResponse {
  message: string;
  groupId?: string;
}

export function GroupJoin() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const api = useApi();

  const codeFromUrl = searchParams.get('code') || '';
  const [inviteCode, setInviteCode] = useState(codeFromUrl);
  const [isJoining, setIsJoining] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<JoinResponse | null>(null);

  useEffect(() => {
    if (codeFromUrl) {
      handleJoin(codeFromUrl);
    }
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleJoin = async (code: string) => {
    if (!code.trim()) {
      setError('Invite code is required');
      return;
    }
    setIsJoining(true);
    setError(null);
    try {
      const data = await api.post<JoinResponse>('/groups/join', { inviteCode: code.trim() });
      setSuccess(data);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to join group';
      if (message.includes('already a member')) {
        setError('You are already a member of this group');
      } else if (message.includes('expired')) {
        setError('This invite code has expired');
      } else if (message.includes('not found')) {
        setError('Invalid invite code or group not found');
      } else {
        setError(message);
      }
    } finally {
      setIsJoining(false);
    }
  };

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    handleJoin(inviteCode);
  };

  if (isJoining) {
    return (
      <div className="flex flex-col items-center justify-center py-16">
        <LoadingSpinner size="lg" />
        <p className="mt-4 text-gray-600 dark:text-stone-400">Joining group...</p>
      </div>
    );
  }

  if (success) {
    return (
      <div className="max-w-md mx-auto py-16">
        <Card>
          <div className="text-center space-y-4">
            <div className="w-12 h-12 mx-auto bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center">
              <svg className="w-6 h-6 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Successfully joined the group!</h2>
            <p className="text-sm text-gray-600 dark:text-stone-400">{success.message}</p>
            <Button onClick={() => navigate('/groups')}>
              Go to Groups
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="max-w-md mx-auto py-16">
      <Card>
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Join a Group</h2>
        {error && (
          <div className="mb-4 p-3 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800">
            <p className="text-sm text-red-700 dark:text-red-400">{error}</p>
          </div>
        )}
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="inviteCode" className="block text-sm font-medium text-gray-700 dark:text-stone-300 mb-1">
              Invite Code
            </label>
            <input
              id="inviteCode"
              type="text"
              value={inviteCode}
              onChange={(e) => setInviteCode(e.target.value)}
              placeholder="Enter invite code"
              className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-lg bg-white dark:bg-stone-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-amber-500"
            />
          </div>
          <div className="flex gap-3 justify-end">
            <Button variant="secondary" type="button" onClick={() => navigate('/groups')}>
              Cancel
            </Button>
            <Button type="submit">
              Join Group
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
