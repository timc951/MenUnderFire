import { Group } from '../../types';
import { Card } from '../common/Card';

interface GroupCardProps {
  group: Group;
  onClick?: () => void;
}

export function GroupCard({ group, onClick }: GroupCardProps) {
  return (
    <Card onClick={onClick} className="space-y-2">
      <div className="flex items-center justify-between">
        <h3 className="font-semibold text-gray-900 dark:text-white">{group.name}</h3>
        {group.role && (
          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
            group.role === 'leader' ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400' : 'bg-stone-100 text-stone-700 dark:bg-stone-700 dark:text-stone-300'
          }`}>
            {group.role}
          </span>
        )}
      </div>
      <p className="text-sm text-gray-600 dark:text-stone-400">{group.description}</p>
      {group.memberCount !== undefined && (
        <p className="text-xs text-gray-400 dark:text-stone-500">{group.memberCount} members</p>
      )}
    </Card>
  );
}
