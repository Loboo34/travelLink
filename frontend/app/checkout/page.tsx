'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { format } from 'date-fns';
import {
  Plane,
  Building2,
  Compass,
  Luggage,
  CreditCard,
  Smartphone,
  Building,
  Plus,
  Trash2,
  Check,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Spinner } from '@/components/ui/spinner';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { Header } from '@/components/layout/header';
import { useBookingStore, useAuthStore } from '@/lib/stores';
import type { PassengerDetails, PaymentMethod } from '@/lib/types';

const paymentMethods: { value: PaymentMethod; label: string; icon: React.ReactNode }[] = [
  { value: 'card', label: 'Credit/Debit Card', icon: <CreditCard className="h-5 w-5" /> },
  { value: 'mpesa', label: 'M-Pesa', icon: <Smartphone className="h-5 w-5" /> },
  { value: 'bank', label: 'Bank Transfer', icon: <Building className="h-5 w-5" /> },
];

export default function CheckoutPage() {
  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-background">
        <Header />
        <CheckoutContent />
      </div>
    </ProtectedRoute>
  );
}

function CheckoutContent() {
  const router = useRouter();
  const { user } = useAuthStore();
  const {
    activeBookingType,
    selectedFlight,
    selectedOffer,
    selectedAccommodation,
    selectedRoomType,
    selectedActivity,
    selectedTimeSlot,
    selectedPackage,
    paymentMethod,
    setPaymentMethod,
    setPassengerDetails,
    setGuestDetails,
    setParticipantDetails,
    setTravelersDetails,
    setCheckInOut,
    submitBooking,
    isBooking,
    error,
    bookingResult,
    clearBooking,
  } = useBookingStore();

  const [travelers, setTravelers] = useState<PassengerDetails[]>([
    {
      firstName: user?.firstName || '',
      lastName: user?.lastName || '',
      email: user?.email || '',
      phone: user?.phoneNumber || '',
      dateOfBirth: '',
      passportNumber: '',
      nationality: user?.nationality || '',
    },
  ]);

  const [checkIn, setCheckIn] = useState('');
  const [checkOut, setCheckOut] = useState('');

  const addTraveler = () => {
    setTravelers([
      ...travelers,
      {
        firstName: '',
        lastName: '',
        email: '',
        phone: '',
        dateOfBirth: '',
        passportNumber: '',
        nationality: '',
      },
    ]);
  };

  const removeTraveler = (index: number) => {
    if (travelers.length > 1) {
      setTravelers(travelers.filter((_, i) => i !== index));
    }
  };

  const updateTraveler = (
    index: number,
    field: keyof PassengerDetails,
    value: string
  ) => {
    const updated = [...travelers];
    updated[index] = { ...updated[index], [field]: value };
    setTravelers(updated);
  };

  const handleSubmit = async () => {
    // Set traveler details based on booking type
    if (activeBookingType === 'flight') {
      setPassengerDetails(travelers);
    } else if (activeBookingType === 'accommodation') {
      setGuestDetails(
        travelers.map((t) => ({
          firstName: t.firstName,
          lastName: t.lastName,
          email: t.email,
          phone: t.phone,
        }))
      );
      setCheckInOut(checkIn, checkOut);
    } else if (activeBookingType === 'activity') {
      setParticipantDetails(
        travelers.map((t) => ({
          firstName: t.firstName,
          lastName: t.lastName,
          email: t.email,
          phone: t.phone,
        }))
      );
    } else if (activeBookingType === 'package') {
      setTravelersDetails(travelers);
    }

    try {
      await submitBooking();
    } catch {
      // Error handled by store
    }
  };

  // Redirect if no booking selected
  if (!activeBookingType) {
    return (
      <main className="mx-auto max-w-4xl px-4 py-8 sm:px-6 lg:px-8">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <p className="text-lg font-medium">No booking selected</p>
            <p className="mt-1 text-muted-foreground">
              Please select something to book first
            </p>
            <Button className="mt-4" onClick={() => router.push('/')}>
              Browse Options
            </Button>
          </CardContent>
        </Card>
      </main>
    );
  }

  // Show success message
  if (bookingResult) {
    return (
      <main className="mx-auto max-w-4xl px-4 py-8 sm:px-6 lg:px-8">
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
              <Check className="h-8 w-8 text-green-600" />
            </div>
            <h2 className="mt-4 text-2xl font-bold">Booking Confirmed!</h2>
            <p className="mt-2 text-muted-foreground">
              Your booking has been successfully submitted.
            </p>
            <p className="mt-1 text-sm text-muted-foreground">
              Booking ID: {bookingResult.id}
            </p>
            <div className="mt-6 flex gap-4">
              <Button
                onClick={() => {
                  clearBooking();
                  router.push('/');
                }}
              >
                Back to Home
              </Button>
            </div>
          </CardContent>
        </Card>
      </main>
    );
  }

  const getBookingIcon = () => {
    switch (activeBookingType) {
      case 'flight':
        return <Plane className="h-6 w-6" />;
      case 'accommodation':
        return <Building2 className="h-6 w-6" />;
      case 'activity':
        return <Compass className="h-6 w-6" />;
      case 'package':
        return <Luggage className="h-6 w-6" />;
    }
  };

  const getBookingTitle = () => {
    switch (activeBookingType) {
      case 'flight':
        return selectedFlight?.flightNumber || 'Flight Booking';
      case 'accommodation':
        return selectedAccommodation?.name || 'Accommodation Booking';
      case 'activity':
        return selectedActivity?.title || 'Activity Booking';
      case 'package':
        return selectedPackage?.name || 'Package Booking';
    }
  };

  const getBookingDetails = () => {
    switch (activeBookingType) {
      case 'flight':
        return selectedFlight ? (
          <div className="space-y-2">
            <p>
              {selectedFlight.originID} → {selectedFlight.destinationID}
            </p>
            <p className="text-sm text-muted-foreground">
              {format(new Date(selectedFlight.departureTime), 'MMM d, yyyy HH:mm')}
            </p>
            {selectedOffer && (
              <p className="text-xl font-bold text-primary">
                ${selectedOffer.priceTotal.toLocaleString()}
              </p>
            )}
          </div>
        ) : null;
      case 'accommodation':
        return selectedAccommodation ? (
          <div className="space-y-2">
            <p>{selectedAccommodation.address.city}, {selectedAccommodation.address.country}</p>
            {selectedRoomType && (
              <>
                <p className="text-sm text-muted-foreground">
                  {selectedRoomType.name}
                </p>
                <p className="text-xl font-bold text-primary">
                  ${selectedRoomType.basePrice}/night
                </p>
              </>
            )}
          </div>
        ) : null;
      case 'activity':
        return selectedActivity ? (
          <div className="space-y-2">
            <p>{selectedActivity.city}, {selectedActivity.country}</p>
            {selectedTimeSlot && (
              <>
                <p className="text-sm text-muted-foreground">
                  {format(new Date(selectedTimeSlot.startTime), 'MMM d, yyyy HH:mm')}
                </p>
                <p className="text-xl font-bold text-primary">
                  ${selectedTimeSlot.price}/person
                </p>
              </>
            )}
          </div>
        ) : null;
      case 'package':
        return selectedPackage ? (
          <div className="space-y-2">
            <p>{selectedPackage.destination}</p>
            <p className="text-sm text-muted-foreground">
              {selectedPackage.durationDays} days
            </p>
            <p className="text-xl font-bold text-primary">
              {selectedPackage.currency} {selectedPackage.basePrice.toLocaleString()}
            </p>
          </div>
        ) : null;
    }
  };

  return (
    <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
      <h1 className="mb-8 text-3xl font-bold">Complete Your Booking</h1>

      {error && (
        <Alert variant="destructive" className="mb-6">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="grid gap-8 lg:grid-cols-3">
        {/* Booking Form */}
        <div className="lg:col-span-2 space-y-6">
          {/* Traveler Details */}
          <Card>
            <CardHeader>
              <CardTitle>Traveler Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {travelers.map((traveler, index) => (
                <div key={index} className="space-y-4">
                  {index > 0 && <Separator />}
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium">Traveler {index + 1}</h4>
                    {index > 0 && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => removeTraveler(index)}
                      >
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    )}
                  </div>
                  <div className="grid gap-4 sm:grid-cols-2">
                    <div className="space-y-2">
                      <Label>First Name</Label>
                      <Input
                        value={traveler.firstName}
                        onChange={(e) =>
                          updateTraveler(index, 'firstName', e.target.value)
                        }
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Last Name</Label>
                      <Input
                        value={traveler.lastName}
                        onChange={(e) =>
                          updateTraveler(index, 'lastName', e.target.value)
                        }
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Email</Label>
                      <Input
                        type="email"
                        value={traveler.email}
                        onChange={(e) =>
                          updateTraveler(index, 'email', e.target.value)
                        }
                        required
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Phone</Label>
                      <Input
                        type="tel"
                        value={traveler.phone}
                        onChange={(e) =>
                          updateTraveler(index, 'phone', e.target.value)
                        }
                        required
                      />
                    </div>
                    {activeBookingType === 'flight' && (
                      <>
                        <div className="space-y-2">
                          <Label>Date of Birth</Label>
                          <Input
                            type="date"
                            value={traveler.dateOfBirth}
                            onChange={(e) =>
                              updateTraveler(index, 'dateOfBirth', e.target.value)
                            }
                            required
                          />
                        </div>
                        <div className="space-y-2">
                          <Label>Passport Number</Label>
                          <Input
                            value={traveler.passportNumber}
                            onChange={(e) =>
                              updateTraveler(index, 'passportNumber', e.target.value)
                            }
                          />
                        </div>
                      </>
                    )}
                  </div>
                </div>
              ))}
              <Button variant="outline" onClick={addTraveler} className="w-full">
                <Plus className="mr-2 h-4 w-4" />
                Add Another Traveler
              </Button>
            </CardContent>
          </Card>

          {/* Check-in/out for accommodations */}
          {activeBookingType === 'accommodation' && (
            <Card>
              <CardHeader>
                <CardTitle>Stay Dates</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label>Check-in Date</Label>
                    <Input
                      type="date"
                      value={checkIn}
                      onChange={(e) => setCheckIn(e.target.value)}
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Check-out Date</Label>
                    <Input
                      type="date"
                      value={checkOut}
                      onChange={(e) => setCheckOut(e.target.value)}
                      required
                    />
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Payment Method */}
          <Card>
            <CardHeader>
              <CardTitle>Payment Method</CardTitle>
            </CardHeader>
            <CardContent>
              <RadioGroup
                value={paymentMethod}
                onValueChange={(value) =>
                  setPaymentMethod(value as PaymentMethod)
                }
                className="space-y-3"
              >
                {paymentMethods.map((method) => (
                  <div
                    key={method.value}
                    className={`flex items-center space-x-3 rounded-lg border p-4 transition-colors ${
                      paymentMethod === method.value
                        ? 'border-primary bg-primary/5'
                        : 'hover:bg-muted/50'
                    }`}
                  >
                    <RadioGroupItem value={method.value} id={method.value} />
                    <Label
                      htmlFor={method.value}
                      className="flex flex-1 cursor-pointer items-center gap-3"
                    >
                      <span className="text-muted-foreground">
                        {method.icon}
                      </span>
                      {method.label}
                    </Label>
                  </div>
                ))}
              </RadioGroup>
            </CardContent>
          </Card>
        </div>

        {/* Booking Summary */}
        <div className="lg:col-span-1">
          <Card className="sticky top-24">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                {getBookingIcon()}
                Booking Summary
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h4 className="font-semibold">{getBookingTitle()}</h4>
                {getBookingDetails()}
              </div>

              <Separator />

              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Travelers</span>
                  <span>{travelers.length}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Payment</span>
                  <span className="capitalize">{paymentMethod}</span>
                </div>
              </div>

              <Button
                className="w-full"
                size="lg"
                onClick={handleSubmit}
                disabled={isBooking}
              >
                {isBooking ? (
                  <>
                    <Spinner size="sm" className="mr-2" />
                    Processing...
                  </>
                ) : (
                  'Confirm Booking'
                )}
              </Button>
              <p className="text-center text-xs text-muted-foreground">
                By confirming, you agree to our terms and conditions
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </main>
  );
}
