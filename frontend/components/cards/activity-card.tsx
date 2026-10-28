'use client';

import Link from 'next/link';
import Image from 'next/image';
import { MapPin, Star, Clock, Compass } from 'lucide-react';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import type { Activity } from '@/lib/types';

interface ActivityCardProps {
  activity: Activity;
}

export function ActivityCard({ activity }: ActivityCardProps) {
  const hours = Math.floor(activity.durationMinutes / 60);
  const minutes = activity.durationMinutes % 60;
  const durationText =
    hours > 0
      ? `${hours}h${minutes > 0 ? ` ${minutes}m` : ''}`
      : `${minutes}m`;

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-lg">
      <div className="relative aspect-[4/3] overflow-hidden">
        {activity.images?.[0] ? (
          <Image
            src={activity.images[0]}
            alt={activity.title}
            fill
            className="object-cover transition-transform hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-muted">
            <Compass className="h-12 w-12 text-muted-foreground" />
          </div>
        )}
        {activity.categories?.[0] && (
          <Badge className="absolute right-3 top-3">
            {activity.categories[0]}
          </Badge>
        )}
      </div>
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between gap-2">
          <h3 className="line-clamp-1 text-lg font-semibold">
            {activity.title}
          </h3>
          {activity.rating > 0 && (
            <div className="flex items-center gap-1 text-sm">
              <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
              <span className="font-medium">{activity.rating.toFixed(1)}</span>
              <span className="text-muted-foreground">
                ({activity.reviewCount})
              </span>
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="pb-3">
        <div className="flex items-center gap-1 text-sm text-muted-foreground">
          <MapPin className="h-4 w-4" />
          <span>
            {activity.city}, {activity.country}
          </span>
        </div>
        <p className="mt-2 line-clamp-2 text-sm text-muted-foreground">
          {activity.description}
        </p>
        <div className="mt-3 flex items-center gap-4 text-sm">
          <div className="flex items-center gap-1 text-muted-foreground">
            <Clock className="h-4 w-4" />
            <span>{durationText}</span>
          </div>
        </div>
        <div className="mt-3 flex flex-wrap gap-1">
          {activity.categories?.slice(0, 2).map((category) => (
            <Badge key={category} variant="secondary" className="text-xs">
              {category}
            </Badge>
          ))}
        </div>
      </CardContent>
      <CardFooter className="border-t bg-muted/30 pt-4">
        <Button asChild className="w-full">
          <Link href={`/activities/${activity.id}`}>View Details</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
