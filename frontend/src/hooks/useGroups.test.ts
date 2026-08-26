import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useGroups } from './useGroups';

const mockApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('./useApi', () => ({
  useApi: () => mockApi,
}));

describe('useGroups', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockApi.get.mockResolvedValue([
      { id: '1', name: 'Fitness Group', description: 'Stay fit', inviteCode: 'ABC', createdBy: 'u1', createdAt: '2024-01-01T00:00:00Z', memberCount: 5, role: 'leader' },
      { id: '2', name: 'Study Group', description: 'Study hard', inviteCode: 'DEF', createdBy: 'u2', createdAt: '2024-01-02T00:00:00Z', memberCount: 3, role: 'member' },
    ]);
  });

  it('should fetch groups on mount', async () => {
    const { result } = renderHook(() => useGroups());

    await waitFor(() => {
      expect(result.current.groups).toHaveLength(2);
    });
    expect(mockApi.get).toHaveBeenCalledWith('/groups');
  });

  it('should handle loading state', () => {
    const { result } = renderHook(() => useGroups());
    expect(result.current.isLoading).toBe(true);
  });

  it('should handle error state', async () => {
    mockApi.get.mockRejectedValue(new Error('Network error'));
    const { result } = renderHook(() => useGroups());

    await waitFor(() => {
      expect(result.current.error).toBe('Network error');
    });
  });

  it('should create group and refetch', async () => {
    mockApi.post.mockResolvedValue({ id: 'new' });
    const { result } = renderHook(() => useGroups());

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await result.current.createGroup({ name: 'New', description: 'Desc', organizationId: 'org1' });

    expect(mockApi.post).toHaveBeenCalledWith('/groups', { name: 'New', description: 'Desc', organizationId: 'org1' });
  });

  it('should join group and refetch', async () => {
    mockApi.post.mockResolvedValue({});
    const { result } = renderHook(() => useGroups());

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await result.current.joinGroup('ABC123');

    expect(mockApi.post).toHaveBeenCalledWith('/groups/join', { inviteCode: 'ABC123' });
  });
});
