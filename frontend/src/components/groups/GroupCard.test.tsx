import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { GroupCard } from './GroupCard';
import { Group } from '../../types';

describe('GroupCard', () => {
  const mockGroup: Group = {
    id: '1',
    name: 'Fitness Group',
    description: 'Stay fit together',
    inviteCode: 'ABC123',
    createdBy: 'user-1',
    createdAt: '2024-01-01T00:00:00Z',
    memberCount: 5,
    role: 'leader',
  };

  it('should display group name and description', () => {
    render(<GroupCard group={mockGroup} />);
    expect(screen.getByText('Fitness Group')).toBeInTheDocument();
    expect(screen.getByText('Stay fit together')).toBeInTheDocument();
  });

  it('should display role badge', () => {
    render(<GroupCard group={mockGroup} />);
    expect(screen.getByText('leader')).toBeInTheDocument();
  });

  it('should display member count', () => {
    render(<GroupCard group={mockGroup} />);
    expect(screen.getByText('5 members')).toBeInTheDocument();
  });

  it('should be clickable when onClick is provided', async () => {
    const onClick = vi.fn();
    render(<GroupCard group={mockGroup} onClick={onClick} />);
    await userEvent.click(screen.getByRole('button'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('should show member role with primary color', () => {
    const memberGroup: Group = { ...mockGroup, role: 'member' };
    render(<GroupCard group={memberGroup} />);
    expect(screen.getByText('member')).toHaveClass('bg-stone-100');
  });
});
