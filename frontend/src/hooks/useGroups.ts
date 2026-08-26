import { useState, useEffect, useCallback } from 'react';
import { useApi } from './useApi';
import { Group, CreateGroupRequest } from '../types';

export function useGroups() {
  const api = useApi();
  const [groups, setGroups] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchGroups = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await api.get<Group[]>('/groups');
      setGroups(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load groups');
    } finally {
      setIsLoading(false);
    }
  }, [api]);

  useEffect(() => {
    fetchGroups();
  }, [fetchGroups]);

  const createGroup = async (data: CreateGroupRequest): Promise<void> => {
    await api.post('/groups', data);
    await fetchGroups();
  };

  const joinGroup = async (inviteCode: string): Promise<void> => {
    await api.post('/groups/join', { inviteCode });
    await fetchGroups();
  };

  return { groups, isLoading, error, createGroup, joinGroup, refetch: fetchGroups };
}
