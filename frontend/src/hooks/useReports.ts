import { useState, useEffect, useCallback } from 'react';
import { useApi } from './useApi';
import { Report, CreateReportRequest } from '../types';

export function useReports(groupId: string) {
  const api = useApi();
  const [reports, setReports] = useState<Report[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchReports = useCallback(async () => {
    if (!groupId) return;
    setIsLoading(true);
    setError(null);
    try {
      const data = await api.get<Report[]>(`/reports?groupId=${groupId}`);
      setReports(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load reports');
    } finally {
      setIsLoading(false);
    }
  }, [api, groupId]);

  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  const createReport = async (data: CreateReportRequest): Promise<unknown> => {
    const result = await api.post('/reports', data);
    await fetchReports();
    return result;
  };

  return { reports, isLoading, error, createReport, refetch: fetchReports };
}
