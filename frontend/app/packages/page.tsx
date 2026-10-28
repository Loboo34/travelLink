'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Search, Luggage } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { PackageCard } from '@/components/cards';
import { packagesApi } from '@/lib/api';
import type { Package } from '@/lib/types';

export default function PackagesPage() {
  const [packages, setPackages] = useState<Package[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchPackages() {
      try {
        const data = await packagesApi.getAll();
        setPackages(data);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : 'Failed to load packages'
        );
      } finally {
        setIsLoading(false);
      }
    }
    fetchPackages();
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">
              Travel Packages
            </h1>
            <p className="mt-1 text-muted-foreground">
              Curated travel experiences with everything included
            </p>
          </div>
          <Button asChild>
            <Link href="/search/packages">
              <Search className="mr-2 h-4 w-4" />
              Search Packages
            </Link>
          </Button>
        </div>

        {/* Content */}
        {isLoading ? (
          <div className="flex min-h-[40vh] items-center justify-center">
            <Spinner size="lg" />
          </div>
        ) : error ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <p className="text-destructive">{error}</p>
            <Button onClick={() => window.location.reload()}>Try Again</Button>
          </div>
        ) : packages.length === 0 ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Luggage className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-lg font-medium">No packages available</p>
            <p className="text-muted-foreground">
              Check back later for new packages
            </p>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {packages.map((pkg) => (
              <PackageCard key={pkg.id} pkg={pkg} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
