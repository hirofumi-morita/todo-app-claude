export interface User {
  id: number;
  email: string;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

export interface Todo {
  id: number;
  user_id: number;
  title: string;
  description: string;
  completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface TodoRequest {
  title: string;
  description: string;
  completed: boolean;
}
