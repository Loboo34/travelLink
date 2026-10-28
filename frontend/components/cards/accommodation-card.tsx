'use client';

import Link from 'next/link';
import Image from 'next/image';
import { MapPin, Star, Bed } from 'lucide-react';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Accommodation } from '@/lib/types';

interface AccommodationCardProps {
  accommodation: Accommodation;
}

export function AccommodationCard({ accommodation }: AccommodationCardProps) {
  const lowestPrice = accommodation.roomTypes?.length
    ? Math.min(...accommodation.roomTypes.map((r) => r.basePrice))
    : 0;

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-lg">
      <div className="relative aspect-[4/3] overflow-hidden">
        {accommodation.images?.[0] ? (
          <Image
            src={accommodation.images[0]}
            alt={accommodation.name}
            fill
            className="object-cover transition-transform hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-muted">
            <Bed className="h-12 w-12 text-muted-foreground" />
          </div>
        )}
        <Badge className="absolute right-3 top-3">{accommodation.type}</Badge>
      </div>
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between gap-2">
          <h3 className="line-clamp-1 text-lg font-semibold">
            {accommodation.name}
          </h3>
          {accommodation.rating > 0 && (
            <div className="flex items-center gap-1 text-sm">
              <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
              <span className="font-medium">{accommodation.rating.toFixed(1)}</span>
              <span className="text-muted-foreground">
                ({accommodation.reviews})
              </span>
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="pb-3">
        <div className="flex items-center gap-1 text-sm text-muted-foreground">
          <MapPin className="h-4 w-4" />
          <span className="line-clamp-1">
            {accommodation.address.city}, {accommodation.address.country}
          </span>
        </div>
        <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
          {accommodation.description}
        </p>
        <div className="mt-3 flex flex-wrap gap-1">
          {accommodation.amenities?.slice(0, 3).map((amenity) => (
            <Badge key={amenity} variant="secondary" className="text-xs">
              {amenity}
            </Badge>
          ))}
          {accommodation.amenities?.length > 3 && (
            <Badge variant="secondary" className="text-xs">
              +{accommodation.amenities.length - 3} more
            </Badge>
          )}
        </div>
      </CardContent>
      <CardFooter className="flex items-center justify-between border-t bg-muted/30 pt-4">
        <div>
          <p className="text-sm text-muted-foreground">From</p>
          <p className="text-xl font-bold text-primary">
            ${lowestPrice}
            <span className="text-sm font-normal text-muted-foreground">
              /night
            </span>
          </p>
        </div>
        <Button asChild>
          <Link href={`/accommodations/${accommodation.id}`}>View Details</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
