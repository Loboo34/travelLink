'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/lib/stores/auth-store';
import { AdminSidebar } from '@/components/admin/admin-sidebar';
import { Spinner } from '@/components/ui/spinner';

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated, user, isLoading, isInitialized } = useAuthStore();
  const router = useRouter();
  console.log('AdminLayout - Auth State:', { user, isInitialized });

  useEffect(() => {
    if (!isInitialized) return;

    if (!isLoading && !isAuthenticated) {
      router.push('/login?redirect=/admin');
    } else if (!isLoading && user && user.role !== 'admin') {
      router.push('/');
    }
  }, [isInitialized, isAuthenticated, isLoading, user, router]);

  if (!isInitialized || isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-background">
        <Spinner className="h-8 w-8 text-primary" />
      </div>
    );
  }

  if (!isAuthenticated || (user && user.role !== "admin")) {
    return (
      <div className="flex h-screen items-center justify-center bg-background">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-foreground">Access Denied</h1>
          <p className="mt-2 text-muted-foreground">
            You don&apos;t have permission to access this page.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <AdminSidebar />
      <main className="pl-16 lg:pl-64 transition-all duration-300">
        <div className="p-6">{children}</div>
      </main>
    </div>
  );
}
