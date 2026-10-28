import { create } from 'zustand';
import {
  flightsApi,
  accommodationsApi,
  activitiesApi,
  packagesApi,
} from '../api';
import type {
  Flight,
  Accommodation,
  Activity,
  Package,
  Offer,
  TimeSlot,
  RoomType,
  PaymentMethod,
  PassengerDetails,
  GuestDetails,
  ParticipantDetails,
  TravelerDetails,
  Booking,
} from '../types';

type BookingType = 'flight' | 'accommodation' | 'activity' | 'package';

interface BookingState {
  // Active booking type
  activeBookingType: BookingType | null;

  // Selected items
  selectedFlight: Flight | null;
  selectedOffer: Offer | null;
  selectedAccommodation: Accommodation | null;
  selectedRoomType: RoomType | null;
  selectedActivity: Activity | null;
  selectedTimeSlot: TimeSlot | null;
  selectedPackage: Package | null;

  // Traveler details
  passengerDetails: PassengerDetails[];
  guestDetails: GuestDetails[];
  participantDetails: ParticipantDetails[];
  travelersDetails: TravelerDetails[];

  // Booking specifics
  selectedSeats: string[];
  rooms: number;
  checkIn: string;
  checkOut: string;
  participants: number;
  componentSelection: string[];

  // Payment
  paymentMethod: PaymentMethod;
  currency: string;

  // State
  isBooking: boolean;
  bookingResult: Booking | null;
  error: string | null;

  // Actions
  setActiveBookingType: (type: BookingType) => void;

  // Flight booking
  setFlightBooking: (flight: Flight, offer?: Offer) => void;
  setPassengerDetails: (passengers: PassengerDetails[]) => void;
  setSelectedSeats: (seats: string[]) => void;

  // Accommodation booking
  setAccommodationBooking: (
    accommodation: Accommodation,
    roomType: RoomType
  ) => void;
  setGuestDetails: (guests: GuestDetails[]) => void;
  setRooms: (rooms: number) => void;
  setCheckInOut: (checkIn: string, checkOut: string) => void;

  // Activity booking
  setActivityBooking: (activity: Activity, timeSlot: TimeSlot) => void;
  setParticipantDetails: (participants: ParticipantDetails[]) => void;
  setParticipants: (count: number) => void;

  // Package booking
  setPackageBooking: (pkg: Package) => void;
  setTravelersDetails: (travelers: TravelerDetails[]) => void;
  setComponentSelection: (components: string[]) => void;

  // Payment
  setPaymentMethod: (method: PaymentMethod) => void;
  setCurrency: (currency: string) => void;

  // Submit booking
  submitBooking: () => Promise<Booking>;

  // Utility
  clearBooking: () => void;
  clearError: () => void;
}

const initialState = {
  activeBookingType: null as BookingType | null,
  selectedFlight: null,
  selectedOffer: null,
  selectedAccommodation: null,
  selectedRoomType: null,
  selectedActivity: null,
  selectedTimeSlot: null,
  selectedPackage: null,
  passengerDetails: [],
  guestDetails: [],
  participantDetails: [],
  travelersDetails: [],
  selectedSeats: [],
  rooms: 1,
  checkIn: '',
  checkOut: '',
  participants: 1,
  componentSelection: [],
  paymentMethod: 'card' as PaymentMethod,
  currency: 'USD',
  isBooking: false,
  bookingResult: null,
  error: null,
};

export const useBookingStore = create<BookingState>((set, get) => ({
  ...initialState,

  setActiveBookingType: (type) => set({ activeBookingType: type }),

  // Flight booking
  setFlightBooking: (flight, offer) =>
    set({
      activeBookingType: 'flight',
      selectedFlight: flight,
      selectedOffer: offer || null,
    }),

  setPassengerDetails: (passengers) => set({ passengerDetails: passengers }),
  setSelectedSeats: (seats) => set({ selectedSeats: seats }),

  // Accommodation booking
  setAccommodationBooking: (accommodation, roomType) =>
    set({
      activeBookingType: 'accommodation',
      selectedAccommodation: accommodation,
      selectedRoomType: roomType,
    }),

  setGuestDetails: (guests) => set({ guestDetails: guests }),
  setRooms: (rooms) => set({ rooms }),
  setCheckInOut: (checkIn, checkOut) => set({ checkIn, checkOut }),

  // Activity booking
  setActivityBooking: (activity, timeSlot) =>
    set({
      activeBookingType: 'activity',
      selectedActivity: activity,
      selectedTimeSlot: timeSlot,
    }),

  setParticipantDetails: (participants) =>
    set({ participantDetails: participants }),
  setParticipants: (count) => set({ participants: count }),

  // Package booking
  setPackageBooking: (pkg) =>
    set({
      activeBookingType: 'package',
      selectedPackage: pkg,
    }),

  setTravelersDetails: (travelers) => set({ travelersDetails: travelers }),
  setComponentSelection: (components) =>
    set({ componentSelection: components }),

  // Payment
  setPaymentMethod: (method) => set({ paymentMethod: method }),
  setCurrency: (currency) => set({ currency }),

  // Submit booking
  submitBooking: async () => {
    const state = get();
    set({ isBooking: true, error: null });

    try {
      let booking: Booking;

      switch (state.activeBookingType) {
        case 'flight':
          if (!state.selectedFlight) throw new Error('No flight selected');
          booking = await flightsApi.book({
            flightID: state.selectedFlight.id,
            passengerDetails: state.passengerDetails,
            selectedSeats: state.selectedSeats,
            paymentMethod: state.paymentMethod,
            currency: state.currency,
          });
          break;

        case 'accommodation':
          if (!state.selectedAccommodation || !state.selectedRoomType)
            throw new Error('No accommodation selected');
          booking = await accommodationsApi.book({
            accommodationID: state.selectedAccommodation.id,
            guestDetails: state.guestDetails,
            roomTypeID: state.selectedRoomType.id,
            rooms: state.rooms,
            checkIn: state.checkIn,
            checkOut: state.checkOut,
            paymentMethod: state.paymentMethod,
            currency: state.currency,
          });
          break;

        case 'activity':
          if (!state.selectedActivity || !state.selectedTimeSlot)
            throw new Error('No activity selected');
          booking = await activitiesApi.book({
            activityID: state.selectedActivity.id,
            participantDetails: state.participantDetails,
            timeSlotID: state.selectedTimeSlot.id,
            participants: state.participants,
            paymentMethod: state.paymentMethod,
            currency: state.currency,
          });
          break;

        case 'package':
          if (!state.selectedPackage) throw new Error('No package selected');
          booking = await packagesApi.book({
            packageID: state.selectedPackage.id,
            travelersDetails: state.travelersDetails,
            componentSelection: state.componentSelection,
            paymentMethod: state.paymentMethod,
            currency: state.currency,
          });
          break;

        default:
          throw new Error('No booking type selected');
      }

      set({ bookingResult: booking, isBooking: false });
      return booking;
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Booking failed';
      set({ error: message, isBooking: false });
      throw error;
    }
  },

  clearBooking: () => set(initialState),
  clearError: () => set({ error: null }),
}));
