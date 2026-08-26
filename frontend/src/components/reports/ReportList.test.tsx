import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { ReportList } from './ReportList';
import { Report } from '../../types';

describe('ReportList', () => {
  const mockReports: Report[] = [
    { id: '1', title: 'Report 1', content: 'Content 1', groupId: '1', reporterName: 'User 1', isAnonymous: false, createdAt: '2024-01-15T10:00:00Z' },
    { id: '2', title: 'Report 2', content: 'Content 2', groupId: '1', reporterName: null, isAnonymous: true, createdAt: '2024-01-14T10:00:00Z' },
  ];

  it('should render list of reports', () => {
    render(<ReportList reports={mockReports} />);
    expect(screen.getByText('Report 1')).toBeInTheDocument();
    expect(screen.getByText('Report 2')).toBeInTheDocument();
  });

  it('should show loading spinner when loading', () => {
    render(<ReportList reports={[]} isLoading={true} />);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('should show error message', () => {
    render(<ReportList reports={[]} error="Failed to load reports" />);
    expect(screen.getByText('Failed to load reports')).toBeInTheDocument();
  });

  it('should show empty message when no reports', () => {
    render(<ReportList reports={[]} />);
    expect(screen.getByText(/no reports yet/i)).toBeInTheDocument();
  });
});
