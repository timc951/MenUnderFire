import { Group } from '../../types';
import { GroupCard } from './GroupCard';
import { LoadingSpinner } from '../common/LoadingSpinner';

interface GroupListProps {
  groups: Group[];
  isLoading?: boolean;
  error?: string | null;
  onGroupClick?: (group: Group) => void;
}

export function GroupList({ groups, isLoading, error, onGroupClick }: GroupListProps) {
  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <p className="text-red-600">{error}</p>
      </div>
    );
  }

  if (groups.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-500">No groups yet. Create or join one to get started!</p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {groups.map((group) => (
        <GroupCard
          key={group.id}
          group={group}
          onClick={onGroupClick ? () => onGroupClick(group) : undefined}
        />
      ))}
    </div>
  );
}
