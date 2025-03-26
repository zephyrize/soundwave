import React from 'react';
import { useUserStore } from '../store/userStore';

interface PermissionGuardProps {
  permission: string;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

const PermissionGuard: React.FC<PermissionGuardProps> = ({
  permission,
  children,
  fallback = null,
}) => {
  const user = useUserStore((state) => state.user);
  
  if (!user?.permissions.includes(permission)) {
    return <>{fallback}</>;
  }
  
  return <>{children}</>;
};

export default PermissionGuard; 