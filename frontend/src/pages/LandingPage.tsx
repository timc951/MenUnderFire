import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { PagesMenu } from '../components/common/PagesMenu';

export function LandingPage() {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading, login, logout } = useAuth();

  const handleLogout = () => {
    logout();
  };

  const handleGetStarted = () => {
    if (isAuthenticated) {
      navigate('/dashboard');
    } else {
      login();
    }
  };

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
              <a className="hover:text-white transition" href="#how">How it works</a>
              <button className="hover:text-white transition" onClick={() => navigate('/page/about')}>About</button>
              <a className="hover:text-white transition" href="#faq">FAQ</a>
            </nav>

            <div className="flex items-center gap-3">
              {isLoading ? (
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
                    onClick={handleGetStarted}
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
        {/* Hero */}
        <section className="mx-auto max-w-6xl px-4 sm:px-6 pt-14 pb-10">
          <div className="grid gap-10 lg:grid-cols-12 lg:items-center">
            <div className="lg:col-span-7">
              {/* Wordmark */}
              <img
                src="https://content.menunderfire.com/MenUnderFireStone2.webp"
                alt="Men Under Fire"
                width={1275}
                height={149}
                fetchPriority="high"
                className="h-16 sm:h-20 w-auto drop-shadow-[0_10px_30px_rgba(0,0,0,0.55)]"
              />

              {/* Message as H1 */}
              <h1 className="mt-6 text-4xl sm:text-5xl font-bold tracking-tight">
                Forge discipline.
                <span className="text-amber-300"> Build brotherhood.</span>
              </h1>

              <p className="mt-5 text-lg text-white/75 max-w-2xl leading-relaxed">
                Join a focused community built around accountability, growth, and spiritual strength.
                Weekly check-ins, encouragement, and a path that's actually sustainable.
              </p>

              <div className="mt-7 flex flex-wrap items-center gap-3">
                <button
                  onClick={handleGetStarted}
                  className="inline-flex items-center justify-center rounded-xl bg-amber-500 px-6 py-3 font-semibold text-black hover:bg-amber-400 transition shadow-[0_0_0_1px_rgba(245,158,11,0.35),0_12px_35px_-15px_rgba(245,158,11,0.55)]"
                >
                  {isAuthenticated ? 'Go to Dashboard' : 'Join a Group'}
                </button>
                <a
                  href="#how"
                  className="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/5 px-6 py-3 font-semibold text-white/90 hover:bg-white/10 transition"
                >
                  See how it works
                </a>
              </div>

              <div className="mt-8 flex flex-wrap items-center gap-6 text-sm text-white/65">
                <div className="flex items-center gap-2">
                  <span className="inline-block h-2 w-2 rounded-full bg-amber-400"></span>
                  Weekly accountability
                </div>
                <div className="flex items-center gap-2">
                  <span className="inline-block h-2 w-2 rounded-full bg-amber-400"></span>
                  Private groups
                </div>
                <div className="flex items-center gap-2">
                  <span className="inline-block h-2 w-2 rounded-full bg-amber-400"></span>
                  Zero spam / focused feed
                </div>
              </div>
            </div>

            {/* Right "preview" card */}
            <div className="lg:col-span-5">
              <div className="rounded-2xl border border-white/10 bg-black/35 backdrop-blur p-6 shadow-[0_30px_80px_-40px_rgba(0,0,0,0.8)]">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-semibold text-white/80">This week's focus</p>
                  <span className="rounded-full border border-amber-400/30 bg-amber-400/10 px-3 py-1 text-xs text-amber-200">
                    New
                  </span>
                </div>

                <div className="mt-5 space-y-4">
                  <div className="rounded-xl border border-white/10 bg-white/5 p-4">
                    <p className="text-sm text-white/70">Goal</p>
                    <p className="mt-1 font-semibold">Daily prayer + 20 min strength work</p>
                  </div>
                  <div className="rounded-xl border border-white/10 bg-white/5 p-4">
                    <p className="text-sm text-white/70">Check-in prompt</p>
                    <p className="mt-1 font-semibold">What did you win this week? What needs prayer?</p>
                  </div>
                  <div className="rounded-xl border border-white/10 bg-white/5 p-4">
                    <p className="text-sm text-white/70">Encouragement</p>
                    <p className="mt-1 font-semibold">You're not behind. You're being forged.</p>
                  </div>
                </div>

                <div className="mt-6 flex gap-3">
                  <button
                    onClick={() => isAuthenticated ? navigate('/dashboard') : navigate('/page/about')}
                    className="flex-1 rounded-xl bg-white/5 border border-white/10 px-4 py-3 text-sm font-semibold text-white/85 hover:bg-white/10 transition"
                  >
                    {isAuthenticated ? 'View groups' : 'Learn more'}
                  </button>
                  <button
                    onClick={handleGetStarted}
                    className="flex-1 rounded-xl bg-amber-500 px-4 py-3 text-sm font-semibold text-black hover:bg-amber-400 transition"
                  >
                    Start today
                  </button>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Feature cards */}
        <section className="mx-auto max-w-6xl px-4 sm:px-6 pb-14">
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div className="group rounded-2xl border border-white/10 bg-black/35 backdrop-blur p-7 hover:bg-black/45 transition hover:-translate-y-0.5">
              <div className="flex items-center gap-3">
                <div className="rounded-xl bg-amber-500/15 border border-amber-400/20 p-2">
                  <svg className="h-5 w-5 text-amber-200" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path strokeLinecap="round" d="M12 2l3 7 7 3-7 3-3 7-3-7-7-3 7-3 3-7z" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold">Accountability that sticks</h3>
              </div>
              <p className="mt-4 text-white/70 leading-relaxed">
                Weekly check-ins and clear goals—simple enough to keep, strong enough to change you.
              </p>
            </div>

            <div className="group rounded-2xl border border-white/10 bg-black/35 backdrop-blur p-7 hover:bg-black/45 transition hover:-translate-y-0.5">
              <div className="flex items-center gap-3">
                <div className="rounded-xl bg-amber-500/15 border border-amber-400/20 p-2">
                  <svg className="h-5 w-5 text-amber-200" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path strokeLinecap="round" d="M4 12a8 8 0 0 1 16 0c0 6-8 10-8 10S4 18 4 12z" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold">Brotherhood & encouragement</h3>
              </div>
              <p className="mt-4 text-white/70 leading-relaxed">
                A focused community that builds up—no noise, no drama, just progress and prayer.
              </p>
            </div>

            <div className="group rounded-2xl border border-white/10 bg-black/35 backdrop-blur p-7 hover:bg-black/45 transition hover:-translate-y-0.5">
              <div className="flex items-center gap-3">
                <div className="rounded-xl bg-amber-500/15 border border-amber-400/20 p-2">
                  <svg className="h-5 w-5 text-amber-200" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path strokeLinecap="round" d="M12 3v18M3 12h18" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold">Simple, repeatable habits</h3>
              </div>
              <p className="mt-4 text-white/70 leading-relaxed">
                Faith, fitness, and discipline—tracked with just enough structure to keep you moving.
              </p>
            </div>
          </div>
        </section>

        {/* How it works */}
        <section id="how" className="mx-auto max-w-6xl px-4 sm:px-6 py-16">
          <div className="flex items-end justify-between gap-6 flex-wrap">
            <div>
              <h2 className="text-3xl sm:text-4xl font-bold tracking-tight">How it works</h2>
              <p className="mt-3 text-white/70 max-w-2xl">Three steps. Real momentum. No complicated dashboards.</p>
            </div>
            <button onClick={() => navigate('/page/about')} className="text-sm font-semibold text-amber-200 hover:text-amber-100 transition">
              Learn more →
            </button>
          </div>

          <div className="mt-10 grid gap-4 lg:grid-cols-3">
            <div className="rounded-2xl border border-white/10 bg-black/50 backdrop-blur p-7">
              <p className="text-xs text-white/60">STEP 1</p>
              <h3 className="mt-2 text-xl font-semibold">Join a group</h3>
              <p className="mt-3 text-white/70 leading-relaxed">Pick a cadence and a focus (faith, fitness, discipline).</p>
            </div>
            <div className="rounded-2xl border border-white/10 bg-black/50 backdrop-blur p-7">
              <p className="text-xs text-white/60">STEP 2</p>
              <h3 className="mt-2 text-xl font-semibold">Set your weekly goal</h3>
              <p className="mt-3 text-white/70 leading-relaxed">Make it measurable. Make it honest. Make it doable.</p>
            </div>
            <div className="rounded-2xl border border-white/10 bg-black/50 backdrop-blur p-7">
              <p className="text-xs text-white/60">STEP 3</p>
              <h3 className="mt-2 text-xl font-semibold">Check in + encourage</h3>
              <p className="mt-3 text-white/70 leading-relaxed">Share wins, struggles, and prayer requests. Stay in the fight.</p>
            </div>
          </div>
        </section>

        {/* CTA */}
        <section id="groups" className="mx-auto max-w-6xl px-4 sm:px-6 pb-16">
          <div className="rounded-3xl border border-amber-400/20 bg-gradient-to-r from-amber-500/15 via-black/40 to-black/50 backdrop-blur p-10">
            <h2 className="text-3xl sm:text-4xl font-bold tracking-tight">Ready to be forged?</h2>
            <p className="mt-3 text-white/70 max-w-2xl">
              Join a group and start with one small, repeatable win this week.
            </p>
            <div className="mt-7">
              <button
                onClick={handleGetStarted}
                className="rounded-xl bg-amber-500 px-6 py-3 font-semibold text-black hover:bg-amber-400 transition"
              >
                Get Started
              </button>
            </div>
          </div>
        </section>

        {/* FAQ */}
        <section id="faq" className="mx-auto max-w-6xl px-4 sm:px-6 pb-16">
          <h2 className="text-3xl font-bold tracking-tight">FAQ</h2>
          <div className="mt-6 grid gap-4 lg:grid-cols-2">
            <div className="rounded-2xl border border-white/10 bg-black/50 backdrop-blur p-6">
              <h3 className="font-semibold">Is this only for men?</h3>
              <p className="mt-2 text-white/70 leading-relaxed">
                You can expand to women/co-ed with subdomains or separate communities while keeping the same system.
              </p>
            </div>
            <div className="rounded-2xl border border-white/10 bg-black/50 backdrop-blur p-6">
              <h3 className="font-semibold">How private are groups?</h3>
              <p className="mt-2 text-white/70 leading-relaxed">
                Keep groups private-by-default and require invite or approval to join.
              </p>
            </div>
          </div>
        </section>
      </main>

      {/* Fixed Footer */}
      <footer className="relative z-20 flex-shrink-0 border-t border-white/10 bg-black/60 backdrop-blur">
        <div className="mx-auto max-w-6xl px-4 sm:px-6 py-4 flex flex-col sm:flex-row gap-4 sm:items-center sm:justify-between">
          <div className="flex items-center gap-3">
            <img
              src="https://content.menunderfire.com/MenUnderFireIcon.webp"
              alt="Men Under Fire"
              width={36}
              height={36}
              loading="lazy"
              decoding="async"
              className="h-9 w-9 rounded-xl"
            />
            <div className="text-sm text-white/70">
              <div className="font-semibold text-white/85">Men Under Fire</div>
              <div className="text-white/55">Accountability • Brotherhood • Faith</div>
            </div>
          </div>

          <div className="flex flex-wrap gap-x-6 gap-y-2 text-sm text-white/65">
            <a className="hover:text-white transition" href="#faq">FAQ</a>
            <button className="hover:text-white transition" onClick={() => navigate('/page/about')}>About</button>
            <button className="hover:text-white transition" onClick={() => navigate('/terms')}>Terms</button>
            <button className="hover:text-white transition" onClick={() => navigate('/privacy')}>Privacy</button>
          </div>
        </div>
      </footer>
    </div>
  );
}
