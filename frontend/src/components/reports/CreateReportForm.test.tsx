import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { CreateReportForm } from './CreateReportForm';

describe('CreateReportForm', () => {
  const mockGroups = [
    { id: '1', name: 'Fitness Group' },
    { id: '2', name: 'Study Group' },
  ];

  it('should render all form fields', () => {
    render(<CreateReportForm groups={mockGroups} />);
    expect(screen.getByLabelText(/group/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/content/i)).toBeInTheDocument();
    expect(screen.getByRole('switch', { name: /anonymous/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument();
  });

  it('should populate group dropdown with user groups', () => {
    render(<CreateReportForm groups={mockGroups} />);
    expect(screen.getByText('Fitness Group')).toBeInTheDocument();
    expect(screen.getByText('Study Group')).toBeInTheDocument();
  });

  it('should validate required fields', async () => {
    render(<CreateReportForm groups={mockGroups} />);
    await userEvent.click(screen.getByRole('button', { name: /submit/i }));
    expect(screen.getByText(/title is required/i)).toBeInTheDocument();
    expect(screen.getByText(/content is required/i)).toBeInTheDocument();
  });

  it('should submit form with correct data', async () => {
    const onSubmit = vi.fn().mockResolvedValue({ id: '123' });
    render(<CreateReportForm groups={mockGroups} onSubmit={onSubmit} />);

    await userEvent.selectOptions(screen.getByLabelText(/group/i), '1');
    await userEvent.type(screen.getByLabelText(/title/i), 'Weekly Report');
    await userEvent.type(screen.getByLabelText(/content/i), 'Did great this week!');
    await userEvent.click(screen.getByRole('switch', { name: /anonymous/i }));
    await userEvent.click(screen.getByRole('button', { name: /submit/i }));

    expect(onSubmit).toHaveBeenCalledWith({
      groupId: '1',
      title: 'Weekly Report',
      content: 'Did great this week!',
      isAnonymousToGroup: true,
    });
  });

  it('should show loading state while submitting', async () => {
    const onSubmit = vi.fn(() => new Promise((resolve) => setTimeout(resolve, 100)));
    render(<CreateReportForm groups={mockGroups} onSubmit={onSubmit} />);

    await userEvent.selectOptions(screen.getByLabelText(/group/i), '1');
    await userEvent.type(screen.getByLabelText(/title/i), 'Test');
    await userEvent.type(screen.getByLabelText(/content/i), 'Content');
    await userEvent.click(screen.getByRole('button', { name: /submit/i }));

    expect(screen.getByRole('button', { name: /submitting/i })).toBeDisabled();
  });

  it('should show success message after submission', async () => {
    const onSubmit = vi.fn().mockResolvedValue({ id: '123' });
    render(<CreateReportForm groups={mockGroups} onSubmit={onSubmit} />);

    await userEvent.selectOptions(screen.getByLabelText(/group/i), '1');
    await userEvent.type(screen.getByLabelText(/title/i), 'Test');
    await userEvent.type(screen.getByLabelText(/content/i), 'Content');
    await userEvent.click(screen.getByRole('button', { name: /submit/i }));

    await waitFor(() => {
      expect(screen.getByText(/report submitted/i)).toBeInTheDocument();
    });
  });

  it('should show error message on failure', async () => {
    const onSubmit = vi.fn().mockRejectedValue(new Error('Network error'));
    render(<CreateReportForm groups={mockGroups} onSubmit={onSubmit} />);

    await userEvent.selectOptions(screen.getByLabelText(/group/i), '1');
    await userEvent.type(screen.getByLabelText(/title/i), 'Test');
    await userEvent.type(screen.getByLabelText(/content/i), 'Content');
    await userEvent.click(screen.getByRole('button', { name: /submit/i }));

    await waitFor(() => {
      expect(screen.getByText(/failed to submit/i)).toBeInTheDocument();
    });
  });
});
