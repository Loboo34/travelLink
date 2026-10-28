// User types
export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  gender?: string;
  dateOfBirth?: string;
  nationality?: string;
  phoneNumber?: string;
  profilePic?: string;
  role: string;
  isActive: boolean;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Address type
export interface Address {
  street: string;
  city: string;
  country: string;
  zipCode: string;
}

// Location type
export interface Location {
  latitude: number;
  longitude: number;
  type?: string;
  coordinates?: number[];
}

// Accommodation types
export interface RoomType {
  id: string;
  name: string;
  description: string;
  bedrooms: number;
  bathrooms: number;
  maxGuests: number;
  isEntirePlace?: boolean;
  amenities: string[];
  images: string[];
  basePrice: number;
}

export interface Accommodation {
  id: string;
  hostID?: string;
  name: string;
  description: string;
  type: string;
  propertyType?: string;
  address: Address;
  location: Location;
  amenities: string[];
  images: string[];
  roomTypes: RoomType[];
  roomType?: RoomType[];
  rating: number;
  reviews: number;
  checkInTime: string | number;
  checkOutTime: string | number;
  isActive?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface AccommodationAvailability {
  id: string;
  accommodationID: string;
  roomTypeID: string;
  date: string;
  totalRooms: number;
  reservedRooms: number;
  pricePerNight: number;
  currency?: string;
  isActive: boolean;
  createdAt?: string;
  updatedAt?: string;
}

// Flight types
export interface FlightCabinClass {
  id?: string;
  classType: 'economy' | 'business' | 'first';
  totalSeats: number;
  reservedSeats: number;
  price: number;
}

export interface PlaneCabinClass {
  name: string;
  type: 'economy' | 'business' | 'first';
  seatCapacity: number;
}

export interface FlightSegment {
  id?: string;
  originAirportID?: string;
  destinationAirportID?: string;
  departureAirport?: string;
  arrivalAirport?: string;
  departureTime: string;
  arrivalTime: string;
  flightNumber: string;
  airline?: string;
  airlineID?: string;
}

export interface PlaneModel {
  id: string;
  name: string;
  model: string;
  cabinClasses: PlaneCabinClass[];
}

export interface Flight {
  id: string;
  originID: string;
  destinationID: string;
  departureTime: string;
  arrivalTime: string;
  airlineID: string;
  flightNumber: string;
  cabinClass: FlightCabinClass[];
  stops: number;
  segments: Array<string | FlightSegment>; // Backend usually returns segment IDs.
  planeID: string;
  status: string;
}

export interface Airline {
  id: string;
  name: string;
  code?: string;
  country?: string;
}

export interface Airport {
  id: string;
  name: string;
  code: string;
  city?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
}

// Offer types
export interface BaggageAllowance {
  pieces: number;
  weightKg: number;
}

export interface Offer {
  id: string;
  flightID: string;
  providerReference: string;
  provider: string;
  oneway?: boolean;
  oneWay?: boolean;
  segments?: string[];
  priceTotal: number;
  baggageAllowance?: BaggageAllowance[];
  lastTicketingDate?: string;
  bookableSeats: number;
  expiresAt?: string;
  isActive: boolean;
  flight?: Flight;
}

// Activity types
export interface TimeSlot {
  id: string;
  activityID: string;
  startTime: string;
  endTime: string;
  availableSpots: number;
  price: number;
}

export interface Activity {
  id: string;
  title: string;
  description: string;
  city: string;
  country: string;
  location: Location;
  meetingPoint: string;
  categories: string[];
  durationMinutes: number;
  inclusions: string[];
  exclusions: string[];
  images: string[];
  rating: number;
  reviewCount: number;
  timeSlots?: TimeSlot[];
}

// Package types
export interface PackageComponent {
  type: string;
  name: string;
  description: string;
  included: boolean;
}

export interface Package {
  id: string;
  name: string;
  description: string;
  destination: string;
  durationDays: number;
  startDateFrom: string;
  startDateTo: string;
  maxTravelers: number;
  basePrice: number;
  currency: string;
  includedComponents: PackageComponent[];
  tags: string[];
  images: string[];
  rating: number;
}

// Booking types
export type PaymentMethod = 'card' | 'mpesa' | 'bank';

export interface PassengerDetails {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dateOfBirth: string;
  passportNumber?: string;
  nationality?: string;
}

export interface GuestDetails {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
}

export interface ParticipantDetails {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  age?: number;
}

export interface TravelerDetails {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dateOfBirth: string;
  passportNumber?: string;
}

export interface FlightBookingParams {
  flightID: string;
  passengerDetails: PassengerDetails[];
  selectedSeats: string[];
  paymentMethod: PaymentMethod;
  currency: string;
}

export interface AccommodationBookingParams {
  accommodationID: string;
  guestDetails: GuestDetails[];
  roomTypeID: string;
  rooms: number;
  checkIn: string;
  checkOut: string;
  paymentMethod: PaymentMethod;
  currency: string;
}

export interface ActivityBookingParams {
  activityID: string;
  participantDetails: ParticipantDetails[];
  timeSlotID: string;
  participants: number;
  paymentMethod: PaymentMethod;
  currency: string;
}

export interface PackageBookingParams {
  packageID: string;
  travelersDetails: TravelerDetails[];
  componentSelection: string[];
  paymentMethod: PaymentMethod;
  currency: string;
}

export interface CancelBookingParams {
  bookingID: string;
  reason: string;
}

export interface Booking {
  id: string;
  type: 'flight' | 'accommodation' | 'activity' | 'package';
  status: string;
  createdAt: string;
  totalPrice: number;
  currency: string;
}

// Review types
export interface Review {
  id: string;
  userID: string;
  referenceID: string;
  reviewFor: 'flight' | 'accommodation' | 'activity' | 'package';
  rating: number;
  comment: string;
  createdAt: string;
  user?: {
    firstName: string;
    lastName: string;
    profilePic?: string;
  };
}

export interface CreateReviewParams {
  referenceID: string;
  reviewFor: 'flight' | 'accommodation' | 'activity' | 'package';
  rating: number;
  comment: string;
}

// Search filter types
export interface FlightSearchFilters {
  originCode?: string;
  destinationCode?: string;
  departureTime?: string;
  returnDate?: string;
  adults?: number;
  children?: number;
  infants?: number;
  cabinClass?: string;
  sortBy?: string;
  page?: number;
  pageSize?: number;
}

export interface AccommodationSearchFilters {
  city?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
  radiusKm?: number;
  checkInDate?: string;
  checkOutDate?: string;
  adults?: number;
  children?: number;
  infants?: number;
  totalRooms?: number;
  propertyType?: string;
  sortBy?: string;
  page?: number;
  pageSize?: number;
}

export interface ActivitySearchFilters {
  city?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
  radiusKm?: number;
  date?: string;
  adults?: number;
  children?: number;
  infants?: number;
  maxDurationMinutes?: number;
  category?: string;
  sortBy?: string;
  page?: number;
  pageSize?: number;
}

export interface PackageSearchFilters {
  destination?: string;
  startDate?: string;
  endDate?: string;
  adults?: number;
  children?: number;
  infants?: number;
  maxBudget?: number;
  tag?: string;
  component?: string;
  sortBy?: string;
  page?: number;
  pageSize?: number;
}

// Pagination response
export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}
