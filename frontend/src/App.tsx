import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './auth/AuthProvider';
import { ProtectedRoute } from './auth/ProtectedRoute';
import { Layout } from './components/layout/Layout';
import { Dashboard } from './pages/Dashboard';
import { LandingPage } from './pages/LandingPage';
import { Groups } from './pages/Groups';
import { GroupDetail } from './pages/GroupDetail';
import { GroupJoin } from './pages/GroupJoin';
import { Profile } from './pages/Profile';
import { Settings } from './pages/Settings';
import { ApiKeys } from './pages/ApiKeys';
import { Organizations } from './pages/Organizations';
import { OrganizationDetail } from './pages/OrganizationDetail';
import { AcceptInvitationPage } from './pages/AcceptInvitation';
import { SitePages } from './pages/SitePages';
import { SitePageEdit } from './pages/SitePageEdit';
import { SitePageView } from './pages/SitePageView';
import { PagePreview } from './pages/PagePreview';
import { TermsPage } from './pages/TermsPage';
import { PrivacyPage } from './pages/PrivacyPage';
import { Forms } from './pages/Forms';
import { FormBuilder } from './pages/FormBuilder';
import { Feedback } from './pages/Feedback';
import { Invites } from './pages/Invites';
import { PageHitsAdmin } from './pages/PageHitsAdmin';
import { ThemeProvider } from './theme';
import { usePageTracking } from './hooks/usePageTracking';

function PageTracker() {
  usePageTracking();
  return null;
}

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
      <BrowserRouter>
        <PageTracker />
        <Routes>
          {/* Public landing page */}
          <Route
            path="/"
            element={<LandingPage />}
          />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Layout>
                  <Dashboard />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/groups"
            element={
              <ProtectedRoute>
                <Layout>
                  <Groups />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/groups/join"
            element={
              <ProtectedRoute>
                <Layout>
                  <GroupJoin />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/groups/:id"
            element={
              <ProtectedRoute>
                <Layout>
                  <GroupDetail />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Layout>
                  <Profile />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/settings"
            element={
              <ProtectedRoute>
                <Layout>
                  <Settings />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/api-keys"
            element={
              <ProtectedRoute>
                <Layout>
                  <ApiKeys />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/organizations"
            element={
              <ProtectedRoute>
                <Layout>
                  <Organizations />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/organizations/:id"
            element={
              <ProtectedRoute>
                <Layout>
                  <OrganizationDetail />
                </Layout>
              </ProtectedRoute>
            }
          />
          {/* Public route - no authentication required */}
          <Route
            path="/invite/:token"
            element={<AcceptInvitationPage />}
          />
          {/* Public site pages */}
          <Route
            path="/page/:slug"
            element={<SitePageView />}
          />
          {/* Page preview (public - accessed via temporary draft ID) */}
          <Route
            path="/page/preview/:draftId"
            element={<PagePreview />}
          />
          {/* Terms and Privacy pages */}
          <Route
            path="/terms"
            element={<TermsPage />}
          />
          <Route
            path="/privacy"
            element={<PrivacyPage />}
          />
          {/* Site Admin only - page management */}
          <Route
            path="/admin/pages"
            element={
              <ProtectedRoute>
                <Layout>
                  <SitePages />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/pages/:id"
            element={
              <ProtectedRoute>
                <Layout>
                  <SitePageEdit />
                </Layout>
              </ProtectedRoute>
            }
          />
          {/* Site Admin only - page hit analytics */}
          <Route
            path="/admin/hits"
            element={
              <ProtectedRoute>
                <Layout>
                  <PageHitsAdmin />
                </Layout>
              </ProtectedRoute>
            }
          />
          {/* Feedback */}
          <Route
            path="/feedback"
            element={
              <ProtectedRoute>
                <Layout>
                  <Feedback />
                </Layout>
              </ProtectedRoute>
            }
          />
          {/* Invitations */}
          <Route
            path="/invites"
            element={
              <ProtectedRoute>
                <Layout>
                  <Invites />
                </Layout>
              </ProtectedRoute>
            }
          />
          {/* Form builder routes */}
          <Route
            path="/organizations/:orgId/forms"
            element={
              <ProtectedRoute>
                <Layout>
                  <Forms />
                </Layout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/forms/:formId"
            element={
              <ProtectedRoute>
                <Layout>
                  <FormBuilder />
                </Layout>
              </ProtectedRoute>
            }
          />
        </Routes>
      </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
