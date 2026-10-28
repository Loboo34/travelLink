import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import type {
  User,
  AuthResponse,
  Flight,
  Airline,
  Airport,
  Offer,
  Accommodation,
  AccommodationAvailability,
  Activity,
  Package,
  TimeSlot,
  Review,
  Booking,
  FlightSearchFilters,
  AccommodationSearchFilters,
  ActivitySearchFilters,
  PackageSearchFilters,
  FlightBookingParams,
  AccommodationBookingParams,
  ActivityBookingParams,
  PackageBookingParams,
  CancelBookingParams,
  CreateReviewParams,
  PaginatedResponse,
  PlaneModel,
} from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Create axios instance
const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

function normalizeAccommodation(accommodation: Accommodation): Accommodation {
  const rawLocation = accommodation.location as
    | (Accommodation['location'] & { coordinates?: number[] })
    | undefined;
  const coordinates = rawLocation?.coordinates;

  return {
    ...accommodation,
    type: accommodation.type || accommodation.propertyType || '',
    roomTypes: accommodation.roomTypes || accommodation.roomType || [],
    location: {
      ...rawLocation,
      latitude: rawLocation?.latitude ?? coordinates?.[1] ?? 0,
      longitude: rawLocation?.longitude ?? coordinates?.[0] ?? 0,
    },
    checkInTime: String(accommodation.checkInTime ?? ''),
    checkOutTime: String(accommodation.checkOutTime ?? ''),
  };
}

// Token management
let authToken: string | null = null;

export const setAuthToken = (token: string | null) => {
  authToken = token;
  if (typeof window !== 'undefined') {
    if (token) {
      localStorage.setItem('auth_token', token);
    } else {
      localStorage.removeItem('auth_token');
    }
  }
};

export const getAuthToken = (): string | null => {
  if (authToken) return authToken;
  if (typeof window !== 'undefined') {
    authToken = localStorage.getItem('auth_token');
  }
  return authToken;
};

// Request interceptor for auth token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAuthToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      setAuthToken(null);
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authApi = {
  register: async (data: {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
  }): Promise<AuthResponse> => {
    const response = await api.post('/auth/register', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to register');
    }
    return response.data.data;
  },

  login: async (data: { email: string; password: string }): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', data);
    console.log('Login response:', response.data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to login');
    }
    return response.data.data;
  },

  getProfile: async (): Promise<User> => {
    const response = await api.get('/auth/profile');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to get profile');
    }
    return response.data.data;
  },

  updateProfile: async (data: Partial<User>): Promise<User> => {
    const response = await api.put('/user/profile', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update profile');
    }
    return response.data.data;
  },

  updatePassword: async (data: {
    currentPassword: string;
    newPassword: string;
  }): Promise<void> => {
    const response = await api.put('/user/password', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update password');
    }
  },

  updateProfilePic: async (formData: FormData): Promise<User> => {
    const response = await api.put('/auth/profilepic', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update profile picture');
    }
    return response.data.data;
  },
};

// Flights API
export const flightsApi = {
  getAll: async (): Promise<Flight[]> => {
    const response = await api.get('/flights');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch flights');
    }
    return response.data.data;
  },

  getById: async (id: string): Promise<Flight> => {
    const response = await api.get(`/flight/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch flight');
    }
    return response.data.data;
  },

  search: async (filters: FlightSearchFilters): Promise<PaginatedResponse<Flight>> => {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, String(value));
      }
    });
    const response = await api.get(`/flight/search?${params}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to search flights');
    }
    return response.data.data;
  },

  book: async (params: FlightBookingParams): Promise<Booking> => {
    const response = await api.post('/flight/book', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to book flight');
    }
    return response.data.data;
  },

  cancel: async (params: CancelBookingParams): Promise<void> => {
    const response = await api.post('/flight/cancel', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to cancel flight');
    }
  },
};

// Offers API
export const offersApi = {
  getAll: async (): Promise<Offer[]> => {
    const response = await api.get('/offers');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch offers');
    }
    return response.data.data;
  },

  getById: async (id: string): Promise<Offer> => {
    const response = await api.get(`/offer/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch offer');
    }
    return response.data.data;
  },
};

// Accommodations API
export const accommodationsApi = {
  getAll: async (): Promise<Accommodation[]> => {
    const response = await api.get('/accommodations');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch accommodations');
    }
    return (response.data.data as Accommodation[]).map(normalizeAccommodation);
  },

  getById: async (id: string): Promise<Accommodation> => {
    const response = await api.get(`/accommodation/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch accommodation');
    }
    return normalizeAccommodation(response.data.data);
  },

  search: async (
    filters: AccommodationSearchFilters
  ): Promise<PaginatedResponse<Accommodation>> => {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, String(value));
      }
    });
    const response = await api.get(`/accommodation/search?${params}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to search accommodations');
    }
    return {
      ...response.data.data,
      data: (response.data.data.data as Accommodation[]).map(normalizeAccommodation),
    };
  },

  book: async (params: AccommodationBookingParams): Promise<Booking> => {
    const response = await api.post('/accommodation/book', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to book accommodation');
    }
    return response.data.data;
  },

  cancel: async (params: CancelBookingParams): Promise<void> => {
    const response = await api.post('/accommodation/cancel', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to cancel accommodation');
    }
  },
};

// Activities API
export const activitiesApi = {
  getAll: async (): Promise<Activity[]> => {
    const response = await api.get('/activitties');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch activities');
    }
    return response.data.data;
  },

  getById: async (id: string): Promise<Activity> => {
    const response = await api.get(`/activity/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch activity');
    }
    return response.data.data;
  },

  search: async (filters: ActivitySearchFilters): Promise<PaginatedResponse<Activity>> => {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, String(value));
      }
    });
    const response = await api.get(`/activity/search?${params}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to search activities');
    }
    return response.data.data;
  },

  book: async (params: ActivityBookingParams): Promise<Booking> => {
    const response = await api.post('/activity/book', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to book activity');
    }
    return response.data.data;
  },

  cancel: async (params: CancelBookingParams): Promise<void> => {
    const response = await api.post('/activity/cancel', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to cancel activity');
    }
  },
};

// Time Slots API
export const timeSlotsApi = {
  getAll: async (): Promise<TimeSlot[]> => {
    const response = await api.get('/timeslots');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch time slots');
    }
    return response.data.data;
  },

  getById: async (id: string): Promise<TimeSlot> => {
    const response = await api.get(`/timeslot/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch time slot');
    }
    return response.data.data;
  },
};

// Packages API
export const packagesApi = {
  getAll: async (): Promise<Package[]> => {
    const response = await api.get('/packages');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch packages');
    }
    return response.data.data;
  },

  getById: async (id: string): Promise<Package> => {
    const response = await api.get(`/package/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch package');
    }
    return response.data.data;
  },

  search: async (filters: PackageSearchFilters): Promise<PaginatedResponse<Package>> => {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, String(value));
      }
    });
    const response = await api.get(`/package/search?${params}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to search packages');
    }
    return response.data.data;
  },

  book: async (params: PackageBookingParams): Promise<Booking> => {
    const response = await api.post('/package/book', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to book package');
    }
    return response.data.data;
  },

  cancel: async (params: CancelBookingParams): Promise<void> => {
    const response = await api.post('/package/cancel', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to cancel package');
    }
  },
};

// Reviews API
export const reviewsApi = {
  get: async (referenceID: string, reviewFor: string): Promise<Review[]> => {
    const response = await api.get(
      `/review/get?referenceID=${referenceID}&reviewFor=${reviewFor}`
    );
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch reviews');
    }
    return response.data.data;
  },

  create: async (params: CreateReviewParams): Promise<Review> => {
    const response = await api.post('/review/create', params);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to create review');
    }
    return response.data.data;
  },

  update: async (
    reviewID: string,
    data: { rating: number; comment: string }
  ): Promise<Review> => {
    const response = await api.put('/review/update', { reviewID, ...data });
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update review');
    }
    return response.data.data;
  },

  delete: async (reviewID: string): Promise<void> => {
    const response = await api.delete(`/review/delete/${reviewID}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to delete review');
    }
  },
};

  // Admin API
export const adminApi = {
  // Dashboard stats
  getStats: async (): Promise<{
    totalFlights: number;
    totalAccommodations: number;
    totalActivities: number;
    totalPackages: number;
    totalOffers: number;
    totalBookings: number;
    totalUsers: number;
    revenue: number;
  }> => {
    const response = await api.get('/admin/stats');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch admin stats');
    }
    return response.data.data;
  },

  // Users management
  getUsers: async (): Promise<User[]> => {
    const response = await api.get('/admin/users');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch users');
    }
    return response.data.data;
  },

  updateUserRole: async (userId: string, role: string): Promise<User> => {
    const response = await api.put(`/admin/user/${userId}/role`, { role });
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update user role');
    }
    return response.data.data;
  },

  toggleUserStatus: async (userId: string): Promise<User> => {
    const response = await api.put(`/admin/user/${userId}/toggle-status`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to toggle user status');
    }
    return response.data.data;
  },

  // Flights management
  createFlight: async (data: Omit<Flight, 'id'>): Promise<Flight> => {
    const response = await api.post('/flight/add', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to create flight');
    }
    return response.data.data;
  },

  updateFlight: async (id: string, data: Partial<Flight>): Promise<Flight> => {
    const response = await api.put(`/flight/update/${id}`, data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update flight');
    }
    return response.data.data;
  },

  updateFlightStatus: async (id: string, status: string): Promise<{ status: string }> => {
    const response = await api.patch(`/flight/status/${id}`, { status });
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update flight status');
    }
    return response.data.data;
  },

  deleteFlight: async (id: string): Promise<void> => {
    const response = await api.delete(`/flight/delete/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to delete flight');
    }
  },

  // Accommodations management
  createAccommodation: async (data: Omit<Accommodation, 'id'>): Promise<Accommodation> => {
    const response = await api.post('/accommodation/add', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to create accommodation');
    }
    return response.data.data;
  },

  updateAccommodation: async (id: string, data: Partial<Accommodation>): Promise<Accommodation> => {
    const response = await api.patch(`/accommodation/update/${id}`, data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update accommodation');
    }
    return response.data.data;
  },

  getAccommodationAvailabilities: async (): Promise<AccommodationAvailability[]> => {
    const response = await api.get('/available');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch accommodation availability');
    }
    return response.data.data;
  },

  createAccommodationAvailability: async (
    data: Pick<
      AccommodationAvailability,
      'accommodationID' | 'roomTypeID' | 'date' | 'totalRooms' | 'pricePerNight'
    >
  ): Promise<AccommodationAvailability> => {
    const response = await api.post('/availability/add', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to create accommodation availability');
    }
    return response.data.data;
  },

  updateAccommodationAvailabilityStatus: async (
    id: string,
    isActive: boolean
  ): Promise<{ isActive: boolean }> => {
    const response = await api.patch(`/availability/status/${id}`, { isActive });
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update accommodation availability');
    }
    return response.data.data;
  },

  deleteAccommodation: async (id: string): Promise<void> => {
    const response = await api.delete(`/accommodation/delete/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to delete accommodation');
    }
  },

  // Activities management
  createActivity: async (data: Omit<Activity, 'id'>): Promise<Activity> => {
    const response = await api.post('/activity/add', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to create activity');
    }
    return response.data.data;
  },

  updateActivity: async (id: string, data: Partial<Activity>): Promise<Activity> => {
    const response = await api.put(`/activity/update/${id}`, data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update activity');
    }
    return response.data.data;
  },

  deleteActivity: async (id: string): Promise<void> => {
    const response = await api.delete(`/activity/delete/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to delete activity');
    }
  },

  // Packages management
  createPackage: async (data: Omit<Package, 'id'>): Promise<Package> => {
    const response = await api.post('/package/add', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to create package');
    }
    return response.data.data;
  },

  updatePackage: async (id: string, data: Partial<Package>): Promise<Package> => {
    const response = await api.put(`/package/update/${id}`, data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update package');
    }
    return response.data.data;
  },

  deletePackage: async (id: string): Promise<void> => {
    const response = await api.delete(`/package/delete/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to delete package');
    }
  },

  // Offers management
  createOffer: async (data: Omit<Offer, 'id'>): Promise<Offer> => {
    const response = await api.post('/offer/create', data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to create offer');
    }
    return response.data.data;
  },

  updateOffer: async (id: string, data: Partial<Offer>): Promise<Offer> => {
    const response = await api.put(`/offer/update/${id}`, data);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update offer');
    }
    return response.data.data;
  },

  updateOfferStatus: async (id: string, isActive: boolean): Promise<{ isActive: boolean }> => {
    const response = await api.patch(`/offer/isactive/${id}`, { isActive });
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update offer status');
    }
    return response.data.data;
  },

  deleteOffer: async (id: string): Promise<void> => {
    const response = await api.delete(`/offer/delete/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to delete offer');
    }
  },

  // Bookings management
  getBookings: async (): Promise<Booking[]> => {
    const response = await api.get('/admin/bookings');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch bookings');
    }
    return response.data.data as Booking[];
  },

  updateBookingStatus: async (id: string, status: string): Promise<Booking> => {
    const response = await api.put(`/admin/booking/${id}/status`, { status });
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to update booking status');
    }
    return response.data.data as Booking;
  },

  cancelBooking: async (id: string): Promise<void> => {
    const response = await api.post(`/admin/booking/${id}/cancel`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to cancel booking');
    }
  },

  // Airlines & Airports
  getAirlines: async (): Promise<Airline[]> => {
    const response = await api.get('/admin/airlines');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch airlines');
    }
    return response.data.data as Airline[];
  },

  getAirports: async (): Promise<Airport[]> => {
    const response = await api.get('/admin/airports');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch airports');
    }
    return response.data.data as Airport[];
  },

  //planes
   getPlanes: async (): Promise<PlaneModel[]> => {
    const response = await api.get('/admin/planes');
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch planes');
    }
    return response.data.data as PlaneModel[];
  },

  getPlane: async (id: string): Promise<PlaneModel> => {
    const response = await api.get(`/admin/plane/${id}`);
    if (!response.data.success) {
      throw new Error(response.data.error?.message || 'Failed to fetch plane');
    }
    return response.data.data as PlaneModel;
  },
};

export default api;
