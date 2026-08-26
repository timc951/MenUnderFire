// Minimal Markdown renderer for site page content.
//
// The output is injected via dangerouslySetInnerHTML, so every path that emits
// HTML must assume the input is hostile. Source text is HTML-escaped up front;
// the only place a caller-supplied value reaches an HTML *attribute* is the
// link href, which is scheme-checked and quote-escaped by sanitizeHref.

const ALLOWED_PROTOCOLS = ['http:', 'https:', 'mailto:'];

// & < > are already escaped by the time a URL reaches here, so quotes are all
// that stand between a href value and attribute injection.
function escapeQuotes(value: string): string {
  return value.replace(/"/g, '&quot;').replace(/'/g, '&#39;');
}

function sanitizeHref(raw: string): string {
  const trimmed = raw.trim();

  // Same-document and site-relative links carry no scheme to validate.
  if (trimmed.startsWith('/') || trimmed.startsWith('#')) {
    return escapeQuotes(trimmed);
  }

  let parsed: URL;
  try {
    parsed = new URL(trimmed);
  } catch {
    return '#';
  }

  if (!ALLOWED_PROTOCOLS.includes(parsed.protocol)) {
    return '#';
  }

  return escapeQuotes(trimmed);
}

export function renderMarkdown(content: string): string {
  let html = content
    // Escape HTML first
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    // Headers
    .replace(/^### (.*$)/gm, '<h3 class="text-xl font-semibold mt-6 mb-3 text-white">$1</h3>')
    .replace(/^## (.*$)/gm, '<h2 class="text-2xl font-semibold mt-8 mb-4 text-white">$1</h2>')
    .replace(/^# (.*$)/gm, '<h1 class="text-3xl font-bold mt-8 mb-4 text-white">$1</h1>')
    // Bold and italic
    .replace(/\*\*\*(.*?)\*\*\*/g, '<strong><em>$1</em></strong>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    // Links
    .replace(
      /\[([^\]]+)\]\(([^)]+)\)/g,
      (_match, text: string, href: string) =>
        `<a href="${sanitizeHref(href)}" class="text-amber-400 hover:underline" target="_blank" rel="noopener noreferrer">${text}</a>`
    )
    // Unordered lists
    .replace(/^\s*[-*]\s+(.*)$/gm, '<li class="ml-4">$1</li>')
    // Ordered lists
    .replace(/^\s*\d+\.\s+(.*)$/gm, '<li class="ml-4">$1</li>')
    // Code blocks
    .replace(/```([\s\S]*?)```/g, '<pre class="bg-stone-800 p-4 rounded-md overflow-x-auto my-4"><code>$1</code></pre>')
    // Inline code
    .replace(/`([^`]+)`/g, '<code class="bg-stone-800 px-1 rounded">$1</code>')
    // Horizontal rule
    .replace(/^---$/gm, '<hr class="my-8 border-stone-600" />')
    // Blockquotes
    .replace(/^>\s+(.*)$/gm, '<blockquote class="border-l-4 border-amber-500 pl-4 italic my-4">$1</blockquote>')
    // Paragraphs (double newlines)
    .replace(/\n\n/g, '</p><p class="mb-4">')
    // Single newlines in lists
    .replace(/<\/li>\n<li/g, '</li><li');

  // Wrap lists
  html = html.replace(/(<li[^>]*>.*<\/li>\n?)+/g, '<ul class="list-disc list-inside mb-4">$&</ul>');

  // Wrap in paragraph if not already wrapped
  if (!html.startsWith('<')) {
    html = '<p class="mb-4">' + html + '</p>';
  }

  return html;
}
