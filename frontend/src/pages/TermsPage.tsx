import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Button } from '../components/common/Button';
import { PagesMenu } from '../components/common/PagesMenu';

export function TermsPage() {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading, login, logout } = useAuth();

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="h-screen bg-black relative flex flex-col">
      {/* Background image with overlay */}
      <div
        className="absolute inset-0 bg-cover bg-top bg-no-repeat"
        style={{ backgroundImage: 'url(https://content.menunderfire.com/poster.webp)' }}
      >
        <div className="absolute inset-0 bg-black/60" />
      </div>

      {/* Content */}
      <div className="relative z-10 flex flex-col h-full">
        {/* Header */}
        <header className="flex-shrink-0">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <PagesMenu />
                <img src="https://content.menunderfire.com/logo.webp" alt="Men Under Fire" width={48} height={48} className="w-12 h-12" />
                <span className="text-2xl font-bold text-amber-500">Men Under Fire</span>
              </div>
              <div className="flex items-center space-x-4">
                {isLoading ? (
                  <div className="w-24 h-10 bg-amber-600/30 rounded animate-pulse" />
                ) : isAuthenticated ? (
                  <>
                    <Button onClick={() => navigate('/dashboard')} className="bg-amber-600 hover:bg-amber-700 border-amber-600">
                      Go to Dashboard
                    </Button>
                    <Button
                      variant="secondary"
                      onClick={handleLogout}
                      className="bg-stone-800/80 hover:bg-stone-700/80 text-white border-stone-600"
                    >
                      Logout
                    </Button>
                  </>
                ) : (
                  <Button onClick={() => login()} className="bg-amber-600 hover:bg-amber-700 border-amber-600">
                    Sign In
                  </Button>
                )}
              </div>
            </div>
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 pb-4 flex flex-col overflow-y-auto">
          <div className="w-full max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 flex-1 flex flex-col">
            <article className="bg-stone-900/80 rounded-xl p-4 sm:p-6 lg:p-8 border border-stone-700/50 backdrop-blur-sm flex-1">
              <h1 className="text-3xl font-bold text-white mb-6">Terms of Use</h1>

              <div className="prose prose-invert max-w-none text-stone-300 space-y-6">
                <p>
                  Welcome to Men Under Fire. By accessing and using this platform, you agree to be bound by these Terms of Use.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">1. Acceptance of Terms</h2>
                <p>
                  By creating an account or using our services, you acknowledge that you have read, understood, and agree to these terms.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">2. User Conduct</h2>
                <p>
                  Users are expected to treat all members with respect and maintain confidentiality within their accountability groups.
                  Harassment, discrimination, or sharing of private information outside of groups is prohibited.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">3. Account Responsibilities</h2>
                <p>
                  You are responsible for maintaining the security of your account credentials and for all activities that occur under your account.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">4. Content Guidelines</h2>
                <p>
                  All content shared within groups should be appropriate and supportive. Content that is illegal, harmful, or violates the privacy of others is not permitted.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">5. Modifications</h2>
                <p>
                  We reserve the right to modify these terms at any time. Continued use of the platform after changes constitutes acceptance of the new terms.
                </p>

                <p className="text-stone-500 mt-8 text-sm">
                  Last updated: {new Date().toLocaleDateString()}
                </p>
              </div>
            </article>
          </div>
        </main>

        {/* Footer */}
        <footer className="flex-shrink-0 border-t border-stone-800">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <p className="text-center text-stone-500 text-sm">
              &copy; {new Date().getFullYear()} Men Under Fire. All rights reserved.
            </p>
          </div>
        </footer>
      </div>
    </div>
  );
}
