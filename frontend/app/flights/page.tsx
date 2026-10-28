'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Search, Plane } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { FlightCard } from '@/components/cards';
import { flightsApi } from '@/lib/api';
import type { Flight } from '@/lib/types';

export default function FlightsPage() {
  const [flights, setFlights] = useState<Flight[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchFlights() {
      try {
        const data = await flightsApi.getAll();
        setFlights(data || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load flights');
      } finally {
        setIsLoading(false);
      }
    }
    fetchFlights();
  }, []);

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Flights</h1>
            <p className="mt-1 text-muted-foreground">
              Discover and book flights to destinations worldwide
            </p>
          </div>
          <Button asChild>
            <Link href="/search/flights">
              <Search className="mr-2 h-4 w-4" />
              Search Flights
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
        ) : flights.length === 0 ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Plane className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-lg font-medium">No flights available</p>
            <p className="text-muted-foreground">
              Check back later for new flights
            </p>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {flights.map((flight) => (
              <FlightCard key={flight.id} flight={flight} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
