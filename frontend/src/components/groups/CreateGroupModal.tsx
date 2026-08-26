import { useState, FormEvent, useMemo, useEffect } from 'react';
import { Modal } from '../common/Modal';
import { Input } from '../common/Input';
import { Button } from '../common/Button';
import { useOrganizations, Organization } from '../../hooks/useOrganizations';
import { useRbac } from '../../rbac';

interface CreateGroupModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: { name: string; description: string; organizationId: string }) => Promise<void>;
}

export function CreateGroupModal({ isOpen, onClose, onSubmit }: CreateGroupModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [organizationId, setOrganizationId] = useState('');
  const [orgSearchTerm, setOrgSearchTerm] = useState('');
  const [isOrgDropdownOpen, setIsOrgDropdownOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const { organizations, isLoading: orgsLoading } = useOrganizations();
  const rbac = useRbac();
  const isSiteAdmin = rbac.isSiteAdmin();
  const isOrgAdmin = rbac.isOrgAdmin();
  const adminOrgIds = rbac.permissions?.adminOfOrganizationIds || [];

  // For org admins, filter to only organizations they admin
  const availableOrganizations = useMemo(() => {
    if (isSiteAdmin) return organizations;
    return organizations.filter(org => adminOrgIds.includes(org.id));
  }, [organizations, isSiteAdmin, adminOrgIds]);

  // Auto-select if there's only one organization available
  useEffect(() => {
    if (availableOrganizations.length === 1 && !organizationId) {
      setOrganizationId(availableOrganizations[0].id);
    }
  }, [availableOrganizations, organizationId]);

  // Filter organizations based on search term
  const filteredOrganizations = useMemo(() => {
    if (!orgSearchTerm.trim()) return availableOrganizations;
    const searchLower = orgSearchTerm.toLowerCase();
    return availableOrganizations.filter(org =>
      org.name.toLowerCase().includes(searchLower)
    );
  }, [availableOrganizations, orgSearchTerm]);

  // Get selected organization name
  const selectedOrg = organizations.find(org => org.id === organizationId);

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    if (!name.trim()) newErrors.name = 'Group name is required';
    if (!organizationId) newErrors.organization = 'Please select an organization';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Show org dropdown if site admin, or org admin with multiple orgs
  const showOrgDropdown = isSiteAdmin || (isOrgAdmin && availableOrganizations.length > 1);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    setIsSubmitting(true);
    try {
      await onSubmit({ name, description, organizationId });
      setName('');
      setDescription('');
      setOrganizationId('');
      setOrgSearchTerm('');
      onClose();
    } catch {
      setErrors({ submit: 'Failed to create group' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSelectOrg = (org: Organization) => {
    setOrganizationId(org.id);
    setOrgSearchTerm('');
    setIsOrgDropdownOpen(false);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Create Group">
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Group Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          error={errors.name}
          placeholder="Enter group name"
        />

        {/* Organization dropdown - shown for site admins or org admins with multiple orgs */}
        {showOrgDropdown && (
          <div className="space-y-1">
            <label htmlFor="organization" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              Organization
            </label>
            <div className="relative">
              <div
                className={`block w-full rounded-md border ${errors.organization ? 'border-red-500' : 'border-gray-300 dark:border-stone-600'} bg-white dark:bg-stone-700 px-3 py-2 shadow-sm cursor-pointer`}
                onClick={() => setIsOrgDropdownOpen(!isOrgDropdownOpen)}
              >
                {selectedOrg ? (
                  <span className="text-gray-900 dark:text-white">{selectedOrg.name}</span>
                ) : (
                  <span className="text-gray-400">Select an organization...</span>
                )}
              </div>

              {isOrgDropdownOpen && (
                <div className="absolute z-10 mt-1 w-full bg-white dark:bg-stone-700 border border-gray-300 dark:border-stone-600 rounded-md shadow-lg max-h-60 overflow-hidden">
                  <div className="p-2 border-b border-gray-200 dark:border-stone-700">
                    <input
                      type="text"
                      value={orgSearchTerm}
                      onChange={(e) => setOrgSearchTerm(e.target.value)}
                      placeholder="Search organizations..."
                      className="block w-full rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-600 px-3 py-2 text-sm text-gray-900 dark:text-white placeholder-gray-400 focus:border-primary-500 focus:ring-1 focus:ring-primary-500"
                      onClick={(e) => e.stopPropagation()}
                    />
                  </div>
                  <ul className="max-h-40 overflow-y-auto">
                    {orgsLoading ? (
                      <li className="px-3 py-2 text-gray-500">Loading...</li>
                    ) : filteredOrganizations.length === 0 ? (
                      <li className="px-3 py-2 text-gray-500">No organizations found</li>
                    ) : (
                      filteredOrganizations.map((org) => (
                        <li
                          key={org.id}
                          className={`px-3 py-2 cursor-pointer hover:bg-gray-100 dark:hover:bg-stone-600 ${org.id === organizationId ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-400' : 'text-gray-900 dark:text-white'}`}
                          onClick={() => handleSelectOrg(org)}
                        >
                          {org.name}
                        </li>
                      ))
                    )}
                  </ul>
                </div>
              )}
            </div>
            {errors.organization && <p className="text-sm text-red-600">{errors.organization}</p>}
          </div>
        )}

        <div className="space-y-1">
          <label htmlFor="group-description" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            Description
          </label>
          <textarea
            id="group-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={3}
            className="block w-full rounded-md border border-gray-300 dark:border-stone-600 bg-white dark:bg-stone-700 px-3 py-2 shadow-sm text-gray-900 dark:text-white focus:border-primary-500 focus:ring-1 focus:ring-primary-500"
            placeholder="Describe the group's purpose"
          />
        </div>
        {errors.submit && <p className="text-sm text-red-600">{errors.submit}</p>}
        <div className="flex gap-3 justify-end">
          <Button variant="secondary" type="button" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" isLoading={isSubmitting}>
            Create Group
          </Button>
        </div>
      </form>
    </Modal>
  );
}
