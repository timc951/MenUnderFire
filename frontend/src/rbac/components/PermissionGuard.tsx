import { ReactNode } from 'react';
import { useRbac } from '../RbacContext';
import { CapabilityCheck, ScopedCheck } from '../types';

type CheckType = CapabilityCheck | ScopedCheck;

interface PermissionGuardProps {
  check: CheckType;
  scopeId?: string;
  fallback?: ReactNode;
  children: ReactNode;
}

export function PermissionGuard({
  check,
  scopeId,
  fallback = null,
  children,
}: PermissionGuardProps) {
  const rbac = useRbac();

  if (rbac.loading) {
    return null;
  }

  // Get the check function from the rbac context
  const checkFn = rbac[check as keyof typeof rbac];

  if (typeof checkFn !== 'function') {
    console.warn(`PermissionGuard: Unknown check "${check}"`);
    return <>{fallback}</>;
  }

  // Call the check function, passing scopeId if provided
  const hasPermission = scopeId ? (checkFn as (id: string) => boolean)(scopeId) : (checkFn as () => boolean)();

  return hasPermission ? <>{children}</> : <>{fallback}</>;
}
