import { describe, it, expect } from 'vitest';
import { renderMarkdown } from './markdown';

describe('renderMarkdown', () => {
  describe('escaping', () => {
    it('escapes raw HTML', () => {
      const html = renderMarkdown('<script>alert(1)</script>');
      expect(html).not.toContain('<script>');
      expect(html).toContain('&lt;script&gt;');
    });
  });

  describe('link hrefs', () => {
    it('renders an ordinary https link', () => {
      const html = renderMarkdown('[docs](https://example.com/a?x=1)');
      expect(html).toContain('href="https://example.com/a?x=1"');
    });

    it('allows site-relative and fragment links', () => {
      expect(renderMarkdown('[a](/about)')).toContain('href="/about"');
      expect(renderMarkdown('[a](#top)')).toContain('href="#top"');
    });

    it('rejects a javascript: url', () => {
      const html = renderMarkdown('[click](javascript:alert%281%29)');
      expect(html).not.toContain('javascript:');
      expect(html).toContain('href="#"');
    });

    it('rejects a data: url', () => {
      const html = renderMarkdown('[click](data:text/html,<h1>hi</h1>)');
      expect(html).not.toContain('data:text/html');
    });

    it('does not allow a quote in the url to inject an attribute', () => {
      const html = renderMarkdown('[click](" onmouseover="alert(1)');
      expect(html).not.toMatch(/\sonmouseover=/);
    });

    it('keeps the link text escaped', () => {
      const html = renderMarkdown('[<img src=x onerror=alert(1)>](https://example.com)');
      expect(html).not.toContain('<img');
    });
  });
});
