import { NavLink, useLocation } from 'react-router-dom';
import { useEffect } from 'react';
import { useRbac } from '../../rbac';

interface NavItem {
  path: string;
  label: string;
  icon: string;
  requiresAdmin?: boolean;
  requiresSiteAdmin?: boolean;
  requiresAnyLeader?: boolean;
  hidden?: boolean;
}

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export function Sidebar({ isOpen, onClose }: SidebarProps) {
  const rbac = useRbac();
  const location = useLocation();

  // Close sidebar on route change (mobile)
  useEffect(() => {
    onClose();
  }, [location.pathname]); // eslint-disable-line react-hooks/exhaustive-deps

  // Check if user can see Organizations (Site Admin or Org Admin)
  const isSiteAdmin = rbac.isSiteAdmin();
  const isOrgAdmin = rbac.isOrgAdmin();
  const isGroupOwner = rbac.isGroupOwner();
  const canSeeOrganizations = isSiteAdmin || isOrgAdmin;
  const canSeeInvitations = isSiteAdmin || isOrgAdmin || isGroupOwner;

  // Use singular "Organization" for org admins, plural for site admins
  const organizationLabel = isSiteAdmin ? 'Organizations' : 'Organization';

  const navItems: NavItem[] = [
    { path: '/dashboard', label: 'Dashboard', icon: '📊' },
    { path: '/organizations', label: organizationLabel, icon: '🏢', requiresAdmin: true },
    { path: '/groups', label: 'Groups', icon: '👥' },
    { path: '/invites', label: 'Invitations', icon: '✉️', requiresAnyLeader: true },
    { path: '/profile', label: 'Profile', icon: '👤' },
    { path: '/feedback', label: 'Feedback', icon: '💬' },
    { path: '/admin/hits', label: 'Site Analytics', icon: '📈', requiresSiteAdmin: true },
    { path: '/admin/pages', label: 'Manage Pages', icon: '📄', requiresSiteAdmin: true },
    { path: '/settings', label: 'Settings', icon: '⚙️' },
    { path: '/api-keys', label: 'API Keys', icon: '🔑', hidden: true }, // Hidden for now
  ];

  const visibleNavItems = navItems.filter((item) => {
    if (item.hidden) return false;
    if (item.requiresSiteAdmin) {
      return isSiteAdmin;
    }
    if (item.requiresAnyLeader) {
      return canSeeInvitations;
    }
    if (item.requiresAdmin) {
      return canSeeOrganizations;
    }
    return true;
  });

  return (
    <>
      {/* Mobile overlay backdrop */}
      {isOpen && (
        <div
          className="fixed inset-0 z-30 bg-black/50 md:hidden"
          onClick={onClose}
        />
      )}

      {/* Sidebar */}
      <nav
        className={`
          fixed inset-y-0 left-0 z-40 w-64 bg-stone-50 dark:bg-stone-800 border-r border-stone-200 dark:border-stone-700 p-4
          transform transition-transform duration-200 ease-in-out overflow-y-auto shrink-0
          ${isOpen ? 'translate-x-0' : '-translate-x-full'}
          md:relative md:translate-x-0 md:transition-none
        `}
        aria-label="Main navigation"
      >
        {/* Mobile close button */}
        <div className="flex items-center justify-between mb-4 md:hidden">
          <span className="text-sm font-semibold text-stone-700 dark:text-stone-300">Menu</span>
          <button
            onClick={onClose}
            className="p-1.5 rounded-md text-stone-600 dark:text-stone-300 hover:bg-stone-200 dark:hover:bg-stone-700"
            aria-label="Close menu"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <ul className="space-y-2">
          {visibleNavItems.map((item) => (
            <li key={item.path}>
              <NavLink
                to={item.path}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-primary-50 dark:bg-primary-900/30 text-primary-700 dark:text-primary-400'
                      : 'text-stone-600 dark:text-stone-300 hover:bg-stone-100 dark:hover:bg-stone-700 hover:text-stone-900 dark:hover:text-white'
                  }`
                }
              >
                <span>{item.icon}</span>
                <span>{item.label}</span>
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </>
  );
}
