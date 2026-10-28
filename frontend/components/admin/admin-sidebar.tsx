'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  LayoutDashboard,
  Plane,
  Building2,
  Ticket,
  Package,
  Tag,
  Users,
  CalendarCheck,
  Settings,
  ChevronLeft,
  ChevronRight,
  LogOut,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/lib/stores/auth-store';
import { useState } from 'react';

const navigation = [
  { name: 'Dashboard', href: '/admin', icon: LayoutDashboard },
  { name: 'Flights', href: '/admin/flights', icon: Plane },
  { name: 'Accommodations', href: '/admin/accommodations', icon: Building2 },
  { name: 'Activities', href: '/admin/activities', icon: Ticket },
  { name: 'Packages', href: '/admin/packages', icon: Package },
  { name: 'Offers', href: '/admin/offers', icon: Tag },
  { name: 'Users', href: '/admin/users', icon: Users },
  { name: 'Bookings', href: '/admin/bookings', icon: CalendarCheck },
];

export function AdminSidebar() {
  const pathname = usePathname();
  const { logout, user } = useAuthStore();
  const [collapsed, setCollapsed] = useState(false);

  return (
    <aside
      className={cn(
        'fixed left-0 top-0 z-40 h-screen bg-deep-ocean text-white transition-all duration-300',
        collapsed ? 'w-16' : 'w-64'
      )}
    >
      <div className="flex h-full flex-col">
        {/* Logo */}
        <div className="flex h-16 items-center justify-between border-b border-navy-600 px-4">
          {!collapsed && (
            <Link href="/admin" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-tropical-teal">
                <Settings className="h-4 w-4 text-white" />
              </div>
              <span className="font-bold">Admin Panel</span>
            </Link>
          )}
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setCollapsed(!collapsed)}
            className="h-8 w-8 text-white hover:bg-white/10"
          >
            {collapsed ? (
              <ChevronRight className="h-4 w-4" />
            ) : (
              <ChevronLeft className="h-4 w-4" />
            )}
          </Button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 space-y-1 overflow-y-auto p-3">
          {navigation.map((item) => {
            const Icon = item.icon;
            const isActive =
              item.href === '/admin'
                ? pathname === '/admin'
                : pathname.startsWith(item.href);

            return (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                  isActive
                    ? 'bg-tropical-teal text-white'
                    : 'text-white/70 hover:bg-white/10 hover:text-white',
                  collapsed && 'justify-center px-2'
                )}
                title={collapsed ? item.name : undefined}
              >
                <Icon className="h-5 w-5 flex-shrink-0" />
                {!collapsed && <span>{item.name}</span>}
              </Link>
            );
          })}
        </nav>

        {/* User section */}
        <div className="border-t border-navy-600 p-3">
          {!collapsed && user && (
            <div className="mb-3 rounded-lg bg-navy-600 p-3">
              <p className="text-sm font-medium text-white">
                {user.firstName} {user.lastName}
              </p>
              <p className="text-xs text-white/60">{user.email}</p>
              <p className="mt-1 inline-block rounded bg-tropical-teal/20 px-2 py-0.5 text-xs text-tropical-teal">
                {user.role}
              </p>
            </div>
          )}

          <div className="flex gap-2">
            <Button
              variant="ghost"
              size={collapsed ? 'icon' : 'default'}
              className={cn(
                'text-white/70 hover:bg-white/10 hover:text-white',
                !collapsed && 'w-full justify-start'
              )}
              asChild
            >
              <Link href="/">
                <ChevronLeft className="h-4 w-4" />
                {!collapsed && <span className="ml-2">Back to Site</span>}
              </Link>
            </Button>
          </div>

          <Button
            variant="ghost"
            size={collapsed ? 'icon' : 'default'}
            onClick={logout}
            className={cn(
              'mt-2 text-red-400 hover:bg-red-500/10 hover:text-red-300',
              !collapsed && 'w-full justify-start'
            )}
          >
            <LogOut className="h-4 w-4" />
            {!collapsed && <span className="ml-2">Logout</span>}
          </Button>
        </div>
      </div>
    </aside>
  );
}
