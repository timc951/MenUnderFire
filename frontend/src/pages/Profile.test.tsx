import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Profile } from './Profile';

vi.mock('../hooks/useAuth', () => ({
  useAuth: () => ({
    user: {
      id: 'user-1',
      name: 'Test User',
      email: 'test@example.com',
      picture: 'https://example.com/photo.jpg',
    },
  }),
}));

describe('Profile', () => {
  it('should display user name', () => {
    render(<Profile />);
    expect(screen.getByText('Test User')).toBeInTheDocument();
  });

  it('should display user email', () => {
    render(<Profile />);
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('should display user avatar', () => {
    render(<Profile />);
    expect(screen.getByAltText('Test User')).toHaveAttribute('src', 'https://example.com/photo.jpg');
  });

  it('should display page title', () => {
    render(<Profile />);
    expect(screen.getByText('Profile')).toBeInTheDocument();
  });
});
