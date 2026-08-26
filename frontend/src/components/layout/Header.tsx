import { Link } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useRbac, RoleBadge } from '../../rbac';
import { RoleType } from '../../rbac/types';

function getUserRole(
  isSiteAdmin: boolean,
  isOrgAdmin: boolean,
  isGroupOwner: boolean
): RoleType {
  if (isSiteAdmin) return 'site_admin';
  if (isOrgAdmin) return 'org_admin';
  if (isGroupOwner) return 'group_owner';
  return 'member';
}

interface HeaderProps {
  onToggleSidebar?: () => void;
}

export function Header({ onToggleSidebar }: HeaderProps) {
  const { user, logout } = useAuth();
  const { isSiteAdmin, isOrgAdmin, isGroupOwner, loading } = useRbac();

  const role = getUserRole(isSiteAdmin(), isOrgAdmin(), isGroupOwner());

  return (
    <header className="bg-stone-50 dark:bg-stone-800 border-b border-stone-200 dark:border-stone-700 px-4 md:px-6 py-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          {onToggleSidebar && (
            <button
              onClick={onToggleSidebar}
              className="md:hidden p-1.5 rounded-md text-stone-600 dark:text-stone-300 hover:bg-stone-200 dark:hover:bg-stone-700"
              aria-label="Toggle menu"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          )}
          <Link to="/" className="flex items-center gap-3 hover:opacity-80 transition-opacity">
            <img src="https://content.menunderfire.com/logo.webp" alt="Men Under Fire" width={32} height={32} className="w-8 h-8" />
            <h1 className="text-xl font-bold text-primary-600 dark:text-primary-500 hidden sm:block">Men Under Fire</h1>
          </Link>
        </div>
        <div className="flex items-center gap-2 md:gap-4">
          {user && (
            <>
              <span className="text-sm text-stone-600 dark:text-stone-300 hidden sm:inline">{user.name}</span>
              {!loading && <span className="hidden sm:inline"><RoleBadge role={role} /></span>}
              <button
                onClick={logout}
                className="text-sm text-stone-500 dark:text-stone-400 hover:text-stone-700 dark:hover:text-stone-200"
              >
                Logout
              </button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
