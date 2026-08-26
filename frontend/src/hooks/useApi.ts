import { useMemo } from 'react';
import { useAuth } from './useAuth';
import { createApiClient } from '../services/api';

export function useApi() {
  const { getToken } = useAuth();
  return useMemo(() => createApiClient(getToken), [getToken]);
}
