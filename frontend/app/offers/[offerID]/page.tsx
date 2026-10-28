'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { format } from 'date-fns';
import {
  Plane,
  Luggage,
  Clock,
  Tag,
  ChevronLeft,
  AlertCircle,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Spinner } from '@/components/ui/spinner';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Header } from '@/components/layout/header';
import { offersApi, flightsApi } from '@/lib/api';
import { useBookingStore, useAuthStore } from '@/lib/stores';
import type { Offer, Flight } from '@/lib/types';

export default function OfferDetailPage() {
  const params = useParams();
  const router = useRouter();
  const offerID = params.offerID as string;
  const { isAuthenticated } = useAuthStore();
  const { setFlightBooking } = useBookingStore();

  const [offer, setOffer] = useState<Offer | null>(null);
  const [flight, setFlight] = useState<Flight | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const offerData = await offersApi.getById(offerID);
        setOffer(offerData);
        if (offerData.flightID) {
          const flightData = await flightsApi.getById(offerData.flightID);
          setFlight(flightData);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load offer');
      } finally {
        setIsLoading(false);
      }
    }
    if (offerID) fetchData();
  }, [offerID]);

  const handleBook = () => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }
    if (flight && offer) {
      setFlightBooking(flight, offer);
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

  if (error || !offer) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4">
          <p className="text-destructive">{error || 'Booking option not found'}</p>
          <Button onClick={() => router.back()}>Go Back</Button>
        </div>
      </div>
    );
  }

  const expiresAt = offer.expiresAt ? new Date(offer.expiresAt) : null;
  const isExpired = expiresAt !== null && expiresAt < new Date();
  const isExpiringSoon =
    expiresAt !== null &&
    !isExpired &&
    expiresAt.getTime() - Date.now() < 24 * 60 * 60 * 1000;
  const baggage = offer.baggageAllowance?.[0];
  const isOneWay = offer.oneway ?? offer.oneWay ?? true;
  const cabinLabel =
    flight && flight.cabinClass.length > 0
      ? flight.cabinClass.map((cabin) => cabin.classType).join(', ')
      : 'Not set';

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-4xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Back Button */}
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ChevronLeft className="mr-2 h-4 w-4" />
          Back to Booking Options
        </Button>

        {/* Expiry Warning */}
        {isExpired && (
          <Alert variant="destructive" className="mb-6">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              This booking option has expired and is no longer available.
            </AlertDescription>
          </Alert>
        )}
        {isExpiringSoon && !isExpired && (
          <Alert className="mb-6 border-yellow-500 bg-yellow-50 text-yellow-700">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              This booking option expires soon. Book now to secure this fare.
            </AlertDescription>
          </Alert>
        )}

        {/* Booking Option Header */}
        <Card className="mb-6">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
                  <Tag className="h-6 w-6 text-primary" />
                </div>
                <div>
                  <CardTitle className="text-2xl">{offer.provider}</CardTitle>
                  <p className="text-muted-foreground">
                    Ref: {offer.providerReference}
                  </p>
                </div>
              </div>
              <Badge
                variant={offer.isActive && !isExpired ? 'default' : 'secondary'}
                className="text-sm"
              >
                {isExpired ? 'Expired' : offer.isActive ? 'Active' : 'Inactive'}
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            {/* Price Section */}
            <div className="mb-6 rounded-xl bg-primary/5 p-6 text-center">
              <p className="text-sm text-muted-foreground">Total Price</p>
              <p className="text-5xl font-bold text-primary">
                ${offer.priceTotal.toLocaleString()}
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                {isOneWay ? 'One Way' : 'Round Trip'}
              </p>
            </div>

            {/* Booking Option Details */}
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  <Plane className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Trip Type</p>
                  <p className="font-medium">
                    {isOneWay ? 'One Way' : 'Round Trip'}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  <Luggage className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">
                    Baggage Allowance
                  </p>
                  <p className="font-medium">
                    {baggage?.pieces || 0} piece(s), {baggage?.weightKg || 0}kg
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  <Clock className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Expires</p>
                  <p
                    className={`font-medium ${isExpired ? 'text-destructive' : isExpiringSoon ? 'text-yellow-600' : ''}`}
                  >
                    {expiresAt ? format(expiresAt, 'MMM d, yyyy HH:mm') : 'Not set'}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-muted">
                  <Tag className="h-5 w-5 text-muted-foreground" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">
                    Available Seats
                  </p>
                  <p className="font-medium">{offer.bookableSeats} seats</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Flight Details */}
        {flight && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>Flight Details</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between py-4">
                <div className="text-center">
                  <p className="text-3xl font-bold">
                    {format(new Date(flight.departureTime), 'HH:mm')}
                  </p>
                  <p className="mt-1 font-medium">{flight.originID}</p>
                  <p className="text-sm text-muted-foreground">
                    {format(new Date(flight.departureTime), 'MMM d, yyyy')}
                  </p>
                </div>
                <div className="flex flex-1 flex-col items-center px-4">
                  <p className="text-sm text-muted-foreground">
                    {flight.flightNumber}
                  </p>
                  <div className="relative mt-2 w-full">
                    <div className="h-0.5 w-full bg-border" />
                    <Plane className="absolute left-1/2 top-1/2 h-4 w-4 -translate-x-1/2 -translate-y-1/2 text-primary" />
                  </div>
                  <p className="mt-2 text-sm text-muted-foreground">
                    {flight.stops === 0
                      ? 'Direct'
                      : `${flight.stops} stop${flight.stops > 1 ? 's' : ''}`}
                  </p>
                </div>
                <div className="text-center">
                  <p className="text-3xl font-bold">
                    {format(new Date(flight.arrivalTime), 'HH:mm')}
                  </p>
                  <p className="mt-1 font-medium">{flight.destinationID}</p>
                  <p className="text-sm text-muted-foreground">
                    {format(new Date(flight.arrivalTime), 'MMM d, yyyy')}
                  </p>
                </div>
              </div>

              <Separator className="my-4" />

              <div className="flex flex-wrap gap-4 text-sm">
                <Badge variant="outline" className="capitalize">
                  {cabinLabel}
                </Badge>
                <Badge variant="outline">{flight.airlineID}</Badge>
                <Badge
                  variant={
                    flight.status === 'scheduled' ? 'default' : 'secondary'
                  }
                >
                  {flight.status}
                </Badge>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Book Button */}
        <Card>
          <CardContent className="flex items-center justify-between py-6">
            <div>
              <p className="text-2xl font-bold text-primary">
                ${offer.priceTotal.toLocaleString()}
              </p>
              <p className="text-sm text-muted-foreground">
                {offer.bookableSeats} seats remaining
              </p>
            </div>
            <Button
              size="lg"
              onClick={handleBook}
              disabled={isExpired || !offer.isActive}
            >
              {isExpired ? 'Option Expired' : 'Book This Fare'}
            </Button>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
