'use client';

import { useEffect, useState } from 'react';
import {
  CalendarCheck,
  Plane,
  Building2,
  Ticket,
  Package,
  DollarSign,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Spinner } from '@/components/ui/spinner';
import { DataTable, Column } from '@/components/admin/data-table';
import { useAdminStore } from '@/lib/stores/admin-store';
import type { Booking } from '@/lib/types';
import { toast } from 'sonner';
import { format } from 'date-fns';

const typeIcons = {
  flight: Plane,
  accommodation: Building2,
  activity: Ticket,
  package: Package,
};

const columns: Column<Booking>[] = [
  {
    key: 'id',
    header: 'Booking ID',
    render: (booking) => (
      <span className="font-mono text-sm">{booking.id.slice(0, 8)}...</span>
    ),
  },
  {
    key: 'type',
    header: 'Type',
    sortable: true,
    render: (booking) => {
      const Icon = typeIcons[booking.type] || CalendarCheck;
      return (
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
            <Icon className="h-4 w-4 text-primary" />
          </div>
          <span className="capitalize">{booking.type}</span>
        </div>
      );
    },
  },
  {
    key: 'createdAt',
    header: 'Date',
    sortable: true,
    render: (booking) => format(new Date(booking.createdAt), 'MMM d, yyyy HH:mm'),
  },
  {
    key: 'totalPrice',
    header: 'Amount',
    sortable: true,
    render: (booking) => (
      <span className="font-semibold text-primary">
        {booking.currency} {booking.totalPrice.toLocaleString()}
      </span>
    ),
  },
  {
    key: 'status',
    header: 'Status',
    sortable: true,
    render: (booking) => {
      const statusColors: Record<string, string> = {
        confirmed: 'bg-success-50 text-success-400',
        pending: 'bg-amber-50 text-amber-600',
        cancelled: 'bg-danger-50 text-danger-400',
        completed: 'bg-info-50 text-info-400',
      };
      return (
        <span
          className={`rounded-full px-2 py-1 text-xs font-medium capitalize ${
            statusColors[booking.status] || 'bg-muted text-muted-foreground'
          }`}
        >
          {booking.status}
        </span>
      );
    },
  },
];

export default function AdminBookingsPage() {
  const { bookings, isLoading, fetchBookings, updateBookingStatus } =
    useAdminStore();

  const [isStatusDialogOpen, setIsStatusDialogOpen] = useState(false);
  const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null);
  const [newStatus, setNewStatus] = useState('pending');

  useEffect(() => {
    fetchBookings();
  }, [fetchBookings]);

  const handleEditStatus = (booking: Booking) => {
    setSelectedBooking(booking);
    setNewStatus(booking.status);
    setIsStatusDialogOpen(true);
  };

  const handleUpdateStatus = async () => {
    if (selectedBooking) {
      try {
        await updateBookingStatus(selectedBooking.id, newStatus);
        toast.success('Booking status updated successfully');
        setIsStatusDialogOpen(false);
      } catch {
        toast.error('Failed to update booking status');
      }
    }
  };

  const confirmedBookings = bookings.filter((b) => b.status === 'confirmed');
  const pendingBookings = bookings.filter((b) => b.status === 'pending');
  const totalRevenue = bookings
    .filter((b) => b.status === 'confirmed' || b.status === 'completed')
    .reduce((sum, b) => sum + b.totalPrice, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-foreground">Bookings</h1>
        <p className="mt-1 text-muted-foreground">
          Manage and track all bookings
        </p>
      </div>

      {/* Stats */}
      <div className="grid gap-4 sm:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Bookings
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2 text-2xl font-bold">
              <CalendarCheck className="h-5 w-5 text-muted-foreground" />
              {bookings.length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Confirmed
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-success-400">
              {confirmedBookings.length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Pending
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-600">
              {pendingBookings.length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Revenue
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-1 text-2xl font-bold text-primary">
              <DollarSign className="h-5 w-5" />
              {totalRevenue.toLocaleString()}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Breakdown by Type */}
      <div className="grid gap-4 sm:grid-cols-4">
        {Object.entries(typeIcons).map(([type, Icon]) => {
          const count = bookings.filter((b) => b.type === type).length;
          return (
            <Card key={type}>
              <CardContent className="flex items-center gap-3 p-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                  <Icon className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{count}</p>
                  <p className="text-sm capitalize text-muted-foreground">
                    {type} bookings
                  </p>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Table */}
      {isLoading && bookings.length === 0 ? (
        <div className="flex h-64 items-center justify-center">
          <Spinner className="h-8 w-8 text-primary" />
        </div>
      ) : (
        <DataTable
          data={bookings}
          columns={columns}
          searchPlaceholder="Search bookings..."
          searchKey="id"
          onEdit={handleEditStatus}
        />
      )}

      {/* Status Dialog */}
      <Dialog open={isStatusDialogOpen} onOpenChange={setIsStatusDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Update Booking Status</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="rounded-lg border border-border p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">Booking ID</p>
                  <p className="font-mono text-sm text-muted-foreground">
                    {selectedBooking?.id}
                  </p>
                </div>
                <div className="text-right">
                  <p className="font-medium">Amount</p>
                  <p className="text-primary">
                    {selectedBooking?.currency}{' '}
                    {selectedBooking?.totalPrice.toLocaleString()}
                  </p>
                </div>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="status">Status</Label>
              <Select value={newStatus} onValueChange={setNewStatus}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="confirmed">Confirmed</SelectItem>
                  <SelectItem value="completed">Completed</SelectItem>
                  <SelectItem value="cancelled">Cancelled</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsStatusDialogOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleUpdateStatus} disabled={isLoading}>
              {isLoading ? 'Updating...' : 'Update Status'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
