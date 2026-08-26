import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useReports } from './useReports';

const mockApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('./useApi', () => ({
  useApi: () => mockApi,
}));

describe('useReports', () => {
  const mockReports = [
    { id: 'r1', title: 'Report 1', content: 'Content 1', groupId: '1', reporterName: 'John Doe', isAnonymous: false, createdAt: '2024-01-15T10:00:00Z' },
    { id: 'r2', title: 'Report 2', content: 'Content 2', groupId: '1', reporterName: null, isAnonymous: true, createdAt: '2024-01-14T09:00:00Z' },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    mockApi.get.mockResolvedValue(mockReports);
  });

  it('should fetch reports for group', async () => {
    const { result } = renderHook(() => useReports('group-123'));

    await waitFor(() => {
      expect(result.current.reports).toHaveLength(2);
    });
    expect(mockApi.get).toHaveBeenCalledWith('/reports?groupId=group-123');
  });

  it('should handle loading state', () => {
    const { result } = renderHook(() => useReports('group-123'));
    expect(result.current.isLoading).toBe(true);
  });

  it('should handle error state', async () => {
    mockApi.get.mockRejectedValue(new Error('Server error'));
    const { result } = renderHook(() => useReports('group-123'));

    await waitFor(() => {
      expect(result.current.error).toBe('Server error');
    });
  });

  it('should create report and refetch', async () => {
    mockApi.post.mockResolvedValue({ id: 'new-report' });
    const { result } = renderHook(() => useReports('group-123'));

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await result.current.createReport({
      groupId: 'group-123',
      title: 'New Report',
      content: 'Test content',
      isAnonymousToGroup: false,
    });

    expect(mockApi.post).toHaveBeenCalledWith('/reports', {
      groupId: 'group-123',
      title: 'New Report',
      content: 'Test content',
      isAnonymousToGroup: false,
    });
  });

  it('should not fetch if groupId is empty', async () => {
    renderHook(() => useReports(''));
    expect(mockApi.get).not.toHaveBeenCalled();
  });

  it('should show anonymous reports without reporter name', async () => {
    const { result } = renderHook(() => useReports('group-123'));

    await waitFor(() => {
      const anonymousReport = result.current.reports.find(r => r.isAnonymous);
      expect(anonymousReport?.reporterName).toBeNull();
    });
  });
});
