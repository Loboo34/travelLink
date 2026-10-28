'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Spinner } from '@/components/ui/spinner';

export default function AdminOffersPage() {
  const router = useRouter();

  useEffect(() => {
    router.replace('/admin/flights');
  }, [router]);

  return (
    <div className="flex h-64 items-center justify-center">
      <Spinner className="h-8 w-8 text-primary" />
    </div>
  );
}
