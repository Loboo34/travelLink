'use client';

import { useEffect } from 'react';
import {
  Plane,
  Building2,
  Ticket,
  Package,
  Tag,
  Users,
  CalendarCheck,
  DollarSign,
  TrendingUp,
  ArrowUpRight,
  ArrowDownRight,
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Spinner } from '@/components/ui/spinner';
import { useAdminStore } from '@/lib/stores/admin-store';

const statCards = [
  {
    name: 'Total Flights',
    key: 'totalFlights' as const,
    icon: Plane,
    color: 'bg-info-400',
    trend: '+12%',
    trendUp: true,
  },
  {
    name: 'Accommodations',
    key: 'totalAccommodations' as const,
    icon: Building2,
    color: 'bg-success-400',
    trend: '+8%',
    trendUp: true,
  },
  {
    name: 'Activities',
    key: 'totalActivities' as const,
    icon: Ticket,
    color: 'bg-amber-400',
    trend: '+15%',
    trendUp: true,
  },
  {
    name: 'Packages',
    key: 'totalPackages' as const,
    icon: Package,
    color: 'bg-teal-500',
    trend: '+5%',
    trendUp: true,
  },
  {
    name: 'Active Options',
    key: 'totalOffers' as const,
    icon: Tag,
    color: 'bg-danger-400',
    trend: '-3%',
    trendUp: false,
  },
  {
    name: 'Total Users',
    key: 'totalUsers' as const,
    icon: Users,
    color: 'bg-navy-500',
    trend: '+22%',
    trendUp: true,
  },
  {
    name: 'Bookings',
    key: 'totalBookings' as const,
    icon: CalendarCheck,
    color: 'bg-teal-600',
    trend: '+18%',
    trendUp: true,
  },
  {
    name: 'Revenue',
    key: 'revenue' as const,
    icon: DollarSign,
    color: 'bg-success-400',
    trend: '+25%',
    trendUp: true,
    isCurrency: true,
  },
];

export default function AdminDashboard() {
  const { stats, isLoading, fetchStats } = useAdminStore();

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  if (isLoading && !stats) {
    return (
      <div className="flex h-[60vh] items-center justify-center">
        <Spinner className="h-8 w-8 text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
        <p className="mt-1 text-muted-foreground">
          Welcome to TravelLink Admin Panel
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => {
          const Icon = stat.icon;
          const value = stats?.[stat.key] ?? 0;
          const displayValue = stat.isCurrency
            ? `$${value.toLocaleString()}`
            : value.toLocaleString();

          return (
            <Card key={stat.key} className="overflow-hidden">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.name}
                </CardTitle>
                <div
                  className={`flex h-10 w-10 items-center justify-center rounded-lg ${stat.color}`}
                >
                  <Icon className="h-5 w-5 text-white" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-foreground">
                  {displayValue}
                </div>
                <div className="mt-1 flex items-center text-sm">
                  {stat.trendUp ? (
                    <ArrowUpRight className="h-4 w-4 text-success-400" />
                  ) : (
                    <ArrowDownRight className="h-4 w-4 text-danger-400" />
                  )}
                  <span
                    className={
                      stat.trendUp ? 'text-success-400' : 'text-danger-400'
                    }
                  >
                    {stat.trend}
                  </span>
                  <span className="ml-1 text-muted-foreground">
                    from last month
                  </span>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Quick Actions */}
      <div>
        <h2 className="mb-4 text-xl font-semibold text-foreground">
          Quick Actions
        </h2>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <QuickActionCard
            title="Add New Flight"
            description="Create a new flight listing"
            href="/admin/flights?action=new"
            icon={Plane}
          />
          <QuickActionCard
            title="Add Accommodation"
            description="List a new property"
            href="/admin/accommodations?action=new"
            icon={Building2}
          />
          <QuickActionCard
            title="Create Activity"
            description="Add a new activity"
            href="/admin/activities?action=new"
            icon={Ticket}
          />
          <QuickActionCard
            title="Build Package"
            description="Create a travel package"
            href="/admin/packages?action=new"
            icon={Package}
          />
        </div>
      </div>

      {/* Recent Activity Placeholder */}
      <div>
        <h2 className="mb-4 text-xl font-semibold text-foreground">
          Recent Activity
        </h2>
        <Card>
          <CardContent className="py-8">
            <div className="flex flex-col items-center justify-center text-center">
              <TrendingUp className="h-12 w-12 text-muted-foreground/50" />
              <h3 className="mt-4 text-lg font-medium text-foreground">
                Activity Feed
              </h3>
              <p className="mt-1 text-sm text-muted-foreground">
                Recent bookings, reviews, and user activity will appear here.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function QuickActionCard({
  title,
  description,
  href,
  icon: Icon,
}: {
  title: string;
  description: string;
  href: string;
  icon: React.ElementType;
}) {
  return (
    <a href={href}>
      <Card className="cursor-pointer transition-all hover:border-primary hover:shadow-md">
        <CardContent className="flex items-center gap-4 p-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
            <Icon className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h3 className="font-medium text-foreground">{title}</h3>
            <p className="text-sm text-muted-foreground">{description}</p>
          </div>
        </CardContent>
      </Card>
    </a>
  );
}
