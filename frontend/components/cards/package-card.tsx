'use client';

import Link from 'next/link';
import Image from 'next/image';
import { format } from 'date-fns';
import { MapPin, Star, Calendar, Users, Luggage } from 'lucide-react';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Package } from '@/lib/types';

interface PackageCardProps {
  pkg: Package;
}

export function PackageCard({ pkg }: PackageCardProps) {
  const startFrom = new Date(pkg.startDateFrom);
  const startTo = new Date(pkg.startDateTo);

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-lg">
      <div className="relative aspect-[16/9] overflow-hidden">
        {pkg.images?.[0] ? (
          <Image
            src={pkg.images[0]}
            alt={pkg.name}
            fill
            className="object-cover transition-transform hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-muted">
            <Luggage className="h-12 w-12 text-muted-foreground" />
          </div>
        )}
        <div className="absolute left-3 top-3 flex flex-wrap gap-1">
          {pkg.tags?.slice(0, 2).map((tag) => (
            <Badge key={tag} className="bg-primary/90">
              {tag}
            </Badge>
          ))}
        </div>
      </div>
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between gap-2">
          <h3 className="line-clamp-1 text-lg font-semibold">{pkg.name}</h3>
          {pkg.rating > 0 && (
            <div className="flex items-center gap-1 text-sm">
              <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
              <span className="font-medium">{pkg.rating.toFixed(1)}</span>
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-3 pb-3">
        <div className="flex items-center gap-1 text-sm text-muted-foreground">
          <MapPin className="h-4 w-4" />
          <span>{pkg.destination}</span>
        </div>
        <p className="line-clamp-2 text-sm text-muted-foreground">
          {pkg.description}
        </p>
        <div className="flex flex-wrap gap-3 text-sm">
          <div className="flex items-center gap-1 text-muted-foreground">
            <Calendar className="h-4 w-4" />
            <span>{pkg.durationDays} days</span>
          </div>
          <div className="flex items-center gap-1 text-muted-foreground">
            <Users className="h-4 w-4" />
            <span>Up to {pkg.maxTravelers}</span>
          </div>
        </div>
        <p className="text-xs text-muted-foreground">
          Available: {format(startFrom, 'MMM d')} - {format(startTo, 'MMM d, yyyy')}
        </p>
        <div className="flex flex-wrap gap-1">
          {pkg.includedComponents?.slice(0, 3).map((component) => (
            <Badge key={component.type} variant="outline" className="text-xs">
              {component.type}
            </Badge>
          ))}
        </div>
      </CardContent>
      <CardFooter className="flex items-center justify-between border-t bg-muted/30 pt-4">
        <div>
          <p className="text-sm text-muted-foreground">From</p>
          <p className="text-xl font-bold text-primary">
            {pkg.currency} {pkg.basePrice.toLocaleString()}
          </p>
        </div>
        <Button asChild>
          <Link href={`/packages/${pkg.id}`}>View Package</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
