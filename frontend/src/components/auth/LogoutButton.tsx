import { useAuth } from '../../hooks/useAuth';
import { Button } from '../common/Button';

export function LogoutButton() {
  const { logout } = useAuth();

  return (
    <Button onClick={logout} variant="secondary">
      Log Out
    </Button>
  );
}
