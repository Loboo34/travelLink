'use client';

import { useState } from 'react';
import { Search, Building2, SlidersHorizontal } from 'lucide-react';
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
import { AccommodationCard } from '@/components/cards';
import { useSearchStore } from '@/lib/stores';

const propertyTypes = [
  { value: 'all', label: 'All Types' },
  { value: 'hotel', label: 'Hotel' },
  { value: 'resort', label: 'Resort' },
  { value: 'apartment', label: 'Apartment' },
  { value: 'villa', label: 'Villa' },
  { value: 'hostel', label: 'Hostel' },
  { value: 'guesthouse', label: 'Guesthouse' },
];

const sortOptions = [
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
  { value: 'rating_desc', label: 'Rating: Highest First' },
  { value: 'reviews_desc', label: 'Most Reviewed' },
];

export default function AccommodationSearchPage() {
  const {
    accommodationFilters,
    accommodationResults,
    isSearchingAccommodations,
    setAccommodationFilters,
    searchAccommodations,
  } = useSearchStore();

  const [hasSearched, setHasSearched] = useState(false);

  const handleSearch = async () => {
    setHasSearched(true);
    await searchAccommodations();
  };

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">
            Search Accommodations
          </h1>
          <p className="mt-1 text-muted-foreground">
            Find the perfect place to stay
          </p>
        </div>

        {/* Search Form */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Building2 className="h-5 w-5" />
              Accommodation Search
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <div className="space-y-2">
                <Label htmlFor="city">City</Label>
                <Input
                  id="city"
                  placeholder="e.g., Nairobi"
                  value={accommodationFilters.city || ''}
                  onChange={(e) =>
                    setAccommodationFilters({ city: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="country">Country</Label>
                <Input
                  id="country"
                  placeholder="e.g., Kenya"
                  value={accommodationFilters.country || ''}
                  onChange={(e) =>
                    setAccommodationFilters({ country: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="checkIn">Check-in Date</Label>
                <Input
                  id="checkIn"
                  type="date"
                  value={accommodationFilters.checkInDate || ''}
                  onChange={(e) =>
                    setAccommodationFilters({ checkInDate: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="checkOut">Check-out Date</Label>
                <Input
                  id="checkOut"
                  type="date"
                  value={accommodationFilters.checkOutDate || ''}
                  onChange={(e) =>
                    setAccommodationFilters({ checkOutDate: e.target.value })
                  }
                />
              </div>
            </div>

            <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-6">
              <div className="space-y-2">
                <Label htmlFor="adults">Adults</Label>
                <Input
                  id="adults"
                  type="number"
                  min={1}
                  max={20}
                  value={accommodationFilters.adults || 1}
                  onChange={(e) =>
                    setAccommodationFilters({
                      adults: parseInt(e.target.value) || 1,
                    })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="children">Children</Label>
                <Input
                  id="children"
                  type="number"
                  min={0}
                  max={10}
                  value={accommodationFilters.children || 0}
                  onChange={(e) =>
                    setAccommodationFilters({
                      children: parseInt(e.target.value) || 0,
                    })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="rooms">Rooms</Label>
                <Input
                  id="rooms"
                  type="number"
                  min={1}
                  max={10}
                  value={accommodationFilters.totalRooms || 1}
                  onChange={(e) =>
                    setAccommodationFilters({
                      totalRooms: parseInt(e.target.value) || 1,
                    })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="propertyType">Property Type</Label>
                <Select
                  value={accommodationFilters.propertyType || 'all'}
                  onValueChange={(value) =>
                    setAccommodationFilters({
                      propertyType: value === 'all' ? undefined : value,
                    })
                  }
                >
                  <SelectTrigger id="propertyType">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {propertyTypes.map((type) => (
                      <SelectItem key={type.value} value={type.value}>
                        {type.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="col-span-2 flex items-end">
                <Button
                  className="w-full"
                  onClick={handleSearch}
                  disabled={isSearchingAccommodations}
                >
                  {isSearchingAccommodations ? (
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
            {hasSearched && accommodationResults
              ? `${accommodationResults.total || accommodationResults.data?.length || 0} accommodations found`
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
                    value={accommodationFilters.sortBy || 'price_asc'}
                    onValueChange={(value) =>
                      setAccommodationFilters({ sortBy: value })
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
                <div className="space-y-2">
                  <Label>Search Radius (km)</Label>
                  <Input
                    type="number"
                    min={1}
                    max={100}
                    value={accommodationFilters.radiusKm || 10}
                    onChange={(e) =>
                      setAccommodationFilters({
                        radiusKm: parseInt(e.target.value) || 10,
                      })
                    }
                  />
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>

        {isSearchingAccommodations ? (
          <div className="flex min-h-[40vh] items-center justify-center">
            <Spinner size="lg" />
          </div>
        ) : hasSearched && accommodationResults?.data?.length === 0 ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Building2 className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-lg font-medium">No accommodations found</p>
            <p className="text-muted-foreground">
              Try adjusting your search criteria
            </p>
          </div>
        ) : accommodationResults?.data?.length ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {accommodationResults.data.map((accommodation) => (
              <AccommodationCard
                key={accommodation.id}
                accommodation={accommodation}
              />
            ))}
          </div>
        ) : null}
      </main>
    </div>
  );
}
