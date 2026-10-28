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
  FlightSearchFilters,
  AccommodationSearchFilters,
  ActivitySearchFilters,
  PackageSearchFilters,
  PaginatedResponse,
} from '../types';

interface SearchState {
  // Flight search
  flightFilters: FlightSearchFilters;
  flightResults: PaginatedResponse<Flight> | null;
  isSearchingFlights: boolean;

  // Accommodation search
  accommodationFilters: AccommodationSearchFilters;
  accommodationResults: PaginatedResponse<Accommodation> | null;
  isSearchingAccommodations: boolean;

  // Activity search
  activityFilters: ActivitySearchFilters;
  activityResults: PaginatedResponse<Activity> | null;
  isSearchingActivities: boolean;

  // Package search
  packageFilters: PackageSearchFilters;
  packageResults: PaginatedResponse<Package> | null;
  isSearchingPackages: boolean;

  // Selected items
  selectedFlight: Flight | null;
  selectedAccommodation: Accommodation | null;
  selectedActivity: Activity | null;
  selectedPackage: Package | null;

  // Error state
  error: string | null;

  // Actions
  setFlightFilters: (filters: Partial<FlightSearchFilters>) => void;
  searchFlights: () => Promise<void>;
  setAccommodationFilters: (filters: Partial<AccommodationSearchFilters>) => void;
  searchAccommodations: () => Promise<void>;
  setActivityFilters: (filters: Partial<ActivitySearchFilters>) => void;
  searchActivities: () => Promise<void>;
  setPackageFilters: (filters: Partial<PackageSearchFilters>) => void;
  searchPackages: () => Promise<void>;

  setSelectedFlight: (flight: Flight | null) => void;
  setSelectedAccommodation: (accommodation: Accommodation | null) => void;
  setSelectedActivity: (activity: Activity | null) => void;
  setSelectedPackage: (pkg: Package | null) => void;

  clearError: () => void;
  resetSearch: () => void;
}

const initialFlightFilters: FlightSearchFilters = {
  adults: 1,
  children: 0,
  infants: 0,
  cabinClass: 'economy',
  page: 1,
  pageSize: 10,
};

const initialAccommodationFilters: AccommodationSearchFilters = {
  adults: 1,
  children: 0,
  infants: 0,
  totalRooms: 1,
  page: 1,
  pageSize: 10,
};

const initialActivityFilters: ActivitySearchFilters = {
  adults: 1,
  children: 0,
  infants: 0,
  page: 1,
  pageSize: 10,
};

const initialPackageFilters: PackageSearchFilters = {
  adults: 1,
  children: 0,
  infants: 0,
  page: 1,
  pageSize: 10,
};

export const useSearchStore = create<SearchState>((set, get) => ({
  // Initial state
  flightFilters: initialFlightFilters,
  flightResults: null,
  isSearchingFlights: false,

  accommodationFilters: initialAccommodationFilters,
  accommodationResults: null,
  isSearchingAccommodations: false,

  activityFilters: initialActivityFilters,
  activityResults: null,
  isSearchingActivities: false,

  packageFilters: initialPackageFilters,
  packageResults: null,
  isSearchingPackages: false,

  selectedFlight: null,
  selectedAccommodation: null,
  selectedActivity: null,
  selectedPackage: null,

  error: null,

  // Flight actions
  setFlightFilters: (filters) =>
    set((state) => ({
      flightFilters: { ...state.flightFilters, ...filters },
    })),

  searchFlights: async () => {
    set({ isSearchingFlights: true, error: null });
    try {
      const results = await flightsApi.search(get().flightFilters);
      set({ flightResults: results, isSearchingFlights: false });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Search failed';
      set({ error: message, isSearchingFlights: false });
    }
  },

  // Accommodation actions
  setAccommodationFilters: (filters) =>
    set((state) => ({
      accommodationFilters: { ...state.accommodationFilters, ...filters },
    })),

  searchAccommodations: async () => {
    set({ isSearchingAccommodations: true, error: null });
    try {
      const results = await accommodationsApi.search(get().accommodationFilters);
      set({ accommodationResults: results, isSearchingAccommodations: false });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Search failed';
      set({ error: message, isSearchingAccommodations: false });
    }
  },

  // Activity actions
  setActivityFilters: (filters) =>
    set((state) => ({
      activityFilters: { ...state.activityFilters, ...filters },
    })),

  searchActivities: async () => {
    set({ isSearchingActivities: true, error: null });
    try {
      const results = await activitiesApi.search(get().activityFilters);
      set({ activityResults: results, isSearchingActivities: false });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Search failed';
      set({ error: message, isSearchingActivities: false });
    }
  },

  // Package actions
  setPackageFilters: (filters) =>
    set((state) => ({
      packageFilters: { ...state.packageFilters, ...filters },
    })),

  searchPackages: async () => {
    set({ isSearchingPackages: true, error: null });
    try {
      const results = await packagesApi.search(get().packageFilters);
      set({ packageResults: results, isSearchingPackages: false });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Search failed';
      set({ error: message, isSearchingPackages: false });
    }
  },

  // Selection actions
  setSelectedFlight: (flight) => set({ selectedFlight: flight }),
  setSelectedAccommodation: (accommodation) =>
    set({ selectedAccommodation: accommodation }),
  setSelectedActivity: (activity) => set({ selectedActivity: activity }),
  setSelectedPackage: (pkg) => set({ selectedPackage: pkg }),

  // Utility actions
  clearError: () => set({ error: null }),

  resetSearch: () =>
    set({
      flightFilters: initialFlightFilters,
      flightResults: null,
      accommodationFilters: initialAccommodationFilters,
      accommodationResults: null,
      activityFilters: initialActivityFilters,
      activityResults: null,
      packageFilters: initialPackageFilters,
      packageResults: null,
      selectedFlight: null,
      selectedAccommodation: null,
      selectedActivity: null,
      selectedPackage: null,
      error: null,
    }),
}));
