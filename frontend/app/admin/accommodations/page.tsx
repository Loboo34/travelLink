'use client';

import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import { Bed, Building2, CalendarClock, Eye, Plus, Star, Trash2, X } from 'lucide-react';
import { format } from 'date-fns';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Spinner } from '@/components/ui/spinner';
import { Switch } from '@/components/ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';
import { DataTable, Column } from '@/components/admin/data-table';
import { DeleteDialog } from '@/components/admin/delete-dialog';
import { useAdminStore } from '@/lib/stores/admin-store';
import type { Accommodation, AccommodationAvailability, RoomType } from '@/lib/types';

const PROPERTY_TYPES = ['Hotel', 'AirBNB', 'Villa', 'Resort', 'Guesthouse', 'Apartment'];

const emptyRoom: RoomType = {
  id: '',
  name: '',
  description: '',
  bedrooms: 1,
  bathrooms: 1,
  maxGuests: 2,
  isEntirePlace: false,
  amenities: [],
  images: [],
  basePrice: 0,
};

const emptyAccommodationForm = {
  name: '',
  description: '',
  propertyType: 'Hotel',
  street: '',
  city: '',
  country: '',
  zipCode: '',
  latitude: 0,
  longitude: 0,
  checkInTime: '14:00',
  checkOutTime: '11:00',
  amenities: '',
  images: '',
  rooms: [{ ...emptyRoom }],
};

const emptyAvailabilityForm = {
  accommodationID: '',
  roomTypeID: '',
  date: '',
  totalRooms: 1,
  pricePerNight: 0,
};

function getRooms(accommodation?: Accommodation | null) {
  return accommodation?.roomTypes || accommodation?.roomType || [];
}

function formatDate(value?: string) {
  if (!value) return 'Not set';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? 'Invalid date' : format(date, 'MMM d, yyyy');
}

function getTimeNumber(value: string) {
  const [hours = '0', minutes = '0'] = value.split(':');
  return Number.parseInt(hours, 10) * 100 + Number.parseInt(minutes, 10);
}

function toTimeInput(value?: string | number) {
  if (!value) return '';
  if (typeof value === 'string' && value.includes(':')) return value.slice(0, 5);
  const numeric = Number.parseInt(String(value), 10);
  if (Number.isNaN(numeric)) return '';
  return `${String(Math.floor(numeric / 100)).padStart(2, '0')}:${String(numeric % 100).padStart(2, '0')}`;
}

export default function AdminAccommodationsPage() {
  const searchParams = useSearchParams();
  const {
    accommodations,
    accommodationAvailabilities,
    isLoading,
    fetchAccommodations,
    fetchAccommodationAvailabilities,
    createAccommodation,
    updateAccommodation,
    deleteAccommodation,
    createAccommodationAvailability,
    updateAccommodationAvailabilityStatus,
  } = useAdminStore();

  const [activeTab, setActiveTab] = useState('properties');
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isAvailabilityFormOpen, setIsAvailabilityFormOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<Accommodation | null>(null);
  const [selectedDetails, setSelectedDetails] = useState<Accommodation | null>(null);
  const [formData, setFormData] = useState(emptyAccommodationForm);
  const [availabilityForm, setAvailabilityForm] = useState(emptyAvailabilityForm);

  const accommodationList = Array.isArray(accommodations) ? accommodations : [];
  const availabilityList = Array.isArray(accommodationAvailabilities)
    ? accommodationAvailabilities
    : [];

  const accommodationById = useMemo(
    () => new Map(accommodationList.map((accommodation) => [accommodation.id, accommodation])),
    [accommodationList]
  );

  const selectedAvailabilityAccommodation = accommodationById.get(
    availabilityForm.accommodationID
  );
  const selectedAvailabilityRooms = getRooms(selectedAvailabilityAccommodation);

  const columns: Column<Accommodation>[] = [
    {
      key: 'name',
      header: 'Property',
      sortable: true,
      render: (acc) => (
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded bg-primary/10">
            <Building2 className="h-5 w-5 text-primary" />
          </div>
          <div>
            <p className="font-medium">{acc.name}</p>
            <p className="text-sm text-muted-foreground">{acc.type || acc.propertyType}</p>
          </div>
        </div>
      ),
    },
    {
      key: 'location',
      header: 'Location',
      render: (acc) => `${acc.address?.city || '-'}, ${acc.address?.country || '-'}`,
    },
    {
      key: 'rating',
      header: 'Rating',
      sortable: true,
      render: (acc) => (
        <div className="flex items-center gap-1">
          <Star className="h-4 w-4 fill-amber-400 text-amber-400" />
          <span className="font-medium">{(acc.rating || 0).toFixed(1)}</span>
          <span className="text-muted-foreground">({acc.reviews || 0})</span>
        </div>
      ),
    },
    {
      key: 'roomTypes',
      header: 'Rooms',
      render: (acc) => (
        <span className="rounded-full bg-primary/10 px-2 py-1 text-xs font-medium text-primary">
          {getRooms(acc).length} types
        </span>
      ),
    },
    {
      key: 'isActive',
      header: 'Status',
      render: (acc) => (
        <span className="rounded-full bg-muted px-2 py-1 text-xs font-medium">
          {acc.isActive === false ? 'Inactive' : 'Listed'}
        </span>
      ),
    },
  ];

  const availabilityColumns: Column<AccommodationAvailability>[] = [
    {
      key: 'accommodationID',
      header: 'Property',
      render: (availability) => {
        const accommodation = accommodationById.get(availability.accommodationID);
        const room = getRooms(accommodation).find((item) => item.id === availability.roomTypeID);

        return (
          <div>
            <p className="font-medium">{accommodation?.name || availability.accommodationID}</p>
            <p className="text-sm text-muted-foreground">{room?.name || availability.roomTypeID}</p>
          </div>
        );
      },
    },
    {
      key: 'date',
      header: 'Date',
      sortable: true,
      render: (availability) => formatDate(availability.date),
    },
    {
      key: 'totalRooms',
      header: 'Rooms',
      sortable: true,
      render: (availability) => (
        <span>
          {Math.max(availability.totalRooms - availability.reservedRooms, 0)} of{' '}
          {availability.totalRooms} available
        </span>
      ),
    },
    {
      key: 'pricePerNight',
      header: 'Price',
      sortable: true,
      render: (availability) => (
        <span className="font-semibold text-primary">
          {(availability.currency || 'USD')} {availability.pricePerNight.toLocaleString()}
        </span>
      ),
    },
    {
      key: 'isActive',
      header: 'Available',
      render: (availability) => (
        <Switch
          checked={availability.isActive}
          onCheckedChange={(checked) => handleToggleAvailability(availability, checked)}
          aria-label="Toggle accommodation availability"
        />
      ),
    },
  ];

  useEffect(() => {
    fetchAccommodations();
    fetchAccommodationAvailabilities();
  }, [fetchAccommodations, fetchAccommodationAvailabilities]);

  useEffect(() => {
    if (searchParams.get('action') === 'new') {
      setActiveTab('properties');
      handleCreate();
    }
  }, [searchParams]);

  function handleCreate() {
    setSelectedItem(null);
    setFormData(emptyAccommodationForm);
    setIsFormOpen(true);
  }

  function handleEdit(item: Accommodation) {
    setSelectedItem(item);
    setFormData({
      name: item.name,
      description: item.description,
      propertyType: item.type || item.propertyType || 'Hotel',
      street: item.address?.street || '',
      city: item.address?.city || '',
      country: item.address?.country || '',
      zipCode: item.address?.zipCode || '',
      latitude: item.location?.latitude || item.location?.coordinates?.[1] || 0,
      longitude: item.location?.longitude || item.location?.coordinates?.[0] || 0,
      checkInTime: toTimeInput(item.checkInTime) || '14:00',
      checkOutTime: toTimeInput(item.checkOutTime) || '11:00',
      amenities: (item.amenities || []).join(', '),
      images: (item.images || []).join(', '),
      rooms: getRooms(item).length ? getRooms(item) : [{ ...emptyRoom }],
    });
    setIsFormOpen(true);
  }

  function handleView(item: Accommodation) {
    setSelectedDetails(item);
    setIsDetailsOpen(true);
  }

  function handleDelete(item: Accommodation) {
    setSelectedItem(item);
    setIsDeleteOpen(true);
  }

  function handleCreateAvailability(accommodation?: Accommodation) {
    const rooms = getRooms(accommodation);
    setAvailabilityForm({
      ...emptyAvailabilityForm,
      accommodationID: accommodation?.id || '',
      roomTypeID: rooms[0]?.id || '',
      pricePerNight: rooms[0]?.basePrice || 0,
    });
    setActiveTab('availability');
    setIsAvailabilityFormOpen(true);
  }

  function handleUpdateRoom(index: number, field: keyof RoomType, value: string | number | boolean) {
    setFormData((current) => {
      const rooms = [...current.rooms];
      rooms[index] = { ...rooms[index], [field]: value };
      return { ...current, rooms };
    });
  }

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault();

    const rooms = formData.rooms
      .filter((room) => room.name && room.basePrice > 0)
      .map((room) => ({
        ...room,
        id: room.id || undefined,
        amenities:
          typeof room.amenities === 'string'
            ? String(room.amenities).split(',').map((item) => item.trim()).filter(Boolean)
            : room.amenities,
      }));

    const payload = {
      propertyType: formData.propertyType,
      type: formData.propertyType,
      name: formData.name,
      description: formData.description,
      address: {
        street: formData.street,
        city: formData.city,
        country: formData.country,
        zipCode: formData.zipCode,
      },
      location: {
        type: 'Point',
        coordinates: [formData.longitude, formData.latitude],
        latitude: formData.latitude,
        longitude: formData.longitude,
      },
      checkInTime: getTimeNumber(formData.checkInTime),
      checkOutTime: getTimeNumber(formData.checkOutTime),
      amenities: formData.amenities.split(',').map((item) => item.trim()).filter(Boolean),
      images: formData.images.split(',').map((item) => item.trim()).filter(Boolean),
      roomType: rooms,
      roomTypes: rooms,
      rating: selectedItem?.rating || 0,
      reviews: selectedItem?.reviews || 0,
    };

    try {
      if (selectedItem) {
        await updateAccommodation(selectedItem.id, payload as Partial<Accommodation>);
        toast.success('Accommodation updated successfully');
      } else {
        await createAccommodation(payload as Omit<Accommodation, 'id'>);
        toast.success('Accommodation created successfully');
      }
      setIsFormOpen(false);
    } catch (error) {
      console.error('Accommodation save error:', error);
      toast.error('Failed to save accommodation');
    }
  }

  async function handleAvailabilitySubmit(event: React.FormEvent) {
    event.preventDefault();

    try {
      await createAccommodationAvailability({
        ...availabilityForm,
        date: new Date(availabilityForm.date).toISOString(),
      });
      toast.success('Accommodation availability created');
      setIsAvailabilityFormOpen(false);
    } catch (error) {
      console.error('Availability save error:', error);
      toast.error('Failed to create availability');
    }
  }

  async function handleToggleAvailability(
    availability: AccommodationAvailability,
    isActive: boolean
  ) {
    try {
      await updateAccommodationAvailabilityStatus(availability.id, isActive);
      toast.success(`Availability ${isActive ? 'activated' : 'deactivated'}`);
    } catch {
      toast.error('Failed to update availability');
    }
  }

  async function handleConfirmDelete() {
    if (!selectedItem) return;

    try {
      await deleteAccommodation(selectedItem.id);
      toast.success('Accommodation deleted successfully');
      setIsDeleteOpen(false);
    } catch {
      toast.error('Failed to delete accommodation');
    }
  }

  const activeAvailabilityCount = availabilityList.filter((item) => item.isActive).length;
  const roomsCount = accommodationList.reduce((sum, item) => sum + getRooms(item).length, 0);
  const averageRating =
    accommodationList.length > 0
      ? accommodationList.reduce((sum, item) => sum + (item.rating || 0), 0) /
        accommodationList.length
      : 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Accommodations</h1>
          <p className="mt-1 text-muted-foreground">
            Manage properties, room details, and bookable availability
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => handleCreateAvailability()} className="gap-2">
            <CalendarClock className="h-4 w-4" />
            Add Availability
          </Button>
          <Button onClick={handleCreate} className="gap-2">
            <Plus className="h-4 w-4" />
            Add Accommodation
          </Button>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Properties
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{accommodationList.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Room Types
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{roomsCount}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Active Availability
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{activeAvailabilityCount}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Avg Rating
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-1 text-2xl font-bold">
              <Star className="h-5 w-5 fill-amber-400 text-amber-400" />
              {averageRating.toFixed(1)}
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="properties">Properties</TabsTrigger>
          <TabsTrigger value="availability">Availability</TabsTrigger>
        </TabsList>
        <TabsContent value="properties" className="mt-4">
          {isLoading && accommodationList.length === 0 ? (
            <div className="flex h-64 items-center justify-center">
              <Spinner className="h-8 w-8 text-primary" />
            </div>
          ) : (
            <DataTable
              data={accommodationList}
              columns={columns}
              searchPlaceholder="Search accommodations..."
              searchKey="name"
              onView={handleView}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          )}
        </TabsContent>
        <TabsContent value="availability" className="mt-4">
          <DataTable
            data={availabilityList}
            columns={availabilityColumns}
            searchPlaceholder="Search by availability ID..."
            searchKey="accommodationID"
          />
        </TabsContent>
      </Tabs>

      <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
        <DialogContent className="max-h-[90vh] max-w-4xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Building2 className="h-5 w-5" />
              {selectedItem ? 'Edit Accommodation' : 'Create Accommodation'}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="name">Property Name</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(event) => setFormData({ ...formData, name: event.target.value })}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Property Type</Label>
                <Select
                  value={formData.propertyType}
                  onValueChange={(value) => setFormData({ ...formData, propertyType: value })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {PROPERTY_TYPES.map((type) => (
                      <SelectItem key={type} value={type}>
                        {type}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(event) =>
                    setFormData({ ...formData, description: event.target.value })
                  }
                  rows={3}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="street">Street Address</Label>
                <Input
                  id="street"
                  value={formData.street}
                  onChange={(event) => setFormData({ ...formData, street: event.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="city">City</Label>
                <Input
                  id="city"
                  value={formData.city}
                  onChange={(event) => setFormData({ ...formData, city: event.target.value })}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="country">Country</Label>
                <Input
                  id="country"
                  value={formData.country}
                  onChange={(event) => setFormData({ ...formData, country: event.target.value })}
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="zipCode">Zip Code</Label>
                <Input
                  id="zipCode"
                  value={formData.zipCode}
                  onChange={(event) => setFormData({ ...formData, zipCode: event.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="latitude">Latitude</Label>
                <Input
                  id="latitude"
                  type="number"
                  step="any"
                  value={formData.latitude}
                  onChange={(event) =>
                    setFormData({ ...formData, latitude: Number.parseFloat(event.target.value) || 0 })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="longitude">Longitude</Label>
                <Input
                  id="longitude"
                  type="number"
                  step="any"
                  value={formData.longitude}
                  onChange={(event) =>
                    setFormData({ ...formData, longitude: Number.parseFloat(event.target.value) || 0 })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="checkInTime">Check-in Time</Label>
                <Input
                  id="checkInTime"
                  type="time"
                  value={formData.checkInTime}
                  onChange={(event) =>
                    setFormData({ ...formData, checkInTime: event.target.value })
                  }
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="checkOutTime">Check-out Time</Label>
                <Input
                  id="checkOutTime"
                  type="time"
                  value={formData.checkOutTime}
                  onChange={(event) =>
                    setFormData({ ...formData, checkOutTime: event.target.value })
                  }
                />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="amenities">Amenities</Label>
                <Input
                  id="amenities"
                  value={formData.amenities}
                  onChange={(event) => setFormData({ ...formData, amenities: event.target.value })}
                  placeholder="wifi, parking, breakfast"
                />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="images">Image URLs</Label>
                <Input
                  id="images"
                  value={formData.images}
                  onChange={(event) => setFormData({ ...formData, images: event.target.value })}
                  placeholder="https://example.com/room.jpg, https://example.com/lobby.jpg"
                />
              </div>
            </div>

            <div>
              <div className="mb-3 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Bed className="h-4 w-4 text-muted-foreground" />
                  <h3 className="text-sm font-semibold">Room Types</h3>
                </div>
                {!selectedItem && (
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => setFormData({ ...formData, rooms: [...formData.rooms, { ...emptyRoom }] })}
                  >
                    <Plus className="mr-1 h-4 w-4" />
                    Add Room
                  </Button>
                )}
              </div>
              <div className="space-y-3">
                {formData.rooms.map((room, index) => (
                  <div key={index} className="rounded-md border p-3">
                    <div className="grid gap-3 sm:grid-cols-3">
                      <div className="space-y-2">
                        <Label>Name</Label>
                        <Input
                          value={room.name}
                          disabled={Boolean(selectedItem)}
                          onChange={(event) => handleUpdateRoom(index, 'name', event.target.value)}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label>Base Price</Label>
                        <Input
                          type="number"
                          min="0"
                          value={room.basePrice}
                          disabled={Boolean(selectedItem)}
                          onChange={(event) =>
                            handleUpdateRoom(index, 'basePrice', Number.parseInt(event.target.value, 10) || 0)
                          }
                        />
                      </div>
                      <div className="space-y-2">
                        <Label>Max Guests</Label>
                        <Input
                          type="number"
                          min="1"
                          value={room.maxGuests}
                          disabled={Boolean(selectedItem)}
                          onChange={(event) =>
                            handleUpdateRoom(index, 'maxGuests', Number.parseInt(event.target.value, 10) || 1)
                          }
                        />
                      </div>
                      <div className="space-y-2">
                        <Label>Bedrooms</Label>
                        <Input
                          type="number"
                          min="0"
                          value={room.bedrooms}
                          disabled={Boolean(selectedItem)}
                          onChange={(event) =>
                            handleUpdateRoom(index, 'bedrooms', Number.parseInt(event.target.value, 10) || 0)
                          }
                        />
                      </div>
                      <div className="space-y-2">
                        <Label>Bathrooms</Label>
                        <Input
                          type="number"
                          min="0"
                          value={room.bathrooms}
                          disabled={Boolean(selectedItem)}
                          onChange={(event) =>
                            handleUpdateRoom(index, 'bathrooms', Number.parseInt(event.target.value, 10) || 0)
                          }
                        />
                      </div>
                      <div className="flex items-center gap-3 rounded-md border p-3">
                        <Switch
                          checked={Boolean(room.isEntirePlace)}
                          disabled={Boolean(selectedItem)}
                          onCheckedChange={(checked) =>
                            handleUpdateRoom(index, 'isEntirePlace', checked)
                          }
                        />
                        <Label>Entire Place</Label>
                      </div>
                      <div className="space-y-2 sm:col-span-3">
                        <Label>Description</Label>
                        <Textarea
                          value={room.description}
                          disabled={Boolean(selectedItem)}
                          onChange={(event) =>
                            handleUpdateRoom(index, 'description', event.target.value)
                          }
                          rows={2}
                        />
                      </div>
                    </div>
                    {!selectedItem && formData.rooms.length > 1 && (
                      <Button
                        type="button"
                        variant="ghost"
                        className="mt-3 text-destructive hover:text-destructive"
                        onClick={() =>
                          setFormData({
                            ...formData,
                            rooms: formData.rooms.filter((_, roomIndex) => roomIndex !== index),
                          })
                        }
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        Remove Room
                      </Button>
                    )}
                  </div>
                ))}
              </div>
            </div>

            <div className="flex justify-end gap-2 border-t pt-4">
              <Button type="button" variant="outline" onClick={() => setIsFormOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Saving...' : selectedItem ? 'Update Accommodation' : 'Create Accommodation'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={isAvailabilityFormOpen} onOpenChange={setIsAvailabilityFormOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <CalendarClock className="h-5 w-5" />
              Add Accommodation Availability
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleAvailabilitySubmit} className="space-y-4">
            <div className="space-y-2">
              <Label>Accommodation</Label>
              <Select
                value={availabilityForm.accommodationID}
                onValueChange={(value) => {
                  const accommodation = accommodationById.get(value);
                  const firstRoom = getRooms(accommodation)[0];
                  setAvailabilityForm({
                    ...availabilityForm,
                    accommodationID: value,
                    roomTypeID: firstRoom?.id || '',
                    pricePerNight: firstRoom?.basePrice || availabilityForm.pricePerNight,
                  });
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select accommodation" />
                </SelectTrigger>
                <SelectContent>
                  {accommodationList.map((accommodation) => (
                    <SelectItem key={accommodation.id} value={accommodation.id}>
                      {accommodation.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Room Type</Label>
              <Select
                value={availabilityForm.roomTypeID}
                onValueChange={(value) => {
                  const room = selectedAvailabilityRooms.find((item) => item.id === value);
                  setAvailabilityForm({
                    ...availabilityForm,
                    roomTypeID: value,
                    pricePerNight: room?.basePrice || availabilityForm.pricePerNight,
                  });
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select room type" />
                </SelectTrigger>
                <SelectContent>
                  {selectedAvailabilityRooms.map((room) => (
                    <SelectItem key={room.id} value={room.id}>
                      {room.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="space-y-2">
                <Label htmlFor="date">Date</Label>
                <Input
                  id="date"
                  type="date"
                  value={availabilityForm.date}
                  onChange={(event) =>
                    setAvailabilityForm({ ...availabilityForm, date: event.target.value })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="totalRooms">Total Rooms</Label>
                <Input
                  id="totalRooms"
                  type="number"
                  min="1"
                  value={availabilityForm.totalRooms}
                  onChange={(event) =>
                    setAvailabilityForm({
                      ...availabilityForm,
                      totalRooms: Number.parseInt(event.target.value, 10) || 1,
                    })
                  }
                  required
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="pricePerNight">Price Per Night</Label>
                <Input
                  id="pricePerNight"
                  type="number"
                  min="1"
                  value={availabilityForm.pricePerNight}
                  onChange={(event) =>
                    setAvailabilityForm({
                      ...availabilityForm,
                      pricePerNight: Number.parseInt(event.target.value, 10) || 0,
                    })
                  }
                  required
                />
              </div>
            </div>
            <div className="flex justify-end gap-2 border-t pt-4">
              <Button type="button" variant="outline" onClick={() => setIsAvailabilityFormOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading || !availabilityForm.roomTypeID}>
                {isLoading ? 'Saving...' : 'Create Availability'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={isDetailsOpen} onOpenChange={setIsDetailsOpen}>
        <DialogContent className="max-h-[90vh] max-w-3xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Eye className="h-5 w-5" />
              Accommodation Details
            </DialogTitle>
          </DialogHeader>
          {selectedDetails && (
            <div className="space-y-5">
              <div>
                <h2 className="text-xl font-semibold">{selectedDetails.name}</h2>
                <p className="text-sm text-muted-foreground">
                  {selectedDetails.address?.street}, {selectedDetails.address?.city},{' '}
                  {selectedDetails.address?.country}
                </p>
              </div>
              <p className="text-sm leading-relaxed text-muted-foreground">
                {selectedDetails.description || 'No description provided.'}
              </p>
              <div className="grid gap-3 sm:grid-cols-3">
                <div className="rounded-md border p-3">
                  <p className="text-xs text-muted-foreground">Property Type</p>
                  <p className="font-medium">{selectedDetails.type || selectedDetails.propertyType}</p>
                </div>
                <div className="rounded-md border p-3">
                  <p className="text-xs text-muted-foreground">Check-in</p>
                  <p className="font-medium">{toTimeInput(selectedDetails.checkInTime) || 'Not set'}</p>
                </div>
                <div className="rounded-md border p-3">
                  <p className="text-xs text-muted-foreground">Check-out</p>
                  <p className="font-medium">{toTimeInput(selectedDetails.checkOutTime) || 'Not set'}</p>
                </div>
              </div>
              <div>
                <h3 className="mb-2 text-sm font-semibold">Amenities</h3>
                <div className="flex flex-wrap gap-2">
                  {(selectedDetails.amenities || []).map((amenity) => (
                    <span key={amenity} className="rounded-full bg-muted px-2 py-1 text-xs">
                      {amenity}
                    </span>
                  ))}
                </div>
              </div>
              <div>
                <h3 className="mb-2 text-sm font-semibold">Room Types</h3>
                <div className="space-y-2">
                  {getRooms(selectedDetails).map((room) => (
                    <div key={room.id || room.name} className="rounded-md border p-3">
                      <div className="flex items-start justify-between gap-3">
                        <div>
                          <p className="font-medium">{room.name}</p>
                          <p className="text-sm text-muted-foreground">{room.description}</p>
                        </div>
                        <p className="font-semibold text-primary">
                          USD {room.basePrice.toLocaleString()}
                        </p>
                      </div>
                      <p className="mt-2 text-xs text-muted-foreground">
                        {room.bedrooms} bed, {room.bathrooms} bath, {room.maxGuests} guests
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <DeleteDialog
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
        onConfirm={handleConfirmDelete}
        title="Delete Accommodation"
        description={`Are you sure you want to delete "${selectedItem?.name}"? This action cannot be undone.`}
        isLoading={isLoading}
      />
    </div>
  );
}
