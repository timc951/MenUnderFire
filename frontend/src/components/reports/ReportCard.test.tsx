import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect } from 'vitest';
import { ReportCard } from './ReportCard';
import { Report } from '../../types';

describe('ReportCard', () => {
  const mockReport: Report = {
    id: '123',
    title: 'Weekly Check-in',
    content: 'I exercised 4 times this week.',
    groupId: '1',
    reporterName: 'John Doe',
    isAnonymous: false,
    createdAt: '2024-01-15T10:00:00Z',
  };

  it('should display report title and content', () => {
    render(<ReportCard report={mockReport} />);
    expect(screen.getByText('Weekly Check-in')).toBeInTheDocument();
    expect(screen.getByText(/exercised 4 times/i)).toBeInTheDocument();
  });

  it('should display reporter name for non-anonymous report', () => {
    render(<ReportCard report={mockReport} />);
    expect(screen.getByText('John Doe')).toBeInTheDocument();
  });

  it('should display "Anonymous" for anonymous report', () => {
    const anonymousReport: Report = { ...mockReport, reporterName: null, isAnonymous: true };
    render(<ReportCard report={anonymousReport} />);
    expect(screen.getByTestId('anonymous-badge')).toBeInTheDocument();
    expect(screen.getByTestId('anonymous-badge')).toHaveTextContent('Anonymous');
  });

  it('should format date correctly', () => {
    render(<ReportCard report={mockReport} />);
    expect(screen.getByText(/Jan 15, 2024/i)).toBeInTheDocument();
  });

  it('should expand to show full content when clicked', async () => {
    const longReport: Report = { ...mockReport, content: 'A'.repeat(500) };
    render(<ReportCard report={longReport} />);
    expect(screen.getByText(/show more/i)).toBeInTheDocument();
    await userEvent.click(screen.getByText(/show more/i));
    expect(screen.getByText(/show less/i)).toBeInTheDocument();
  });

  it('should not show expand button for short content', () => {
    render(<ReportCard report={mockReport} />);
    expect(screen.queryByText(/show more/i)).not.toBeInTheDocument();
  });
});
