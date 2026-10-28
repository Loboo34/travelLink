'use client';

import { useState } from 'react';
import { Search, Compass, SlidersHorizontal } from 'lucide-react';
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
import { ActivityCard } from '@/components/cards';
import { useSearchStore } from '@/lib/stores';

const categories = [
  { value: 'all', label: 'All Categories' },
  { value: 'adventure', label: 'Adventure' },
  { value: 'cultural', label: 'Cultural' },
  { value: 'nature', label: 'Nature & Wildlife' },
  { value: 'food', label: 'Food & Drink' },
  { value: 'water', label: 'Water Activities' },
  { value: 'tours', label: 'Tours & Sightseeing' },
];

const sortOptions = [
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
  { value: 'rating_desc', label: 'Rating: Highest First' },
  { value: 'duration_asc', label: 'Duration: Shortest' },
];

export default function ActivitySearchPage() {
  const {
    activityFilters,
    activityResults,
    isSearchingActivities,
    setActivityFilters,
    searchActivities,
  } = useSearchStore();

  const [hasSearched, setHasSearched] = useState(false);

  const handleSearch = async () => {
    setHasSearched(true);
    await searchActivities();
  };

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">
            Search Activities
          </h1>
          <p className="mt-1 text-muted-foreground">
            Discover exciting experiences and tours
          </p>
        </div>

        {/* Search Form */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Compass className="h-5 w-5" />
              Activity Search
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <div className="space-y-2">
                <Label htmlFor="city">City</Label>
                <Input
                  id="city"
                  placeholder="e.g., Nairobi"
                  value={activityFilters.city || ''}
                  onChange={(e) =>
                    setActivityFilters({ city: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="country">Country</Label>
                <Input
                  id="country"
                  placeholder="e.g., Kenya"
                  value={activityFilters.country || ''}
                  onChange={(e) =>
                    setActivityFilters({ country: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="date">Date</Label>
                <Input
                  id="date"
                  type="date"
                  value={activityFilters.date || ''}
                  onChange={(e) =>
                    setActivityFilters({ date: e.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="category">Category</Label>
                <Select
                  value={activityFilters.category || 'all'}
                  onValueChange={(value) =>
                    setActivityFilters({
                      category: value === 'all' ? undefined : value,
                    })
                  }
                >
                  <SelectTrigger id="category">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map((cat) => (
                      <SelectItem key={cat.value} value={cat.value}>
                        {cat.label}
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
                  value={activityFilters.adults || 1}
                  onChange={(e) =>
                    setActivityFilters({
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
                  value={activityFilters.children || 0}
                  onChange={(e) =>
                    setActivityFilters({
                      children: parseInt(e.target.value) || 0,
                    })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="duration">Max Duration (min)</Label>
                <Input
                  id="duration"
                  type="number"
                  min={30}
                  placeholder="Any duration"
                  value={activityFilters.maxDurationMinutes || ''}
                  onChange={(e) =>
                    setActivityFilters({
                      maxDurationMinutes: e.target.value
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
                  disabled={isSearchingActivities}
                >
                  {isSearchingActivities ? (
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
            {hasSearched && activityResults
              ? `${activityResults.total || activityResults.data?.length || 0} activities found`
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
                    value={activityFilters.sortBy || 'price_asc'}
                    onValueChange={(value) =>
                      setActivityFilters({ sortBy: value })
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
                    value={activityFilters.radiusKm || 10}
                    onChange={(e) =>
                      setActivityFilters({
                        radiusKm: parseInt(e.target.value) || 10,
                      })
                    }
                  />
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>

        {isSearchingActivities ? (
          <div className="flex min-h-[40vh] items-center justify-center">
            <Spinner size="lg" />
          </div>
        ) : hasSearched && activityResults?.data?.length === 0 ? (
          <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Compass className="h-8 w-8 text-muted-foreground" />
            </div>
            <p className="text-lg font-medium">No activities found</p>
            <p className="text-muted-foreground">
              Try adjusting your search criteria
            </p>
          </div>
        ) : activityResults?.data?.length ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {activityResults.data.map((activity) => (
              <ActivityCard key={activity.id} activity={activity} />
            ))}
          </div>
        ) : null}
      </main>
    </div>
  );
}
