import { create } from 'zustand';
import {
  adminApi,
  flightsApi,
  accommodationsApi,
  activitiesApi,
  packagesApi,
  offersApi,
} from '../api';
import type {
  User,
  Flight,
  Accommodation,
  AccommodationAvailability,
  Activity,
  Package,
  Offer,
  Booking,
  Airline,
  Airport,
  PlaneModel,
} from '../types';

interface AdminStats {
  totalFlights: number;
  totalAccommodations: number;
  totalActivities: number;
  totalPackages: number;
  totalOffers: number;
  totalBookings: number;
  totalUsers: number;
  revenue: number;
}

interface AdminState {
  // Data
  stats: AdminStats | null;
  users: User[];
  flights: Flight[];
  accommodations: Accommodation[];
  accommodationAvailabilities: AccommodationAvailability[];
  activities: Activity[];
  packages: Package[];
  offers: Offer[];
  bookings: Booking[];
  airlines: Airline[];
  airports: Airport[];
  planes: PlaneModel[];

  // Loading states
  isLoading: boolean;
  error: string | null;

  // Actions
  fetchStats: () => Promise<void>;
  fetchUsers: () => Promise<void>;
  fetchFlights: () => Promise<void>;
  fetchAirlines: () => Promise<void>;
  fetchAirports: () => Promise<void>;
  fetchPlanes: () => Promise<void>;
  fetchAccommodations: () => Promise<void>;
  fetchAccommodationAvailabilities: () => Promise<void>;
  fetchActivities: () => Promise<void>;
  fetchPackages: () => Promise<void>;
  fetchOffers: () => Promise<void>;
  fetchBookings: () => Promise<void>;

  // User actions
  updateUserRole: (userId: string, role: string) => Promise<void>;
  toggleUserStatus: (userId: string) => Promise<void>;

  // Flight actions
  createFlight: (data: Omit<Flight, 'id'>) => Promise<void>;
  updateFlight: (id: string, data: Partial<Flight>) => Promise<void>;
  updateFlightStatus: (id: string, status: string) => Promise<void>;
  deleteFlight: (id: string) => Promise<void>;

  // Accommodation actions
  createAccommodation: (data: Omit<Accommodation, 'id'>) => Promise<void>;
  updateAccommodation: (id: string, data: Partial<Accommodation>) => Promise<void>;
  deleteAccommodation: (id: string) => Promise<void>;
  createAccommodationAvailability: (
    data: Pick<
      AccommodationAvailability,
      'accommodationID' | 'roomTypeID' | 'date' | 'totalRooms' | 'pricePerNight'
    >
  ) => Promise<void>;
  updateAccommodationAvailabilityStatus: (id: string, isActive: boolean) => Promise<void>;

  // Activity actions
  createActivity: (data: Omit<Activity, 'id'>) => Promise<void>;
  updateActivity: (id: string, data: Partial<Activity>) => Promise<void>;
  deleteActivity: (id: string) => Promise<void>;

  // Package actions
  createPackage: (data: Omit<Package, 'id'>) => Promise<void>;
  updatePackage: (id: string, data: Partial<Package>) => Promise<void>;
  deletePackage: (id: string) => Promise<void>;

  // Offer actions
  createOffer: (data: Omit<Offer, 'id'>) => Promise<void>;
  updateOffer: (id: string, data: Partial<Offer>) => Promise<void>;
  updateOfferStatus: (id: string, isActive: boolean) => Promise<void>;
  deleteOffer: (id: string) => Promise<void>;

  // Booking actions
  updateBookingStatus: (id: string, status: string) => Promise<void>;

  clearError: () => void;
}

export const useAdminStore = create<AdminState>((set, get) => ({
  stats: null,
  users: [],
  flights: [],
  airlines: [],
  airports: [],
  planes: [],
  accommodations: [],
  accommodationAvailabilities: [],
  activities: [],
  packages: [],
  offers: [],
  bookings: [],
  isLoading: false,
  error: null,

  fetchStats: async () => {
    set({ isLoading: true, error: null });
    try {
      const [
        flights,
        accommodations,
        activities,
        packages,
        offers,
        users,
        bookings,
      ] = await Promise.all([
        flightsApi.getAll().catch(() => []),
        accommodationsApi.getAll().catch(() => []),
        activitiesApi.getAll().catch(() => []),
        packagesApi.getAll().catch(() => []),
        offersApi.getAll().catch(() => []),
        adminApi.getUsers().catch(() => []),
        adminApi.getBookings().catch(() => []),
      ]);

      const revenue = bookings.reduce(
        (sum, booking) => sum + (booking.totalPrice ?? 0),
        0
      );

      set({
        stats: {
          totalFlights: flights.length,
          totalAccommodations: accommodations.length,
          totalActivities: activities.length,
          totalPackages: packages.length,
          totalOffers: offers.length,
          totalBookings: bookings.length,
          totalUsers: users.length,
          revenue,
        },
        isLoading: false,
      });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to load dashboard stats',
        isLoading: false,
      });
    }
  },

  fetchUsers: async () => {
    set({ isLoading: true, error: null });
    try {
      const users = await adminApi.getUsers();
      set({ users, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch users', isLoading: false });
    }
  },

  fetchAirlines: async () => {
    set({ isLoading: true, error: null });
    try {
      const airlines = await adminApi.getAirlines();
      set({ airlines, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch airlines', isLoading: false });
    }
  },

  fetchAirports: async () => {
    set({ isLoading: true, error: null });
    try {
      const airports = await adminApi.getAirports();
      set({ airports, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch airports', isLoading: false });
    }
  },

  fetchFlights: async () => {
    set({ isLoading: true, error: null });
    try {
      const flights = await flightsApi.getAll();
      set({ flights, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch flights', isLoading: false });
    }
  },

  fetchAccommodations: async () => {
    set({ isLoading: true, error: null });
    try {
      const accommodations = await accommodationsApi.getAll();
      set({ accommodations, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch accommodations', isLoading: false });
    }
  },

  fetchPlanes: async () => {
    set({ isLoading: true, error: null });
    try {
      const planes = await adminApi.getPlanes();
      set({ planes, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch planes', isLoading: false });
    }
  },

  fetchActivities: async () => {
    set({ isLoading: true, error: null });
    try {
      const activities = await activitiesApi.getAll();
      set({ activities, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch activities', isLoading: false });
    }
  },

  fetchPackages: async () => {
    set({ isLoading: true, error: null });
    try {
      const packages = await packagesApi.getAll();
      set({ packages, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch packages', isLoading: false });
    }
  },

  fetchOffers: async () => {
    set({ isLoading: true, error: null });
    try {
      const offers = await offersApi.getAll();
      set({ offers, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch offers', isLoading: false });
    }
  },

  fetchBookings: async () => {
    set({ isLoading: true, error: null });
    try {
      const bookings = await adminApi.getBookings();
      set({ bookings, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch bookings', isLoading: false });
    }
  },

  updateUserRole: async (userId: string, role: string) => {
    try {
      await adminApi.updateUserRole(userId, role);
      await get().fetchUsers();
    } catch (error) {
      set({ error: 'Failed to update user role' });
    }
  },

  toggleUserStatus: async (userId: string) => {
    try {
      await adminApi.toggleUserStatus(userId);
      await get().fetchUsers();
    } catch (error) {
      set({ error: 'Failed to toggle user status' });
    }
  },

  createFlight: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.createFlight(data);
      console.log('Flight created successfully:', data);
      await get().fetchFlights();
    } catch (error) {
      set({ error: 'Failed to create flight', isLoading: false });
      throw error;
    }
  },

  updateFlight: async (id, data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updateFlight(id, data);
      await get().fetchFlights();
    } catch (error) {
      set({ error: 'Failed to update flight', isLoading: false });
      throw error;
    }
  },

  fetchAccommodationAvailabilities: async () => {
    set({ isLoading: true, error: null });
    try {
      const accommodationAvailabilities = await adminApi.getAccommodationAvailabilities();
      set({ accommodationAvailabilities, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to fetch accommodation availability', isLoading: false });
    }
  },

  updateFlightStatus: async (id, status) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updateFlightStatus(id, status);
      await get().fetchFlights();
    } catch (error) {
      set({ error: 'Failed to update flight status', isLoading: false });
      throw error;
    }
  },

  deleteFlight: async (id) => {
    try {
      await adminApi.deleteFlight(id);
      set((state) => ({
        flights: state.flights.filter((f) => f.id !== id),
      }));
    } catch (error) {
      set({ error: 'Failed to delete flight' });
      throw error;
    }
  },

  createAccommodation: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.createAccommodation(data);
      await get().fetchAccommodations();
    } catch (error) {
      set({ error: 'Failed to create accommodation', isLoading: false });
      throw error;
    }
  },

  updateAccommodation: async (id, data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updateAccommodation(id, data);
      await get().fetchAccommodations();
    } catch (error) {
      set({ error: 'Failed to update accommodation', isLoading: false });
      throw error;
    }
  },

  deleteAccommodation: async (id) => {
    try {
      await adminApi.deleteAccommodation(id);
      set((state) => ({
        accommodations: state.accommodations.filter((a) => a.id !== id),
      }));
    } catch (error) {
      set({ error: 'Failed to delete accommodation' });
      throw error;
    }
  },

  createActivity: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.createActivity(data);
      await get().fetchActivities();
    } catch (error) {
      set({ error: 'Failed to create activity', isLoading: false });
      throw error;
    }
  },

  updateActivity: async (id, data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updateActivity(id, data);
      await get().fetchActivities();
    } catch (error) {
      set({ error: 'Failed to update activity', isLoading: false });
      throw error;
    }
  },

  deleteActivity: async (id) => {
    try {
      await adminApi.deleteActivity(id);
      set((state) => ({
        activities: state.activities.filter((a) => a.id !== id),
      }));
    } catch (error) {
      set({ error: 'Failed to delete activity' });
      throw error;
    }
  },

  createPackage: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.createPackage(data);
      await get().fetchPackages();
    } catch (error) {
      set({ error: 'Failed to create package', isLoading: false });
      throw error;
    }
  },

  updatePackage: async (id, data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updatePackage(id, data);
      await get().fetchPackages();
    } catch (error) {
      set({ error: 'Failed to update package', isLoading: false });
      throw error;
    }
  },

  deletePackage: async (id) => {
    try {
      await adminApi.deletePackage(id);
      set((state) => ({
        packages: state.packages.filter((p) => p.id !== id),
      }));
    } catch (error) {
      set({ error: 'Failed to delete package' });
      throw error;
    }
  },

  createOffer: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.createOffer(data);
      await get().fetchOffers();
    } catch (error) {
      set({ error: 'Failed to create offer', isLoading: false });
      throw error;
    }
  },

  updateOffer: async (id, data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updateOffer(id, data);
      await get().fetchOffers();
    } catch (error) {
      set({ error: 'Failed to update offer', isLoading: false });
      throw error;
    }
  },

  createAccommodationAvailability: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.createAccommodationAvailability(data);
      await get().fetchAccommodationAvailabilities();
    } catch (error) {
      set({ error: 'Failed to create accommodation availability', isLoading: false });
      throw error;
    }
  },

  updateAccommodationAvailabilityStatus: async (id, isActive) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updateAccommodationAvailabilityStatus(id, isActive);
      await get().fetchAccommodationAvailabilities();
    } catch (error) {
      set({ error: 'Failed to update accommodation availability', isLoading: false });
      throw error;
    }
  },

  updateOfferStatus: async (id, isActive) => {
    set({ isLoading: true, error: null });
    try {
      await adminApi.updateOfferStatus(id, isActive);
      await get().fetchOffers();
    } catch (error) {
      set({ error: 'Failed to update offer status', isLoading: false });
      throw error;
    }
  },

  deleteOffer: async (id) => {
    try {
      await adminApi.deleteOffer(id);
      set((state) => ({
        offers: state.offers.filter((o) => o.id !== id),
      }));
    } catch (error) {
      set({ error: 'Failed to delete offer' });
      throw error;
    }
  },

  updateBookingStatus: async (id, status) => {
    try {
      await adminApi.updateBookingStatus(id, status);
      await get().fetchBookings();
    } catch (error) {
      set({ error: 'Failed to update booking status' });
      throw error;
    }
  },

  clearError: () => set({ error: null }),
}));
