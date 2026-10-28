'use client';

import Link from 'next/link';
import { Plane, Building2, Compass, Luggage, ArrowRight, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Header } from '@/components/layout/header';

const categories = [
  {
    title: 'Flights',
    description: 'Find and book flights to destinations worldwide',
    icon: Plane,
    href: '/flights',
    searchHref: '/search/flights',
    color: 'bg-blue-500/10 text-blue-600',
  },
  {
    title: 'Accommodations',
    description: 'Discover hotels, apartments, and unique stays',
    icon: Building2,
    href: '/accommodations',
    searchHref: '/search/accommodations',
    color: 'bg-green-500/10 text-green-600',
  },
  {
    title: 'Activities',
    description: 'Explore tours, experiences, and adventures',
    icon: Compass,
    href: '/activities',
    searchHref: '/search/activities',
    color: 'bg-orange-500/10 text-orange-600',
  },
  {
    title: 'Packages',
    description: 'All-inclusive travel packages for easy planning',
    icon: Luggage,
    href: '/packages',
    searchHref: '/search/packages',
    color: 'bg-purple-500/10 text-purple-600',
  },
];

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      
      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gradient-to-br from-primary/10 via-background to-accent/10">
        <div className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8 lg:py-32">
          <div className="text-center">
            <h1 className="text-balance text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
              Your Gateway to
              <span className="text-primary"> Amazing Travel</span>
            </h1>
            <p className="mx-auto mt-6 max-w-2xl text-pretty text-lg text-muted-foreground">
              Discover incredible destinations, book flights, find perfect accommodations,
              and create unforgettable experiences with TravelLink.
            </p>
            <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
              <Button size="lg" asChild>
                <Link href="/search/flights">
                  <Search className="mr-2 h-5 w-5" />
                  Search Flights
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild>
                <Link href="/packages">
                  Explore Packages
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
        
        {/* Decorative elements */}
        <div className="absolute -left-20 -top-20 h-72 w-72 rounded-full bg-primary/5 blur-3xl" />
        <div className="absolute -bottom-20 -right-20 h-72 w-72 rounded-full bg-accent/10 blur-3xl" />
      </section>

      {/* Categories Section */}
      <section className="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
        <div className="text-center">
          <h2 className="text-3xl font-bold text-foreground">
            What are you looking for?
          </h2>
          <p className="mt-2 text-muted-foreground">
            Browse our categories to find exactly what you need
          </p>
        </div>

        <div className="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {categories.map((category) => {
            const Icon = category.icon;
            return (
              <Card
                key={category.title}
                className="group relative overflow-hidden transition-all hover:shadow-lg"
              >
                <CardHeader>
                  <div
                    className={`mb-2 inline-flex h-12 w-12 items-center justify-center rounded-lg ${category.color}`}
                  >
                    <Icon className="h-6 w-6" />
                  </div>
                  <CardTitle className="text-xl">{category.title}</CardTitle>
                  <CardDescription>{category.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex gap-2">
                    <Button variant="outline" size="sm" asChild>
                      <Link href={category.href}>Browse</Link>
                    </Button>
                    <Button size="sm" asChild>
                      <Link href={category.searchHref}>
                        <Search className="mr-1 h-4 w-4" />
                        Search
                      </Link>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      </section>

      {/* Features Section */}
      <section className="border-t bg-muted/30">
        <div className="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-foreground">
              Why Choose TravelLink?
            </h2>
            <p className="mt-2 text-muted-foreground">
              We make travel planning simple and enjoyable
            </p>
          </div>

          <div className="mt-12 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <div className="text-center">
              <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-primary/10">
                <Search className="h-7 w-7 text-primary" />
              </div>
              <h3 className="mt-4 text-lg font-semibold">Easy Search</h3>
              <p className="mt-2 text-sm text-muted-foreground">
                Find flights, hotels, and activities with our powerful search tools
              </p>
            </div>
            <div className="text-center">
              <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-primary/10">
                <Plane className="h-7 w-7 text-primary" />
              </div>
              <h3 className="mt-4 text-lg font-semibold">Best Prices</h3>
              <p className="mt-2 text-sm text-muted-foreground">
                Compare fares from multiple providers before booking
              </p>
            </div>
            <div className="text-center">
              <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-primary/10">
                <Building2 className="h-7 w-7 text-primary" />
              </div>
              <h3 className="mt-4 text-lg font-semibold">Verified Listings</h3>
              <p className="mt-2 text-sm text-muted-foreground">
                All our listings are verified for quality and authenticity
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-primary">
        <div className="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-primary-foreground">
              Ready to start your journey?
            </h2>
            <p className="mt-2 text-primary-foreground/80">
              Create an account and start exploring amazing destinations today
            </p>
            <div className="mt-8 flex justify-center gap-4">
              <Button size="lg" variant="secondary" asChild>
                <Link href="/register">Get Started</Link>
              </Button>
              <Button
                size="lg"
                variant="outline"
                className="border-primary-foreground/20 text-primary-foreground hover:bg-primary-foreground/10"
                asChild
              >
                <Link href="/offers">View Flight Fares</Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t">
        <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
          <div className="flex flex-col items-center justify-between gap-6 md:flex-row">
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                <Plane className="h-4 w-4 text-primary-foreground" />
              </div>
              <span className="text-lg font-bold">TravelLink</span>
            </div>
            <nav className="flex flex-wrap justify-center gap-6">
              <Link
                href="/flights"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Flights
              </Link>
              <Link
                href="/accommodations"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Accommodations
              </Link>
              <Link
                href="/activities"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Activities
              </Link>
              <Link
                href="/packages"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Packages
              </Link>
              <Link
                href="/offers"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Offers
              </Link>
            </nav>
            <p className="text-sm text-muted-foreground">
              © {new Date().getFullYear()} TravelLink. All rights reserved.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
