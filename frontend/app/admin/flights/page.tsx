'use client';

import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'next/navigation';
import {
  BadgeDollarSign,
  CalendarClock,
  Luggage,
  Plane,
  Plus,
  Route,
  Tag,
  Trash2,
  X,
} from 'lucide-react';
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
import { DataTable, Column } from '@/components/admin/data-table';
import { DeleteDialog } from '@/components/admin/delete-dialog';
import { useAdminStore } from '@/lib/stores/admin-store';
import type { BaggageAllowance, Flight, FlightCabinClass, Offer } from '@/lib/types';

const CABIN_CLASS_OPTIONS = [
  { value: 'economy', label: 'Economy' },
  { value: 'business', label: 'Business' },
  { value: 'first', label: 'First Class' },
] as const;

const FLIGHT_STATUSES = ['scheduled', 'delayed', 'cancelled', 'departed', 'arrived'] as const;

type FlightStatus = (typeof FLIGHT_STATUSES)[number];

interface FlightFormData {
  flightNumber: string;
  originID: string;
  destinationID: string;
  departureTime: string;
  arrivalTime: string;
  airlineID: string;
  cabinClass: FlightCabinClass[];
  segments: string[];
  planeID: string;
  stops: number;
  status: FlightStatus;
}

interface OfferFormData {
  flightID: string;
  providerReference: string;
  provider: string;
  oneway: boolean;
  segments: string[];
  priceTotal: number;
  baggageAllowance: BaggageAllowance[];
  lastTicketingDate: string;
  bookableSeats: number;
  expiresAt: string;
}

const emptyFlightForm: FlightFormData = {
  flightNumber: '',
  originID: '',
  destinationID: '',
  departureTime: '',
  arrivalTime: '',
  airlineID: '',
  cabinClass: [],
  segments: [],
  planeID: '',
  stops: 0,
  status: 'scheduled',
};

const emptyOfferForm: OfferFormData = {
  flightID: '',
  providerReference: '',
  provider: '',
  oneway: true,
  segments: [],
  priceTotal: 0,
  baggageAllowance: [{ pieces: 1, weightKg: 23 }],
  lastTicketingDate: '',
  bookableSeats: 1,
  expiresAt: '',
};

function toDateTimeLocal(value?: string) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return date.toISOString().slice(0, 16);
}

function toIsoOrUndefined(value: string) {
  return value ? new Date(value).toISOString() : undefined;
}

function getFlightTicketingDates(flight?: Flight) {
  if (!flight?.departureTime) {
    return {
      expiresAt: '',
      lastTicketingDate: '',
    };
  }

  const departureTime = new Date(flight.departureTime);
  const lastTicketingDate = new Date(departureTime.getTime() - 30 * 60 * 1000);

  return {
    expiresAt: toDateTimeLocal(departureTime.toISOString()),
    lastTicketingDate: toDateTimeLocal(lastTicketingDate.toISOString()),
  };
}

function formatDateTime(value?: string) {
  if (!value) return 'Not set';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? 'Invalid date' : format(date, 'MMM d, yyyy HH:mm');
}

function segmentToId(segment: string | { id?: string }) {
  return typeof segment === 'string' ? segment : segment.id || '';
}

function getAvailableSeats(cabin?: FlightCabinClass) {
  if (!cabin) return 0;
  return Math.max(cabin.totalSeats - cabin.reservedSeats, 0);
}

function getFlightFareSummary(flight?: Flight) {
  const cabinClasses = flight?.cabinClass || [];
  const prices = cabinClasses.map((cabin) => cabin.price).filter((price) => price > 0);
  return {
    priceTotal: prices.length > 0 ? Math.min(...prices) : 0,
    bookableSeats: cabinClasses.reduce((sum, cabin) => sum + getAvailableSeats(cabin), 0),
  };
}

export default function AdminFlightsPage() {
  const searchParams = useSearchParams();
  const {
    flights,
    offers,
    isLoading,
    fetchFlights,
    fetchOffers,
    createFlight,
    updateFlight,
    updateFlightStatus,
    deleteFlight,
    createOffer,
    updateOffer,
    updateOfferStatus,
    deleteOffer,
    airlines,
    airports,
    planes,
    fetchAirlines,
    fetchAirports,
    fetchPlanes,
  } = useAdminStore();

  const [activeTab, setActiveTab] = useState('flights');
  const [isFlightFormOpen, setIsFlightFormOpen] = useState(false);
  const [isOfferFormOpen, setIsOfferFormOpen] = useState(false);
  const [isFlightDeleteOpen, setIsFlightDeleteOpen] = useState(false);
  const [isOfferDeleteOpen, setIsOfferDeleteOpen] = useState(false);
  const [selectedFlight, setSelectedFlight] = useState<Flight | null>(null);
  const [selectedOffer, setSelectedOffer] = useState<Offer | null>(null);
  const [flightFormErrors, setFlightFormErrors] = useState<Record<string, string>>({});
  const [offerFormErrors, setOfferFormErrors] = useState<Record<string, string>>({});
  const [flightForm, setFlightForm] = useState<FlightFormData>(emptyFlightForm);
  const [offerForm, setOfferForm] = useState<OfferFormData>(emptyOfferForm);

  const flightList = Array.isArray(flights) ? flights : [];
  const offerList = Array.isArray(offers) ? offers : [];
  const airlineList = Array.isArray(airlines) ? airlines : [];
  const airportList = Array.isArray(airports) ? airports : [];
  const planeList = Array.isArray(planes) ? planes : [];

  const airportById = useMemo(
    () => new Map(airportList.map((airport) => [airport.id, airport])),
    [airportList]
  );
  const airlineById = useMemo(
    () => new Map(airlineList.map((airline) => [airline.id, airline])),
    [airlineList]
  );
  const planeById = useMemo(
    () => new Map(planeList.map((plane) => [plane.id, plane])),
    [planeList]
  );
  const flightById = useMemo(
    () => new Map(flightList.map((flight) => [flight.id, flight])),
    [flightList]
  );

  useEffect(() => {
    fetchFlights();
    fetchOffers();
    fetchAirlines();
    fetchAirports();
    fetchPlanes();
  }, [fetchFlights, fetchOffers, fetchAirlines, fetchAirports, fetchPlanes]);

  useEffect(() => {
    if (searchParams.get('action') === 'new') {
      setActiveTab('flights');
      handleCreateFlight();
    }
  }, [searchParams]);

  const getAirportLabel = (id: string) => {
    const airport = airportById.get(id);
    if (!airport) return id;
    return `${airport.code} ${airport.city || airport.name}`;
  };

  const getFlightLabel = (flight: Flight) => {
    return `${flight.flightNumber} - ${getAirportLabel(flight.originID)} to ${getAirportLabel(
      flight.destinationID
    )}`;
  };

  const flightColumns: Column<Flight>[] = [
    {
      key: 'flightNumber',
      header: 'Flight',
      sortable: true,
      render: (flight) => (
        <div>
          <p className="font-mono font-medium">{flight.flightNumber}</p>
          <p className="text-sm text-muted-foreground">
            {airlineById.get(flight.airlineID)?.name || flight.airlineID}
          </p>
        </div>
      ),
    },
    {
      key: 'route',
      header: 'Route',
      render: (flight) => (
        <span>
          {getAirportLabel(flight.originID)} to {getAirportLabel(flight.destinationID)}
        </span>
      ),
    },
    {
      key: 'departureTime',
      header: 'Schedule',
      sortable: true,
      render: (flight) => (
        <div className="text-sm">
          <p>Departs {formatDateTime(flight.departureTime)}</p>
          <p className="text-muted-foreground">Arrives {formatDateTime(flight.arrivalTime)}</p>
        </div>
      ),
    },
    {
      key: 'planeID',
      header: 'Plane',
      render: (flight) => {
        const plane = planeById.get(flight.planeID);
        return plane ? `${plane.name} ${plane.model}` : flight.planeID;
      },
    },
    {
      key: 'cabinClass',
      header: 'Classes',
      render: (flight) => (
        <div className="flex flex-wrap gap-1">
          {flight.cabinClass.map((cabin, index) => (
            <span
              key={`${cabin.classType}-${index}`}
              className="rounded px-2 py-1 text-xs font-medium bg-primary/10 text-primary"
            >
              {cabin.classType} ${cabin.price.toLocaleString()}
            </span>
          ))}
        </div>
      ),
    },
    {
      key: 'status',
      header: 'Status',
      render: (flight) => (
        <span className="rounded-full bg-muted px-2 py-1 text-xs font-medium capitalize">
          {flight.status}
        </span>
      ),
    },
  ];

  const offerColumns: Column<Offer>[] = [
    {
      key: 'provider',
      header: 'Booking Option',
      sortable: true,
      render: (offer) => (
        <div>
          <p className="font-medium">{offer.provider}</p>
          <p className="text-sm text-muted-foreground">{offer.providerReference}</p>
        </div>
      ),
    },
    {
      key: 'flightID',
      header: 'Flight',
      render: (offer) => {
        const flight = flightById.get(offer.flightID);
        return flight ? getFlightLabel(flight) : offer.flightID;
      },
    },
    {
      key: 'priceTotal',
      header: 'Price',
      sortable: true,
      render: (offer) => (
        <span className="font-semibold text-primary">
          ${offer.priceTotal.toLocaleString()}
        </span>
      ),
    },
    {
      key: 'bookableSeats',
      header: 'Seats',
      sortable: true,
      render: (offer) => <span>{offer.bookableSeats}</span>,
    },
    {
      key: 'expiresAt',
      header: 'Expires',
      render: (offer) => formatDateTime(offer.expiresAt),
    },
    {
      key: 'isActive',
      header: 'Active',
      render: (offer) => (
        <Switch
          checked={offer.isActive}
          onCheckedChange={(checked) => handleToggleOfferStatus(offer, checked)}
          aria-label={`Toggle ${offer.provider} booking option`}
        />
      ),
    },
  ];

  function handleCreateFlight() {
    setSelectedFlight(null);
    setFlightForm(emptyFlightForm);
    setFlightFormErrors({});
    setIsFlightFormOpen(true);
  }

  function handleEditFlight(flight: Flight) {
    setSelectedFlight(flight);
    setFlightForm({
      flightNumber: flight.flightNumber,
      originID: flight.originID,
      destinationID: flight.destinationID,
      departureTime: toDateTimeLocal(flight.departureTime),
      arrivalTime: toDateTimeLocal(flight.arrivalTime),
      airlineID: flight.airlineID,
      cabinClass: flight.cabinClass,
      segments: (flight.segments || []).map(segmentToId).filter(Boolean),
      planeID: flight.planeID,
      stops: flight.stops ?? 0,
      status: FLIGHT_STATUSES.includes(flight.status as FlightStatus)
        ? (flight.status as FlightStatus)
        : 'scheduled',
    });
    setFlightFormErrors({});
    setIsFlightFormOpen(true);
  }

  function handleDeleteFlight(flight: Flight) {
    setSelectedFlight(flight);
    setIsFlightDeleteOpen(true);
  }

  function handleCreateOffer(flight?: Flight) {
    const fareSummary = getFlightFareSummary(flight);
    const ticketingDates = getFlightTicketingDates(flight);
    setSelectedOffer(null);
    setOfferForm({
      ...emptyOfferForm,
      flightID: flight?.id || '',
      segments: flight?.segments ? flight.segments.map(segmentToId).filter(Boolean) : [],
      ...fareSummary,
      ...ticketingDates,
    });
    setOfferFormErrors({});
    setActiveTab('offers');
    setIsOfferFormOpen(true);
  }

  function handleEditOffer(offer: Offer) {
    const flight = flightById.get(offer.flightID);
    const fareSummary = getFlightFareSummary(flight);
    setSelectedOffer(offer);
    setOfferForm({
      flightID: offer.flightID,
      providerReference: offer.providerReference,
      provider: offer.provider,
      oneway: offer.oneway ?? offer.oneWay ?? true,
      segments: (offer.segments || []).map(segmentToId).filter(Boolean),
      priceTotal: fareSummary.priceTotal || offer.priceTotal,
      baggageAllowance:
        offer.baggageAllowance && offer.baggageAllowance.length > 0
          ? offer.baggageAllowance
          : [{ pieces: 1, weightKg: 23 }],
      lastTicketingDate: toDateTimeLocal(offer.lastTicketingDate),
      bookableSeats: fareSummary.bookableSeats || offer.bookableSeats,
      expiresAt: toDateTimeLocal(offer.expiresAt),
    });
    setOfferFormErrors({});
    setIsOfferFormOpen(true);
  }

  function handleDeleteOffer(offer: Offer) {
    setSelectedOffer(offer);
    setIsOfferDeleteOpen(true);
  }

  function validateFlightForm() {
    const errors: Record<string, string> = {};

    if (!flightForm.flightNumber) errors.flightNumber = 'Flight number is required';
    if (!flightForm.originID) errors.originID = 'Origin airport is required';
    if (!flightForm.destinationID) errors.destinationID = 'Destination airport is required';
    if (!flightForm.departureTime) errors.departureTime = 'Departure time is required';
    if (!flightForm.arrivalTime) errors.arrivalTime = 'Arrival time is required';
    if (!flightForm.airlineID) errors.airlineID = 'Airline is required';
    if (!flightForm.planeID) errors.planeID = 'Plane model is required';
    if (!selectedFlight && flightForm.cabinClass.length === 0) {
      errors.cabinClass = 'At least one cabin class is required';
    }

    flightForm.cabinClass.forEach((cabin, index) => {
      if (!cabin.totalSeats || cabin.totalSeats <= 0) {
        errors[`cabin_${index}_seats`] = 'Total seats must be greater than 0';
      }
      if (!cabin.price || cabin.price <= 0) {
        errors[`cabin_${index}_price`] = 'Price must be greater than 0';
      }
    });

    setFlightFormErrors(errors);
    return Object.keys(errors).length === 0;
  }

  function validateOfferForm() {
    const errors: Record<string, string> = {};

    if (!offerForm.flightID) errors.flightID = 'Flight is required';
    if (selectedOfferFlight && selectedOfferFlight.cabinClass.length === 0) {
      errors.flightID = 'Selected flight must have cabin classes';
    }
    if (!offerForm.providerReference) errors.providerReference = 'Provider reference is required';
    if (!offerForm.provider) errors.provider = 'Provider is required';
    if (!offerForm.priceTotal || offerForm.priceTotal <= 0) {
      errors.priceTotal = 'Price must be greater than 0';
    }
    if (!offerForm.bookableSeats || offerForm.bookableSeats <= 0) {
      errors.bookableSeats = 'Bookable seats must be greater than 0';
    }
    const maxSeats = getFlightFareSummary(selectedOfferFlight).bookableSeats;
    if (maxSeats > 0 && offerForm.bookableSeats > maxSeats) {
      errors.bookableSeats = `Bookable seats cannot exceed ${maxSeats}`;
    }
    if (!selectedOffer && !offerForm.expiresAt) errors.expiresAt = 'Expiry date is required';

    offerForm.baggageAllowance.forEach((baggage, index) => {
      if (baggage.pieces < 0) errors[`baggage_${index}_pieces`] = 'Pieces cannot be negative';
      if (baggage.weightKg < 0) errors[`baggage_${index}_weight`] = 'Weight cannot be negative';
    });

    setOfferFormErrors(errors);
    return Object.keys(errors).length === 0;
  }

  function handleAddCabinClass() {
    setFlightForm((current) => ({
      ...current,
      cabinClass: [
        ...current.cabinClass,
        { classType: 'economy', totalSeats: 100, reservedSeats: 0, price: 10000 },
      ],
    }));
  }

  function handleUpdateCabinClass(
    index: number,
    field: keyof FlightCabinClass,
    value: string | number
  ) {
    setFlightForm((current) => {
      const cabinClass = [...current.cabinClass];
      cabinClass[index] = { ...cabinClass[index], [field]: value };
      return { ...current, cabinClass };
    });
  }

  function handleAddFlightSegment() {
    setFlightForm((current) => ({ ...current, segments: [...current.segments, ''] }));
  }

  function handleUpdateFlightSegment(index: number, value: string) {
    setFlightForm((current) => {
      const segments = [...current.segments];
      segments[index] = value;
      return { ...current, segments };
    });
  }

  function handleAddOfferSegment() {
    setOfferForm((current) => ({ ...current, segments: [...current.segments, ''] }));
  }

  function handleUpdateOfferSegment(index: number, value: string) {
    setOfferForm((current) => {
      const segments = [...current.segments];
      segments[index] = value;
      return { ...current, segments };
    });
  }

  function handleAddBaggage() {
    setOfferForm((current) => ({
      ...current,
      baggageAllowance: [...current.baggageAllowance, { pieces: 1, weightKg: 23 }],
    }));
  }

  function handleUpdateBaggage(index: number, field: keyof BaggageAllowance, value: number) {
    setOfferForm((current) => {
      const baggageAllowance = [...current.baggageAllowance];
      baggageAllowance[index] = { ...baggageAllowance[index], [field]: value };
      return { ...current, baggageAllowance };
    });
  }

  function handleOfferFlightChange(flightID: string) {
    const flight = flightById.get(flightID);
    const fareSummary = getFlightFareSummary(flight);
    const ticketingDates = getFlightTicketingDates(flight);

    setOfferForm({
      ...offerForm,
      flightID,
      segments: flight?.segments ? flight.segments.map(segmentToId).filter(Boolean) : [],
      ...fareSummary,
      ...ticketingDates,
    });
  }

  async function handleFlightSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (!validateFlightForm()) {
      toast.error('Please fill in the required flight fields');
      return;
    }

    try {
      const departureTime = new Date(flightForm.departureTime).toISOString();
      const arrivalTime = new Date(flightForm.arrivalTime).toISOString();

      if (selectedFlight) {
        await updateFlight(selectedFlight.id, {
          departureTime,
          arrivalTime,
          stops: flightForm.stops,
        });
        if (selectedFlight.status !== flightForm.status) {
          await updateFlightStatus(selectedFlight.id, flightForm.status);
        }
        toast.success('Flight updated successfully');
      } else {
        await createFlight({
          flightNumber: flightForm.flightNumber,
          originID: flightForm.originID,
          destinationID: flightForm.destinationID,
          departureTime,
          arrivalTime,
          airlineID: flightForm.airlineID,
          cabinClass: flightForm.cabinClass,
          segments: flightForm.segments.filter(Boolean),
          planeID: flightForm.planeID,
          stops: flightForm.stops,
          status: flightForm.status,
        });
        toast.success('Flight created successfully');
      }

      setIsFlightFormOpen(false);
    } catch (error) {
      console.error('Flight save error:', error);
      toast.error('Failed to save flight');
    }
  }

  async function handleOfferSubmit(event: React.FormEvent) {
    event.preventDefault();
    if (!validateOfferForm()) {
      toast.error('Please fill in the required booking option fields');
      return;
    }

    const payload = {
      flightID: offerForm.flightID,
      providerReference: offerForm.providerReference,
      provider: offerForm.provider,
      oneway: offerForm.oneway,
      segments: offerForm.segments.filter(Boolean),
      priceTotal: offerForm.priceTotal,
      baggageAllowance: offerForm.baggageAllowance,
      lastTicketingDate: toIsoOrUndefined(offerForm.lastTicketingDate),
      bookableSeats: offerForm.bookableSeats,
      expiresAt: toIsoOrUndefined(offerForm.expiresAt),
      isActive: selectedOffer?.isActive ?? true,
    };

    try {
      if (selectedOffer) {
        await updateOffer(selectedOffer.id, {
          priceTotal: payload.priceTotal,
          oneway: payload.oneway,
          bookableSeats: payload.bookableSeats,
        });
        toast.success('Booking option updated successfully');
      } else {
        await createOffer(payload as Omit<Offer, 'id'>);
        toast.success('Booking option created successfully');
      }
      setIsOfferFormOpen(false);
    } catch (error) {
      console.error('Booking option save error:', error);
      toast.error('Failed to save booking option');
    }
  }

  async function handleConfirmFlightDelete() {
    if (!selectedFlight) return;
    try {
      await deleteFlight(selectedFlight.id);
      toast.success('Flight deleted successfully');
      setIsFlightDeleteOpen(false);
    } catch {
      toast.error('Failed to delete flight');
    }
  }

  async function handleConfirmOfferDelete() {
    if (!selectedOffer) return;
    try {
      await deleteOffer(selectedOffer.id);
      toast.success('Booking option deleted successfully');
      setIsOfferDeleteOpen(false);
    } catch {
      toast.error('Failed to delete booking option');
    }
  }

  async function handleToggleOfferStatus(offer: Offer, checked: boolean) {
    try {
      await updateOfferStatus(offer.id, checked);
      toast.success(`Booking option ${checked ? 'activated' : 'deactivated'}`);
    } catch {
      toast.error('Failed to update booking option status');
    }
  }

  const selectedPlane = planeById.get(flightForm.planeID);
  const activeOffers = offerList.filter((offer) => offer.isActive).length;
  const selectedOfferFlight = flightById.get(offerForm.flightID);
  const selectedOfferFareSummary = getFlightFareSummary(selectedOfferFlight);
  const maxBookableSeats = selectedOfferFareSummary.bookableSeats;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Flights</h1>
          <p className="mt-1 text-muted-foreground">
            Manage flight schedules and the booking options customers see
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" onClick={() => handleCreateOffer()} className="gap-2">
            <Tag className="h-4 w-4" />
            Add Booking Option
          </Button>
          <Button onClick={handleCreateFlight} className="gap-2">
            <Plus className="h-4 w-4" />
            Add Flight
          </Button>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Total Flights
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{flightList.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Scheduled
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-primary">
              {flightList.filter((flight) => flight.status === 'scheduled').length}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Booking Options
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{offerList.length}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Active Options
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600">{activeOffers}</div>
          </CardContent>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-4">
        <TabsList>
          <TabsTrigger value="flights" className="gap-2">
            <Plane className="h-4 w-4" />
            Flights
          </TabsTrigger>
          <TabsTrigger value="offers" className="gap-2">
            <BadgeDollarSign className="h-4 w-4" />
            Booking Options
          </TabsTrigger>
        </TabsList>

        <TabsContent value="flights" className="space-y-4">
          {isLoading && flightList.length === 0 ? (
            <div className="flex h-64 items-center justify-center">
              <Spinner className="h-8 w-8 text-primary" />
            </div>
          ) : (
            <DataTable
              data={flightList}
              columns={flightColumns}
              searchPlaceholder="Search flights..."
              searchKey="flightNumber"
              onView={handleCreateOffer}
              onEdit={handleEditFlight}
              onDelete={handleDeleteFlight}
            />
          )}
        </TabsContent>

        <TabsContent value="offers" className="space-y-4">
          {isLoading && offerList.length === 0 ? (
            <div className="flex h-64 items-center justify-center">
              <Spinner className="h-8 w-8 text-primary" />
            </div>
          ) : (
            <DataTable
              data={offerList}
              columns={offerColumns}
              searchPlaceholder="Search booking options..."
              searchKey="provider"
              onEdit={handleEditOffer}
              onDelete={handleDeleteOffer}
            />
          )}
        </TabsContent>
      </Tabs>

      <Dialog open={isFlightFormOpen} onOpenChange={setIsFlightFormOpen}>
        <DialogContent className="max-h-[90vh] max-w-4xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Plane className="h-5 w-5" />
              {selectedFlight ? 'Edit Flight' : 'Create Flight'}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleFlightSubmit} className="space-y-6">
            {selectedFlight && (
              <div className="rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground">
                This backend update route accepts schedule, stops, and status changes for existing
                flights. Create a new flight when route, airline, plane, or cabin inventory changes.
              </div>
            )}

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="flightNumber">Flight Number</Label>
                <Input
                  id="flightNumber"
                  value={flightForm.flightNumber}
                  disabled={Boolean(selectedFlight)}
                  onChange={(event) =>
                    setFlightForm({ ...flightForm, flightNumber: event.target.value })
                  }
                  className={flightFormErrors.flightNumber ? 'border-red-500' : ''}
                  placeholder="KQ101"
                />
                {flightFormErrors.flightNumber && (
                  <p className="text-xs text-red-500">{flightFormErrors.flightNumber}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label>Status</Label>
                <Select
                  value={flightForm.status}
                  onValueChange={(value) =>
                    setFlightForm({ ...flightForm, status: value as FlightStatus })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {FLIGHT_STATUSES.map((status) => (
                      <SelectItem key={status} value={status}>
                        {status}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Airline</Label>
                <Select
                  value={flightForm.airlineID}
                  disabled={Boolean(selectedFlight)}
                  onValueChange={(value) => setFlightForm({ ...flightForm, airlineID: value })}
                >
                  <SelectTrigger className={flightFormErrors.airlineID ? 'border-red-500' : ''}>
                    <SelectValue placeholder="Select airline" />
                  </SelectTrigger>
                  <SelectContent>
                    {airlineList.map((airline) => (
                      <SelectItem value={airline.id} key={airline.id}>
                        {airline.name} {airline.code ? `(${airline.code})` : ''}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {flightFormErrors.airlineID && (
                  <p className="text-xs text-red-500">{flightFormErrors.airlineID}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label>Plane Model</Label>
                <Select
                  value={flightForm.planeID}
                  disabled={Boolean(selectedFlight)}
                  onValueChange={(value) => {
                    const plane = planeById.get(value);
                    setFlightForm({
                      ...flightForm,
                      planeID: value,
                      cabinClass: plane
                        ? plane.cabinClasses.map((cabin) => ({
                            classType: cabin.type,
                            totalSeats: cabin.seatCapacity,
                            reservedSeats: 0,
                            price: 10000,
                          }))
                        : flightForm.cabinClass,
                    });
                  }}
                >
                  <SelectTrigger className={flightFormErrors.planeID ? 'border-red-500' : ''}>
                    <SelectValue placeholder="Select plane model" />
                  </SelectTrigger>
                  <SelectContent>
                    {planeList.map((plane) => (
                      <SelectItem value={plane.id} key={plane.id}>
                        {plane.name} {plane.model}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {flightFormErrors.planeID && (
                  <p className="text-xs text-red-500">{flightFormErrors.planeID}</p>
                )}
                {selectedPlane && !selectedFlight && (
                  <p className="text-xs text-muted-foreground">
                    {selectedPlane.cabinClasses.length} cabin class(es) loaded from this model.
                  </p>
                )}
              </div>
            </div>

            <div>
              <div className="mb-3 flex items-center gap-2">
                <Route className="h-4 w-4 text-muted-foreground" />
                <h3 className="text-sm font-semibold">Route</h3>
              </div>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label>Origin</Label>
                  <Select
                    value={flightForm.originID}
                    disabled={Boolean(selectedFlight)}
                    onValueChange={(value) => setFlightForm({ ...flightForm, originID: value })}
                  >
                    <SelectTrigger className={flightFormErrors.originID ? 'border-red-500' : ''}>
                      <SelectValue placeholder="Select origin airport" />
                    </SelectTrigger>
                    <SelectContent>
                      {airportList.map((airport) => (
                        <SelectItem value={airport.id} key={airport.id}>
                          {airport.code} - {airport.city || airport.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {flightFormErrors.originID && (
                    <p className="text-xs text-red-500">{flightFormErrors.originID}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label>Destination</Label>
                  <Select
                    value={flightForm.destinationID}
                    disabled={Boolean(selectedFlight)}
                    onValueChange={(value) =>
                      setFlightForm({ ...flightForm, destinationID: value })
                    }
                  >
                    <SelectTrigger
                      className={flightFormErrors.destinationID ? 'border-red-500' : ''}
                    >
                      <SelectValue placeholder="Select destination airport" />
                    </SelectTrigger>
                    <SelectContent>
                      {airportList.map((airport) => (
                        <SelectItem value={airport.id} key={airport.id}>
                          {airport.code} - {airport.city || airport.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {flightFormErrors.destinationID && (
                    <p className="text-xs text-red-500">{flightFormErrors.destinationID}</p>
                  )}
                </div>
              </div>
            </div>

            <div>
              <div className="mb-3 flex items-center gap-2">
                <CalendarClock className="h-4 w-4 text-muted-foreground" />
                <h3 className="text-sm font-semibold">Schedule</h3>
              </div>
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label htmlFor="departureTime">Departure Time</Label>
                  <Input
                    id="departureTime"
                    type="datetime-local"
                    value={flightForm.departureTime}
                    onChange={(event) =>
                      setFlightForm({ ...flightForm, departureTime: event.target.value })
                    }
                    className={flightFormErrors.departureTime ? 'border-red-500' : ''}
                  />
                  {flightFormErrors.departureTime && (
                    <p className="text-xs text-red-500">{flightFormErrors.departureTime}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="arrivalTime">Arrival Time</Label>
                  <Input
                    id="arrivalTime"
                    type="datetime-local"
                    value={flightForm.arrivalTime}
                    onChange={(event) =>
                      setFlightForm({ ...flightForm, arrivalTime: event.target.value })
                    }
                    className={flightFormErrors.arrivalTime ? 'border-red-500' : ''}
                  />
                  {flightFormErrors.arrivalTime && (
                    <p className="text-xs text-red-500">{flightFormErrors.arrivalTime}</p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label htmlFor="stops">Stops</Label>
                  <Input
                    id="stops"
                    type="number"
                    min="0"
                    value={flightForm.stops}
                    onChange={(event) =>
                      setFlightForm({
                        ...flightForm,
                        stops: Number.parseInt(event.target.value, 10) || 0,
                      })
                    }
                  />
                </div>
              </div>
            </div>

            <div>
              <div className="mb-3 flex items-center justify-between">
                <h3 className="text-sm font-semibold">Cabin Classes</h3>
                {!selectedFlight && (
                  <Button type="button" size="sm" variant="outline" onClick={handleAddCabinClass}>
                    <Plus className="mr-1 h-4 w-4" />
                    Add Class
                  </Button>
                )}
              </div>
              {flightFormErrors.cabinClass && (
                <p className="mb-3 text-xs text-red-500">{flightFormErrors.cabinClass}</p>
              )}
              <div className="space-y-3">
                {flightForm.cabinClass.map((cabin, index) => (
                  <div key={`${cabin.classType}-${index}`} className="rounded-md border p-4">
                    <div className="mb-3 flex items-center justify-between">
                      <h4 className="text-sm font-medium">Cabin Class {index + 1}</h4>
                      {!selectedFlight && (
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          onClick={() =>
                            setFlightForm({
                              ...flightForm,
                              cabinClass: flightForm.cabinClass.filter((_, i) => i !== index),
                            })
                          }
                        >
                          <Trash2 className="h-4 w-4 text-red-500" />
                        </Button>
                      )}
                    </div>
                    <div className="grid gap-3 sm:grid-cols-4">
                      <div className="space-y-2">
                        <Label>Class Type</Label>
                        <Select
                          value={cabin.classType}
                          disabled={Boolean(selectedFlight)}
                          onValueChange={(value) =>
                            handleUpdateCabinClass(
                              index,
                              'classType',
                              value as FlightCabinClass['classType']
                            )
                          }
                        >
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {CABIN_CLASS_OPTIONS.map((option) => (
                              <SelectItem key={option.value} value={option.value}>
                                {option.label}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="space-y-2">
                        <Label>Total Seats</Label>
                        <Input
                          type="number"
                          min="1"
                          disabled={Boolean(selectedFlight)}
                          value={cabin.totalSeats}
                          onChange={(event) =>
                            handleUpdateCabinClass(
                              index,
                              'totalSeats',
                              Number.parseInt(event.target.value, 10) || 0
                            )
                          }
                          className={flightFormErrors[`cabin_${index}_seats`] ? 'border-red-500' : ''}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label>Reserved Seats</Label>
                        <Input type="number" value={cabin.reservedSeats} disabled />
                      </div>
                      <div className="space-y-2">
                        <Label>Price</Label>
                        <Input
                          type="number"
                          min="1"
                          disabled={Boolean(selectedFlight)}
                          value={cabin.price}
                          onChange={(event) =>
                            handleUpdateCabinClass(
                              index,
                              'price',
                              Number.parseInt(event.target.value, 10) || 0
                            )
                          }
                          className={flightFormErrors[`cabin_${index}_price`] ? 'border-red-500' : ''}
                        />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div>
              <div className="mb-3 flex items-center justify-between">
                <h3 className="text-sm font-semibold">Segments</h3>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  disabled={Boolean(selectedFlight)}
                  onClick={handleAddFlightSegment}
                >
                  <Plus className="mr-1 h-4 w-4" />
                  Add Segment
                </Button>
              </div>
              <div className="space-y-2">
                {flightForm.segments.map((segment, index) => (
                  <div key={index} className="flex gap-2">
                    <Input
                      value={segment}
                      disabled={Boolean(selectedFlight)}
                      onChange={(event) => handleUpdateFlightSegment(index, event.target.value)}
                      placeholder="Segment ObjectID"
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      disabled={Boolean(selectedFlight)}
                      onClick={() =>
                        setFlightForm({
                          ...flightForm,
                          segments: flightForm.segments.filter((_, i) => i !== index),
                        })
                      }
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex justify-end gap-2 border-t pt-4">
              <Button type="button" variant="outline" onClick={() => setIsFlightFormOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Saving...' : selectedFlight ? 'Update Flight' : 'Create Flight'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={isOfferFormOpen} onOpenChange={setIsOfferFormOpen}>
        <DialogContent className="max-h-[90vh] max-w-3xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Tag className="h-5 w-5" />
              {selectedOffer ? 'Edit Booking Option' : 'Create Booking Option'}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleOfferSubmit} className="space-y-6">
            {selectedOffer && (
              <div className="rounded-md border bg-muted/30 p-3 text-sm text-muted-foreground">
                Existing booking options can update price, trip type, and bookable seats through
                the current backend route.
              </div>
            )}

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2 sm:col-span-2">
                <Label>Flight</Label>
                <Select
                  value={offerForm.flightID}
                  disabled={Boolean(selectedOffer)}
                  onValueChange={handleOfferFlightChange}
                >
                  <SelectTrigger className={offerFormErrors.flightID ? 'border-red-500' : ''}>
                    <SelectValue placeholder="Select flight" />
                  </SelectTrigger>
                  <SelectContent>
                    {flightList.map((flight) => (
                      <SelectItem value={flight.id} key={flight.id}>
                        {getFlightLabel(flight)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {selectedOfferFlight && (
                  <p className="text-xs text-muted-foreground">
                    Departs {formatDateTime(selectedOfferFlight.departureTime)}.
                  </p>
                )}
                {offerFormErrors.flightID && (
                  <p className="text-xs text-red-500">{offerFormErrors.flightID}</p>
                )}
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label>Cabin Classes</Label>
                {selectedOfferFlight?.cabinClass?.length ? (
                  <div className="grid gap-3 sm:grid-cols-3">
                    {selectedOfferFlight.cabinClass.map((cabin, index) => {
                      const availableSeats = getAvailableSeats(cabin);

                      return (
                        <div
                          key={`${cabin.classType}-${index}`}
                          className="rounded-md border border-border bg-background p-3"
                        >
                          <span className="block text-sm font-semibold capitalize">
                            {cabin.classType}
                          </span>
                          <span className="mt-1 block text-sm">
                            ${cabin.price.toLocaleString()}
                          </span>
                          <span className="mt-2 block text-xs text-muted-foreground">
                            {availableSeats} bookable of {cabin.totalSeats} seats
                          </span>
                          <span className="block text-xs text-muted-foreground">
                            {cabin.reservedSeats} reserved
                          </span>
                        </div>
                      );
                    })}
                  </div>
                ) : null}
                {selectedOfferFlight && selectedOfferFlight.cabinClass.length === 0 && (
                  <p className="text-xs text-muted-foreground">
                    This flight does not have cabin pricing yet.
                  </p>
                )}
                {!selectedOfferFlight && (
                  <p className="text-xs text-muted-foreground">
                    Select a flight to view its cabin classes.
                  </p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="provider">Provider</Label>
                <Input
                  id="provider"
                  value={offerForm.provider}
                  disabled={Boolean(selectedOffer)}
                  onChange={(event) => setOfferForm({ ...offerForm, provider: event.target.value })}
                  className={offerFormErrors.provider ? 'border-red-500' : ''}
                  placeholder="Amadeus"
                />
                {offerFormErrors.provider && (
                  <p className="text-xs text-red-500">{offerFormErrors.provider}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="providerReference">Provider Reference</Label>
                <Input
                  id="providerReference"
                  value={offerForm.providerReference}
                  disabled={Boolean(selectedOffer)}
                  onChange={(event) =>
                    setOfferForm({ ...offerForm, providerReference: event.target.value })
                  }
                  className={offerFormErrors.providerReference ? 'border-red-500' : ''}
                  placeholder="OFF-2026-001"
                />
                {offerFormErrors.providerReference && (
                  <p className="text-xs text-red-500">{offerFormErrors.providerReference}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="priceTotal">Total Price</Label>
                <Input
                  id="priceTotal"
                  type="number"
                  min="1"
                  value={offerForm.priceTotal}
                  readOnly
                  className={offerFormErrors.priceTotal ? 'border-red-500' : ''}
                />
                <p className="text-xs text-muted-foreground">
                  Filled from the selected flight cabin class.
                </p>
                {offerFormErrors.priceTotal && (
                  <p className="text-xs text-red-500">{offerFormErrors.priceTotal}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="bookableSeats">Bookable Seats</Label>
                <Input
                  id="bookableSeats"
                  type="number"
                  min="1"
                  value={offerForm.bookableSeats}
                  readOnly
                  className={offerFormErrors.bookableSeats ? 'border-red-500' : ''}
                />
                <p className="text-xs text-muted-foreground">
                  Sum of all available seats across the flight cabin classes.
                </p>
                {offerFormErrors.bookableSeats && (
                  <p className="text-xs text-red-500">{offerFormErrors.bookableSeats}</p>
                )}
              </div>
              <div className="flex items-center gap-3 rounded-md border p-3">
                <Switch
                  checked={offerForm.oneway}
                  onCheckedChange={(checked) => setOfferForm({ ...offerForm, oneway: checked })}
                  aria-label="One way booking option"
                />
                <div>
                  <Label>One Way</Label>
                  <p className="text-xs text-muted-foreground">
                    Turn off for a return or multi-leg offer.
                  </p>
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="expiresAt">Expires At</Label>
                <Input
                  id="expiresAt"
                  type="datetime-local"
                  value={offerForm.expiresAt}
                  readOnly
                  className={offerFormErrors.expiresAt ? 'border-red-500' : ''}
                />
                <p className="text-xs text-muted-foreground">
                  Filled from the selected flight departure time.
                </p>
                {offerFormErrors.expiresAt && (
                  <p className="text-xs text-red-500">{offerFormErrors.expiresAt}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastTicketingDate">Last Ticketing Date</Label>
                <Input
                  id="lastTicketingDate"
                  type="datetime-local"
                  value={offerForm.lastTicketingDate}
                  readOnly
                />
                <p className="text-xs text-muted-foreground">
                  Automatically set to 30 minutes before departure.
                </p>
              </div>
            </div>

            <div>
              <div className="mb-3 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <Luggage className="h-4 w-4 text-muted-foreground" />
                  <h3 className="text-sm font-semibold">Baggage Allowance</h3>
                </div>
                {!selectedOffer && (
                  <Button type="button" size="sm" variant="outline" onClick={handleAddBaggage}>
                    <Plus className="mr-1 h-4 w-4" />
                    Add Baggage
                  </Button>
                )}
              </div>
              <div className="space-y-3">
                {offerForm.baggageAllowance.map((baggage, index) => (
                  <div key={index} className="grid gap-3 rounded-md border p-3 sm:grid-cols-[1fr_1fr_auto]">
                    <div className="space-y-2">
                      <Label>Pieces</Label>
                      <Input
                        type="number"
                        min="0"
                        disabled={Boolean(selectedOffer)}
                        value={baggage.pieces}
                        onChange={(event) =>
                          handleUpdateBaggage(
                            index,
                            'pieces',
                            Number.parseInt(event.target.value, 10) || 0
                          )
                        }
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Weight Kg</Label>
                      <Input
                        type="number"
                        min="0"
                        disabled={Boolean(selectedOffer)}
                        value={baggage.weightKg}
                        onChange={(event) =>
                          handleUpdateBaggage(
                            index,
                            'weightKg',
                            Number.parseInt(event.target.value, 10) || 0
                          )
                        }
                      />
                    </div>
                    {!selectedOffer && (
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="self-end"
                        onClick={() =>
                          setOfferForm({
                            ...offerForm,
                            baggageAllowance: offerForm.baggageAllowance.filter(
                              (_, i) => i !== index
                            ),
                          })
                        }
                      >
                        <Trash2 className="h-4 w-4 text-red-500" />
                      </Button>
                    )}
                  </div>
                ))}
              </div>
            </div>

            <div>
              <div className="mb-3 flex items-center justify-between">
                <h3 className="text-sm font-semibold">Booking Segments</h3>
                {!selectedOffer && (
                  <Button type="button" size="sm" variant="outline" onClick={handleAddOfferSegment}>
                    <Plus className="mr-1 h-4 w-4" />
                    Add Segment
                  </Button>
                )}
              </div>
              <div className="space-y-2">
                {offerForm.segments.map((segment, index) => (
                  <div key={index} className="flex gap-2">
                    <Input
                      value={segment}
                      disabled={Boolean(selectedOffer)}
                      onChange={(event) => handleUpdateOfferSegment(index, event.target.value)}
                      placeholder="Segment ObjectID"
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      disabled={Boolean(selectedOffer)}
                      onClick={() =>
                        setOfferForm({
                          ...offerForm,
                          segments: offerForm.segments.filter((_, i) => i !== index),
                        })
                      }
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex justify-end gap-2 border-t pt-4">
              <Button type="button" variant="outline" onClick={() => setIsOfferFormOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading
                  ? 'Saving...'
                  : selectedOffer
                    ? 'Update Booking Option'
                    : 'Create Booking Option'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <DeleteDialog
        open={isFlightDeleteOpen}
        onOpenChange={setIsFlightDeleteOpen}
        onConfirm={handleConfirmFlightDelete}
        title="Delete Flight"
        description={`Are you sure you want to delete flight ${selectedFlight?.flightNumber}? This action cannot be undone.`}
        isLoading={isLoading}
      />

      <DeleteDialog
        open={isOfferDeleteOpen}
        onOpenChange={setIsOfferDeleteOpen}
        onConfirm={handleConfirmOfferDelete}
        title="Delete Booking Option"
        description={`Are you sure you want to delete booking option ${selectedOffer?.providerReference}? This action cannot be undone.`}
        isLoading={isLoading}
      />
    </div>
  );
}
