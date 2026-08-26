import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '../components/common/Button';
import { renderMarkdown } from '../utils/markdown';

interface PageDraft {
  id: string;
  pageId?: string;
  title: string;
  content: string;
  createdAt: string;
  expiresAt: string;
}

// Simple markdown renderer - converts basic markdown to HTML

export function PagePreview() {
  const { draftId } = useParams<{ draftId: string }>();
  const navigate = useNavigate();
  const [draft, setDraft] = useState<PageDraft | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDraft = async () => {
      if (!draftId) return;
      setIsLoading(true);
      setError(null);
      try {
        const response = await fetch(`/api/pages/preview/${draftId}`);
        if (!response.ok) {
          if (response.status === 404) {
            throw new Error('Preview not found or has expired');
          }
          throw new Error('Failed to load preview');
        }
        const data = await response.json();
        setDraft(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load preview');
      } finally {
        setIsLoading(false);
      }
    };
    fetchDraft();
  }, [draftId]);

  const formatExpiresAt = (expiresAt: string) => {
    const expires = new Date(expiresAt);
    const now = new Date();
    const diffMs = expires.getTime() - now.getTime();
    const diffMins = Math.round(diffMs / 60000);
    if (diffMins <= 0) return 'Expired';
    if (diffMins < 60) return `${diffMins} minute${diffMins === 1 ? '' : 's'}`;
    const diffHours = Math.round(diffMins / 60);
    return `${diffHours} hour${diffHours === 1 ? '' : 's'}`;
  };

  if (isLoading) {
    return (
      <div className="h-screen text-white antialiased relative flex flex-col overflow-hidden">
        <div
          className="fixed inset-0 bg-top bg-cover"
          style={{ backgroundImage: 'url(https://content.menunderfire.com/poster.webp)' }}
          aria-hidden="true"
        />
        <div className="fixed inset-0 bg-gradient-to-b from-black/85 via-black/55 to-black/80" aria-hidden="true" />
        <div className="relative z-10 flex-1 flex items-center justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-500" />
        </div>
      </div>
    );
  }

  if (error || !draft) {
    return (
      <div className="h-screen text-white antialiased relative flex flex-col overflow-hidden">
        <div
          className="fixed inset-0 bg-top bg-cover"
          style={{ backgroundImage: 'url(https://content.menunderfire.com/poster.webp)' }}
          aria-hidden="true"
        />
        <div className="fixed inset-0 bg-gradient-to-b from-black/85 via-black/55 to-black/80" aria-hidden="true" />
        <div className="relative z-10 flex-1 flex items-center justify-center">
          <div className="text-center">
            <h1 className="text-4xl font-bold text-white mb-4">Preview Not Available</h1>
            <p className="text-stone-300 mb-8">
              {error || 'This preview has expired or does not exist.'}
            </p>
            <Button onClick={() => window.close()} className="bg-amber-600 hover:bg-amber-700 border-amber-600">
              Close Preview
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen text-white antialiased relative flex flex-col overflow-hidden">
      {/* Fixed background that covers the entire page */}
      <div
        className="fixed inset-0 bg-top bg-cover"
        style={{ backgroundImage: 'url(https://content.menunderfire.com/poster.webp)' }}
        aria-hidden="true"
      />

      {/* Gradient overlay - also fixed */}
      <div className="fixed inset-0 bg-gradient-to-b from-black/85 via-black/55 to-black/80" aria-hidden="true" />

      {/* Preview Banner */}
      <div className="relative z-30 bg-amber-500 text-black py-2 px-4">
        <div className="max-w-6xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-3">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            <span className="font-semibold">Preview Mode</span>
            <span className="text-amber-800">
              This is a preview of unsaved changes. Expires in {formatExpiresAt(draft.expiresAt)}.
            </span>
          </div>
          <button
            onClick={() => window.close()}
            className="px-3 py-1 bg-black/20 hover:bg-black/30 rounded text-sm font-medium transition"
          >
            Close Preview
          </button>
        </div>
      </div>

      {/* Fixed Header */}
      <header className="relative z-20 flex-shrink-0">
        <div className="mx-auto max-w-6xl px-4 sm:px-6">
          <div className="flex h-16 items-center justify-between border-b border-white/10 bg-black/25 backdrop-blur supports-[backdrop-filter]:bg-black/20">
            <div className="flex items-center gap-3">
              {/* Icon-only brand */}
              <button onClick={() => navigate('/')} className="group inline-flex items-center gap-2" aria-label="Home">
                <img
                  src="https://content.menunderfire.com/logo.webp"
                  alt="Men Under Fire"
                  width={36}
                  height={36}
                  className="h-9 w-9 rounded-xl group-hover:opacity-80 transition"
                />
                <span className="hidden sm:inline text-xs tracking-wide text-white/60 group-hover:text-white/80 transition">
                  Home
                </span>
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Scrollable Content Area */}
      <main className="relative z-10 flex-1 overflow-y-auto scrollbar-thin scrollbar-thumb-stone-600 scrollbar-track-transparent [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-stone-600 [&::-webkit-scrollbar-thumb]:rounded-full">
        {/* Title image */}
        <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 pt-6 mb-6">
          <img
            src="https://content.menunderfire.com/MenUnderFireStone2.webp"
            alt="Men Under Fire"
            width={1275}
            height={149}
            className="mx-auto max-w-full h-auto"
            style={{ maxHeight: '120px' }}
          />
        </div>
        <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 pb-8">
          <article className="bg-stone-900/80 rounded-xl p-4 sm:p-6 lg:p-8 border border-stone-700/50 backdrop-blur-sm">
            {/* Show preview title */}
            <h1 className="text-3xl font-bold text-white mb-6">{draft.title}</h1>
            <div
              className="prose prose-invert max-w-none text-stone-300"
              dangerouslySetInnerHTML={{ __html: renderMarkdown(draft.content) }}
            />
          </article>
        </div>
      </main>

      {/* Fixed Footer */}
      <footer className="relative z-20 flex-shrink-0 border-t border-white/10 bg-black/60 backdrop-blur">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex flex-col sm:flex-row items-center justify-center gap-2 sm:gap-4">
            <p className="text-stone-400 text-sm">
              &copy; {new Date().getFullYear()} Men Under Fire. All rights reserved.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
