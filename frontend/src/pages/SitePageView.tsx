import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { useSitePage } from '../hooks/useSitePages';
import { Button } from '../components/common/Button';
import { PagesMenu } from '../components/common/PagesMenu';
import { renderMarkdown } from '../utils/markdown';

// Simple markdown renderer - converts basic markdown to HTML

export function SitePageView() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading, login, logout } = useAuth();
  const { page, isLoading, error } = useSitePage(slug || '');

  const handleLogout = () => {
    logout();
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

  if (error || !page) {
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
            <h1 className="text-4xl font-bold text-white mb-4">Page Not Found</h1>
            <p className="text-stone-300 mb-8">
              The page you're looking for doesn't exist or is not published.
            </p>
            <Button onClick={() => navigate('/')} className="bg-amber-600 hover:bg-amber-700 border-amber-600">
              Go Home
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

      {/* Fixed Header */}
      <header className="relative z-20 flex-shrink-0">
        <div className="mx-auto max-w-6xl px-4 sm:px-6">
          <div className="flex h-16 items-center justify-between border-b border-white/10 bg-black/25 backdrop-blur supports-[backdrop-filter]:bg-black/20">
            <div className="flex items-center gap-3">
              {/* Pages Menu (hamburger) */}
              <PagesMenu />

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

            <nav className="hidden md:flex items-center gap-7 text-sm text-white/75">
              <button className="hover:text-white transition" onClick={() => navigate('/')}>How it works</button>
              <button className="hover:text-white transition" onClick={() => navigate('/page/about')}>About</button>
              <button className="hover:text-white transition" onClick={() => navigate('/#faq')}>FAQ</button>
            </nav>

            <div className="flex items-center gap-3">
              {authLoading ? (
                <div className="w-24 h-10 bg-amber-500/30 rounded-lg animate-pulse" />
              ) : isAuthenticated ? (
                <>
                  <button
                    onClick={() => navigate('/dashboard')}
                    className="hidden sm:inline-flex rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm text-white/90 hover:bg-white/10 transition"
                  >
                    Dashboard
                  </button>
                  <button
                    onClick={handleLogout}
                    className="inline-flex rounded-lg bg-amber-500 px-4 py-2 text-sm font-semibold text-black hover:bg-amber-400 transition shadow-[0_0_0_1px_rgba(245,158,11,0.35),0_10px_30px_-12px_rgba(245,158,11,0.55)]"
                  >
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <button
                    onClick={() => login()}
                    className="hidden sm:inline-flex rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm text-white/90 hover:bg-white/10 transition"
                  >
                    Sign in
                  </button>
                  <button
                    onClick={() => login()}
                    className="inline-flex rounded-lg bg-amber-500 px-4 py-2 text-sm font-semibold text-black hover:bg-amber-400 transition shadow-[0_0_0_1px_rgba(245,158,11,0.35),0_10px_30px_-12px_rgba(245,158,11,0.55)]"
                  >
                    Get Started
                  </button>
                </>
              )}
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
            <div
              className="prose prose-invert max-w-none text-stone-300"
              dangerouslySetInnerHTML={{ __html: renderMarkdown(page.content) }}
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
            <div className="flex items-center gap-4 text-sm">
              <button
                onClick={() => navigate('/terms')}
                className="text-stone-400 hover:text-amber-400 transition-colors"
              >
                Terms of Use
              </button>
              <span className="text-stone-600">|</span>
              <button
                onClick={() => navigate('/privacy')}
                className="text-stone-400 hover:text-amber-400 transition-colors"
              >
                Privacy Policy
              </button>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
