import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { JoinGroupModal } from './JoinGroupModal';

describe('JoinGroupModal', () => {
  const defaultProps = {
    isOpen: true,
    onClose: vi.fn(),
    onSubmit: vi.fn().mockResolvedValue(undefined),
  };

  it('should render invite code input', () => {
    render(<JoinGroupModal {...defaultProps} />);
    expect(screen.getByLabelText(/invite code/i)).toBeInTheDocument();
  });

  it('should validate required invite code', async () => {
    render(<JoinGroupModal {...defaultProps} />);
    await userEvent.click(screen.getByRole('button', { name: /join group/i }));
    expect(screen.getByText(/invite code is required/i)).toBeInTheDocument();
  });

  it('should submit invite code', async () => {
    render(<JoinGroupModal {...defaultProps} />);
    await userEvent.type(screen.getByLabelText(/invite code/i), 'ABC123');
    await userEvent.click(screen.getByRole('button', { name: /join group/i }));

    await waitFor(() => {
      expect(defaultProps.onSubmit).toHaveBeenCalledWith('ABC123');
    });
  });

  it('should show error on invalid code', async () => {
    const onSubmit = vi.fn().mockRejectedValue(new Error('Not found'));
    render(<JoinGroupModal {...defaultProps} onSubmit={onSubmit} />);
    await userEvent.type(screen.getByLabelText(/invite code/i), 'INVALID');
    await userEvent.click(screen.getByRole('button', { name: /join group/i }));

    await waitFor(() => {
      expect(screen.getByText(/invalid invite code/i)).toBeInTheDocument();
    });
  });

  it('should call onClose when cancel is clicked', async () => {
    render(<JoinGroupModal {...defaultProps} />);
    await userEvent.click(screen.getByRole('button', { name: /cancel/i }));
    expect(defaultProps.onClose).toHaveBeenCalled();
  });
});
