import { useState, useEffect, FormEvent } from 'react';
import { useRbac } from '../RbacContext';
import { useApi } from '../../hooks/useApi';
import { QRCodeSVG } from 'qrcode.react';

type InvitationType = 'org_admin' | 'group_owner' | 'group_member';

interface Organization {
  id: string;
  name: string;
}

interface Group {
  id: string;
  name: string;
  organizationId?: string;
}

interface InvitationResponse {
  id: string;
  token: string;
  email: string;
  type: string;
  organizationId?: string;
  groupId?: string;
  status: string;
  expiresAt: string;
  createdAt: string;
}

interface CreateInvitationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
  /** Pre-select invitation type */
  defaultType?: InvitationType;
  /** Pre-select organization (for org_admin invitations) */
  defaultOrganizationId?: string;
  /** Pre-select group (for group_owner/group_member invitations) */
  defaultGroupId?: string;
}

export function CreateInvitationModal({
  isOpen,
  onClose,
  onSuccess,
  defaultType,
  defaultOrganizationId,
  defaultGroupId,
}: CreateInvitationModalProps) {
  const rbac = useRbac();
  const api = useApi();

  const [invitationType, setInvitationType] = useState<InvitationType | null>(null);
  const [email, setEmail] = useState('');
  const [organizationId, setOrganizationId] = useState('');
  const [groupId, setGroupId] = useState('');
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [groups, setGroups] = useState<Group[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [createdInvitation, setCreatedInvitation] = useState<InvitationResponse | null>(null);
  const [copied, setCopied] = useState(false);

  // Determine available invitation types based on permissions
  const canInviteOrgAdmins = rbac.canInviteOrgAdmins();
  const canInviteGroupOwners = rbac.canInviteGroupOwners();
  const canInviteMembers = rbac.isSiteAdmin() || rbac.isOrgAdmin() || rbac.isGroupOwner();

  // Check if user can send any invitations
  const hasAnyPermission = canInviteOrgAdmins || canInviteGroupOwners || canInviteMembers;

  // Load organizations and groups when modal opens
  useEffect(() => {
    if (isOpen && hasAnyPermission) {
      // Load organizations if user can invite org admins
      if (canInviteOrgAdmins) {
        api.get<Organization[]>('/organizations').then(setOrganizations).catch(() => {});
      }
      // Load groups if user can invite to groups
      if (canInviteGroupOwners || canInviteMembers) {
        api.get<Group[]>('/groups').then(setGroups).catch(() => {});
      }
    }
  }, [isOpen, api, canInviteOrgAdmins, canInviteGroupOwners, canInviteMembers, hasAnyPermission]);

  // Reset form when modal closes, or set defaults when it opens
  useEffect(() => {
    if (!isOpen) {
      setInvitationType(null);
      setEmail('');
      setOrganizationId('');
      setGroupId('');
      setError(null);
      setCreatedInvitation(null);
      setCopied(false);
    } else {
      // Set defaults when modal opens
      if (defaultType) setInvitationType(defaultType);
      if (defaultOrganizationId) setOrganizationId(defaultOrganizationId);
      if (defaultGroupId) setGroupId(defaultGroupId);
    }
  }, [isOpen, defaultType, defaultOrganizationId, defaultGroupId]);

  // Filter groups based on user permissions
  const availableGroups = groups.filter((group) => {
    if (rbac.isSiteAdmin() || rbac.isOrgAdmin()) return true;
    return rbac.permissions?.ownedGroupIds?.includes(group.id);
  });

  // Generate invite link from token
  const getInviteLink = (token: string) => {
    return `${window.location.origin}/invite/${token}`;
  };

  // Copy invite link to clipboard
  const copyToClipboard = async () => {
    if (!createdInvitation) return;

    const link = getInviteLink(createdInvitation.token);
    try {
      await navigator.clipboard.writeText(link);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = link;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    if (!email || !invitationType) return;

    setIsSubmitting(true);
    setError(null);

    try {
      let endpoint = '';
      let body: Record<string, string> = { email };

      switch (invitationType) {
        case 'org_admin':
          endpoint = '/invitations/org-admin';
          body.organizationId = organizationId;
          break;
        case 'group_owner':
          endpoint = '/invitations/group-owner';
          body.groupId = groupId;
          break;
        case 'group_member':
          endpoint = '/invitations/group-member';
          body.groupId = groupId;
          break;
      }

      const response = await api.post<InvitationResponse>(endpoint, body);
      setCreatedInvitation(response);
      onSuccess?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to send invitation');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    onClose();
  };

  const handleCreateAnother = () => {
    setCreatedInvitation(null);
    setEmail('');
    setCopied(false);
  };

  if (!isOpen) return null;

  return (
    <div
      role="dialog"
      aria-modal="true"
      className="fixed inset-0 z-50 flex items-center justify-center"
    >
      <div className="fixed inset-0 bg-black bg-opacity-50" onClick={handleClose} />
      <div className="relative bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
        <button
          onClick={handleClose}
          aria-label="Close"
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>

        <h2 className="text-xl font-semibold mb-4 text-gray-900 dark:text-white">
          {createdInvitation ? 'Invitation Created' : 'Create Invitation'}
        </h2>

        {!hasAnyPermission ? (
          <div className="text-center py-8">
            <p className="text-gray-600 dark:text-gray-400">You don't have permission to send invitations.</p>
          </div>
        ) : createdInvitation ? (
          // Success state - show invite link
          <div className="space-y-4">
            <div className="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
              <div className="flex items-center mb-2">
                <svg className="w-5 h-5 text-green-600 dark:text-green-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-green-800 dark:text-green-300 font-medium">Invitation sent successfully!</span>
              </div>
              <p className="text-sm text-green-700 dark:text-green-400">
                An invitation has been created for <strong>{createdInvitation.email}</strong>
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Invite Link
              </label>
              <div className="flex items-center space-x-2">
                <input
                  type="text"
                  readOnly
                  value={getInviteLink(createdInvitation.token)}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                />
                <button
                  type="button"
                  onClick={copyToClipboard}
                  className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                    copied
                      ? 'bg-green-600 text-white'
                      : 'bg-primary-600 text-white hover:bg-primary-700'
                  }`}
                >
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
                Share this link with the invitee. The invitation expires on {new Date(createdInvitation.expiresAt).toLocaleDateString()}.
              </p>
            </div>

            <div className="flex justify-center">
              <div className="bg-white p-3 rounded-lg shadow-sm">
                <QRCodeSVG
                  value={getInviteLink(createdInvitation.token)}
                  size={160}
                  level="M"
                  marginSize={0}
                />
                <p className="text-xs text-center text-gray-500 mt-1">Scan to accept invitation</p>
              </div>
            </div>

            <div className="flex justify-end space-x-3 mt-6">
              <button
                type="button"
                onClick={handleCreateAnother}
                className="px-4 py-2 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700"
              >
                Create Another
              </button>
              <button
                type="button"
                onClick={handleClose}
                className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
              >
                Done
              </button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit}>
            {error && (
              <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 rounded-md text-sm">
                {error}
              </div>
            )}

            {/* Show type selection only if no default type is provided */}
            {defaultType ? (
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Invitation Type
                </label>
                <div className="px-3 py-2 bg-gray-100 dark:bg-gray-700 rounded-md text-gray-900 dark:text-white">
                  {defaultType === 'org_admin' && 'Organization Admin'}
                  {defaultType === 'group_owner' && 'Group Owner'}
                  {defaultType === 'group_member' && 'Group Member'}
                </div>
              </div>
            ) : (
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Invitation Type
                </label>
                <div className="space-y-2">
                  {canInviteOrgAdmins && (
                    <label className="flex items-center text-gray-900 dark:text-gray-100">
                      <input
                        type="radio"
                        name="invitationType"
                        value="org_admin"
                        checked={invitationType === 'org_admin'}
                        onChange={() => setInvitationType('org_admin')}
                        className="mr-2"
                        aria-label="Invite Organization Admin"
                      />
                      <span>Organization Admin</span>
                    </label>
                  )}
                  {canInviteGroupOwners && (
                    <label className="flex items-center text-gray-900 dark:text-gray-100">
                      <input
                        type="radio"
                        name="invitationType"
                        value="group_owner"
                        checked={invitationType === 'group_owner'}
                        onChange={() => setInvitationType('group_owner')}
                        className="mr-2"
                        aria-label="Invite Group Owner"
                      />
                      <span>Group Owner</span>
                    </label>
                  )}
                  {canInviteMembers && (
                    <label className="flex items-center text-gray-900 dark:text-gray-100">
                      <input
                        type="radio"
                        name="invitationType"
                        value="group_member"
                        checked={invitationType === 'group_member'}
                        onChange={() => setInvitationType('group_member')}
                        className="mr-2"
                        aria-label="Invite Group Member"
                      />
                      <span>Group Member</span>
                    </label>
                  )}
                </div>
              </div>
            )}

            <div className="mb-4">
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Email Address
              </label>
              <input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                placeholder="user@example.com"
              />
            </div>

            {invitationType === 'org_admin' && (
              <div className="mb-4">
                <label
                  htmlFor="organization"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  Organization
                </label>
                {defaultOrganizationId ? (
                  <div className="px-3 py-2 bg-gray-100 dark:bg-gray-700 rounded-md text-gray-900 dark:text-white">
                    {organizations.find(o => o.id === defaultOrganizationId)?.name || 'Loading...'}
                  </div>
                ) : (
                  <select
                    id="organization"
                    value={organizationId}
                    onChange={(e) => setOrganizationId(e.target.value)}
                    required
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                  >
                    <option value="">Select an organization</option>
                    {organizations.map((org) => (
                      <option key={org.id} value={org.id}>
                        {org.name}
                      </option>
                    ))}
                  </select>
                )}
              </div>
            )}

            {(invitationType === 'group_owner' || invitationType === 'group_member') && (
              <div className="mb-4">
                <label
                  htmlFor="group"
                  className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                >
                  Select Group
                </label>
                <select
                  id="group"
                  aria-label="Select Group"
                  value={groupId}
                  onChange={(e) => setGroupId(e.target.value)}
                  required
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                >
                  <option value="">Select a group</option>
                  {availableGroups.map((group) => (
                    <option key={group.id} value={group.id}>
                      {group.name}
                    </option>
                  ))}
                </select>
              </div>
            )}

            <div className="flex justify-end space-x-3 mt-6">
              <button
                type="button"
                onClick={handleClose}
                className="px-4 py-2 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSubmitting || !invitationType || !email}
                className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? 'Creating...' : 'Create Invitation'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
