import { useAuth } from '../hooks/useAuth';
import { Card } from '../components/common/Card';

export function Profile() {
  const { user } = useAuth();

  if (!user) return null;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Profile</h1>
      <Card className="max-w-md space-y-4">
        <div className="flex items-center gap-4">
          {user.picture && (
            <img
              src={user.picture}
              alt={user.name}
              width={64}
              height={64}
              loading="lazy"
              className="h-16 w-16 rounded-full"
            />
          )}
          <div>
            <h2 className="text-lg font-semibold">{user.name}</h2>
            <p className="text-sm text-gray-500">{user.email}</p>
          </div>
        </div>
      </Card>
    </div>
  );
}
