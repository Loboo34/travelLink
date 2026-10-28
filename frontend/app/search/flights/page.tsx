'use client';

import { useState } from 'react';
import { Search, Plane, SlidersHorizontal } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { FlightCard } from '@/components/cards';
import { useSearchStore } from '@/lib/stores';

const cabinClasses = [
  { value: 'economy', label: 'Economy' },
  { value: 'premium_economy', label: 'Premium Economy' },
  { value: 'business', label: 'Business' },
  { value: 'first', label: 'First Class' },
];

const sortOptions = [
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
  { value: 'duration_asc', label: 'Duration: Shortest' },
  { value: 'departure_asc', label: 'Departure: Earliest' },
];

export default function FlightSearchPage() {
  const {
    flightFilters,
    flightResults,
    isSearchingFlights,
    setFlightFilters,
    searchFlights,
  } = useSearchStore();

  const [hasSearched, setHasSearched] = useState(false);

  const handleSearch = async () => {
    setHasSearched(true);
    await searchFlights();
  };

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Search Flights</h1>
          <p className="mt-1 text-muted-foreground">
            Find the best flights for your journey
          </p>
        </div>

        {/* Search Form */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Plane className="h-5 w-5" />
              Flight Search
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <div className="space-y-2">
                <Label htmlFor="origin">From</Label>
                <Input
                  id="origin"
                  placeholder="Origin airport code (e.g., NBO)"
                  value={flightFilters.originCode || ''}
                  onChange={(e) =>
                    setFlightFilters({ originCode: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="destination">To</Label>
                <Input
                  id="destination"
                  placeholder="Destination code (e.g., LHR)"
                  value={flightFilters.destinationCode || ''}
                  onChange={(e) =>
                    setFlightFilters({ destinationCode: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="departure">Departure Date</Label>
                <Input
                  id="departure"
                  type="date"
                  value={flightFilters.departureTime || ''}
                  onChange={(e) =>
                    setFlightFilters({ departureTime: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="return">Return Date</Label>
                <Input
                  id="return"
                  type="date"
                  value={flightFilters.returnDate || ''}
                  onChange={(e) =>
                    setFlightFilters({ returnDate: e.target.value })
                  }
                />
              </div>
            </div>

            <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-5">
              <div className="space-y-2">
                <Label htmlFor="adults">Adults</Label>
                <Input
                  id="adults"
                  type="number"
                  min={1}
                  max={9}
                  value={flightFilters.adults || 1}
                  onChange={(e) =>
                    setFlightFilters({ adults: parseInt(e.target.value) || 1 })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="children">Children</Label>
                <Input
                  id="children"
                  type="number"
                  min={0}
                  max={9}
                  value={flightFilters.children || 0}
                  onChange={(e) =>
                    setFlightFilters({
                      children: parseInt(e.target.value) || 0,
                    })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="infants">Infants</Label>
                <Input
                  id="infants"
                  type="number"
                  min={0}
                  max={9}
                  value={flightFilters.infants || 0}
                  onChange={(e) =>
                    setFlightFilters({ infants: parseInt(e.target.value) || 0 })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cabin">Cabin Class</Label>
                <Select
                  value={flightFilters.cabinClass || 'economy'}
                  onValueChange={(value) =>
                    setFlightFilters({ cabinClass: value })
                  }
                >
                  <SelectTrigger id="cabin">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {cabinClasses.map((cabin) => (
                      <SelectItem key={cabin.value} value={cabin.value}>
                        {cabin.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="flex items-end">
                <Button
                  className="w-full"
                  onClick={handleSearch}
                  disabled={isSearchingFlights}
                >
                  {isSearchingFlights ? (
                    <>
                      <Spinner size="sm" className="mr-2" />
                      Searching...
                    </>
                  ) : (
                    <>
                      <Search className="mr-2 h-4 w-4" />
                      Search
                    </>
                  )}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Results */}
        <div className="flex items-center justify-between mb-4">
          <p className="text-muted-foreground">
            {hasSearched && flightResults
              ? `${flightResults.total || flightResults.data?.length || 0} flights found`
              : 'Enter your search criteria above'}
          </p>
          <Sheet>
            <SheetTrigger asChild>
              <Button variant="outline" size="sm">
                <SlidersHorizontal className="mr-2 h-4 w-4" />
                Filters
              </Button>
            </SheetTrigger>
            <SheetContent>
              <SheetHeader>
                <SheetTitle>Filter & Sort</SheetTitle>
              </SheetHeader>
              <div className="mt-6 space-y-4">
                <div className="space-y-2">
                  <Label>Sort By</Label>
                  <Select
                    value={flightFilters.sortBy || 'price_asc'}
                    onValueChange={(value) =>
                      setFlightFilters({ sortBy: value })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {sortOptions.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>

        {isSearchingFlights ? (
          <div className="flex min-h-[40vh] items-center justify-center">
            <Spinner size="lg" />
          </div>
        ) : hasSearched && flightResults?.data?.length === 0 ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Plane className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-lg font-medium">No flights found</p>
            <p className="text-muted-foreground">
              Try adjusting your search criteria
            </p>
          </div>
        ) : flightResults?.data?.length ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {flightResults.data.map((flight) => (
              <FlightCard key={flight.id} flight={flight} />
            ))}
          </div>
        ) : null}
      </main>
    </div>
  );
}
