'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Search, Building2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { AccommodationCard } from '@/components/cards';
import { accommodationsApi } from '@/lib/api';
import type { Accommodation } from '@/lib/types';

export default function AccommodationsPage() {
  const [accommodations, setAccommodations] = useState<Accommodation[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchAccommodations() {
      try {
        const data = await accommodationsApi.getAll();
        setAccommodations(data);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : 'Failed to load accommodations'
        );
      } finally {
        setIsLoading(false);
      }
    }
    fetchAccommodations();
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">
              Accommodations
            </h1>
            <p className="mt-1 text-muted-foreground">
              Find the perfect place to stay for your next trip
            </p>
          </div>
          <Button asChild>
            <Link href="/search/accommodations">
              <Search className="mr-2 h-4 w-4" />
              Search Accommodations
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
        ) : accommodations.length === 0 ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Building2 className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-lg font-medium">No accommodations available</p>
            <p className="text-muted-foreground">
              Check back later for new listings
            </p>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {accommodations.map((accommodation) => (
              <AccommodationCard
                key={accommodation.id}
                accommodation={accommodation}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
