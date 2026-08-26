import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { GroupMembers } from './GroupMembers';
import { GroupMembership } from '../../types';

describe('GroupMembers', () => {
  const mockMembers: GroupMembership[] = [
    { id: 'm1', userId: 'u1', groupId: '1', role: 'LEADER', joinedAt: '2024-01-01T00:00:00Z', displayName: 'Leader User' },
    { id: 'm2', userId: 'u2', groupId: '1', role: 'MEMBER', joinedAt: '2024-01-02T00:00:00Z', displayName: 'Member One' },
    { id: 'm3', userId: 'u3', groupId: '1', role: 'MEMBER', joinedAt: '2024-01-03T00:00:00Z', displayName: 'Member Two' },
  ];

  it('should show invite code for leaders', () => {
    render(<GroupMembers members={mockMembers} inviteCode="ABC123" currentUserRole="LEADER" />);
    expect(screen.getByText(/invite code/i)).toBeInTheDocument();
    expect(screen.getByText('ABC123')).toBeInTheDocument();
  });

  it('should hide invite code for regular members', () => {
    render(<GroupMembers members={mockMembers} inviteCode="ABC123" currentUserRole="MEMBER" />);
    expect(screen.queryByText(/invite code/i)).not.toBeInTheDocument();
  });

  it('should show remove button for leaders', () => {
    render(<GroupMembers members={mockMembers} inviteCode="ABC123" currentUserRole="LEADER" />);
    const removeButtons = screen.getAllByRole('button', { name: /remove/i });
    expect(removeButtons.length).toBe(2); // Only non-leader members
  });

  it('should not show remove button for members', () => {
    render(<GroupMembers members={mockMembers} inviteCode="ABC123" currentUserRole="MEMBER" />);
    expect(screen.queryByRole('button', { name: /remove/i })).not.toBeInTheDocument();
  });

  it('should call onRemoveMember when remove is clicked', async () => {
    const onRemove = vi.fn();
    render(<GroupMembers members={mockMembers} inviteCode="ABC123" currentUserRole="LEADER" onRemoveMember={onRemove} />);
    await userEvent.click(screen.getByLabelText(/remove member one/i));
    expect(onRemove).toHaveBeenCalledWith('u2');
  });

  it('should display all member names', () => {
    render(<GroupMembers members={mockMembers} inviteCode="ABC123" currentUserRole="MEMBER" />);
    expect(screen.getByText('Leader User')).toBeInTheDocument();
    expect(screen.getByText('Member One')).toBeInTheDocument();
    expect(screen.getByText('Member Two')).toBeInTheDocument();
  });

  it('should show member count', () => {
    render(<GroupMembers members={mockMembers} inviteCode="ABC123" currentUserRole="MEMBER" />);
    expect(screen.getByText(/members \(3\)/i)).toBeInTheDocument();
  });
});
