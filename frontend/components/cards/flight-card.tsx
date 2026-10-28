'use client';

import Link from 'next/link';
import { format } from 'date-fns';
import { Plane, Clock, ArrowRight } from 'lucide-react';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Flight } from '@/lib/types';

interface FlightCardProps {
  flight: Flight;
}

export function FlightCard({ flight }: FlightCardProps) {
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

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-lg">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
              <Plane className="h-5 w-5 text-primary" />
            </div>
            <div>
              <p className="font-semibold">{flight.flightNumber}</p>
              <p className="text-sm text-muted-foreground">
                {flight.airlineID}
              </p>
            </div>
          </div>
          <Badge
            variant={flight.status === 'scheduled' ? 'default' : 'secondary'}
          >
            {flight.status}
          </Badge>
        </div>
      </CardHeader>
      <CardContent className="pb-3">
        <div className="flex items-center justify-between">
          <div className="text-center">
            <p className="text-2xl font-bold">
              {format(departureTime, 'HH:mm')}
            </p>
            <p className="text-sm text-muted-foreground">{flight.originID}</p>
          </div>
          <div className="flex flex-1 flex-col items-center px-4">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4" />
              {hours}h {minutes}m
            </div>
            <div className="relative mt-2 w-full">
              <div className="h-0.5 w-full bg-border" />
              <ArrowRight className="absolute -right-1 top-1/2 h-4 w-4 -translate-y-1/2 text-primary" />
              {flight.stops > 0 && (
                <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
                  <div className="h-2 w-2 rounded-full bg-primary" />
                </div>
              )}
            </div>
            <p className="mt-2 text-xs text-muted-foreground">
              {flight.stops === 0
                ? 'Direct'
                : `${flight.stops} stop${flight.stops > 1 ? 's' : ''}`}
            </p>
          </div>
          <div className="text-center">
            <p className="text-2xl font-bold">{format(arrivalTime, 'HH:mm')}</p>
            <p className="text-sm text-muted-foreground">
              {flight.destinationID}
            </p>
          </div>
        </div>
        <div className="mt-4 flex items-center justify-between text-sm">
          <span className="text-muted-foreground">
            {format(departureTime, 'EEE, MMM d, yyyy')}
          </span>
          <Badge variant="outline" className="capitalize">
            {cabinLabel}
          </Badge>
        </div>
      </CardContent>
      <CardFooter className="border-t bg-muted/30 pt-4">
        <Button asChild className="w-full">
          <Link href={`/flights/${flight.id}`}>View Details</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
