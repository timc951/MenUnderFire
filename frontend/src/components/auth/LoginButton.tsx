import { useAuth } from '../../hooks/useAuth';
import { Button } from '../common/Button';

export function LoginButton() {
  const { login } = useAuth();

  return (
    <Button onClick={() => login()} variant="primary">
      Log In
    </Button>
  );
}
