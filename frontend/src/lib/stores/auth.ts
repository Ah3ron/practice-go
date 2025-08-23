import { writable } from 'svelte/store';
import type { User } from '$lib/services/api';

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  initialized: boolean;
}

const createAuthStore = () => {
  const initialState: AuthState = {
    user: null,
    token: null,
    isAuthenticated: false,
    loading: true,
    initialized: false
  };

  const { subscribe, set, update } = writable<AuthState>(initialState);

  return {
    subscribe,
    login: (user: User, token: string) => {
      try {
        // Ensure we're in browser environment
        if (typeof window !== 'undefined' && localStorage) {
          localStorage.setItem('token', token);
          localStorage.setItem('user', JSON.stringify(user));
        }

        const newState = {
          user,
          token,
          isAuthenticated: true,
          loading: false,
          initialized: true
        };
        set(newState);

        // Log for debugging
        console.log('Auth state updated:', { token, user: user.username, isAuthenticated: true });
      } catch (error) {
        console.error('Failed to save auth data to localStorage:', error);
        // Still update the store state even if localStorage fails
        set({
          user,
          token,
          isAuthenticated: true,
          loading: false,
          initialized: true
        });
      }
    },
    logout: () => {
      try {
        if (typeof window !== 'undefined' && localStorage) {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
        console.log('User logged out, auth state cleared');
      } catch (error) {
        console.error('Failed to clear localStorage during logout:', error);
      }

      set({
        user: null,
        token: null,
        isAuthenticated: false,
        loading: false,
        initialized: true
      });
    },
    initialize: () => {
      // Prevent multiple initializations
      update((state) => {
        if (state.initialized) {
          console.log('Auth already initialized, skipping');
          return state;
        }

        // Check if we're in browser environment
        if (typeof window === 'undefined' || typeof localStorage === 'undefined') {
          console.log('Not in browser environment, initializing as unauthenticated');
          return {
            user: null,
            token: null,
            isAuthenticated: false,
            loading: false,
            initialized: true
          };
        }

        try {
          const token = localStorage.getItem('token');
          const userStr = localStorage.getItem('user');

          console.log('Checking stored auth data:', {
            hasToken: !!token,
            hasUser: !!userStr,
            tokenLength: token?.length || 0
          });

          if (token && userStr) {
            try {
              const user = JSON.parse(userStr);
              console.log('Successfully restored auth state for user:', user.username);
              return {
                user,
                token,
                isAuthenticated: true,
                loading: false,
                initialized: true
              };
            } catch (parseError) {
              console.error('Failed to parse user data, clearing localStorage:', parseError);
              localStorage.removeItem('token');
              localStorage.removeItem('user');
              return {
                user: null,
                token: null,
                isAuthenticated: false,
                loading: false,
                initialized: true
              };
            }
          } else {
            console.log('No stored auth data found');
            return {
              user: null,
              token: null,
              isAuthenticated: false,
              loading: false,
              initialized: true
            };
          }
        } catch (storageError) {
          console.error('localStorage access error:', storageError);
          return {
            user: null,
            token: null,
            isAuthenticated: false,
            loading: false,
            initialized: true
          };
        }
      });
    },
    setLoading: (loading: boolean) => {
      update((state) => ({ ...state, loading }));
    }
  };
};

export const auth = createAuthStore();
