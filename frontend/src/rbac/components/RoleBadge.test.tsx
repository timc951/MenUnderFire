import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { RoleBadge } from './RoleBadge';

describe('RoleBadge', () => {
  it('should render Site Admin badge', () => {
    render(<RoleBadge role="site_admin" />);
    expect(screen.getByText('Site Admin')).toBeInTheDocument();
  });

  it('should render Org Admin badge', () => {
    render(<RoleBadge role="org_admin" />);
    expect(screen.getByText('Org Admin')).toBeInTheDocument();
  });

  it('should render Group Owner badge', () => {
    render(<RoleBadge role="group_owner" />);
    expect(screen.getByText('Owner')).toBeInTheDocument();
  });

  it('should render Member badge', () => {
    render(<RoleBadge role="member" />);
    expect(screen.getByText('Member')).toBeInTheDocument();
  });

  it('should apply purple styling for Site Admin', () => {
    const { container } = render(<RoleBadge role="site_admin" />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-purple-100');
    expect(badge?.className).toContain('text-purple-800');
  });

  it('should apply indigo styling for Org Admin', () => {
    const { container } = render(<RoleBadge role="org_admin" />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-indigo-100');
    expect(badge?.className).toContain('text-indigo-800');
  });

  it('should apply blue styling for Group Owner', () => {
    const { container } = render(<RoleBadge role="group_owner" />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-blue-100');
    expect(badge?.className).toContain('text-blue-800');
  });

  it('should apply gray styling for Member', () => {
    const { container } = render(<RoleBadge role="member" />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-gray-100');
    expect(badge?.className).toContain('text-gray-800');
  });

  it('should accept additional className prop', () => {
    const { container } = render(<RoleBadge role="member" className="mt-2" />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('mt-2');
  });

  it('should render with proper badge styling', () => {
    const { container } = render(<RoleBadge role="site_admin" />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('inline-flex');
    expect(badge?.className).toContain('items-center');
    expect(badge?.className).toContain('rounded-full');
  });
});
