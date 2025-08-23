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
    initialized: false,
  };

  const { subscribe, set, update } = writable<AuthState>(initialState);

  return {
    subscribe,
    login: (user: User, token: string) => {
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      const newState = {
        user,
        token,
        isAuthenticated: true,
        loading: false,
        initialized: true,
      };
      set(newState);
    },
    logout: () => {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      set({
        user: null,
        token: null,
        isAuthenticated: false,
        loading: false,
        initialized: true,
      });
    },
    initialize: () => {
      // Prevent multiple initializations
      update(state => {
        if (state.initialized) {
          return state;
        }
        
        // Check if we're in browser environment
        if (typeof window === 'undefined' || typeof localStorage === 'undefined') {
          return {
            user: null,
            token: null,
            isAuthenticated: false,
            loading: false,
            initialized: true,
          };
        }
        
        try {
          const token = localStorage.getItem('token');
          const userStr = localStorage.getItem('user');
          
          if (token && userStr) {
            try {
              const user = JSON.parse(userStr);
              return {
                user,
                token,
                isAuthenticated: true,
                loading: false,
                initialized: true,
              };
            } catch (parseError) {
              localStorage.removeItem('token');
              localStorage.removeItem('user');
              return {
                user: null,
                token: null,
                isAuthenticated: false,
                loading: false,
                initialized: true,
              };
            }
          } else {
            return {
              user: null,
              token: null,
              isAuthenticated: false,
              loading: false,
              initialized: true,
            };
          }
        } catch (storageError) {
          return {
            user: null,
            token: null,
            isAuthenticated: false,
            loading: false,
            initialized: true,
          };
        }
      });
    },
    setLoading: (loading: boolean) => {
      update(state => ({ ...state, loading }));
    },
  };
};

export const auth = createAuthStore();