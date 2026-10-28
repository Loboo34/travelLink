'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Image from 'next/image';
import { format } from 'date-fns';
import {
  MapPin,
  Star,
  Calendar,
  Users,
  Luggage,
  Check,
  ChevronLeft,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { ReviewSection } from '@/components/reviews/review-section';
import { packagesApi } from '@/lib/api';
import { useBookingStore, useAuthStore } from '@/lib/stores';
import type { Package } from '@/lib/types';

export default function PackageDetailPage() {
  const params = useParams();
  const router = useRouter();
  const packageID = params.packageID as string;
  const { isAuthenticated } = useAuthStore();
  const { setPackageBooking } = useBookingStore();

  const [pkg, setPkg] = useState<Package | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchPackage() {
      try {
        const data = await packagesApi.getById(packageID);
        setPkg(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load package');
      } finally {
        setIsLoading(false);
      }
    }
    if (packageID) fetchPackage();
  }, [packageID]);

  const handleBook = () => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }
    if (pkg) {
      setPackageBooking(pkg);
      router.push('/checkout');
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="flex min-h-[60vh] items-center justify-center">
          <Spinner size="lg" />
        </div>
      </div>
    );
  }

  if (error || !pkg) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4">
          <p className="text-destructive">{error || 'Package not found'}</p>
          <Button onClick={() => router.back()}>Go Back</Button>
        </div>
      </div>
    );
  }

  const startFrom = new Date(pkg.startDateFrom);
  const startTo = new Date(pkg.startDateTo);

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Back Button */}
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ChevronLeft className="mr-2 h-4 w-4" />
          Back to Packages
        </Button>

        {/* Image Gallery */}
        <div className="mb-8 grid gap-4 md:grid-cols-2">
          <div className="relative aspect-[16/9] overflow-hidden rounded-xl">
            {pkg.images?.[0] ? (
              <Image
                src={pkg.images[0]}
                alt={pkg.name}
                fill
                className="object-cover"
              />
            ) : (
              <div className="flex h-full w-full items-center justify-center bg-muted">
                <Luggage className="h-16 w-16 text-muted-foreground" />
              </div>
            )}
          </div>
          <div className="grid grid-cols-2 gap-4">
            {pkg.images?.slice(1, 5).map((image, index) => (
              <div
                key={index}
                className="relative aspect-[4/3] overflow-hidden rounded-xl"
              >
                <Image
                  src={image}
                  alt={`${pkg.name} ${index + 2}`}
                  fill
                  className="object-cover"
                />
              </div>
            ))}
          </div>
        </div>

        <div className="grid gap-8 lg:grid-cols-3">
          {/* Main Content */}
          <div className="lg:col-span-2">
            {/* Header */}
            <div className="mb-6">
              <div className="mb-2 flex flex-wrap gap-2">
                {pkg.tags?.map((tag) => (
                  <Badge key={tag}>{tag}</Badge>
                ))}
              </div>
              <div className="flex items-start justify-between">
                <h1 className="text-3xl font-bold">{pkg.name}</h1>
                {pkg.rating > 0 && (
                  <div className="flex items-center gap-1">
                    <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                    <span className="text-lg font-semibold">
                      {pkg.rating.toFixed(1)}
                    </span>
                  </div>
                )}
              </div>
              <div className="mt-2 flex items-center gap-1 text-muted-foreground">
                <MapPin className="h-4 w-4" />
                <span>{pkg.destination}</span>
              </div>
            </div>

            <Separator className="my-6" />

            {/* Quick Info */}
            <div className="mb-6 flex flex-wrap gap-6">
              <div className="flex items-center gap-2">
                <Calendar className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Duration</p>
                  <p className="font-medium">{pkg.durationDays} days</p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Users className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Group Size</p>
                  <p className="font-medium">Up to {pkg.maxTravelers} people</p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Calendar className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Available</p>
                  <p className="font-medium">
                    {format(startFrom, 'MMM d')} -{' '}
                    {format(startTo, 'MMM d, yyyy')}
                  </p>
                </div>
              </div>
            </div>

            {/* Description */}
            <div className="mb-6">
              <h2 className="mb-3 text-xl font-semibold">About this package</h2>
              <p className="leading-relaxed text-muted-foreground">
                {pkg.description}
              </p>
            </div>

            <Separator className="my-6" />

            {/* Included Components */}
            {pkg.includedComponents?.length > 0 && (
              <div className="mb-6">
                <h2 className="mb-4 text-xl font-semibold">
                  {"What's Included"}
                </h2>
                <div className="grid gap-4 sm:grid-cols-2">
                  {pkg.includedComponents.map((component, index) => (
                    <Card key={index}>
                      <CardContent className="flex items-start gap-3 p-4">
                        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                          <Check className="h-4 w-4 text-primary" />
                        </div>
                        <div>
                          <p className="font-semibold capitalize">
                            {component.type}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            {component.name}
                          </p>
                          {component.description && (
                            <p className="mt-1 text-xs text-muted-foreground">
                              {component.description}
                            </p>
                          )}
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Booking Card */}
          <div className="lg:col-span-1">
            <Card className="sticky top-24">
              <CardHeader>
                <CardTitle>Book This Package</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="rounded-lg bg-muted/50 p-4">
                  <p className="text-sm text-muted-foreground">Price per person</p>
                  <p className="text-3xl font-bold text-primary">
                    {pkg.currency} {pkg.basePrice.toLocaleString()}
                  </p>
                </div>

                <div className="space-y-2 text-sm">
                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span>{pkg.durationDays} days trip</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Users className="h-4 w-4 text-muted-foreground" />
                    <span>Max {pkg.maxTravelers} travelers</span>
                  </div>
                </div>

                <Button className="w-full" size="lg" onClick={handleBook}>
                  Book Package
                </Button>
                <p className="text-center text-sm text-muted-foreground">
                  Secure your spot with a small deposit
                </p>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Reviews Section */}
        <div className="mt-8">
          <ReviewSection referenceID={packageID} reviewFor="package" />
        </div>
      </main>
    </div>
  );
}
