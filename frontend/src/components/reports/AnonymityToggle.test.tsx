import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { AnonymityToggle } from './AnonymityToggle';

describe('AnonymityToggle', () => {
  it('should render with anonymous OFF by default', () => {
    render(<AnonymityToggle />);
    expect(screen.getByRole('switch')).not.toBeChecked();
    expect(screen.getByText(/your name will be visible/i)).toBeInTheDocument();
  });

  it('should toggle to anonymous ON when clicked', async () => {
    const onChange = vi.fn();
    render(<AnonymityToggle onChange={onChange} />);
    await userEvent.click(screen.getByRole('switch'));
    expect(onChange).toHaveBeenCalledWith(true);
  });

  it('should show explanation when anonymous is ON', () => {
    render(<AnonymityToggle value={true} />);
    expect(screen.getByText(/group members will not see your name/i)).toBeInTheDocument();
    expect(screen.getByText(/leader will still see your identity/i)).toBeInTheDocument();
  });

  it('should show info tooltip on hover', async () => {
    render(<AnonymityToggle />);
    await userEvent.hover(screen.getByTestId('info-icon'));
    expect(screen.getByRole('tooltip')).toBeInTheDocument();
  });

  it('should hide tooltip when not hovered', () => {
    render(<AnonymityToggle />);
    expect(screen.queryByRole('tooltip')).not.toBeInTheDocument();
  });
});
