'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Image from 'next/image';
import {
  MapPin,
  Star,
  Clock,
  Check,
  Users,
  Bed,
  Bath,
  ChevronLeft,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { ReviewSection } from '@/components/reviews/review-section';
import { accommodationsApi } from '@/lib/api';
import { useBookingStore, useAuthStore } from '@/lib/stores';
import type { Accommodation, RoomType } from '@/lib/types';

export default function AccommodationDetailPage() {
  const params = useParams();
  const router = useRouter();
  const accommodationID = params.accommodationID as string;
  const { isAuthenticated } = useAuthStore();
  const { setAccommodationBooking } = useBookingStore();

  const [accommodation, setAccommodation] = useState<Accommodation | null>(
    null
  );
  const [selectedRoom, setSelectedRoom] = useState<RoomType | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchAccommodation() {
      try {
        const data = await accommodationsApi.getById(accommodationID);
        setAccommodation(data);
        if (data.roomTypes?.length > 0) {
          setSelectedRoom(data.roomTypes[0]);
        }
      } catch (err) {
        setError(
          err instanceof Error ? err.message : 'Failed to load accommodation'
        );
      } finally {
        setIsLoading(false);
      }
    }
    if (accommodationID) fetchAccommodation();
  }, [accommodationID]);

  const handleBook = () => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }
    if (accommodation && selectedRoom) {
      setAccommodationBooking(accommodation, selectedRoom);
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

  if (error || !accommodation) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4">
          <p className="text-destructive">
            {error || 'Accommodation not found'}
          </p>
          <Button onClick={() => router.back()}>Go Back</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Back Button */}
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ChevronLeft className="mr-2 h-4 w-4" />
          Back to Accommodations
        </Button>

        {/* Image Gallery */}
        <div className="mb-8 grid gap-4 md:grid-cols-2">
          <div className="relative aspect-[4/3] overflow-hidden rounded-xl">
            {accommodation.images?.[0] ? (
              <Image
                src={accommodation.images[0]}
                alt={accommodation.name}
                fill
                className="object-cover"
              />
            ) : (
              <div className="flex h-full w-full items-center justify-center bg-muted">
                <Bed className="h-16 w-16 text-muted-foreground" />
              </div>
            )}
          </div>
          <div className="grid grid-cols-2 gap-4">
            {accommodation.images?.slice(1, 5).map((image, index) => (
              <div
                key={index}
                className="relative aspect-[4/3] overflow-hidden rounded-xl"
              >
                <Image
                  src={image}
                  alt={`${accommodation.name} ${index + 2}`}
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
              <div className="flex items-start justify-between">
                <div>
                  <Badge className="mb-2">{accommodation.type}</Badge>
                  <h1 className="text-3xl font-bold">{accommodation.name}</h1>
                </div>
                {accommodation.rating > 0 && (
                  <div className="flex items-center gap-1">
                    <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                    <span className="text-lg font-semibold">
                      {accommodation.rating.toFixed(1)}
                    </span>
                    <span className="text-muted-foreground">
                      ({accommodation.reviews} reviews)
                    </span>
                  </div>
                )}
              </div>
              <div className="mt-2 flex items-center gap-1 text-muted-foreground">
                <MapPin className="h-4 w-4" />
                <span>
                  {accommodation.address.street}, {accommodation.address.city},{' '}
                  {accommodation.address.country}
                </span>
              </div>
            </div>

            <Separator className="my-6" />

            {/* Description */}
            <div className="mb-6">
              <h2 className="mb-3 text-xl font-semibold">About this place</h2>
              <p className="leading-relaxed text-muted-foreground">
                {accommodation.description}
              </p>
            </div>

            {/* Check-in/out Times */}
            <div className="mb-6 flex gap-6">
              <div className="flex items-center gap-2">
                <Clock className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Check-in</p>
                  <p className="font-medium">{accommodation.checkInTime}</p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <Clock className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Check-out</p>
                  <p className="font-medium">{accommodation.checkOutTime}</p>
                </div>
              </div>
            </div>

            {/* Amenities */}
            {accommodation.amenities?.length > 0 && (
              <div className="mb-6">
                <h2 className="mb-3 text-xl font-semibold">Amenities</h2>
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
                  {accommodation.amenities.map((amenity) => (
                    <div
                      key={amenity}
                      className="flex items-center gap-2 text-sm"
                    >
                      <Check className="h-4 w-4 text-primary" />
                      <span>{amenity}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <Separator className="my-6" />

            {/* Room Types */}
            {accommodation.roomTypes?.length > 0 && (
              <div className="mb-6">
                <h2 className="mb-4 text-xl font-semibold">Room Types</h2>
                <div className="space-y-4">
                  {accommodation.roomTypes.map((room) => (
                    <Card
                      key={room.id}
                      className={`cursor-pointer transition-all ${
                        selectedRoom?.id === room.id
                          ? 'border-primary ring-2 ring-primary/20'
                          : 'hover:border-primary/50'
                      }`}
                      onClick={() => setSelectedRoom(room)}
                    >
                      <CardContent className="flex items-center gap-4 p-4">
                        {room.images?.[0] && (
                          <div className="relative h-20 w-28 overflow-hidden rounded-lg">
                            <Image
                              src={room.images[0]}
                              alt={room.name}
                              fill
                              className="object-cover"
                            />
                          </div>
                        )}
                        <div className="flex-1">
                          <h3 className="font-semibold">{room.name}</h3>
                          <p className="line-clamp-1 text-sm text-muted-foreground">
                            {room.description}
                          </p>
                          <div className="mt-2 flex items-center gap-4 text-sm text-muted-foreground">
                            <span className="flex items-center gap-1">
                              <Bed className="h-4 w-4" />
                              {room.bedrooms} bed
                            </span>
                            <span className="flex items-center gap-1">
                              <Bath className="h-4 w-4" />
                              {room.bathrooms} bath
                            </span>
                            <span className="flex items-center gap-1">
                              <Users className="h-4 w-4" />
                              {room.maxGuests} guests
                            </span>
                          </div>
                        </div>
                        <div className="text-right">
                          <p className="text-xl font-bold text-primary">
                            ${room.basePrice}
                          </p>
                          <p className="text-sm text-muted-foreground">
                            /night
                          </p>
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
                <CardTitle>Book Your Stay</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {selectedRoom && (
                  <div className="rounded-lg bg-muted/50 p-4">
                    <p className="text-sm text-muted-foreground">
                      Selected Room
                    </p>
                    <p className="font-semibold">{selectedRoom.name}</p>
                    <p className="mt-2 text-2xl font-bold text-primary">
                      ${selectedRoom.basePrice}
                      <span className="text-sm font-normal text-muted-foreground">
                        /night
                      </span>
                    </p>
                  </div>
                )}
                <Button
                  className="w-full"
                  size="lg"
                  onClick={handleBook}
                  disabled={!selectedRoom}
                >
                  Book Now
                </Button>
                <p className="text-center text-sm text-muted-foreground">
                  {"You won't be charged yet"}
                </p>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Reviews Section */}
        <div className="mt-8">
          <ReviewSection
            referenceID={accommodationID}
            reviewFor="accommodation"
          />
        </div>
      </main>
    </div>
  );
}
