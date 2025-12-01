import axios from 'axios';
import { LoginRequest, RegisterRequest, TodoRequest, LoginResponse, User, Todo } from '@/types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  register: async (data: RegisterRequest) => {
    const response = await api.post('/register', data);
    return response.data;
  },

  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/login', data);
    return response.data;
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get<User>('/me');
    return response.data;
  },
};

export const todoAPI = {
  getTodos: async (): Promise<Todo[]> => {
    const response = await api.get<Todo[]>('/todos');
    return response.data;
  },

  getTodo: async (id: number): Promise<Todo> => {
    const response = await api.get<Todo>(`/todos/${id}`);
    return response.data;
  },

  createTodo: async (data: TodoRequest): Promise<Todo> => {
    const response = await api.post<Todo>('/todos', data);
    return response.data;
  },

  updateTodo: async (id: number, data: TodoRequest): Promise<Todo> => {
    const response = await api.put<Todo>(`/todos/${id}`, data);
    return response.data;
  },

  deleteTodo: async (id: number): Promise<void> => {
    await api.delete(`/todos/${id}`);
  },
};

export const adminAPI = {
  getAllUsers: async (): Promise<User[]> => {
    const response = await api.get<User[]>('/admin/users');
    return response.data;
  },

  getUser: async (id: number): Promise<User> => {
    const response = await api.get<User>(`/admin/users/${id}`);
    return response.data;
  },

  deleteUser: async (id: number): Promise<void> => {
    await api.delete(`/admin/users/${id}`);
  },

  updateUserRole: async (id: number, isAdmin: boolean): Promise<User> => {
    const response = await api.put<User>(`/admin/users/${id}/role`, { is_admin: isAdmin });
    return response.data;
  },

  getUserTodos: async (id: number): Promise<Todo[]> => {
    const response = await api.get<Todo[]>(`/admin/users/${id}/todos`);
    return response.data;
  },
};

export default api;
