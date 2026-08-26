import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { LoginButton } from './LoginButton';

const mockLogin = vi.fn();

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => ({ login: mockLogin }),
}));

describe('LoginButton', () => {
  it('should render login button', () => {
    render(<LoginButton />);
    expect(screen.getByRole('button', { name: /log in/i })).toBeInTheDocument();
  });

  it('should call login when clicked', async () => {
    render(<LoginButton />);
    await userEvent.click(screen.getByRole('button', { name: /log in/i }));
    expect(mockLogin).toHaveBeenCalledTimes(1);
  });
});
