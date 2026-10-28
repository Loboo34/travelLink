import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authApi, setAuthToken, getAuthToken } from '../api';
import type { User } from '../types';

interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  isInitialized: boolean;
  isLoading: boolean;
  error: string | null;

  login: (email: string, password: string) => Promise<void>;
  register: (
    firstName: string,
    lastName: string,
    email: string,
    password: string
  ) => Promise<void>;
  logout: () => void;
  fetchProfile: () => Promise<void>;
  updateProfile: (data: Partial<User>) => Promise<void>;
  clearError: () => void;
  initializeAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      isAuthenticated: false,
      isInitialized: false,
      isLoading: false,
      error: null,

      initializeAuth: async () => {
        const token = getAuthToken();
        if (token) {
          set({ token, isAuthenticated: true, isLoading: true });
          try {
            await get().fetchProfile();
          } catch (error) {
            console.error('Failed to initialize auth profile:', error);
          }
          set({ isLoading: false, isInitialized: true });
          return;
        }
        set({ isInitialized: true });
      },

      login: async (email: string, password: string) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authApi.login({ email, password });
          setAuthToken(response.token);
          set({
            token: response.token,
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : 'Login failed';
          set({ error: message, isLoading: false });
          throw error;
        }
      },

      register: async (
        firstName: string,
        lastName: string,
        email: string,
        password: string
      ) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authApi.register({
            firstName,
            lastName,
            email,
            password,
          });
          setAuthToken(response.token);
          set({
            token: response.token,
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : 'Registration failed';
          set({ error: message, isLoading: false });
          throw error;
        }
      },

      logout: () => {
        setAuthToken(null);
        set({
          token: null,
          user: null,
          isAuthenticated: false,
          error: null,
        });
      },

      fetchProfile: async () => {
        try {
          const user = await authApi.getProfile();
          set({ user });
        } catch (error) {
          console.error('Failed to fetch profile:', error);
        }
      },

      updateProfile: async (data: Partial<User>) => {
        set({ isLoading: true, error: null });
        try {
          const user = await authApi.updateProfile(data);
          set({ user, isLoading: false });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : 'Update failed';
          set({ error: message, isLoading: false });
          throw error;
        }
      },

      clearError: () => set({ error: null }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        isAuthenticated: state.isAuthenticated,
        user: state.user,
        isInitialized: state.isInitialized,
      }),
    }
  )
);
