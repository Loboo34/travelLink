'use client';

import Link from 'next/link';
import { format } from 'date-fns';
import { Plane, Luggage, Clock, Tag } from 'lucide-react';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Offer } from '@/lib/types';

interface OfferCardProps {
  offer: Offer;
}

export function OfferCard({ offer }: OfferCardProps) {
  const expiresAt = offer.expiresAt ? new Date(offer.expiresAt) : null;
  const isExpiringSoon =
    expiresAt !== null && expiresAt.getTime() - Date.now() < 24 * 60 * 60 * 1000;
  const baggage = offer.baggageAllowance?.[0];
  const isOneWay = offer.oneway ?? offer.oneWay ?? true;

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-lg">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
              <Tag className="h-5 w-5 text-primary" />
            </div>
            <div>
              <p className="font-semibold">{offer.provider}</p>
              <p className="text-sm text-muted-foreground">
                Ref: {offer.providerReference}
              </p>
            </div>
          </div>
          <Badge variant={offer.isActive ? 'default' : 'secondary'}>
              {offer.isActive ? 'Active' : 'Expired'}
          </Badge>
        </div>
      </CardHeader>
      <CardContent className="space-y-4 pb-3">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Plane className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm">
              {isOneWay ? 'One Way' : 'Round Trip'}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <Luggage className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm">
              {baggage?.pieces || 0} piece(s), {baggage?.weightKg || 0}kg
            </span>
          </div>
        </div>

        <div className="flex items-center justify-between rounded-lg bg-muted/50 px-4 py-3">
          <div>
            <p className="text-sm text-muted-foreground">Total Price</p>
            <p className="text-2xl font-bold text-primary">
              ${offer.priceTotal.toLocaleString()}
            </p>
          </div>
          <div className="text-right">
            <p className="text-sm text-muted-foreground">Seats Available</p>
            <p className="text-lg font-semibold">{offer.bookableSeats}</p>
          </div>
        </div>

        <div className="flex items-center gap-2 text-sm">
          <Clock className="h-4 w-4 text-muted-foreground" />
          <span className={isExpiringSoon ? 'text-destructive' : 'text-muted-foreground'}>
            Expires: {expiresAt ? format(expiresAt, 'MMM d, yyyy HH:mm') : 'Not set'}
          </span>
          {isExpiringSoon && (
            <Badge variant="destructive" className="ml-auto">
              Expiring Soon
            </Badge>
          )}
        </div>
      </CardContent>
      <CardFooter className="border-t bg-muted/30 pt-4">
        <Button asChild className="w-full">
          <Link href={`/offers/${offer.id}`}>View Booking Option</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
