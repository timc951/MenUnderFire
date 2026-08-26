import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { formatDate, formatRelativeTime, truncateText, generateInviteLink } from './helpers';

describe('helpers', () => {
  describe('formatDate', () => {
    it('should format date correctly', () => {
      expect(formatDate('2024-01-15T10:00:00Z')).toBe('Jan 15, 2024');
    });

    it('should handle different months', () => {
      expect(formatDate('2024-06-20T10:00:00Z')).toBe('Jun 20, 2024');
    });
  });

  describe('formatRelativeTime', () => {
    beforeEach(() => {
      vi.useFakeTimers();
      vi.setSystemTime(new Date('2024-01-15T12:00:00Z'));
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it('should return "just now" for < 1 minute', () => {
      expect(formatRelativeTime('2024-01-15T11:59:30Z')).toBe('just now');
    });

    it('should return minutes for < 1 hour', () => {
      expect(formatRelativeTime('2024-01-15T11:30:00Z')).toBe('30m ago');
    });

    it('should return hours for < 24 hours', () => {
      expect(formatRelativeTime('2024-01-15T06:00:00Z')).toBe('6h ago');
    });

    it('should return days for < 7 days', () => {
      expect(formatRelativeTime('2024-01-12T12:00:00Z')).toBe('3d ago');
    });

    it('should return formatted date for >= 7 days', () => {
      expect(formatRelativeTime('2024-01-01T12:00:00Z')).toBe('Jan 1, 2024');
    });
  });

  describe('truncateText', () => {
    it('should not truncate short text', () => {
      expect(truncateText('Hello', 10)).toBe('Hello');
    });

    it('should truncate long text with ellipsis', () => {
      expect(truncateText('This is a long text', 10)).toBe('This is a...');
    });

    it('should handle exact length', () => {
      expect(truncateText('Hello', 5)).toBe('Hello');
    });
  });

  describe('generateInviteLink', () => {
    it('should generate correct invite link', () => {
      expect(generateInviteLink('ABC123')).toContain('/groups/join?code=ABC123');
    });
  });
});
