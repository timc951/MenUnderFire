import { useState, useEffect, useCallback } from 'react';
import { useApi } from '../hooks/useApi';
import { useRbac, CreateInvitationModal } from '../rbac';
import { Button } from '../components/common/Button';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { InvitationListItem } from '../types';

type FilterType = 'ALL' | 'ORG_ADMIN' | 'GROUP_OWNER' | 'GROUP_MEMBER';
type FilterStatus = 'ALL' | 'PENDING' | 'ACCEPTED' | 'EXPIRED' | 'REVOKED';

export function Invites() {
  const api = useApi();
  const rbac = useRbac();
  const isSiteAdmin = rbac.isSiteAdmin();
  const isOrgAdmin = rbac.isOrgAdmin();

  const [invitations, setInvitations] = useState<InvitationListItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [filterType, setFilterType] = useState<FilterType>('ALL');
  const [filterStatus, setFilterStatus] = useState<FilterStatus>('ALL');
  const [searchQuery, setSearchQuery] = useState('');

  const fetchInvitations = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await api.get<InvitationListItem[]>('/invitations');
      setInvitations(data);
    } catch (err) {
      console.error('Failed to load invitations:', err);
    } finally {
      setIsLoading(false);
    }
  }, [api]);

  useEffect(() => {
    fetchInvitations();
  }, [fetchInvitations]);

  const handleDelete = async (id: string) => {
    try {
      await api.delete(`/invitations/${id}`);
      setInvitations((prev) => prev.filter((inv) => inv.id !== id));
    } catch (err) {
      console.error('Failed to revoke invitation:', err);
    }
  };

  const filteredInvitations = invitations.filter((inv) => {
    if (filterType !== 'ALL' && inv.type !== filterType) return false;
    if (filterStatus !== 'ALL' && inv.status !== filterStatus) return false;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      return (
        inv.email.toLowerCase().includes(q) ||
        (inv.organizationName && inv.organizationName.toLowerCase().includes(q)) ||
        (inv.groupName && inv.groupName.toLowerCase().includes(q))
      );
    }
    return true;
  });

  const typeLabel = (t: string) => {
    switch (t) {
      case 'ORG_ADMIN': return 'Org Admin';
      case 'GROUP_OWNER': return 'Group Owner';
      case 'GROUP_MEMBER': return 'Group Member';
      default: return t;
    }
  };

  const typeColor = (t: string) => {
    switch (t) {
      case 'ORG_ADMIN': return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400';
      case 'GROUP_OWNER': return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400';
      case 'GROUP_MEMBER': return 'bg-stone-100 text-stone-700 dark:bg-stone-700 dark:text-stone-300';
      default: return 'bg-stone-100 text-stone-600 dark:bg-stone-700 dark:text-stone-400';
    }
  };

  const statusColor = (status: string) => {
    switch (status) {
      case 'PENDING': return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400';
      case 'ACCEPTED': return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'EXPIRED': return 'bg-stone-100 text-stone-600 dark:bg-stone-700 dark:text-stone-400';
      case 'REVOKED': return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400';
      default: return 'bg-stone-100 text-stone-600 dark:bg-stone-700 dark:text-stone-400';
    }
  };

  const formatDate = (dateStr: string) =>
    new Date(dateStr).toLocaleDateString('en-US', {
      year: 'numeric', month: 'short', day: 'numeric',
    });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Invitations</h1>
          <p className="text-gray-600 dark:text-stone-400">
            {isSiteAdmin
              ? 'Manage all invitations across the platform'
              : isOrgAdmin
                ? 'Manage invitations for your organization and groups'
                : 'Manage invitations for your groups'}
          </p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          New Invitation
        </Button>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search by email, org, or group..."
          className="rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500 w-64"
        />
        <select
          value={filterType}
          onChange={(e) => setFilterType(e.target.value as FilterType)}
          className="rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500"
        >
          <option value="ALL">All Types</option>
          {isSiteAdmin && <option value="ORG_ADMIN">Org Admin</option>}
          <option value="GROUP_OWNER">Group Owner</option>
          <option value="GROUP_MEMBER">Group Member</option>
        </select>
        <select
          value={filterStatus}
          onChange={(e) => setFilterStatus(e.target.value as FilterStatus)}
          className="rounded-md border border-stone-300 dark:border-stone-600 bg-white dark:bg-stone-700 text-gray-900 dark:text-white px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-amber-500"
        >
          <option value="ALL">All Statuses</option>
          <option value="PENDING">Pending</option>
          <option value="ACCEPTED">Accepted</option>
          <option value="EXPIRED">Expired</option>
          <option value="REVOKED">Revoked</option>
        </select>
      </div>

      {isLoading ? (
        <LoadingSpinner size="lg" />
      ) : filteredInvitations.length === 0 ? (
        <div className="text-center py-12 text-gray-500 dark:text-stone-400">
          <p>No invitations found.</p>
          <p className="text-sm mt-1">Click &quot;New Invitation&quot; to invite someone.</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-stone-200 dark:border-stone-700">
                <th className="text-left text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Email</th>
                <th className="text-left text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Type</th>
                <th className="text-left text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Org / Group</th>
                <th className="text-left text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Status</th>
                <th className="text-left text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Invited By</th>
                <th className="text-left text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Created</th>
                <th className="text-left text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Expires</th>
                <th className="text-right text-xs font-medium text-gray-500 dark:text-stone-400 uppercase tracking-wider px-4 py-3">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stone-200 dark:divide-stone-700">
              {filteredInvitations.map((inv) => (
                <tr key={inv.id} className="hover:bg-stone-50 dark:hover:bg-stone-800/50">
                  <td className="px-4 py-3 text-sm text-gray-900 dark:text-white">{inv.email}</td>
                  <td className="px-4 py-3">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${typeColor(inv.type)}`}>
                      {typeLabel(inv.type)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-700 dark:text-stone-300">
                    {inv.organizationName && <div>{inv.organizationName}</div>}
                    {inv.groupName && <div className="text-xs text-gray-500 dark:text-stone-400">{inv.groupName}</div>}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusColor(inv.status)}`}>
                      {inv.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-stone-400">{inv.inviterName || '-'}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-stone-400">{formatDate(inv.createdAt)}</td>
                  <td className="px-4 py-3 text-sm text-gray-600 dark:text-stone-400">
                    <span className={new Date(inv.expiresAt) < new Date() ? 'text-red-600 dark:text-red-400' : ''}>
                      {formatDate(inv.expiresAt)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    {inv.status === 'PENDING' && (
                      <button
                        onClick={() => handleDelete(inv.id)}
                        className="text-xs text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300"
                      >
                        Revoke
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <CreateInvitationModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onSuccess={() => {
          setShowCreateModal(false);
          fetchInvitations();
        }}
      />
    </div>
  );
}
