import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useOrganization } from '../hooks/useOrganizations';
import { useApi } from '../hooks/useApi';
import { CreateInvitationModal } from '../rbac/components/CreateInvitationModal';
import { CreateGroupModal } from '../components/groups/CreateGroupModal';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { Group } from '../types';

interface OrgAdmin {
  id: string;
  userId: string;
  organizationId: string;
  displayName: string;
  email: string;
  joinedAt: string;
}

export function OrganizationDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const api = useApi();
  const { organization, isLoading, error, updateOrganization } = useOrganization(id || '');

  const [isEditing, setIsEditing] = useState(false);
  const [showInviteModal, setShowInviteModal] = useState(false);
  const [showCreateGroupModal, setShowCreateGroupModal] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);

  // Groups and admins state
  const [groups, setGroups] = useState<Group[]>([]);
  const [admins, setAdmins] = useState<OrgAdmin[]>([]);
  const [groupsLoading, setGroupsLoading] = useState(false);
  const [adminsLoading, setAdminsLoading] = useState(false);

  const fetchGroups = useCallback(async () => {
    if (!id) return;
    setGroupsLoading(true);
    try {
      const data = await api.get<Group[]>(`/organizations/${id}/groups`);
      setGroups(data);
    } catch (err) {
      console.error('Failed to load groups:', err);
    } finally {
      setGroupsLoading(false);
    }
  }, [api, id]);

  const fetchAdmins = useCallback(async () => {
    if (!id) return;
    setAdminsLoading(true);
    try {
      const data = await api.get<OrgAdmin[]>(`/organizations/${id}/admins`);
      setAdmins(data);
    } catch (err) {
      console.error('Failed to load admins:', err);
    } finally {
      setAdminsLoading(false);
    }
  }, [api, id]);

  useEffect(() => {
    if (organization) {
      setName(organization.name);
      setDescription(organization.description || '');
      fetchGroups();
      fetchAdmins();
    }
  }, [organization, fetchGroups, fetchAdmins]);

  const handleSave = async () => {
    if (!name.trim()) return;

    setIsSaving(true);
    setSaveError(null);

    try {
      await updateOrganization({
        name: name.trim(),
        description: description.trim() || undefined,
      });
      setIsEditing(false);
    } catch (err) {
      setSaveError(err instanceof Error ? err.message : 'Failed to save changes');
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = () => {
    if (organization) {
      setName(organization.name);
      setDescription(organization.description || '');
    }
    setIsEditing(false);
    setSaveError(null);
  };

  const handleCreateGroup = async (data: { name: string; description: string; organizationId: string }) => {
    await api.post('/groups', {
      name: data.name,
      description: data.description,
      organizationId: id,
    });
    setShowCreateGroupModal(false);
    fetchGroups();
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-amber-600" />
      </div>
    );
  }

  if (error || !organization) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">{error || 'Organization not found'}</p>
        <Button variant="secondary" className="mt-4" onClick={() => navigate('/organizations')}>
          Back to Organizations
        </Button>
      </div>
    );
  }

  const isSiteAdmin = organization.isSiteAdmin;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate('/organizations')}
            className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            aria-label="Back"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {isEditing ? 'Edit Organization' : organization.name}
            </h1>
            {isSiteAdmin && (
              <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400 mt-1">
                Site Admin View
              </span>
            )}
          </div>
        </div>
        {organization.canEdit && !isEditing && (
          <Button onClick={() => setIsEditing(true)}>
            Edit
          </Button>
        )}
      </div>

      <Card>
        {isEditing ? (
          <div className="space-y-4">
            {saveError && (
              <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-md text-sm">
                {saveError}
              </div>
            )}

            <div>
              <label htmlFor="edit-name" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Organization Name
              </label>
              {isSiteAdmin ? (
                <input
                  type="text"
                  id="edit-name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
                />
              ) : (
                <>
                  <p className="px-3 py-2 bg-gray-100 dark:bg-stone-700 text-gray-900 dark:text-white rounded-md">
                    {organization.name}
                  </p>
                  <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    Organization name cannot be changed
                  </p>
                </>
              )}
            </div>

            <div>
              <label htmlFor="edit-description" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Description
              </label>
              <textarea
                id="edit-description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
              />
            </div>

            <div className="flex justify-end space-x-3">
              <Button variant="secondary" onClick={handleCancel} disabled={isSaving}>
                Cancel
              </Button>
              <Button onClick={handleSave} disabled={isSaving || (isSiteAdmin && !name.trim())}>
                {isSaving ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div>
              <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">Name</h3>
              <p className="mt-1 text-gray-900 dark:text-white">{organization.name}</p>
            </div>

            <div>
              <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">Description</h3>
              <p className="mt-1 text-gray-900 dark:text-white">
                {organization.description || <span className="text-gray-400 italic">No description</span>}
              </p>
            </div>

            <div className="pt-4 border-t border-gray-200 dark:border-stone-700">
              <dl className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Created</dt>
                  <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                    {new Date(organization.createdAt).toLocaleDateString()}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Last Updated</dt>
                  <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                    {new Date(organization.updatedAt).toLocaleDateString()}
                  </dd>
                </div>
              </dl>
            </div>
          </div>
        )}
      </Card>

      {isSiteAdmin && (
        <Card>
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Admin Actions
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
            As a Site Admin, you have full control over this organization.
          </p>
          <div className="flex flex-wrap gap-2">
            <Button
              variant="secondary"
              onClick={() => setShowInviteModal(true)}
            >
              Invite Org Admin
            </Button>
          </div>
        </Card>
      )}

      {/* Groups Section */}
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Groups
          </h2>
          <Button onClick={() => setShowCreateGroupModal(true)}>
            Create Group
          </Button>
        </div>
        {groupsLoading ? (
          <div className="flex items-center justify-center py-4">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-amber-600" />
          </div>
        ) : groups.length === 0 ? (
          <p className="text-sm text-gray-500 dark:text-stone-400">
            No groups in this organization yet.
          </p>
        ) : (
          <div className="space-y-2">
            {groups.map((group) => (
              <div
                key={group.id}
                onClick={() => navigate(`/groups/${group.id}`)}
                className="flex items-center justify-between p-3 bg-stone-50 dark:bg-stone-700/50 rounded-lg cursor-pointer hover:bg-stone-100 dark:hover:bg-stone-700 transition-colors"
              >
                <div>
                  <h3 className="font-medium text-gray-900 dark:text-white">{group.name}</h3>
                  {group.description && (
                    <p className="text-sm text-gray-600 dark:text-stone-400 truncate max-w-md">
                      {group.description}
                    </p>
                  )}
                </div>
                <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            ))}
          </div>
        )}
      </Card>

      {/* Forms Section */}
      <Card>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Forms
          </h2>
          <Button onClick={() => navigate(`/organizations/${id}/forms`)}>
            Manage Forms
          </Button>
        </div>
        <p className="text-sm text-gray-500 dark:text-stone-400">
          Create and manage custom forms for data collection.
        </p>
      </Card>

      {/* Organization Admins Section */}
      <Card>
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Organization Admins
        </h2>
        {adminsLoading ? (
          <div className="flex items-center justify-center py-4">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-amber-600" />
          </div>
        ) : admins.length === 0 ? (
          <p className="text-sm text-gray-500 dark:text-stone-400">
            No organization admins yet.
          </p>
        ) : (
          <div className="space-y-2">
            {admins.map((admin) => (
              <div
                key={admin.id}
                className="flex items-center justify-between p-3 bg-stone-50 dark:bg-stone-700/50 rounded-lg"
              >
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">
                    {admin.displayName}
                  </p>
                  <p className="text-sm text-gray-600 dark:text-stone-400">
                    {admin.email}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-stone-500">
                    Joined {new Date(admin.joinedAt).toLocaleDateString()}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>

      <CreateInvitationModal
        isOpen={showInviteModal}
        onClose={() => setShowInviteModal(false)}
        defaultType="org_admin"
        defaultOrganizationId={organization?.id}
      />

      <CreateGroupModal
        isOpen={showCreateGroupModal}
        onClose={() => setShowCreateGroupModal(false)}
        onSubmit={handleCreateGroup}
      />
    </div>
  );
}
