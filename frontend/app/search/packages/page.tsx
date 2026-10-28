'use client';

import { useState } from 'react';
import { Search, Luggage, SlidersHorizontal } from 'lucide-react';
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
import { PackageCard } from '@/components/cards';
import { useSearchStore } from '@/lib/stores';

const tags = [
  { value: 'all', label: 'All Tags' },
  { value: 'adventure', label: 'Adventure' },
  { value: 'romantic', label: 'Romantic' },
  { value: 'family', label: 'Family' },
  { value: 'luxury', label: 'Luxury' },
  { value: 'budget', label: 'Budget' },
  { value: 'safari', label: 'Safari' },
  { value: 'beach', label: 'Beach' },
];

const sortOptions = [
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
  { value: 'rating_desc', label: 'Rating: Highest First' },
  { value: 'duration_asc', label: 'Duration: Shortest' },
];

export default function PackageSearchPage() {
  const {
    packageFilters,
    packageResults,
    isSearchingPackages,
    setPackageFilters,
    searchPackages,
  } = useSearchStore();

  const [hasSearched, setHasSearched] = useState(false);

  const handleSearch = async () => {
    setHasSearched(true);
    await searchPackages();
  };

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">
            Search Packages
          </h1>
          <p className="mt-1 text-muted-foreground">
            Find curated travel packages for your dream vacation
          </p>
        </div>

        {/* Search Form */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Luggage className="h-5 w-5" />
              Package Search
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <div className="space-y-2">
                <Label htmlFor="destination">Destination</Label>
                <Input
                  id="destination"
                  placeholder="e.g., Maasai Mara"
                  value={packageFilters.destination || ''}
                  onChange={(e) =>
                    setPackageFilters({ destination: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="startDate">Start Date</Label>
                <Input
                  id="startDate"
                  type="date"
                  value={packageFilters.startDate || ''}
                  onChange={(e) =>
                    setPackageFilters({ startDate: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="endDate">End Date</Label>
                <Input
                  id="endDate"
                  type="date"
                  value={packageFilters.endDate || ''}
                  onChange={(e) =>
                    setPackageFilters({ endDate: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="tag">Package Type</Label>
                <Select
                  value={packageFilters.tag || 'all'}
                  onValueChange={(value) =>
                    setPackageFilters({
                      tag: value === 'all' ? undefined : value,
                    })
                  }
                >
                  <SelectTrigger id="tag">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {tags.map((tag) => (
                      <SelectItem key={tag.value} value={tag.value}>
                        {tag.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-5">
              <div className="space-y-2">
                <Label htmlFor="adults">Adults</Label>
                <Input
                  id="adults"
                  type="number"
                  min={1}
                  max={20}
                  value={packageFilters.adults || 1}
                  onChange={(e) =>
                    setPackageFilters({
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
                  value={packageFilters.children || 0}
                  onChange={(e) =>
                    setPackageFilters({
                      children: parseInt(e.target.value) || 0,
                    })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="budget">Max Budget (USD)</Label>
                <Input
                  id="budget"
                  type="number"
                  min={0}
                  placeholder="Any budget"
                  value={packageFilters.maxBudget || ''}
                  onChange={(e) =>
                    setPackageFilters({
                      maxBudget: e.target.value
                        ? parseInt(e.target.value)
                        : undefined,
                    })
                  }
                />
              </div>
              <div className="col-span-2 flex items-end">
                <Button
                  className="w-full"
                  onClick={handleSearch}
                  disabled={isSearchingPackages}
                >
                  {isSearchingPackages ? (
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
            {hasSearched && packageResults
              ? `${packageResults.total || packageResults.data?.length || 0} packages found`
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
                    value={packageFilters.sortBy || 'price_asc'}
                    onValueChange={(value) =>
                      setPackageFilters({ sortBy: value })
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
                  <Label>Include Component</Label>
                  <Select
                    value={packageFilters.component || 'all'}
                    onValueChange={(value) =>
                      setPackageFilters({
                        component: value === 'all' ? undefined : value,
                      })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All</SelectItem>
                      <SelectItem value="flight">With Flights</SelectItem>
                      <SelectItem value="accommodation">
                        With Accommodation
                      </SelectItem>
                      <SelectItem value="activity">With Activities</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>

        {isSearchingPackages ? (
          <div className="flex min-h-[40vh] items-center justify-center">
            <Spinner size="lg" />
          </div>
        ) : hasSearched && packageResults?.data?.length === 0 ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Luggage className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-lg font-medium">No packages found</p>
            <p className="text-muted-foreground">
              Try adjusting your search criteria
            </p>
          </div>
        ) : packageResults?.data?.length ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {packageResults.data.map((pkg) => (
              <PackageCard key={pkg.id} pkg={pkg} />
            ))}
          </div>
        ) : null}
      </main>
    </div>
  );
}
