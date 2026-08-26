import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { GroupDetail } from './GroupDetail';

const mockApi = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

vi.mock('../hooks/useApi', () => ({
  useApi: () => mockApi,
}));

const renderWithRouter = () => {
  return render(
    <MemoryRouter initialEntries={['/groups/1']}>
      <Routes>
        <Route path="/groups/:id" element={<GroupDetail />} />
      </Routes>
    </MemoryRouter>
  );
};

describe('GroupDetail', () => {
  beforeEach(() => {
    mockApi.get.mockImplementation((path: string) => {
      if (path === '/groups/1') {
        return Promise.resolve({
          id: '1',
          name: 'Fitness Group',
          description: 'Stay fit',
          organizationId: 'org-1',
          inviteCode: 'ABC123',
          createdAt: '2024-01-01',
          role: 'LEADER',
          members: [
            { id: 'm1', displayName: 'Leader', email: 'leader@test.com', role: 'LEADER', joinedAt: '2024-01-01' },
          ],
        });
      }
      if (path === '/groups/1/messages') {
        return Promise.resolve([]);
      }
      if (path === '/groups/1/form-reports') {
        return Promise.resolve([
          { id: 'msg-1', formId: 'form-1', formName: 'Weekly Check-in', content: 'Please fill out', createdAt: '2024-01-15T10:00:00Z', senderName: 'Leader', senderId: 'm1', groupId: '1', notifyMembers: true },
        ]);
      }
      return Promise.resolve([]);
    });
  });

  it('should display group name', async () => {
    renderWithRouter();
    expect(await screen.findByText('Fitness Group')).toBeInTheDocument();
  });

  it('should display group description', async () => {
    renderWithRouter();
    expect(await screen.findByText('Stay fit')).toBeInTheDocument();
  });

  it('should show form reports when reports tab clicked', async () => {
    renderWithRouter();
    await screen.findByText('Fitness Group');
    await userEvent.click(screen.getByText('Reports'));
    await waitFor(() => {
      expect(screen.getByText('Weekly Check-in')).toBeInTheDocument();
    });
  });

  it('should switch to members tab', async () => {
    renderWithRouter();
    await screen.findByText('Fitness Group');
    await userEvent.click(screen.getByText('Members'));
    expect(screen.getByText('Leader')).toBeInTheDocument();
  });

  it('should show send form button for leaders', async () => {
    renderWithRouter();
    await screen.findByText('Fitness Group');
    expect(screen.getByRole('button', { name: /send form/i })).toBeInTheDocument();
  });

  it('should show empty state when no form reports', async () => {
    mockApi.get.mockImplementation((path: string) => {
      if (path === '/groups/1') {
        return Promise.resolve({
          id: '1',
          name: 'Fitness Group',
          description: 'Stay fit',
          organizationId: 'org-1',
          inviteCode: 'ABC123',
          createdAt: '2024-01-01',
          role: 'LEADER',
          members: [
            { id: 'm1', displayName: 'Leader', email: 'leader@test.com', role: 'LEADER', joinedAt: '2024-01-01' },
          ],
        });
      }
      if (path === '/groups/1/form-reports') {
        return Promise.resolve([]);
      }
      return Promise.resolve([]);
    });
    renderWithRouter();
    await screen.findByText('Fitness Group');
    await userEvent.click(screen.getByText('Reports'));
    await waitFor(() => {
      expect(screen.getByText(/no form reports yet/i)).toBeInTheDocument();
    });
  });

  it('should show loading spinner while fetching', () => {
    mockApi.get.mockImplementation(() => new Promise(() => {})); // never resolves
    renderWithRouter();
    expect(screen.getByRole('status')).toBeInTheDocument();
  });
});
