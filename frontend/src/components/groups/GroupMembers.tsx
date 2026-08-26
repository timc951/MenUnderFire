import { GroupMembership } from '../../types';
import { Button } from '../common/Button';

interface GroupMembersProps {
  members: GroupMembership[];
  inviteCode?: string | null;
  currentUserRole: string;
  onRemoveMember?: (userId: string) => void;
  onChangeRole?: (userId: string, newRole: string) => void;
}

const ROLE_LEVEL: Record<string, number> = {
  OWNER: 4,
  LEADER: 3,
  MODERATOR: 2,
  MEMBER: 1,
};

const ROLE_BADGE_STYLES: Record<string, string> = {
  OWNER: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
  LEADER: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
  MODERATOR: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  MEMBER: 'bg-stone-100 text-stone-600 dark:bg-stone-700 dark:text-stone-300',
  ADMIN: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
};

export function GroupMembers({ members, inviteCode, currentUserRole, onRemoveMember, onChangeRole }: GroupMembersProps) {
  const isLeader = currentUserRole === 'LEADER' || currentUserRole === 'OWNER' || currentUserRole === 'ADMIN';
  const currentLevel = ROLE_LEVEL[currentUserRole] ?? 0;

  const canRemove = (memberRole: string) => {
    const targetLevel = ROLE_LEVEL[memberRole] ?? 0;
    return currentLevel > targetLevel;
  };

  const canChangeRole = (memberRole: string) => {
    if (!onChangeRole) return false;
    const targetLevel = ROLE_LEVEL[memberRole] ?? 0;
    if (currentUserRole === 'OWNER') return true;
    return currentLevel > targetLevel;
  };

  const getAvailableRoles = (memberRole: string) => {
    const allRoles = ['OWNER', 'LEADER', 'MODERATOR', 'MEMBER'];
    return allRoles.filter(r => {
      if (r === memberRole) return false;
      if (currentUserRole === 'OWNER') return true;
      return (ROLE_LEVEL[r] ?? 0) < currentLevel;
    });
  };

  return (
    <div className="space-y-4">
      {isLeader && inviteCode && (
        <div className="bg-stone-50 dark:bg-stone-700/50 rounded-md p-3">
          <p className="text-sm font-medium text-gray-700 dark:text-stone-300">Invite Code</p>
          <p className="text-lg font-mono font-bold text-amber-600 dark:text-amber-400">{inviteCode}</p>
        </div>
      )}

      <div className="space-y-2">
        <h3 className="text-sm font-medium text-gray-700 dark:text-stone-300">Members ({members.length})</h3>
        <ul className="divide-y divide-gray-200 dark:divide-stone-600">
          {members.map((member) => (
            <li key={member.id} className="flex items-center justify-between py-3">
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-3">
                  <span className="text-sm font-medium text-gray-900 dark:text-white">{member.displayName}</span>
                  <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                    ROLE_BADGE_STYLES[member.role] || ROLE_BADGE_STYLES.MEMBER
                  }`}>
                    {member.role}
                  </span>
                </div>
                {member.email && (
                  <span className="text-xs text-gray-500 dark:text-stone-400">{member.email}</span>
                )}
              </div>
              <div className="flex items-center gap-2">
                {canChangeRole(member.role) && (
                  <select
                    value=""
                    onChange={(e) => {
                      if (e.target.value) {
                        onChangeRole?.(member.userId, e.target.value);
                        e.target.value = '';
                      }
                    }}
                    className="text-xs border border-gray-300 dark:border-stone-600 rounded px-1 py-1 bg-white dark:bg-stone-700 text-gray-700 dark:text-stone-300"
                    aria-label={`Change role for ${member.displayName}`}
                  >
                    <option value="">Change role</option>
                    {getAvailableRoles(member.role).map(r => (
                      <option key={r} value={r}>{r}</option>
                    ))}
                  </select>
                )}
                {canRemove(member.role) && (
                  <Button
                    variant="danger"
                    size="sm"
                    onClick={() => onRemoveMember?.(member.userId)}
                    aria-label={`Remove ${member.displayName}`}
                  >
                    Remove
                  </Button>
                )}
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
