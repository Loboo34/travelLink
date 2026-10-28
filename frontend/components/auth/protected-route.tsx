'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/stores';
import { Spinner } from '@/components/ui/spinner';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const router = useRouter();
  const { isAuthenticated, token, isInitialized } = useAuthStore();

  useEffect(() => {
    if (!isInitialized) return;

    // Check for token in localStorage as well
    const storedToken =
      typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;

    if (!isAuthenticated && !token && !storedToken) {
      router.push('/login');
    }
  }, [isInitialized, isAuthenticated, token, router]);

  if (!isInitialized || (!isAuthenticated && !token)) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <Spinner className="h-8 w-8 text-primary" />
      </div>
    );
  }

  return <>{children}</>;
}
