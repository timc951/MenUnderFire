import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useOrganizations, useOrganization, Organization } from '../hooks/useOrganizations';
import { useRbac } from '../rbac';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { CreateInvitationModal } from '../rbac/components/CreateInvitationModal';

interface CreateOrganizationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: { name: string; description?: string }) => Promise<void>;
}

function CreateOrganizationModal({ isOpen, onClose, onSubmit }: CreateOrganizationModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await onSubmit({ name: name.trim(), description: description.trim() || undefined });
      setName('');
      setDescription('');
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create organization');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div
      role="dialog"
      aria-modal="true"
      className="fixed inset-0 z-50 flex items-center justify-center"
    >
      <div className="fixed inset-0 bg-black bg-opacity-50" onClick={onClose} />
      <div className="relative bg-white dark:bg-stone-800 rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
        <button
          onClick={onClose}
          aria-label="Close"
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-white">
          Create Organization
        </h2>

        <form onSubmit={handleSubmit}>
          {error && (
            <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-md text-sm">
              {error}
            </div>
          )}

          <div className="mb-4">
            <label htmlFor="org-name" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Organization Name
            </label>
            <input
              type="text"
              id="org-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
              placeholder="Enter organization name"
            />
          </div>

          <div className="mb-4">
            <label htmlFor="org-description" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Description (optional)
            </label>
            <textarea
              id="org-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
              placeholder="Enter organization description"
            />
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-700 dark:text-stone-300 border border-gray-300 dark:border-stone-600 rounded-md hover:bg-gray-50 dark:hover:bg-stone-700"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !name.trim()}
              className="px-4 py-2 bg-amber-600 text-white rounded-md hover:bg-amber-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? 'Creating...' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

interface OrganizationCardProps {
  organization: Organization;
  onClick: () => void;
}

function OrganizationCard({ organization, onClick }: OrganizationCardProps) {
  return (
    <Card
      className="cursor-pointer hover:shadow-md transition-shadow"
      onClick={onClick}
    >
      <div className="flex items-start justify-between">
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
            {organization.name}
          </h3>
          {organization.description && (
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
              {organization.description}
            </p>
          )}
          <p className="mt-2 text-xs text-gray-500 dark:text-gray-500">
            Created {new Date(organization.createdAt).toLocaleDateString()}
          </p>
        </div>
        <div className="ml-4">
          <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>
    </Card>
  );
}

// Single organization view for org admins (name read-only, description editable)
function OrgAdminOrganizationView({ organizationId }: { organizationId: string }) {
  const { organization, isLoading, error, updateOrganization } = useOrganization(organizationId);
  const [isEditing, setIsEditing] = useState(false);
  const [description, setDescription] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);
  const [showInviteModal, setShowInviteModal] = useState(false);

  useEffect(() => {
    if (organization) {
      setDescription(organization.description || '');
    }
  }, [organization]);

  const handleSave = async () => {
    if (!organization) return;

    setIsSaving(true);
    setSaveError(null);

    try {
      await updateOrganization({
        name: organization.name, // Keep the existing name (can't edit)
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
      setDescription(organization.description || '');
    }
    setIsEditing(false);
    setSaveError(null);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" />
      </div>
    );
  }

  if (error || !organization) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">{error || 'Organization not found'}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Organization</h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            View and manage your organization
          </p>
        </div>
        {!isEditing && (
          <Button onClick={() => setIsEditing(true)}>
            Edit Description
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
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Organization Name
              </label>
              <p className="px-3 py-2 bg-gray-100 dark:bg-stone-700 text-gray-900 dark:text-white rounded-md">
                {organization.name}
              </p>
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Organization name cannot be changed
              </p>
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
                className="w-full px-3 py-2 border border-gray-300 dark:border-stone-600 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white dark:bg-stone-700 text-gray-900 dark:text-white"
              />
            </div>

            <div className="flex justify-end space-x-3">
              <Button variant="secondary" onClick={handleCancel} disabled={isSaving}>
                Cancel
              </Button>
              <Button onClick={handleSave} disabled={isSaving}>
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

            <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
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

      <Card>
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Admin Actions
        </h2>
        <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
          As an Organization Admin, you can invite new members to your organization.
        </p>
        <div className="flex flex-wrap gap-2">
          <Button
            variant="secondary"
            onClick={() => setShowInviteModal(true)}
          >
            Invite Group Owner
          </Button>
        </div>
      </Card>

      <CreateInvitationModal
        isOpen={showInviteModal}
        onClose={() => setShowInviteModal(false)}
        defaultType="group_owner"
        defaultOrganizationId={organization?.id}
      />
    </div>
  );
}

export function Organizations() {
  const navigate = useNavigate();
  const rbac = useRbac();
  const { organizations, isLoading, error, createOrganization } = useOrganizations();
  const [showCreateModal, setShowCreateModal] = useState(false);

  const isSiteAdmin = rbac.isSiteAdmin();
  const isOrgAdmin = rbac.isOrgAdmin();

  // If user is neither site admin nor org admin, they shouldn't see this page
  if (!isSiteAdmin && !isOrgAdmin) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
            Access Denied
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            You don't have permission to view this page.
          </p>
        </div>
      </div>
    );
  }

  // For org admins (non-site-admins), show their single organization directly
  if (isOrgAdmin && !isSiteAdmin && !isLoading && organizations.length === 1) {
    return <OrgAdminOrganizationView organizationId={organizations[0].id} />;
  }

  const handleOrganizationClick = (org: Organization) => {
    navigate(`/organizations/${org.id}`);
  };

  const handleCreateOrganization = async (data: { name: string; description?: string }) => {
    await createOrganization(data);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Organizations</h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            {isSiteAdmin
              ? 'Manage all organizations in the system'
              : 'View and manage your organizations'}
          </p>
        </div>
        {isSiteAdmin && (
          <Button onClick={() => setShowCreateModal(true)}>
            Create Organization
          </Button>
        )}
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" />
        </div>
      ) : error ? (
        <div className="text-center py-12">
          <p className="text-red-600 dark:text-red-400">{error}</p>
        </div>
      ) : organizations.length === 0 ? (
        <div className="text-center py-12">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
            />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-white">
            No organizations
          </h3>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {isSiteAdmin
              ? 'Get started by creating a new organization.'
              : 'You are not an admin of any organizations yet.'}
          </p>
          {isSiteAdmin && (
            <div className="mt-6">
              <Button onClick={() => setShowCreateModal(true)}>
                Create Organization
              </Button>
            </div>
          )}
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {organizations.map((org) => (
            <OrganizationCard
              key={org.id}
              organization={org}
              onClick={() => handleOrganizationClick(org)}
            />
          ))}
        </div>
      )}

      <CreateOrganizationModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onSubmit={handleCreateOrganization}
      />
    </div>
  );
}
