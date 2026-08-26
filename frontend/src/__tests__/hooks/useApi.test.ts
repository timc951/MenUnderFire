import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { createApiClient, ApiError } from '../../services/api';

describe('API Client', () => {
  const mockGetToken = vi.fn().mockResolvedValue('test-access-token');
  let api: ReturnType<typeof createApiClient>;

  beforeEach(() => {
    api = createApiClient(mockGetToken);
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should include auth token in requests', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      json: () => Promise.resolve({ data: 'test' }),
    });

    await api.get('/users/me');

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining('/users/me'),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer test-access-token',
        }),
      })
    );
  });

  it('should throw ApiError on 401 response', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 401,
      json: () => Promise.resolve({ message: 'Unauthorized' }),
    });

    await expect(api.get('/protected')).rejects.toThrow(ApiError);
    await expect(api.get('/protected')).rejects.toMatchObject({
      status: 401,
    });
  });

  it('should throw ApiError on non-ok response', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: false,
      status: 404,
      json: () => Promise.resolve({ message: 'Not found' }),
    });

    await expect(api.get('/missing')).rejects.toThrow(ApiError);
  });

  it('should handle network errors gracefully', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Network error'));

    await expect(api.get('/any')).rejects.toThrow('Network error');
  });

  it('should send POST with JSON body', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 201,
      json: () => Promise.resolve({ id: '123' }),
    });

    await api.post('/groups', { name: 'Test Group' });

    expect(global.fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ name: 'Test Group' }),
      })
    );
  });

  it('should handle 204 No Content responses', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 204,
      json: () => Promise.reject(new Error('no body')),
    });

    const result = await api.delete('/groups/123/members/456');

    expect(result).toBeUndefined();
  });

  it('should call getToken before each request', async () => {
    (global.fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
      ok: true,
      status: 200,
      json: () => Promise.resolve({}),
    });

    await api.get('/test1');
    await api.get('/test2');

    expect(mockGetToken).toHaveBeenCalledTimes(2);
  });
});
