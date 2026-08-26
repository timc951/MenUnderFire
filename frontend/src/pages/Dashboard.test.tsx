import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import { Dashboard } from './Dashboard';

const mockUseRbac = vi.fn();
vi.mock('../rbac', () => ({
  useRbac: () => mockUseRbac(),
}));

const mockUseDashboardStats = vi.fn();
vi.mock('../hooks/useDashboardStats', () => ({
  useDashboardStats: () => mockUseDashboardStats(),
}));

describe('Dashboard', () => {
  beforeEach(() => {
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => true,
      isOrgAdmin: () => false,
    });
    mockUseDashboardStats.mockReturnValue({
      stats: { organizationCount: 5, groupCount: 12 },
      isLoading: false,
      error: null,
    });
  });

  it('should display page title', () => {
    render(<MemoryRouter><Dashboard /></MemoryRouter>);
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('should display organization count for site admin', () => {
    render(<MemoryRouter><Dashboard /></MemoryRouter>);
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('Organizations')).toBeInTheDocument();
  });

  it('should display group count for site admin', () => {
    render(<MemoryRouter><Dashboard /></MemoryRouter>);
    expect(screen.getByText('12')).toBeInTheDocument();
    expect(screen.getByText('Groups')).toBeInTheDocument();
  });

  it('should show loading state', () => {
    mockUseDashboardStats.mockReturnValue({
      stats: null,
      isLoading: true,
      error: null,
    });
    render(<MemoryRouter><Dashboard /></MemoryRouter>);
    expect(screen.getByText(/loading statistics/i)).toBeInTheDocument();
  });

  it('should show error state', () => {
    mockUseDashboardStats.mockReturnValue({
      stats: null,
      isLoading: false,
      error: 'Failed to load',
    });
    render(<MemoryRouter><Dashboard /></MemoryRouter>);
    expect(screen.getByText('Failed to load')).toBeInTheDocument();
  });

  it('should show welcome message for non-admin users', () => {
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => false,
      isOrgAdmin: () => false,
    });
    mockUseDashboardStats.mockReturnValue({
      stats: null,
      isLoading: false,
      error: null,
    });
    render(<MemoryRouter><Dashboard /></MemoryRouter>);
    expect(screen.getByText(/welcome to men under fire/i)).toBeInTheDocument();
  });

  it('should only show group count for org admin (no organization count)', () => {
    mockUseRbac.mockReturnValue({
      isSiteAdmin: () => false,
      isOrgAdmin: () => true,
    });
    mockUseDashboardStats.mockReturnValue({
      stats: { organizationCount: null, groupCount: 8 },
      isLoading: false,
      error: null,
    });
    render(<MemoryRouter><Dashboard /></MemoryRouter>);
    expect(screen.getByText('8')).toBeInTheDocument();
    expect(screen.getByText('Groups')).toBeInTheDocument();
    expect(screen.queryByText('Organizations')).not.toBeInTheDocument();
  });
});
