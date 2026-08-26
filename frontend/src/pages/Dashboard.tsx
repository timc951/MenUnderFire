import { useRbac } from '../rbac';
import { useDashboardStats } from '../hooks/useDashboardStats';

export function Dashboard() {
  const rbac = useRbac();
  const { stats, isLoading, error } = useDashboardStats();

  const canViewStats = rbac.isSiteAdmin() || rbac.isOrgAdmin();
  // Site admins see org count, org admins only see group count
  const showOrgCount = stats?.organizationCount != null && stats.organizationCount > 0;
  const hasStats = stats && (showOrgCount || stats.groupCount > 0);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
      </div>

      {isLoading ? (
        <div className="text-gray-600 dark:text-gray-400">Loading statistics...</div>
      ) : error ? (
        <div className="text-red-600 dark:text-red-400">{error}</div>
      ) : canViewStats && hasStats ? (
        <div className={`grid grid-cols-1 ${showOrgCount ? 'md:grid-cols-2' : ''} gap-6`}>
          {showOrgCount && (
            <StatCard
              title="Organizations"
              count={stats.organizationCount!}
              description="Total organizations in the system"
              href="/organizations"
            />
          )}
          <StatCard
            title="Groups"
            count={stats.groupCount}
            description={rbac.isSiteAdmin() ? 'Total groups across all organizations' : 'Groups in your organizations'}
            href="/groups"
          />
        </div>
      ) : (
        <div className="text-gray-600 dark:text-gray-400">
          Welcome to Men Under Fire. Use the navigation to access your groups and reports.
        </div>
      )}
    </div>
  );
}

interface StatCardProps {
  title: string;
  count: number;
  description: string;
  href?: string;
}

function StatCard({ title, count, description, href }: StatCardProps) {
  const content = (
    <div className="bg-white dark:bg-stone-800 rounded-lg shadow p-6 border border-stone-200 dark:border-stone-700">
      <h2 className="text-sm font-medium text-stone-500 dark:text-stone-400 uppercase tracking-wide">
        {title}
      </h2>
      <p className="mt-2 text-3xl font-bold text-stone-900 dark:text-white">
        {count}
      </p>
      <p className="mt-1 text-sm text-stone-600 dark:text-stone-400">
        {description}
      </p>
    </div>
  );

  if (href) {
    return (
      <a href={href} className="block hover:ring-2 hover:ring-primary-500 rounded-lg transition-all">
        {content}
      </a>
    );
  }

  return content;
}
