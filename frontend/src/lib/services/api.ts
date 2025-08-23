import axios from 'axios';

const API_BASE_URL = 'http://localhost:3000/api';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export interface ApiResponse<T> {
  status: 'success' | 'error';
  message: string;
  data: T;
}

export interface User {
  id: number;
  username: string;
  email: string;
  names?: string;
  created_at: string;
  updated_at: string;
}

export interface Resource {
  id: number;
  name: string;
  description?: string;
  unit: string;
  quantity: number;
  created_at: string;
  updated_at: string;
}

export interface ResourceHistory {
  id: number;
  resource_id: number;
  action: 'CREATE' | 'UPDATE' | 'DELETE';
  user_id: number;
  old_data?: string;
  new_data?: string;
  timestamp: string;
  description: string;
  resource: Resource;
  user: User;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  names?: string;
}

export interface CreateResourceRequest {
  name: string;
  description?: string;
  unit: string;
  quantity: number;
}

// Auth API
export const authAPI = {
  login: async (data: LoginRequest) => {
    const response = await api.post<ApiResponse<{user: User, token: string}>>('/auth/login', data);
    return response.data;
  },

  register: async (data: RegisterRequest) => {
    const response = await api.post<ApiResponse<User>>('/auth/register', data);
    return response.data;
  },
};

// Users API
export const usersAPI = {
  getAll: async () => {
    const response = await api.get<ApiResponse<User[]>>('/user/');
    return response.data;
  },

  getById: async (id: number) => {
    const response = await api.get<ApiResponse<User>>(`/user/${id}`);
    return response.data;
  },

  create: async (data: RegisterRequest) => {
    const response = await api.post<ApiResponse<User>>('/user/', data);
    return response.data;
  },

  update: async (id: number, data: Partial<RegisterRequest>) => {
    const response = await api.patch<ApiResponse<User>>(`/user/${id}`, data);
    return response.data;
  },

  delete: async (id: number) => {
    const response = await api.delete<ApiResponse<void>>(`/user/${id}`);
    return response.data;
  },
};

// Validation helpers
export function validateResourceId(id: any): id is number {
  return id !== null && id !== undefined && !isNaN(Number(id)) && Number(id) > 0;
}

export function validateResource(resource: any): resource is Resource {
  return resource && 
         typeof resource === 'object' && 
         validateResourceId(resource.id) &&
         typeof resource.name === 'string' &&
         typeof resource.unit === 'string' &&
         typeof resource.quantity === 'number';
}

// Resources API
export const resourcesAPI = {
  getAll: async () => {
    const response = await api.get<ApiResponse<Resource[]>>('/resource/');
    return response.data;
  },

  getById: async (id: number) => {
    const response = await api.get<ApiResponse<Resource>>(`/resource/${id}`);
    return response.data;
  },

  create: async (data: CreateResourceRequest) => {
    const response = await api.post<ApiResponse<Resource>>('/resource/', data);
    return response.data;
  },

  update: async (id: number, data: Partial<CreateResourceRequest>) => {
    const response = await api.put<ApiResponse<Resource>>(`/resource/${id}`, data);
    return response.data;
  },

  delete: async (id: number) => {
    const response = await api.delete<ApiResponse<void>>(`/resource/${id}`);
    return response.data;
  },

  getHistory: async (id: number) => {
    const response = await api.get<ApiResponse<ResourceHistory[]>>(`/resource/${id}/history`);
    return response.data;
  },
};

export default api;