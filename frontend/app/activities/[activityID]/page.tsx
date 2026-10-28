'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Image from 'next/image';
import { format } from 'date-fns';
import {
  MapPin,
  Star,
  Clock,
  Check,
  X,
  Compass,
  ChevronLeft,
  Calendar,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Spinner } from '@/components/ui/spinner';
import { Header } from '@/components/layout/header';
import { ReviewSection } from '@/components/reviews/review-section';
import { activitiesApi, timeSlotsApi } from '@/lib/api';
import { useBookingStore, useAuthStore } from '@/lib/stores';
import type { Activity, TimeSlot } from '@/lib/types';

export default function ActivityDetailPage() {
  const params = useParams();
  const router = useRouter();
  const activityID = params.activityID as string;
  const { isAuthenticated } = useAuthStore();
  const { setActivityBooking } = useBookingStore();

  const [activity, setActivity] = useState<Activity | null>(null);
  const [timeSlots, setTimeSlots] = useState<TimeSlot[]>([]);
  const [selectedSlot, setSelectedSlot] = useState<TimeSlot | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const [activityData, slotsData] = await Promise.all([
          activitiesApi.getById(activityID),
          timeSlotsApi.getAll(),
        ]);
        setActivity(activityData);
        const activitySlots = slotsData.filter(
          (slot) => slot.activityID === activityID
        );
        setTimeSlots(activitySlots);
        if (activitySlots.length > 0) {
          setSelectedSlot(activitySlots[0]);
        }
      } catch (err) {
        setError(
          err instanceof Error ? err.message : 'Failed to load activity'
        );
      } finally {
        setIsLoading(false);
      }
    }
    if (activityID) fetchData();
  }, [activityID]);

  const handleBook = () => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }
    if (activity && selectedSlot) {
      setActivityBooking(activity, selectedSlot);
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

  if (error || !activity) {
    return (
      <div className="min-h-screen bg-background">
        <Header />
        <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4">
          <p className="text-destructive">{error || 'Activity not found'}</p>
          <Button onClick={() => router.back()}>Go Back</Button>
        </div>
      </div>
    );
  }

  const hours = Math.floor(activity.durationMinutes / 60);
  const minutes = activity.durationMinutes % 60;
  const durationText =
    hours > 0
      ? `${hours} hour${hours > 1 ? 's' : ''}${minutes > 0 ? ` ${minutes} minutes` : ''}`
      : `${minutes} minutes`;

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        {/* Back Button */}
        <Button variant="ghost" onClick={() => router.back()} className="mb-6">
          <ChevronLeft className="mr-2 h-4 w-4" />
          Back to Activities
        </Button>

        {/* Image Gallery */}
        <div className="mb-8 grid gap-4 md:grid-cols-2">
          <div className="relative aspect-[4/3] overflow-hidden rounded-xl">
            {activity.images?.[0] ? (
              <Image
                src={activity.images[0]}
                alt={activity.title}
                fill
                className="object-cover"
              />
            ) : (
              <div className="flex h-full w-full items-center justify-center bg-muted">
                <Compass className="h-16 w-16 text-muted-foreground" />
              </div>
            )}
          </div>
          <div className="grid grid-cols-2 gap-4">
            {activity.images?.slice(1, 5).map((image, index) => (
              <div
                key={index}
                className="relative aspect-[4/3] overflow-hidden rounded-xl"
              >
                <Image
                  src={image}
                  alt={`${activity.title} ${index + 2}`}
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
                {activity.categories?.map((category) => (
                  <Badge key={category}>{category}</Badge>
                ))}
              </div>
              <div className="flex items-start justify-between">
                <h1 className="text-3xl font-bold">{activity.title}</h1>
                {activity.rating > 0 && (
                  <div className="flex items-center gap-1">
                    <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                    <span className="text-lg font-semibold">
                      {activity.rating.toFixed(1)}
                    </span>
                    <span className="text-muted-foreground">
                      ({activity.reviewCount} reviews)
                    </span>
                  </div>
                )}
              </div>
              <div className="mt-2 flex items-center gap-1 text-muted-foreground">
                <MapPin className="h-4 w-4" />
                <span>
                  {activity.city}, {activity.country}
                </span>
              </div>
            </div>

            <Separator className="my-6" />

            {/* Quick Info */}
            <div className="mb-6 flex flex-wrap gap-6">
              <div className="flex items-center gap-2">
                <Clock className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Duration</p>
                  <p className="font-medium">{durationText}</p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <MapPin className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Meeting Point</p>
                  <p className="font-medium">{activity.meetingPoint}</p>
                </div>
              </div>
            </div>

            {/* Description */}
            <div className="mb-6">
              <h2 className="mb-3 text-xl font-semibold">About this activity</h2>
              <p className="leading-relaxed text-muted-foreground">
                {activity.description}
              </p>
            </div>

            {/* Inclusions & Exclusions */}
            <div className="mb-6 grid gap-6 sm:grid-cols-2">
              {activity.inclusions?.length > 0 && (
                <div>
                  <h3 className="mb-3 text-lg font-semibold">
                    {"What's included"}
                  </h3>
                  <ul className="space-y-2">
                    {activity.inclusions.map((item, index) => (
                      <li key={index} className="flex items-center gap-2">
                        <Check className="h-4 w-4 text-green-500" />
                        <span className="text-sm">{item}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
              {activity.exclusions?.length > 0 && (
                <div>
                  <h3 className="mb-3 text-lg font-semibold">Not included</h3>
                  <ul className="space-y-2">
                    {activity.exclusions.map((item, index) => (
                      <li key={index} className="flex items-center gap-2">
                        <X className="h-4 w-4 text-destructive" />
                        <span className="text-sm">{item}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>

            <Separator className="my-6" />

            {/* Available Time Slots */}
            {timeSlots.length > 0 && (
              <div className="mb-6">
                <h2 className="mb-4 text-xl font-semibold">
                  Available Time Slots
                </h2>
                <div className="grid gap-3 sm:grid-cols-2">
                  {timeSlots.map((slot) => (
                    <Card
                      key={slot.id}
                      className={`cursor-pointer transition-all ${
                        selectedSlot?.id === slot.id
                          ? 'border-primary ring-2 ring-primary/20'
                          : 'hover:border-primary/50'
                      }`}
                      onClick={() => setSelectedSlot(slot)}
                    >
                      <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <Calendar className="h-4 w-4 text-muted-foreground" />
                            <span className="font-medium">
                              {format(new Date(slot.startTime), 'MMM d, yyyy')}
                            </span>
                          </div>
                          <Badge variant="outline">
                            {slot.availableSpots} spots
                          </Badge>
                        </div>
                        <div className="mt-2 flex items-center justify-between">
                          <span className="text-sm text-muted-foreground">
                            {format(new Date(slot.startTime), 'HH:mm')} -{' '}
                            {format(new Date(slot.endTime), 'HH:mm')}
                          </span>
                          <span className="font-semibold text-primary">
                            ${slot.price}
                          </span>
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
                <CardTitle>Book This Activity</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {selectedSlot && (
                  <div className="rounded-lg bg-muted/50 p-4">
                    <p className="text-sm text-muted-foreground">
                      Selected Time Slot
                    </p>
                    <p className="font-semibold">
                      {format(new Date(selectedSlot.startTime), 'MMM d, yyyy')}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {format(new Date(selectedSlot.startTime), 'HH:mm')} -{' '}
                      {format(new Date(selectedSlot.endTime), 'HH:mm')}
                    </p>
                    <p className="mt-2 text-2xl font-bold text-primary">
                      ${selectedSlot.price}
                      <span className="text-sm font-normal text-muted-foreground">
                        /person
                      </span>
                    </p>
                  </div>
                )}
                <Button
                  className="w-full"
                  size="lg"
                  onClick={handleBook}
                  disabled={!selectedSlot}
                >
                  Book Now
                </Button>
                <p className="text-center text-sm text-muted-foreground">
                  Free cancellation up to 24 hours before
                </p>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Reviews Section */}
        <div className="mt-8">
          <ReviewSection referenceID={activityID} reviewFor="activity" />
        </div>
      </main>
    </div>
  );
}
