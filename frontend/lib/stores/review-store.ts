import { create } from 'zustand';
import { reviewsApi } from '../api';
import type { Review, CreateReviewParams } from '../types';

interface ReviewState {
  reviews: Review[];
  isLoading: boolean;
  error: string | null;

  fetchReviews: (referenceID: string, reviewFor: string) => Promise<void>;
  createReview: (params: CreateReviewParams) => Promise<Review>;
  updateReview: (
    reviewID: string,
    data: { rating: number; comment: string }
  ) => Promise<void>;
  deleteReview: (reviewID: string) => Promise<void>;
  clearReviews: () => void;
  clearError: () => void;
}

export const useReviewStore = create<ReviewState>((set, get) => ({
  reviews: [],
  isLoading: false,
  error: null,

  fetchReviews: async (referenceID: string, reviewFor: string) => {
    set({ isLoading: true, error: null });
    try {
      const reviews = await reviewsApi.get(referenceID, reviewFor);
      set({ reviews, isLoading: false });
    } catch (error) {
      const message =
        error instanceof Error ? error.message : 'Failed to fetch reviews';
      set({ error: message, isLoading: false });
    }
  },

  createReview: async (params: CreateReviewParams) => {
    set({ isLoading: true, error: null });
    try {
      const review = await reviewsApi.create(params);
      set((state) => ({
        reviews: [review, ...state.reviews],
        isLoading: false,
      }));
      return review;
    } catch (error) {
      const message =
        error instanceof Error ? error.message : 'Failed to create review';
      set({ error: message, isLoading: false });
      throw error;
    }
  },

  updateReview: async (
    reviewID: string,
    data: { rating: number; comment: string }
  ) => {
    set({ isLoading: true, error: null });
    try {
      const updatedReview = await reviewsApi.update(reviewID, data);
      set((state) => ({
        reviews: state.reviews.map((r) =>
          r.id === reviewID ? updatedReview : r
        ),
        isLoading: false,
      }));
    } catch (error) {
      const message =
        error instanceof Error ? error.message : 'Failed to update review';
      set({ error: message, isLoading: false });
      throw error;
    }
  },

  deleteReview: async (reviewID: string) => {
    set({ isLoading: true, error: null });
    try {
      await reviewsApi.delete(reviewID);
      set((state) => ({
        reviews: state.reviews.filter((r) => r.id !== reviewID),
        isLoading: false,
      }));
    } catch (error) {
      const message =
        error instanceof Error ? error.message : 'Failed to delete review';
      set({ error: message, isLoading: false });
      throw error;
    }
  },

  clearReviews: () => set({ reviews: [], error: null }),
  clearError: () => set({ error: null }),
}));
