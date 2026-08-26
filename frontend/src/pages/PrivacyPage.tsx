import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Button } from '../components/common/Button';
import { PagesMenu } from '../components/common/PagesMenu';

export function PrivacyPage() {
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
              <h1 className="text-3xl font-bold text-white mb-6">Privacy Policy</h1>

              <div className="prose prose-invert max-w-none text-stone-300 space-y-6">
                <p>
                  Your privacy is important to us. This Privacy Policy explains how Men Under Fire collects, uses, and protects your information.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">1. Information We Collect</h2>
                <p>
                  We collect information you provide directly, including your name, email address, and any content you share within accountability groups.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">2. How We Use Your Information</h2>
                <p>
                  Your information is used to provide and improve our services, facilitate group communications, and send important updates about your account.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">3. Information Sharing</h2>
                <p>
                  We do not sell your personal information. Content shared within groups is only visible to group members. We may share information when required by law.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">4. Data Security</h2>
                <p>
                  We implement appropriate security measures to protect your information. However, no method of transmission over the internet is 100% secure.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">5. Your Rights</h2>
                <p>
                  You have the right to access, update, or delete your personal information. Contact us if you wish to exercise these rights.
                </p>

                <h2 className="text-xl font-semibold text-white mt-8">6. Contact Us</h2>
                <p>
                  If you have questions about this Privacy Policy, please contact us through the platform.
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
