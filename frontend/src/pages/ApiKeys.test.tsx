import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { ApiKeys } from './ApiKeys';

const mockApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('../hooks/useApi', () => ({
  useApi: () => mockApi,
}));

describe('ApiKeys Page', () => {
  beforeEach(() => {
    mockApi.get.mockResolvedValue([
      { id: 'key1', name: 'My Integration Key', keyPrefix: 'ak_abc1', permissions: { reports: ['read', 'write'], groups: ['read'] }, expiresAt: '2025-01-01T00:00:00Z', createdAt: '2024-01-01T00:00:00Z', lastUsedAt: null },
    ]);
  });

  it('should display list of existing API keys', async () => {
    render(<ApiKeys />);
    expect(await screen.findByText('My Integration Key')).toBeInTheDocument();
    expect(screen.getByText(/expires/i)).toBeInTheDocument();
  });

  it('should show "no keys" message when empty', async () => {
    mockApi.get.mockResolvedValue([]);
    render(<ApiKeys />);
    await waitFor(() => {
      expect(screen.getByText(/no api keys/i)).toBeInTheDocument();
    });
  });

  it('should generate new API key', async () => {
    mockApi.post.mockResolvedValue({ id: 'new', name: 'New Key', key: 'ak_test1234567890', permissions: {}, expiresAt: null });
    render(<ApiKeys />);
    await screen.findByText('My Integration Key');

    await userEvent.click(screen.getByRole('button', { name: /generate/i }));
    await userEvent.type(screen.getByLabelText(/name/i), 'New Key');
    await userEvent.click(screen.getByRole('button', { name: /create/i }));

    expect(await screen.findByText(/copy this key now/i)).toBeInTheDocument();
    expect(screen.getByText('ak_test1234567890')).toBeInTheDocument();
  });

  it('should revoke API key after confirmation', async () => {
    mockApi.delete.mockResolvedValue(undefined);
    render(<ApiKeys />);
    await screen.findByText('My Integration Key');

    await userEvent.click(screen.getByRole('button', { name: /revoke/i }));
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: /confirm/i }));
    await waitFor(() => {
      expect(screen.getByText(/key revoked/i)).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    mockApi.get.mockImplementation(() => new Promise(() => {}));
    render(<ApiKeys />);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });
});
