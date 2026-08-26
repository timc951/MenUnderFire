import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useGroups } from '../hooks/useGroups';
import { useRbac } from '../rbac';
import { GroupList } from '../components/groups/GroupList';
import { CreateGroupModal } from '../components/groups/CreateGroupModal';
import { JoinGroupModal } from '../components/groups/JoinGroupModal';
import { Button } from '../components/common/Button';
import { Group } from '../types';

export function Groups() {
  const navigate = useNavigate();
  const rbac = useRbac();
  const { groups, isLoading, error, createGroup, joinGroup } = useGroups();
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showJoinModal, setShowJoinModal] = useState(false);

  const canCreateGroups = rbac.isSiteAdmin() || rbac.isOrgAdmin();

  const handleGroupClick = (group: Group) => {
    navigate(`/groups/${group.id}`);
  };

  const handleCreateGroup = async (data: { name: string; description?: string; organizationId: string }) => {
    await createGroup(data);
    setShowCreateModal(false);
  };

  const handleJoinGroup = async (inviteCode: string) => {
    await joinGroup(inviteCode);
    setShowJoinModal(false);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Groups</h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            View and manage your groups
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="secondary" onClick={() => setShowJoinModal(true)}>
            Join Group
          </Button>
          {canCreateGroups && (
            <Button onClick={() => setShowCreateModal(true)}>
              Create Group
            </Button>
          )}
        </div>
      </div>

      <GroupList
        groups={groups}
        isLoading={isLoading}
        error={error}
        onGroupClick={handleGroupClick}
      />

      <CreateGroupModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        onSubmit={handleCreateGroup}
      />

      <JoinGroupModal
        isOpen={showJoinModal}
        onClose={() => setShowJoinModal(false)}
        onSubmit={handleJoinGroup}
      />
    </div>
  );
}
