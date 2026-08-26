import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { GroupList } from './GroupList';
import { Group } from '../../types';

describe('GroupList', () => {
  const mockGroups: Group[] = [
    { id: '1', name: 'Fitness Group', description: 'Stay fit', inviteCode: 'ABC', createdBy: 'u1', createdAt: '2024-01-01T00:00:00Z', memberCount: 5, role: 'leader' },
    { id: '2', name: 'Study Group', description: 'Study hard', inviteCode: 'DEF', createdBy: 'u2', createdAt: '2024-01-02T00:00:00Z', memberCount: 3, role: 'member' },
  ];

  it('should render list of groups', () => {
    render(<GroupList groups={mockGroups} />);
    expect(screen.getByText('Fitness Group')).toBeInTheDocument();
    expect(screen.getByText('Study Group')).toBeInTheDocument();
  });

  it('should show loading spinner when loading', () => {
    render(<GroupList groups={[]} isLoading={true} />);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('should show error message', () => {
    render(<GroupList groups={[]} error="Failed to load" />);
    expect(screen.getByText('Failed to load')).toBeInTheDocument();
  });

  it('should show empty message when no groups', () => {
    render(<GroupList groups={[]} />);
    expect(screen.getByText(/no groups yet/i)).toBeInTheDocument();
  });

  it('should call onGroupClick when group is clicked', async () => {
    const onGroupClick = vi.fn();
    render(<GroupList groups={mockGroups} onGroupClick={onGroupClick} />);
    await userEvent.click(screen.getAllByRole('button')[0]);
    expect(onGroupClick).toHaveBeenCalledWith(mockGroups[0]);
  });
});
