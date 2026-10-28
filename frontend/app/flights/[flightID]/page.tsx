'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { format } from 'date-fns';
import {
  Plane,
  Clock,
  ArrowRight,
  Calendar,
  Users,
  ChevronLeft,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { ReviewSection } from '@/components/reviews/review-section';
import { flightsApi } from '@/lib/api';
import { useBookingStore, useAuthStore } from '@/lib/stores';
import type { Flight, FlightSegment } from '@/lib/types';

function isFlightSegment(segment: Flight['segments'][number]): segment is FlightSegment {
  return typeof segment !== 'string';
}

export default function FlightDetailPage() {
  const params = useParams();
  const router = useRouter();
  const flightID = params.flightID as string;
  const { isAuthenticated } = useAuthStore();
  const { setFlightBooking } = useBookingStore();

  const [flight, setFlight] = useState<Flight | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchFlight() {
      try {
        const data = await flightsApi.getById(flightID);
        setFlight(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load flight');
      } finally {
        setIsLoading(false);
      }
    }
    if (flightID) fetchFlight();
  }, [flightID]);

  const handleBook = () => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }
    if (flight) {
      setFlightBooking(flight);
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

  if (error || !flight) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4">
          <p className="text-destructive">{error || 'Flight not found'}</p>
          <Button onClick={() => router.back()}>Go Back</Button>
        </div>
      </div>
    );
  }

  const departureTime = new Date(flight.departureTime);
  const arrivalTime = new Date(flight.arrivalTime);
  const duration = Math.round(
    (arrivalTime.getTime() - departureTime.getTime()) / (1000 * 60)
  );
  const hours = Math.floor(duration / 60);
  const minutes = duration % 60;
  const cabinLabel =
    flight.cabinClass.length > 0
      ? flight.cabinClass.map((cabin) => cabin.classType).join(', ')
      : 'Not set';
  const segmentDetails = flight.segments?.filter(isFlightSegment) || [];

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-4xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Back Button */}
        <Button
          variant="ghost"
          onClick={() => router.back()}
          className="mb-6"
        >
          <ChevronLeft className="mr-2 h-4 w-4" />
          Back to Flights
        </Button>

        {/* Flight Header */}
        <Card className="mb-6">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  <Plane className="h-6 w-6 text-primary" />
                </div>
                <div>
                  <CardTitle className="text-2xl">{flight.flightNumber}</CardTitle>
                  <p className="text-muted-foreground">{flight.airlineID}</p>
                </div>
              </div>
              <Badge
                variant={flight.status === 'scheduled' ? 'default' : 'secondary'}
                className="text-sm"
              >
                {flight.status}
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            {/* Flight Route */}
            <div className="flex items-center justify-between py-6">
              <div className="text-center">
                <p className="text-4xl font-bold">
                  {format(departureTime, 'HH:mm')}
                </p>
                <p className="mt-1 text-lg font-medium">{flight.originID}</p>
                <p className="text-sm text-muted-foreground">
                  {format(departureTime, 'EEE, MMM d, yyyy')}
                </p>
              </div>
              <div className="flex flex-1 flex-col items-center px-8">
                <div className="flex items-center gap-2 text-muted-foreground">
                  <Clock className="h-4 w-4" />
                  <span>{hours}h {minutes}m</span>
                </div>
                <div className="relative mt-3 w-full">
                  <div className="h-0.5 w-full bg-border" />
                  <ArrowRight className="absolute -right-1 top-1/2 h-5 w-5 -translate-y-1/2 text-primary" />
                  {flight.stops > 0 && (
                    <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
                      <div className="h-3 w-3 rounded-full border-2 border-primary bg-background" />
                    </div>
                  )}
                </div>
                <p className="mt-3 text-sm text-muted-foreground">
                  {flight.stops === 0
                    ? 'Direct Flight'
                    : `${flight.stops} stop${flight.stops > 1 ? 's' : ''}`}
                </p>
              </div>
              <div className="text-center">
                <p className="text-4xl font-bold">
                  {format(arrivalTime, 'HH:mm')}
                </p>
                <p className="mt-1 text-lg font-medium">{flight.destinationID}</p>
                <p className="text-sm text-muted-foreground">
                  {format(arrivalTime, 'EEE, MMM d, yyyy')}
                </p>
              </div>
            </div>

            <Separator className="my-6" />

            {/* Flight Details */}
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  <Calendar className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Cabin Class</p>
                  <p className="font-medium capitalize">{cabinLabel}</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  <Plane className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Aircraft</p>
                  <p className="font-medium">{flight.planeID || 'N/A'}</p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  <Users className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Flight Number</p>
                  <p className="font-medium">{flight.flightNumber}</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Segments */}
        {segmentDetails.length > 0 && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Flight Segments</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {segmentDetails.map((segment, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between rounded-lg border p-4"
                  >
                    <div>
                      <p className="font-medium">
                        {segment.departureAirport || segment.originAirportID} to{' '}
                        {segment.arrivalAirport || segment.destinationAirportID}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {segment.airline} {segment.flightNumber}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">
                        {format(new Date(segment.departureTime), 'HH:mm')} -{' '}
                        {format(new Date(segment.arrivalTime), 'HH:mm')}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}

        {/* Book Button */}
        <Card className="mb-6">
          <CardContent className="flex items-center justify-between py-6">
            <div>
              <p className="text-sm text-muted-foreground">Ready to book?</p>
              <p className="text-lg font-medium">
                Secure your seat on this flight
              </p>
            </div>
            <Button size="lg" onClick={handleBook}>
              Book Now
            </Button>
          </CardContent>
        </Card>

        {/* Reviews Section */}
        <ReviewSection referenceID={flightID} reviewFor="flight" />
      </main>
    </div>
  );
}
